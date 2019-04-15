package healthcheck

import (
	"context"
	"testing"
	"time"

	"github.com/Scalingo/go-philae/prober"
	"github.com/Scalingo/go-philae/sampleprobe"
	"github.com/Scalingo/link/config"
	"github.com/Scalingo/link/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromChecks(t *testing.T) {
	examples := []struct {
		Name           string
		Checks         []models.Healthcheck
		ExpectedChecks []string
	}{
		{
			Name:           "With no checks",
			Checks:         []models.Healthcheck{},
			ExpectedChecks: []string{},
		}, {
			Name: "With a tcp check",
			Checks: []models.Healthcheck{
				{
					Type: models.TCPHealthCheck,
				},
			},
			ExpectedChecks: []string{"tcp"},
		},
	}

	for _, example := range examples {
		t.Run(example.Name, func(t *testing.T) {
			ctx := context.Background()
			c := config.Config{
				HealthcheckTimeout: 10 * time.Millisecond,
			}
			checker := FromChecks(c, example.Checks)
			probes := checker.prober.Check(ctx).Probes

			require.Equal(t, len(example.ExpectedChecks), len(probes))

			for i := 0; i < len(example.ExpectedChecks); i++ {
				assert.Equal(t, example.ExpectedChecks[i], probes[i].Name)
			}
		})
	}
}

func TestIsHealthy(t *testing.T) {
	examples := []struct {
		Name          string
		Probes        []prober.Probe
		ExpetedResult bool
	}{
		{
			Name:          "With no probe configured",
			Probes:        []prober.Probe{},
			ExpetedResult: true,
		}, {
			Name: "With only failing probes",
			Probes: []prober.Probe{
				sampleprobe.NewSampleProbe("test", false),
				sampleprobe.NewSampleProbe("test-2", false),
			},
			ExpetedResult: false,
		}, {
			Name: "With failing and valid probes",
			Probes: []prober.Probe{
				sampleprobe.NewSampleProbe("test", false),
				sampleprobe.NewSampleProbe("test-2", true),
			},
			ExpetedResult: false,
		}, {
			Name: "With only valid probes",
			Probes: []prober.Probe{
				sampleprobe.NewSampleProbe("test", true),
				sampleprobe.NewSampleProbe("test-2", true),
			},
			ExpetedResult: true,
		},
	}

	for _, example := range examples {
		t.Run(example.Name, func(t *testing.T) {
			prober := prober.NewProber()

			for _, probe := range example.Probes {
				prober.AddProbe(probe)
			}

			checker := checker{
				prober: prober,
			}

			assert.Equal(t, example.ExpetedResult, checker.IsHealthy(context.Background()))

		})
	}

}
