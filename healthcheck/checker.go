package healthcheck

import (
	"context"

	"github.com/Scalingo/go-philae/v4/prober"
	"github.com/Scalingo/go-philae/v4/tcpprobe"
	"github.com/Scalingo/link/v3/api"
	"github.com/Scalingo/link/v3/config"
	"github.com/Scalingo/link/v3/models"
)

type Checker interface {
	IsHealthy(ctx context.Context) (bool, error)
}

type HeathChecker struct {
	prober *prober.Prober
}

func FromChecks(cfg config.Config, checks []models.HealthCheck) HeathChecker {
	prober := prober.NewProber()
	for _, check := range checks {
		if check.Type == api.TCPHealthCheck {
			prober.AddProbe(tcpprobe.NewTCPProbe("tcp", check.Addr(), tcpprobe.TCPOptions{
				Timeout: cfg.HealthCheckTimeout,
			}))
		}
	}
	return HeathChecker{
		prober: prober,
	}
}

func (c HeathChecker) IsHealthy(ctx context.Context) (bool, error) {
	res := c.prober.Check(ctx)
	if res.Error != nil {
		return false, res.Error
	}
	return res.Healthy, nil
}
