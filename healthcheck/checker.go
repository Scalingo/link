package healthcheck

import (
	"time"

	"github.com/Scalingo/go-philae/prober"
	"github.com/Scalingo/go-philae/tcpprobe"
	"github.com/Scalingo/link/models"
)

type Checker interface {
	Check() bool
}

type checker struct {
	prober *prober.Prober
}

func FromChecks(checks []models.Healthcheck) checker {
	prober := prober.NewProber()
	for _, check := range checks {
		switch check.Type {
		case models.TCPHealthCheck:
			prober.AddProbe(tcpprobe.NewTCPProbe("tcp", check.Addr(), tcpprobe.TCPOptions{
				Timeout: 5 * time.Second,
			}))
		}
	}
	return checker{
		prober: prober,
	}
}

func (c checker) Check() bool {
	res := c.prober.Check()
	return res.Healthy
}
