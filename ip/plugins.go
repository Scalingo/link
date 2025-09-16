package ip

import (
	"context"
	"time"

	"github.com/looplab/fsm"

	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/link/v3/config"
	"github.com/Scalingo/link/v3/locker"
)

func (m *EndpointManager) setActivated(ctx context.Context, _ *fsm.Event) {
	log := logger.Get(ctx)
	log.Info("New state: ACTIVATED")
	err := m.plugin.Activate(ctx)
	if err != nil {
		log.WithError(err).Error("Fail to activate endpoint")
	}
}

func (m *EndpointManager) setStandBy(ctx context.Context, _ *fsm.Event) {
	log := logger.Get(ctx)
	log.Info("New state: STANDBY")
	err := m.plugin.Deactivate(ctx)
	if err != nil {
		log.WithError(err).Error("Fail to de-activate endpoint")
	}
}

func (m *EndpointManager) setFailing(ctx context.Context, _ *fsm.Event) {
	log := logger.Get(ctx)
	log.Info("New state: FAILING")

	err := m.plugin.Deactivate(ctx)
	if err != nil {
		log.WithError(err).Error("Fail to de-activate IP")
	}

	err = m.locker.Unlock(ctx)
	if err != nil && err != locker.ErrNotMaster {
		// If we are not master, we can safely ignore this error
		log.WithError(err).Error("Fail to unlock the key")
	}
}

func (m *EndpointManager) startPluginEnsureLoop(ctx context.Context) {
	log := logger.Get(ctx).WithField("process", "plugin_ensure")
	for {
		if m.isStopped() {
			return
		}
		currentState := m.Status()

		hasFailed := false

		if currentState == ACTIVATED {
			log.Debug("Start plugin ensure")
			err := m.plugin.Ensure(ctx)
			if err != nil {
				log.WithError(err).Error("Fail to run plugin ensure")
				hasFailed = true
			}
		}

		timeToSleep := config.RandomDurationAround(m.config.PluginEnsureInterval, 0.25)
		if hasFailed {
			timeToSleep = m.ensureBackoff.NextBackOff()
			log.WithField("next_try_in", timeToSleep).Info("Ensure failed, applying backoff")
		} else {
			m.ensureBackoff.Reset()
		}

		time.Sleep(timeToSleep)
	}
}
