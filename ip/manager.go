package ip

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/looplab/fsm"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/go-utils/retry"
	"github.com/Scalingo/link/v2/config"
	"github.com/Scalingo/link/v2/healthcheck"
	"github.com/Scalingo/link/v2/locker"
	"github.com/Scalingo/link/v2/models"
	"github.com/Scalingo/link/v2/network"
	"github.com/Scalingo/link/v2/watcher"
)

type Manager interface {
	Start(context.Context)
	Stop(context.Context) error
	Failover(context.Context) error
	Status() string
	Endpoint() models.Endpoint
	SetHealthChecks(context.Context, config.Config, []models.HealthCheck)
}

type manager struct {
	networkInterface        network.NetworkInterface
	stateMachine            *fsm.FSM
	endpoint                models.Endpoint
	stopMutex               sync.RWMutex
	locker                  locker.Locker
	checker                 healthcheck.Checker
	checkerMutex            sync.RWMutex
	config                  config.Config
	storage                 models.Storage
	watcher                 watcher.Watcher
	retry                   retry.Retry
	eventChan               chan string
	keepaliveRetry          int
	healthCheckFailingCount int
	stopped                 bool
}

func NewManager(ctx context.Context, cfg config.Config, endpoint models.Endpoint, client *clientv3.Client, storage models.Storage, leaseManager locker.EtcdLeaseManager) (*manager, error) {
	i, err := network.NewNetworkInterfaceFromName(cfg.Interface)
	if err != nil {
		return nil, errors.Wrap(err, "fail to instantiate network interface")
	}

	log := logger.Get(ctx).WithFields(logrus.Fields{
		"ip": endpoint.IP,
	})
	ctx = logger.ToCtx(ctx, log)

	m := &manager{
		networkInterface:        i,
		endpoint:                endpoint,
		locker:                  locker.NewEtcdLocker(cfg, client, leaseManager, endpoint),
		checker:                 healthcheck.FromChecks(cfg, endpoint.Checks),
		config:                  cfg,
		storage:                 storage,
		eventChan:               make(chan string),
		healthCheckFailingCount: 0,
		retry:                   retry.New(retry.WithWaitDuration(10*time.Second), retry.WithMaxAttempts(5)),
	}

	prefix := fmt.Sprintf("%s/ips/%s", models.EtcdLinkDirectory, endpoint.StorableIP())
	m.watcher = watcher.NewWatcher(client, prefix, m.onTopologyChange)

	m.stateMachine = NewStateMachine(ctx, NewStateMachineOpts{
		ActivatedCallback: m.setActivated,
		StandbyCallback:   m.setStandBy,
		FailingCallback:   m.setFailing,
	})
	return m, nil
}

func (m *manager) Start(ctx context.Context) {
	log := logger.Get(ctx).WithFields(m.endpoint.ToLogrusFields())
	log.Info("Starting manager")

	err := m.retry.Do(ctx, func(ctx context.Context) error {
		return m.storage.LinkEndpointWithCurrentHost(ctx, m.endpoint)
	})
	if err != nil {
		log.WithError(err).Error("Fail to link endpoint")
	}

	ctx = logger.ToCtx(ctx, log)
	go m.ipCheckLoop(ctx)    // Will continuously try to get the IP
	go m.healthChecker(ctx)  // Healthchecker
	go m.startArpEnsure(ctx) // ARP Gratuitous announces
	go m.watcher.Start(ctx)  // Start a watcher that will notify us if other hosts are joining or leaving this IP

	for event := range m.eventChan {
		err := m.stateMachine.Event(ctx, event)
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

// Status returns the current state of the state machine
func (m *manager) Status() string {
	return m.stateMachine.Current()
}

// IP returns the ip model linked to this manager
func (m *manager) Endpoint() models.Endpoint {
	return m.endpoint
}

// sendEvent sends an event to the state machine
func (m *manager) sendEvent(status string) {
	if m.isStopped() {
		return
	}
	m.eventChan <- status
}

func (m *manager) SetHealthChecks(ctx context.Context, cfg config.Config, healthchecks []models.HealthCheck) {
	log := logger.Get(ctx)
	log.WithField("healtchchecks", healthchecks).Debug("Set new healthchecks")

	m.endpoint.Checks = healthchecks
	m.checkerMutex.Lock()
	m.checker = healthcheck.FromChecks(cfg, healthchecks)
	m.checkerMutex.Unlock()
}
