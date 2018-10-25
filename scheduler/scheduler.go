package scheduler

import (
	"context"
	"sync"

	"github.com/Scalingo/link/api"
	"github.com/Scalingo/link/config"
	"github.com/Scalingo/link/ip"
	"github.com/Scalingo/link/models"
	"github.com/coreos/etcd/clientv3"
	"github.com/pkg/errors"
	errgo "gopkg.in/errgo.v1"
)

type Scheduler interface {
	Start(context.Context, models.IP) error
	Stop(ctx context.Context, id string, stopper func(context.Context) error) error
	CancelStopping(context.Context, string) error
	Status(string) string
	ConfiguredIPs(ctx context.Context) []api.IP
	GetIP(ctx context.Context, id string) *api.IP
	TryGetLock(ctx context.Context, id string) bool
}

type IPScheduler struct {
	mapMutex   sync.Mutex
	ipManagers map[string]ip.Manager
	etcd       *clientv3.Client
	config     config.Config
}

var (
	ErrNotStopping = errors.New("not stopping")
)

func NewIPScheduler(config config.Config, etcd *clientv3.Client) *IPScheduler {
	return &IPScheduler{
		mapMutex:   sync.Mutex{},
		ipManagers: make(map[string]ip.Manager),
		etcd:       etcd,
		config:     config,
	}
}

func (s *IPScheduler) Status(id string) string {
	s.mapMutex.Lock()
	defer s.mapMutex.Unlock()
	manager, ok := s.ipManagers[id]
	if ok {
		return manager.Status()
	}
	return ""
}

func (s *IPScheduler) Start(ctx context.Context, ipAddr models.IP) error {
	manager, err := ip.NewManager(ctx, s.config, ipAddr, s.etcd)
	if err != nil {
		return errors.Wrap(err, "fail to initialize manager")
	}

	s.mapMutex.Lock()
	s.ipManagers[ipAddr.ID] = manager
	s.mapMutex.Unlock()
	go manager.Start(ctx)

	return nil
}

func (s *IPScheduler) Stop(ctx context.Context, id string, stopper func(context.Context) error) error {
	s.mapMutex.Lock()
	defer s.mapMutex.Unlock()

	manager, ok := s.ipManagers[id]
	if !ok {
		return errors.New("not found")
	}

	manager.Stop(ctx, func(ctx context.Context) error {
		err := stopper(ctx)
		if err != nil {
			return errors.Wrapf(err, "fail to stop the scheduler")
		}
		delete(s.ipManagers, id)
		return nil
	})
	return nil
}

func (s *IPScheduler) CancelStopping(ctx context.Context, id string) error {
	s.mapMutex.Lock()
	defer s.mapMutex.Unlock()

	manager, ok := s.ipManagers[id]
	if !ok {
		return errgo.Notef(ErrNotStopping, "not found")
	}

	manager.CancelStopping(ctx)
	return nil
}

func (s *IPScheduler) ConfiguredIPs(ctx context.Context) []api.IP {
	s.mapMutex.Lock()
	defer s.mapMutex.Unlock()

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
	s.mapMutex.Lock()
	defer s.mapMutex.Unlock()

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
	s.mapMutex.Lock()
	defer s.mapMutex.Unlock()

	manager, ok := s.ipManagers[id]
	if !ok {
		return false
	}

	manager.TryGetLock(ctx)
	return true
}
