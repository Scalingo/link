package ip

import (
	"context"
	"time"

	"github.com/Scalingo/go-utils/logger"
)

func (m *manager) Stop(ctx context.Context, stopper func(context.Context) error) {
	log := logger.Get(ctx)
	m.stopMutex.Lock()
	m.stopper = stopper
	m.stopMutex.Unlock()

	if m.stateMachine.Current() == ACTIVATED {
		err := m.locker.Unlock(ctx)
		if err != nil {
			log.WithError(err).Error("fail to release etcd lease")
			return
		}
	}
}

func (m *manager) CancelStopping(ctx context.Context) bool {
	log := logger.Get(ctx)
	if !m.isStopping() {
		log.Debug("Do not cancel stopping of a non-stopping IP")
		return false
	}
	log.Info("Cancel manager stopping")

	m.stopMutex.Lock()
	defer m.stopMutex.Unlock()
	m.stopper = nil
	return true
}

func (m *manager) TryGetLock(ctx context.Context) {
	m.singleEtcdRun(ctx)
}

func (m *manager) isStopping() bool {
	m.stopMutex.RLock()
	defer m.stopMutex.RUnlock()
	return m.stopper != nil
}

func (m *manager) isStopped() bool {
	m.stopMutex.RLock()
	defer m.stopMutex.RUnlock()
	return m.stopped
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
	for {
		shouldContinue := m.singleEventRun(ctx)
		if !shouldContinue {
			return
		}

		time.Sleep(m.config.KeepAliveInterval)
	}
}

func (m *manager) singleEventRun(ctx context.Context) bool {
	log := logger.Get(ctx).WithField("process", "event_manager")
	if m.isStopping() {
		// Sleeping twice the lease time will ensure that we've lost our lease and another node was elected MASTER.
		// So after this sleep, we can safely remove our IP.

		log.Infof("Stop order received, waiting %s to remove IP", (2 * m.config.LeaseTime()).String())
		m.waitTwiceLeaseTimeOrReallocation(ctx)
		if m.stopOrder(ctx) {
			log.Infof("Removing IP %s", m.ip.IP)
			err := m.networkInterface.RemoveIP(m.ip.IP)
			if err != nil {
				log.WithError(err).Error("fail to remove IP")
			}
			return false
		}
		log.Info("Stop order has been cancelled")
	}

	if m.stateMachine.Current() != FAILING {
		m.singleEtcdRun(ctx)
	}

	return true
}

func (m *manager) waitTwiceLeaseTimeOrReallocation(ctx context.Context) {
	log := logger.Get(ctx)
	timer := time.NewTimer(2 * m.config.LeaseTime())
	defer timer.Stop()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timer.C:
			return
		case <-ticker.C:
			log.Debug("tick wait twice lease time")
			if !m.isStopping() {
				return
			}
			master, err := m.locker.IsMaster(ctx)
			// This can return a key not found error
			// This is likely to happen during re-election and is perfectly normal
			if err == nil && !master {
				log.Debug("Someone else took the lock, beginning premature shutdown")
				return
			}
		}
	}
}

// stopOrder actually handles the stopping. It returns true if it has been stopped, false
// otherwise. It can happen if the current manager stopping has been cancelled.
func (m *manager) stopOrder(ctx context.Context) bool {
	m.stopMutex.Lock()
	defer m.stopMutex.Unlock()

	log := logger.Get(ctx)

	// The stopping might have been cancelled during the two lease time sleep. We execute the
	// stopper function only if it is still in stopping state
	if m.stopper != nil {
		log.Info("Stopping!")
		if m.stateMachine.Current() != FAILING {
			m.sendEvent(DemotedEvent)
		}
		err := m.stopper(ctx)
		if err != nil {
			log.WithError(err).Error("fail to execute the stopper function")
		}
		m.closeEventChan()

		m.stopped = true
		return true
	}
	return false
}

func (m *manager) singleEtcdRun(ctx context.Context) {
	log := logger.Get(ctx)
	err := m.locker.Refresh(ctx)
	if err != nil {
		m.keepaliveRetry++
		log.WithError(err).Info("Fail to refresh lock (retry)")
		if m.keepaliveRetry > m.config.KeepAliveRetry {
			log.WithError(err).Error("Fail to refresh lock")
			m.sendEvent(FaultEvent)
		}
		return
	}
	if m.keepaliveRetry > 0 {
		log.Infof("Lock refreshed after %v retries", m.keepaliveRetry)
		m.keepaliveRetry = 0
	}

	isMaster, err := m.locker.IsMaster(ctx)
	if err != nil {
		log.WithError(err).Error("Fail to check lock")
		// Fault event should only be handled in the Refresh operation where the
		// retry loop is present
		return
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
		healthy := m.checker.IsHealthy(ctx)

		// The eventManager will close the chan when we receive the Stop order and we do not want to send things on a close channel.
		// Since the checker can take up to 5s to run his checks, this check must be done between the health check and sending the results.
		if m.isStopped() {
			return
		}

		if !m.isStopping() {
			m.sendHealthcheckResults(ctx, healthy)
		}

		time.Sleep(m.config.HealthcheckInterval)
	}
}

func (m *manager) sendHealthcheckResults(ctx context.Context, healthy bool) {
	log := logger.Get(ctx)
	if healthy {
		if m.failingCount > 0 {
			log.Infof("healthcheck healthy after %v retries", m.failingCount)
			m.failingCount = 0
		}
		m.sendEvent(HealthCheckSuccessEvent)
	} else {
		m.failingCount++
		log.WithField("failing_count", m.failingCount).Info("healthcheck failed (retry)")
		if m.failingCount >= m.config.FailCountBeforeFailover {
			log.Error("healthcheck failed")
			m.sendEvent(HealthCheckFailEvent)
		}
	}
}
