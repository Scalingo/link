package ip

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/looplab/fsm"
	etcdv3 "go.etcd.io/etcd/client/v3"

	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/go-utils/retry"
	"github.com/Scalingo/link/v2/config"
	"github.com/Scalingo/link/v2/healthcheck"
	"github.com/Scalingo/link/v2/locker"
	"github.com/Scalingo/link/v2/models"
	"github.com/Scalingo/link/v2/plugin"
	"github.com/Scalingo/link/v2/watcher"
)

type Manager interface {
	Start(ctx context.Context)
	Stop(ctx context.Context) error
	Failover(ctx context.Context) error
	Status() string
	Endpoint() models.Endpoint
	ElectionKey(ctx context.Context) string
	SetHealthChecks(ctx context.Context, config config.Config, checks []models.HealthCheck)
}

type EndpointManager struct {
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
	plugin                  plugin.Plugin
	eventChan               chan string
	keepaliveRetry          int
	healthCheckFailingCount int
	stopped                 bool
}

func NewManager(ctx context.Context, cfg config.Config, endpoint models.Endpoint, client *etcdv3.Client, storage models.Storage, leaseManager locker.EtcdLeaseManager, plugin plugin.Plugin) (*EndpointManager, error) {
	log := logger.Get(ctx).WithFields(endpoint.ToLogrusFields())
	ctx = logger.ToCtx(ctx, log)

	m := &EndpointManager{
		endpoint:                endpoint,
		locker:                  locker.NewEtcdLocker(cfg, client, leaseManager, endpoint, plugin.ElectionKey(ctx)),
		checker:                 healthcheck.FromChecks(cfg, endpoint.Checks),
		config:                  cfg,
		storage:                 storage,
		eventChan:               make(chan string),
		healthCheckFailingCount: 0,
		retry:                   retry.New(retry.WithWaitDuration(10*time.Second), retry.WithMaxAttempts(5)),
		plugin:                  plugin,
	}

	prefix := fmt.Sprintf("%s/ips/%s", models.EtcdLinkDirectory, plugin.ElectionKey(ctx))
	m.watcher = watcher.NewWatcher(client, prefix, m.onTopologyChange)

	m.stateMachine = NewStateMachine(ctx, NewStateMachineOpts{
		ActivatedCallback: m.setActivated,
		StandbyCallback:   m.setStandBy,
		FailingCallback:   m.setFailing,
	})
	return m, nil
}

func (m *EndpointManager) Start(ctx context.Context) {
	log := logger.Get(ctx).WithFields(m.endpoint.ToLogrusFields())
	log.Info("Starting manager")

	err := m.retry.Do(ctx, func(ctx context.Context) error {
		return m.storage.LinkEndpointWithCurrentHost(ctx, m.plugin.ElectionKey(ctx))
	})
	if err != nil {
		log.WithError(err).Error("Fail to link endpoint")
	}

	ctx = logger.ToCtx(ctx, log)
	go m.ipCheckLoop(ctx)   // Will continuously try to get the IP
	go m.healthChecker(ctx) // HealthChecker
	go m.startPluginEnsureLoop(ctx)
	go m.watcher.Start(ctx) // Start a watcher that will notify us if other hosts are joining or leaving this IP

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
func (m *EndpointManager) Status() string {
	return m.stateMachine.Current()
}

// IP returns the ip model linked to this manager
func (m *EndpointManager) Endpoint() models.Endpoint {
	return m.endpoint
}

func (m *EndpointManager) ElectionKey(ctx context.Context) string {
	return m.plugin.ElectionKey(ctx)
}

// sendEvent sends an event to the state machine
func (m *EndpointManager) sendEvent(status string) {
	if m.isStopped() {
		return
	}
	m.eventChan <- status
}

func (m *EndpointManager) SetHealthChecks(ctx context.Context, cfg config.Config, healthChecks []models.HealthCheck) {
	log := logger.Get(ctx)
	log.WithField("health_checks", healthChecks).Debug("Set new health checks")

	m.endpoint.Checks = healthChecks
	m.checkerMutex.Lock()
	m.checker = healthcheck.FromChecks(cfg, healthChecks)
	m.checkerMutex.Unlock()
}
