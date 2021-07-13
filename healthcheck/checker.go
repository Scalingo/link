package healthcheck

import (
	"context"

	"github.com/Scalingo/go-philae/v4/prober"
	"github.com/Scalingo/go-philae/v4/tcpprobe"
	"github.com/Scalingo/link/config"
	"github.com/Scalingo/link/models"
)

type Checker interface {
	IsHealthy(ctx context.Context) (bool, error)
}

type checker struct {
	prober *prober.Prober
}

func FromChecks(cfg config.Config, checks []models.Healthcheck) checker {
	prober := prober.NewProber()
	for _, check := range checks {
		switch check.Type {
		case models.TCPHealthCheck:
			prober.AddProbe(tcpprobe.NewTCPProbe("tcp", check.Addr(), tcpprobe.TCPOptions{
				Timeout: cfg.HealthcheckTimeout,
			}))
		}
	}
	return checker{
		prober: prober,
	}
}

func (c checker) IsHealthy(ctx context.Context) (bool, error) {
	res := c.prober.Check(ctx)
	if res.Error != nil {
		return false, res.Error
	}
	return res.Healthy, nil
}
