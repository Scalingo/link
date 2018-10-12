package ip

import (
	"context"
	"time"

	"github.com/Scalingo/go-utils/logger"
)

func (m *manager) Stop(ctx context.Context) {
	log := logger.Get(ctx)
	m.stopMutex.Lock()
	defer m.stopMutex.Unlock()
	m.stopping = true
	if m.stateMachine.Current() == ACTIVATED {
		err := m.locker.Unlock(ctx)
		if err != nil {
			log.WithError(err).Error("fail to release etcd lease")
			return
		}
	}

}

func (m *manager) TryGetLock(ctx context.Context) {
	m.singleEtcdRun(ctx)
}

func (m *manager) isStopping() bool {
	m.stopMutex.RLock()
	defer m.stopMutex.RUnlock()
	return m.stopping
}

func (m *manager) sendEvent(status string) {
	m.messageMutex.Lock()
	defer m.messageMutex.Unlock()
	if m.closed {
		return
	}
	m.eventChan <- status
}

func (m *manager) closeEventChan() {
	m.messageMutex.Lock()
	defer m.messageMutex.Unlock()
	m.closed = true
	close(m.eventChan)
}

func (m *manager) eventManager(ctx context.Context) {
	log := logger.Get(ctx).WithField("process", "event_manager")
	for {
		if m.isStopping() {
			// Sleeping twice the lease time will ensure that we've lost our lease and another node was elected MASTER.
			// So after this sleep, we can safely remove our IP.

			log.Infof("Stop order received, waiting %s to remove IP", (2 * m.config.LeaseTime()).String())
			time.Sleep(2 * m.config.LeaseTime())
			if m.stateMachine.Current() != FAILING {
				m.sendEvent(DemotedEvent)
			}
			log.Info("Stopping!")
			m.closeEventChan()
			return
		}

		if m.stateMachine.Current() != FAILING {
			m.singleEtcdRun(ctx)
		}

		time.Sleep(m.config.KeepAliveInterval)
	}
}

func (m *manager) singleEtcdRun(ctx context.Context) {
	log := logger.Get(ctx)
	err := m.locker.Refresh(ctx)
	if err != nil {
		log.WithError(err).Error("Fail to refresh lock")
		log.Info("FAULT")
		m.sendEvent(FaultEvent)
		return
	}

	isMaster, err := m.locker.IsMaster(ctx)
	if err != nil {
		log.WithError(err).Error("Fail to check lock")
		m.sendEvent(FaultEvent)
	} else {
		if isMaster {
			m.sendEvent(ElectedEvent)
		} else {
			m.sendEvent(DemotedEvent)
		}
	}
}

func (m *manager) healthChecker(ctx context.Context) {
	for {
		healthy := m.checker.IsHealthy()

		// The eventManager will close the chan when we receive the Stop order and we do not want to send things on a close channel.
		// Since the checker can take up to 5s to run his checks, this check must be done between the health check and sending the results.
		if m.isStopping() {
			return
		}

		if healthy {
			m.sendEvent(HealthCheckSuccessEvent)
		} else {
			m.sendEvent(HealthCheckFailEvent)
		}

		time.Sleep(m.config.HealthcheckInterval)
	}
}
