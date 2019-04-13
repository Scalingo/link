package healthcheck

import (
	"context"
	"io/ioutil"

	"github.com/Scalingo/go-philae/prober"
	"github.com/Scalingo/go-philae/tcpprobe"
	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/link/config"
	"github.com/Scalingo/link/models"
	"github.com/sirupsen/logrus"
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
	appLogger := logger.Get(ctx)

	// Custom logger to discard Philae output
	log := logrus.New()
	log.Out = ioutil.Discard
	philaeCtx := logger.ToCtx(context.Background(), log)

	res := c.prober.Check(philaeCtx)
	if !res.Healthy {
		var reasons []string
		for _, probe := range res.Probes {
			if !probe.Healthy {
				reasons = append(reasons, probe.Comment)
			}
		}
		appLogger.WithField("reasons", reasons).Error("Probe failed")
	}
	return res.Healthy
}
