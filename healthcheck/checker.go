package healthcheck

import (
	"context"

	"github.com/Scalingo/go-philae/prober"
	"github.com/Scalingo/go-philae/tcpprobe"
	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/link/config"
	"github.com/Scalingo/link/models"
)

type Checker interface {
	IsHealthy(ctx context.Context) bool
}

type checker struct {
	prober *prober.Prober
}

func FromChecks(config config.Config, checks []models.Healthcheck) checker {
	prober := prober.NewProber()
	for _, check := range checks {
		switch check.Type {
		case models.TCPHealthCheck:
			prober.AddProbe(tcpprobe.NewTCPProbe("tcp", check.Addr(), tcpprobe.TCPOptions{
				Timeout: config.HealthcheckTimeout,
			}))
		}
	}
	return checker{
		prober: prober,
	}
}

func (c checker) IsHealthy(ctx context.Context) bool {
	log := logger.Get(ctx)

	res := c.prober.Check(ctx)
	if res.Error != nil {
		log.WithError(res.Error).Error("Healthcheck failed")

	}
	return res.Healthy
}
