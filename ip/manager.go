package ip

import (
	"context"
	"sync"

	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/link/config"
	"github.com/Scalingo/link/network"
	"github.com/coreos/etcd/clientv3"
	"github.com/looplab/fsm"
	"github.com/pkg/errors"
)

type Manager interface {
	Start(context.Context)
	Stop(context.Context)
}

type manager struct {
	networkInterface network.NetworkInterface
	stateMachine     *fsm.FSM
	ip               string
	etcd             *clientv3.Client
	stopMutex        sync.RWMutex
	stopping         bool
}

func NewManager(ctx context.Context, config config.Config, ip string, client *clientv3.Client) (*manager, error) {
	i, err := network.NewNetworkInterfaceFromName(config.Interface)
	if err != nil {
		return nil, errors.Wrap(err, "fail to instantiate network interface")
	}

	log := logger.Get(ctx).WithField("ip", ip)
	ctx = logger.ToCtx(ctx, log)

	m := &manager{
		ip:               ip,
		networkInterface: i,
		etcd:             client,
	}

	m.newStateMachine(ctx)
	return m, nil
}

func (m *manager) setActivated(ctx context.Context) {
	log := logger.Get(ctx)
	log.Info("New state: ACTIVATED")
	err := m.networkInterface.AddIP(m.ip)
	if err != nil {
		log.WithError(err).Error("Fail to activate IP")
	}
}

func (m *manager) setStandBy(ctx context.Context) {
	log := logger.Get(ctx)
	log.Info("New state: STANDBY")
	err := m.networkInterface.RemoveIP(m.ip)
	if err != nil {
		log.WithError(err).Error("Fail to de-activate IP")
	}
}

func (m *manager) Start(ctx context.Context) {
	log := logger.Get(ctx).WithField("ip", m.ip)
	log.Info("Starting manager")

	ctx = logger.ToCtx(ctx, log)
	eventChan := make(chan string)

	go m.lockRoutine(ctx, eventChan)

	for event := range eventChan {
		err := m.stateMachine.Event(event)
		if err != nil {
			// Ignore NoTransitionError since those just means that we did not change state (which can be normal)
			if _, ok := err.(fsm.NoTransitionError); !ok {
				log.WithError(err).Info("INVALID STATE MACHINE TRANSITION")
				panic("STATE MACHINE HAD SOME ISSUE, STOP!!")
			}
		}
	}
	log.Info("Manager stopped")
}

func (m *manager) Stop(ctx context.Context) {
	m.stopMutex.Lock()
	defer m.stopMutex.Unlock()
	m.stopping = true
}
