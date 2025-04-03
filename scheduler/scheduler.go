package scheduler

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	etcdv3 "go.etcd.io/etcd/client/v3"

	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/link/v2/config"
	"github.com/Scalingo/link/v2/ip"
	"github.com/Scalingo/link/v2/locker"
	"github.com/Scalingo/link/v2/models"
	"github.com/Scalingo/link/v2/plugin"
)

var (
	// ErrEndpointAlreadyAssigned can be sent by Start if there is another endpoint with the same election key
	ErrEndpointAlreadyAssigned = errors.New("An endpoint with the same election key already exists on that host")

	// ErrEndpointNotFound can be sent if an operation has been called on an unregistered Endpoint
	ErrEndpointNotFound = errors.New("Endpoint not found")
)

// Scheduler is the central point of LinK it will keep track all of IPs registered on this node
// however the heavy lifting for a single IP is done in the Manager
type Scheduler interface {
	Start(ctx context.Context, endpoint models.Endpoint) (models.Endpoint, error)
	Stop(ctx context.Context, id string) error
	Failover(ctx context.Context, id string) error
	Status(id string) string
	ConfiguredEndpoints(ctx context.Context) EndpointsWithStatus
	GetEndpoint(ctx context.Context, id string) *EndpointWithStatus
	UpdateEndpoint(ctx context.Context, endpoint models.Endpoint) error
}

// IPScheduler is LinK implementation of the Scheduler Interface
type IPScheduler struct {
	mapMutex         sync.RWMutex
	endpointManagers map[string]ip.Manager
	etcd             *etcdv3.Client
	config           config.Config
	storage          models.Storage
	leaseManager     locker.EtcdLeaseManager
	pluginRegistry   plugin.Registry
}

// NewIPScheduler creates and configures a Scheduler
func NewIPScheduler(config config.Config, etcd *etcdv3.Client, storage models.Storage, leaseManager locker.EtcdLeaseManager, registry plugin.Registry) *IPScheduler {
	return &IPScheduler{
		mapMutex:         sync.RWMutex{},
		endpointManagers: make(map[string]ip.Manager),
		etcd:             etcd,
		config:           config,
		storage:          storage,
		leaseManager:     leaseManager,
		pluginRegistry:   registry,
	}
}

// Status gives you access to the state machine status of a specific IP
func (s *IPScheduler) Status(id string) string {
	s.mapMutex.RLock()
	defer s.mapMutex.RUnlock()
	manager, ok := s.endpointManagers[id]
	if ok {
		return manager.Status()
	}
	return ""
}

// Start schedules a new IP on the host. It launches a new manager for the IP and add it to the tracked IP on this host.
func (s *IPScheduler) Start(ctx context.Context, endpoint models.Endpoint) (models.Endpoint, error) {
	log := logger.Get(ctx)

	log.Info("Initialize Endpoint Plugin")

	plugin, err := s.pluginRegistry.Create(ctx, endpoint)
	if err != nil {
		return endpoint, errors.Wrap(err, "initialize plugin")
	}

	for _, manager := range s.endpointManagers {
		if manager.ElectionKey(ctx) == plugin.ElectionKey(ctx) {
			return endpoint, ErrEndpointAlreadyAssigned
		}
	}

	log.Info("Initialize a new IP manager")

	manager, err := ip.NewManager(ctx, s.config, endpoint, s.etcd, s.storage, s.leaseManager, plugin)
	if err != nil {
		return endpoint, errors.Wrap(err, "fail to initialize manager")
	}

	s.mapMutex.Lock()
	s.endpointManagers[endpoint.ID] = manager
	s.mapMutex.Unlock()
	go manager.Start(ctx)

	return endpoint, nil
}

// Stop the manager of the specified IP and remove it from the tracked IP
func (s *IPScheduler) Stop(ctx context.Context, id string) error {
	s.mapMutex.RLock()
	manager, ok := s.endpointManagers[id]
	s.mapMutex.RUnlock()
	if !ok {
		return ErrEndpointNotFound
	}

	err := manager.Stop(ctx)
	if err != nil {
		return errors.Wrap(err, "fail to stop manager")
	}

	s.mapMutex.Lock()
	defer s.mapMutex.Unlock()
	delete(s.endpointManagers, id)
	return nil
}

// Failover triggers a failover on a specific IP
func (s *IPScheduler) Failover(ctx context.Context, id string) error {
	s.mapMutex.RLock()
	manager, ok := s.endpointManagers[id]
	s.mapMutex.RUnlock()
	if !ok {
		return ErrEndpointNotFound
	}

	err := manager.Failover(ctx)
	if err != nil {
		return errors.Wrapf(err, "fail to failover the IP %v", id)
	}

	return nil
}

// ConfiguredIPs lists all IPs currently tracked by the scheduler
func (s *IPScheduler) ConfiguredEndpoints(ctx context.Context) EndpointsWithStatus {
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

// GetIP fetches basic information about a tracked IP
func (s *IPScheduler) GetEndpoint(ctx context.Context, id string) *EndpointWithStatus {
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

// UpdateIP updates the IP in the scheduler storage, and the health checks in the IP manager.
func (s *IPScheduler) UpdateEndpoint(ctx context.Context, endpoint models.Endpoint) error {
	log := logger.Get(ctx)
	s.mapMutex.RLock()
	manager, ok := s.endpointManagers[endpoint.ID]
	s.mapMutex.RUnlock()
	if !ok {
		log.Info("IP manager not found, skipping the IP update")
		return nil
	}

	err := s.storage.UpdateEndpoint(ctx, endpoint)
	if err != nil {
		return errors.Wrap(err, "fail to update the IP from storage")
	}

	manager.SetHealthChecks(ctx, s.config, endpoint.Checks)

	return nil
}
