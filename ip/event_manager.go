package ip

import (
	"context"
	"time"

	"github.com/Scalingo/go-utils/logger"
)

func (m *manager) Stop(ctx context.Context) {
	m.stopMutex.Lock()
	defer m.stopMutex.Unlock()
	m.stopping = true
}

func (m *manager) isStopping() bool {
	m.stopMutex.RLock()
	defer m.stopMutex.RUnlock()
	return m.stopping
}

func (m *manager) eventManager(ctx context.Context, eventChan chan string) {
	log := logger.Get(ctx).WithField("process", "event_manager")
	for {
		if m.isStopping() {
			// Sleeping twice the lease time will ensure that we've lost our lease and another node was elected MASTER.
			// So after this sleep, we can safely remove our IP.
			log.Info("Stop order received, waiting %s to remove IP", (2 * m.config.LeaseTime()).String())
			time.Sleep(2 * m.config.LeaseTime())
			if m.stateMachine.Current() != FAILING {
				eventChan <- DemotedEvent
			}
			log.Info("Stopping!")
			close(eventChan)
			return
		}

		if m.stateMachine.Current() != FAILING {
			m.singleEtcdRun(ctx, eventChan)
		}

		time.Sleep(m.config.KeepAliveInterval)
	}
}

func (m *manager) singleEtcdRun(ctx context.Context, eventChan chan string) {
	log := logger.Get(ctx)
	err := m.locker.Refresh(ctx)
	if err != nil {
		log.WithError(err).Error("Fail to refresh lock")
		log.Info("FAULT")
		eventChan <- FaultEvent
		return
	}

	isMaster, err := m.locker.IsMaster(ctx)
	if err != nil {
		log.WithError(err).Error("Fail to check lock")
		eventChan <- FaultEvent
	} else {
		if isMaster {
			eventChan <- ElectedEvent
		} else {
			eventChan <- DemotedEvent
		}
	}
}

func (m *manager) healthChecker(ctx context.Context, eventChan chan string) {
	for {
		healthy := m.checker.IsHealthy()

		// The eventManager will close the chan when we receive the Stop order and we do not want to send things on a close channel.
		// Since the checker can take up to 5s to run his checks, this check must be done between the health check and sending the results.
		if m.isStopping() {
			return
		}

		if healthy {
			eventChan <- HealthCheckSuccessEvent
		} else {
			eventChan <- HealthCheckFailEvent
		}
		time.Sleep(m.config.HealthcheckInterval)
	}
}
