package scheduler

import (
	"context"
	"sync"

	"github.com/Scalingo/link/config"
	"github.com/Scalingo/link/ip"
	"github.com/Scalingo/link/models"
	"github.com/coreos/etcd/clientv3"
	"github.com/pkg/errors"
)

type Scheduler interface {
	Start(context.Context, models.IP) error
	Stop(context.Context, string) error
	Status(string) string
}

type IPScheduler struct {
	mapMutex   sync.Mutex
	ipManagers map[string]ip.Manager
	etcd       *clientv3.Client
	config     config.Config
}

func NewIPScheduler(config config.Config, etcd *clientv3.Client) IPScheduler {
	return IPScheduler{
		mapMutex:   sync.Mutex{},
		ipManagers: make(map[string]ip.Manager),
		etcd:       etcd,
		config:     config,
	}
}

func (s IPScheduler) Status(id string) string {
	s.mapMutex.Lock()
	defer s.mapMutex.Unlock()
	manager, ok := s.ipManagers[id]
	if ok {
		return manager.Status()
	}
	return ""
}

func (s IPScheduler) Start(ctx context.Context, ipAddr models.IP) error {
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

func (s IPScheduler) Stop(ctx context.Context, id string) error {
	s.mapMutex.Lock()
	defer s.mapMutex.Unlock()

	manager, ok := s.ipManagers[id]
	if !ok {
		return errors.New("Not found")
	}

	delete(s.ipManagers, id)

	manager.Stop(ctx)
	return nil
}
