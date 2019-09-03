package ip

import (
	"context"
	"sync"

	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/link/config"
	"github.com/Scalingo/link/healthcheck"
	"github.com/Scalingo/link/locker"
	"github.com/Scalingo/link/models"
	"github.com/Scalingo/link/network"
	"github.com/coreos/etcd/clientv3"
	"github.com/looplab/fsm"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Manager interface {
	Start(context.Context)
	Stop(ctx context.Context, stopper func(context.Context) error)
	CancelStopping(context.Context) bool
	Status() string
	IP() models.IP
	TryGetLock(context.Context)
}

type manager struct {
	networkInterface network.NetworkInterface
	stateMachine     *fsm.FSM
	ip               models.IP
	stopMutex        sync.RWMutex
	stopper          func(context.Context) error
	messageMutex     sync.Mutex
	closed           bool
	locker           locker.Locker
	checker          healthcheck.Checker
	config           config.Config
	eventChan        chan string
	failingCount     int
	stopped          bool
}

func NewManager(ctx context.Context, config config.Config, ip models.IP, client *clientv3.Client, storage models.Storage) (*manager, error) {
	i, err := network.NewNetworkInterfaceFromName(config.Interface)
	if err != nil {
		return nil, errors.Wrap(err, "fail to instantiate network interface")
	}

	log := logger.Get(ctx).WithFields(logrus.Fields{
		"ip": ip.IP,
	})
	ctx = logger.ToCtx(ctx, log)

	m := &manager{
		networkInterface: i,
		ip:               ip,
		locker:           locker.NewEtcdLocker(config, client, storage, ip),
		checker:          healthcheck.FromChecks(config, ip.Checks),
		config:           config,
		eventChan:        make(chan string),
		failingCount:     0,
	}

	m.stateMachine = NewStateMachine(ctx, NewStateMachineOpts{
		ActivatedCallback: m.setActivated,
		StandbyCallback:   m.setStandBy,
		FailingCallback:   m.setFailing,
	})
	return m, nil
}

func (m *manager) Start(ctx context.Context) {
	log := logger.Get(ctx).WithField("ip", m.ip.IP)
	log.Info("Starting manager")

	ctx = logger.ToCtx(ctx, log)
	go m.eventManager(ctx)
	go m.healthChecker(ctx)
	go m.startArpEnsure(ctx)

	for event := range m.eventChan {
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

func (u *manager) Status() string {
	return u.stateMachine.Current()
}

func (u *manager) IP() models.IP {
	return u.ip
}
