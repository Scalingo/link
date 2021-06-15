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
		healthy, err := m.checker.IsHealthy(ctx)

		// The eventManager will close the chan when we receive the Stop order and we do not want to send things on a close channel.
		// Since the checker can take up to 5s to run his checks, this check must be done between the health check and sending the results.
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
		if m.failingCount > 0 {
			log.Infof("healthcheck healthy after %v retries", m.failingCount)
			m.failingCount = 0
		}
		m.sendEvent(HealthCheckSuccessEvent)
		return
	}

	m.failingCount++
	if m.failingCount < m.config.FailCountBeforeFailover {
		log.WithField("failing_count", m.failingCount).WithError(err).Info("healthcheck failed (will be retried)")
		return
	}

	if m.failingCount == m.config.FailCountBeforeFailover {
		log.WithError(err).Error("healthcheck failed")
	}

	m.sendEvent(HealthCheckFailEvent)
}
