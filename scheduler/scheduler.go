package scheduler

import (
	"context"
	"sync"

	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/link/api"
	"github.com/Scalingo/link/config"
	"github.com/Scalingo/link/ip"
	"github.com/Scalingo/link/locker"
	"github.com/Scalingo/link/models"
	"github.com/pkg/errors"
	"go.etcd.io/etcd/v3/clientv3"
)

var (
	// ErrIPAlreadyAssigned can be sent by AddIP if the IP has already been assigned to this scheduler
	ErrIPAlreadyAssigned = errors.New("IP already assigned")

	// ErrIPNotFound can be sent if an operation has been called on an unregistered IP
	ErrIPNotFound = errors.New("IP not found")
)

// Scheduler is the central point of LinK it will keep track all of IPs registered on this node
// however the heavy lifting for a single IP is done in the Manager
type Scheduler interface {
	Start(context.Context, models.IP) (models.IP, error)
	Stop(ctx context.Context, id string) error
	Failover(ctx context.Context, id string) error
	Status(string) string
	ConfiguredIPs(ctx context.Context) []api.IP
	GetIP(ctx context.Context, id string) *api.IP
	UpdateIP(context.Context, models.IP) error
}

// IPScheduler is LinK implementation of the Scheduler Interface
type IPScheduler struct {
	mapMutex     sync.RWMutex
	ipManagers   map[string]ip.Manager
	etcd         *clientv3.Client
	config       config.Config
	storage      models.Storage
	leaseManager locker.EtcdLeaseManager
}

// NewIPScheduler creates and configures a Scheduler
func NewIPScheduler(config config.Config, etcd *clientv3.Client, storage models.Storage, leaseManager locker.EtcdLeaseManager) *IPScheduler {
	return &IPScheduler{
		mapMutex:     sync.RWMutex{},
		ipManagers:   make(map[string]ip.Manager),
		etcd:         etcd,
		config:       config,
		storage:      storage,
		leaseManager: leaseManager,
	}
}

// Status gives you access to the state machine status of a specific IP
func (s *IPScheduler) Status(id string) string {
	s.mapMutex.RLock()
	defer s.mapMutex.RUnlock()
	manager, ok := s.ipManagers[id]
	if ok {
		return manager.Status()
	}
	return ""
}

// Start schedules a new IP on the host. It launches a new manager for the IP and add it to the tracked IP on this host.
func (s *IPScheduler) Start(ctx context.Context, ipAddr models.IP) (models.IP, error) {
	log := logger.Get(ctx)
	newIP, err := s.storage.AddIP(ctx, ipAddr)
	if err != nil {
		if errors.Cause(err) != models.ErrIPAlreadyPresent {
			return newIP, errors.Wrap(err, "fail to add IP to storage")
		} else if ipAddr.ID == "" { // If the ID is not empty we are in the startup process
			return newIP, ErrIPAlreadyAssigned
		}
	}
	log = log.WithFields(newIP.ToLogrusFields())
	ctx = logger.ToCtx(ctx, log)
	ipAdded := (err == nil)

	log.Info("Initialize a new IP manager")

	manager, err := ip.NewManager(ctx, s.config, newIP, s.etcd, s.storage, s.leaseManager)
	if err != nil {
		if ipAdded {
			err := s.storage.RemoveIP(ctx, newIP.ID)
			if err != nil {
				log.WithError(err).Error("Fail to remove IP from storage after failed initialization of IP manager")
			}
		}
		return newIP, errors.Wrap(err, "fail to initialize manager")
	}

	s.mapMutex.Lock()
	s.ipManagers[newIP.ID] = manager
	s.mapMutex.Unlock()
	go manager.Start(ctx)

	return newIP, nil
}

// Stop the manager of the specified IP and remove it from the tracked IP
func (s *IPScheduler) Stop(ctx context.Context, id string) error {
	s.mapMutex.RLock()
	manager, ok := s.ipManagers[id]
	s.mapMutex.RUnlock()
	if !ok {
		return ErrIPNotFound
	}

	err := manager.Stop(ctx)
	if err != nil {
		return errors.Wrap(err, "fail to stop manager")
	}

	err = s.storage.RemoveIP(ctx, id)
	if err != nil {
		return errors.Wrap(err, "fail to remove IP from storage")
	}

	s.mapMutex.Lock()
	defer s.mapMutex.Unlock()
	delete(s.ipManagers, id)
	return nil
}

// Failover triggers a failover on a specific IP
func (s *IPScheduler) Failover(ctx context.Context, id string) error {
	s.mapMutex.RLock()
	manager, ok := s.ipManagers[id]
	s.mapMutex.RUnlock()
	if !ok {
		return ErrIPNotFound
	}

	err := manager.Failover(ctx)
	if err != nil {
		return errors.Wrapf(err, "fail to failover the IP %v", id)
	}

	return nil
}

// ConfiguredIPs lists all IPs currently tracked by the scheduler
func (s *IPScheduler) ConfiguredIPs(ctx context.Context) []api.IP {
	s.mapMutex.RLock()
	defer s.mapMutex.RUnlock()

	var ips []api.IP
	for _, manager := range s.ipManagers {
		ips = append(ips, api.IP{
			IP:     manager.IP(),
			Status: manager.Status(),
		})
	}
	return ips
}

// GetIP fetches basic information about a tracked IP
func (s *IPScheduler) GetIP(ctx context.Context, id string) *api.IP {
	s.mapMutex.RLock()
	defer s.mapMutex.RUnlock()

	manager, ok := s.ipManagers[id]
	if !ok {
		return nil
	}
	return &api.IP{
		IP:     manager.IP(),
		Status: manager.Status(),
	}
}

// UpdateIP updates the IP in the scheduler storage, and the healthchecks in the IP manager.
func (s *IPScheduler) UpdateIP(ctx context.Context, ip models.IP) error {
	log := logger.Get(ctx)
	s.mapMutex.RLock()
	manager, ok := s.ipManagers[ip.ID]
	s.mapMutex.RUnlock()
	if !ok {
		log.Info("IP manager not found, skipping the IP update")
		return nil
	}

	err := s.storage.UpdateIP(ctx, ip)
	if err != nil {
		return errors.Wrap(err, "fail to update the IP from storage")
	}

	manager.SetHealthchecks(ctx, s.config, ip.Checks)

	return nil
}
