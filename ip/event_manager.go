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
			log.Info("Stop order received, waiting 10s to remove IP")
			time.Sleep(10 * time.Second)
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

		time.Sleep(3 * time.Second)
	}
}

func (m *manager) singleEtcdRun(ctx context.Context, eventChan chan string) {
	log := logger.Get(ctx)
	err := m.locker.Refresh(ctx)
	if err != nil {
		log.WithError(err).Error("Fail to refresh lock")
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
		check := m.checker.Check()
		if m.isStopping() {
			return
		}

		if check {
			eventChan <- HealthCheckSuccessEvent
		} else {
			eventChan <- HealthCheckFailEvent
		}
		time.Sleep(30 * time.Second)
	}
}
