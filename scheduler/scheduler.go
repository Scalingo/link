package scheduler

import (
	"context"
	"sync"

	stderrors "github.com/pkg/errors"
	etcdv3 "go.etcd.io/etcd/client/v3"

	"github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/link/v3/config"
	"github.com/Scalingo/link/v3/ip"
	"github.com/Scalingo/link/v3/locker"
	"github.com/Scalingo/link/v3/models"
	"github.com/Scalingo/link/v3/plugin"
)

var (
	// ErrEndpointAlreadyAssigned can be sent by Start if there is another endpoint with the same election key
	ErrEndpointAlreadyAssigned = stderrors.New("An endpoint with the same election key already exists on that host")

	// ErrEndpointNotFound can be sent if an operation has been called on an unregistered Endpoint
	ErrEndpointNotFound = stderrors.New("Endpoint not found")
)

// Scheduler is the central point of LinK it will keep track all of endpoints registered on this node
// however the heavy lifting for a single endpoint is done in the Manager
type Scheduler interface {
	Start(ctx context.Context, endpoint models.Endpoint) (models.Endpoint, error)
	Stop(ctx context.Context, id string) error
	Failover(ctx context.Context, id string) error
	Status(id string) string
	ConfiguredEndpoints(ctx context.Context) EndpointsWithStatus
	GetEndpoint(ctx context.Context, id string) *EndpointWithStatus
	EndpointCount() int
	UpdateEndpoint(ctx context.Context, endpoint models.Endpoint) error
}

// EndpointScheduler is LinK implementation of the Scheduler Interface
type EndpointScheduler struct {
	mapMutex         sync.RWMutex
	endpointManagers map[string]ip.Manager
	etcd             *etcdv3.Client
	config           config.Config
	storage          models.Storage
	leaseManager     locker.EtcdLeaseManager
	pluginRegistry   plugin.Registry
}

// NewEndpointScheduler creates and configures a Scheduler
func NewEndpointScheduler(config config.Config, etcd *etcdv3.Client, storage models.Storage, leaseManager locker.EtcdLeaseManager, registry plugin.Registry) *EndpointScheduler {
	return &EndpointScheduler{
		mapMutex:         sync.RWMutex{},
		endpointManagers: make(map[string]ip.Manager),
		etcd:             etcd,
		config:           config,
		storage:          storage,
		leaseManager:     leaseManager,
		pluginRegistry:   registry,
	}
}

// Status gives you access to the state machine status of a specific endpoint
func (s *EndpointScheduler) Status(id string) string {
	s.mapMutex.RLock()
	defer s.mapMutex.RUnlock()
	manager, ok := s.endpointManagers[id]
	if ok {
		return manager.Status()
	}
	return ""
}

// Start schedules a new endpoint on the host. It launches a new manager for the endpoint and add it to the tracked endpoint on this host.
func (s *EndpointScheduler) Start(ctx context.Context, endpoint models.Endpoint) (models.Endpoint, error) {
	log := logger.Get(ctx)

	log.Info("Initialize Endpoint Plugin")

	plugin, err := s.pluginRegistry.Create(ctx, endpoint)
	if err != nil {
		return endpoint, errors.Wrap(ctx, err, "initialize plugin")
	}

	ctx, log = logger.WithFieldToCtx(ctx, "election_key", plugin.ElectionKey(ctx))

	for _, manager := range s.endpointManagers {
		if manager.ElectionKey(ctx) == plugin.ElectionKey(ctx) {
			return endpoint, errors.Wrap(ctx, ErrEndpointAlreadyAssigned, "endpoint already assigned")
		}
	}

	log.Info("Initialize a new endpoint manager")

	manager, err := ip.NewManager(ctx, s.config, endpoint, s.etcd, s.storage, s.leaseManager, plugin)
	if err != nil {
		return endpoint, errors.Wrap(ctx, err, "fail to initialize manager")
	}

	s.mapMutex.Lock()
	s.endpointManagers[endpoint.ID] = manager
	s.mapMutex.Unlock()
	go manager.Start(ctx)

	return endpoint, nil
}

// Stop the manager of the specified endpoint and remove it from the tracked endpoints
func (s *EndpointScheduler) Stop(ctx context.Context, id string) error {
	s.mapMutex.RLock()
	manager, ok := s.endpointManagers[id]
	s.mapMutex.RUnlock()
	if !ok {
		return ErrEndpointNotFound
	}

	err := manager.Stop(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "fail to stop manager")
	}

	s.mapMutex.Lock()
	defer s.mapMutex.Unlock()
	delete(s.endpointManagers, id)
	return nil
}

// Failover triggers a failover on a specific endpoint
func (s *EndpointScheduler) Failover(ctx context.Context, id string) error {
	s.mapMutex.RLock()
	manager, ok := s.endpointManagers[id]
	s.mapMutex.RUnlock()
	if !ok {
		return ErrEndpointNotFound
	}

	err := manager.Failover(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to failover the endpoint %v", id)
	}

	return nil
}

// ConfiguredEndpoints lists all endpoints currently tracked by the scheduler
func (s *EndpointScheduler) ConfiguredEndpoints(ctx context.Context) EndpointsWithStatus {
	s.mapMutex.RLock()
	defer s.mapMutex.RUnlock()

	res := make(EndpointsWithStatus, 0, len(s.endpointManagers))

	for _, manager := range s.endpointManagers {
		res = append(res, EndpointWithStatus{
			Endpoint:    manager.Endpoint(),
			Status:      manager.Status(),
			ElectionKey: manager.ElectionKey(ctx),
		})
	}
	return res
}

// GetEndpoint fetches basic information about a tracked endpoint
func (s *EndpointScheduler) GetEndpoint(ctx context.Context, id string) *EndpointWithStatus {
	s.mapMutex.RLock()
	defer s.mapMutex.RUnlock()

	manager, ok := s.endpointManagers[id]
	if !ok {
		return nil
	}

	return &EndpointWithStatus{
		Endpoint:    manager.Endpoint(),
		Status:      manager.Status(),
		ElectionKey: manager.ElectionKey(ctx),
	}
}

// UpdateEndpoint updates the endpoint in the scheduler storage, and the health checks in the endpoint manager.
func (s *EndpointScheduler) UpdateEndpoint(ctx context.Context, endpoint models.Endpoint) error {
	log := logger.Get(ctx)
	s.mapMutex.RLock()
	manager, ok := s.endpointManagers[endpoint.ID]
	s.mapMutex.RUnlock()
	if !ok {
		log.Info("Endpoint manager not found, skipping the endpoint update")
		return nil
	}

	err := s.storage.UpdateEndpoint(ctx, endpoint)
	if err != nil {
		return errors.Wrap(ctx, err, "fail to update the endpoint from storage")
	}

	manager.SetHealthChecks(ctx, s.config, endpoint.Checks)

	return nil
}

func (s *EndpointScheduler) EndpointCount() int {
	s.mapMutex.RLock()
	defer s.mapMutex.RUnlock()
	return len(s.endpointManagers)
}
