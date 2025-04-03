package ip

import (
	"context"
	"time"

	"github.com/Scalingo/go-utils/logger"
)

func (m *EndpointManager) healthChecker(ctx context.Context) {
	interval := time.Duration(m.Endpoint().HealthCheckInterval) * time.Second
	if interval == 0 {
		interval = m.config.HealthCheckInterval
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

		m.sendHealthCheckResults(ctx, healthy, err)

		time.Sleep(interval)
	}
}

func (m *EndpointManager) sendHealthCheckResults(ctx context.Context, healthy bool, err error) {
	log := logger.Get(ctx)
	if healthy {
		if m.healthCheckFailingCount > 0 {
			log.Infof("Health check healthy after %v retries", m.healthCheckFailingCount)
			m.healthCheckFailingCount = 0
		}
		m.sendEvent(HealthCheckSuccessEvent)
		return
	}

	m.healthCheckFailingCount++
	if m.healthCheckFailingCount < m.config.FailCountBeforeFailover {
		log.WithField("failing_count", m.healthCheckFailingCount).WithError(err).Info("Health check failed (will be retried)")
		return
	}

	if m.healthCheckFailingCount == m.config.FailCountBeforeFailover {
		log.WithError(err).Error("Health check failed")
	}

	m.sendEvent(HealthCheckFailEvent)
}
