package ip

import (
	"context"
	"time"

	"github.com/Scalingo/go-utils/logger"
)

func (m *manager) healthChecker(ctx context.Context) {
	interval := time.Duration(m.IP().HealthcheckInterval) * time.Second
	if interval == 0 {
		interval = m.config.HealthcheckInterval
	}

	for {
		m.checkerMutex.RLock()
		healthy, err := m.checker.IsHealthy(ctx)
		m.checkerMutex.RUnlock()

		// The eventManager closes the channel `eventChan` when we receive the Stop order. We do not want to send anything on a closed channel.
		// Since the checker can take up to 5s to run the checks, we need to check the manager stopped status between the health check and sending the results.
		if m.isStopped() {
			return
		}

		m.sendHealthcheckResults(ctx, healthy, err)

		time.Sleep(interval)
	}
}

func (m *manager) sendHealthcheckResults(ctx context.Context, healthy bool, err error) {
	log := logger.Get(ctx)
	if healthy {
		if m.healthcheckFailingCount > 0 {
			log.Infof("Healthcheck healthy after %v retries", m.healthcheckFailingCount)
			m.healthcheckFailingCount = 0
		}
		m.sendEvent(HealthCheckSuccessEvent)
		return
	}

	m.healthcheckFailingCount++
	if m.healthcheckFailingCount < m.config.FailCountBeforeFailover {
		log.WithField("failing_count", m.healthcheckFailingCount).WithError(err).Info("Healthcheck failed (will be retried)")
		return
	}

	if m.healthcheckFailingCount == m.config.FailCountBeforeFailover {
		log.WithError(err).Error("Healthcheck failed")
	}

	m.sendEvent(HealthCheckFailEvent)
}
