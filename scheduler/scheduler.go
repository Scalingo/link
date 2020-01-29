package scheduler

import (
	"context"
	"sync"

	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/link/api"
	"github.com/Scalingo/link/config"
	"github.com/Scalingo/link/ip"
	"github.com/Scalingo/link/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
)

type Scheduler interface {
	Start(context.Context, models.IP) (models.IP, error)
	Stop(ctx context.Context, id string) error
	Status(string) string
	ConfiguredIPs(ctx context.Context) []api.IP
	GetIP(ctx context.Context, id string) *api.IP
	TryGetLock(ctx context.Context, id string) bool
}

type IPScheduler struct {
	mapMutex   sync.RWMutex
	ipManagers map[string]ip.Manager
	etcd       *clientv3.Client
	config     config.Config
	storage    models.Storage
}

var (
	ErrNotStopping = errors.New("not stopping")
)

func NewIPScheduler(config config.Config, etcd *clientv3.Client, storage models.Storage) *IPScheduler {
	return &IPScheduler{
		mapMutex:   sync.RWMutex{},
		ipManagers: make(map[string]ip.Manager),
		etcd:       etcd,
		config:     config,
		storage:    storage,
	}
}

func (s *IPScheduler) Status(id string) string {
	s.mapMutex.RLock()
	defer s.mapMutex.RUnlock()
	manager, ok := s.ipManagers[id]
	if ok {
		return manager.Status()
	}
	return ""
}

func (s *IPScheduler) Start(ctx context.Context, ipAddr models.IP) (models.IP, error) {
	log := logger.Get(ctx)
	newIP, err := s.storage.AddIP(ctx, ipAddr)
	if err != nil && errors.Cause(err) != models.ErrIPAlreadyPresent {
		return newIP, errors.Wrap(err, "fail to add IP to storage")
	}
	log = log.WithFields(logrus.Fields{
		"ip": newIP.IP,
		"id": newIP.ID,
	})
	ctx = logger.ToCtx(ctx, log)
	ipAdded := (err == nil)

	s.mapMutex.RLock()
	manager, ok := s.ipManagers[newIP.ID]
	s.mapMutex.RUnlock()
	log.WithFields(logrus.Fields{
		"ip_added":      ipAdded,
		"manager_found": ok,
	}).Debug("")
	// If the interface has the IP, it might be in stopping state. We just want to cancel the
	// stopping
	if ok {
		if manager.CancelStopping(ctx) {
			return newIP, nil
		}
		return newIP, ErrNotStopping
	}
	log.Info("Initialize a new IP manager")

	manager, err = ip.NewManager(ctx, s.config, newIP, s.etcd, s.storage)
	if err != nil {
		if ipAdded {
			err := s.storage.RemoveIP(ctx, newIP.ID)
			if err != nil {
				log.WithError(err).Error("fail to remove IP from storage after failed intialization of IP manager")
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

func (s *IPScheduler) Stop(ctx context.Context, id string) error {
	s.mapMutex.RLock()
	manager, ok := s.ipManagers[id]
	s.mapMutex.RUnlock()
	if !ok {
		return errors.New("not found")
	}

	manager.Stop(ctx, func(ctx context.Context) error {
		s.mapMutex.Lock()
		defer s.mapMutex.Unlock()
		err := s.storage.RemoveIP(ctx, id)
		if err != nil {
			return errors.Wrap(err, "fail to remove IP from storage")
		}
		delete(s.ipManagers, id)
		return nil
	})
	return nil
}

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

func (s *IPScheduler) TryGetLock(ctx context.Context, id string) bool {
	s.mapMutex.RLock()
	manager, ok := s.ipManagers[id]
	s.mapMutex.RUnlock()
	if !ok {
		return false
	}

	manager.TryGetLock(ctx)
	return true
}
