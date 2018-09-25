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
	IsHealthy() bool
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

func (c checker) IsHealthy() bool {
	ctx := context.Background()

	// Custom logger to discard Philae output
	log := logrus.New()
	log.Out = ioutil.Discard
	ctx = logger.ToCtx(ctx, log)

	res := c.prober.Check(ctx)
	return res.Healthy
}
