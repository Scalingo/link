package healthcheck

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Scalingo/go-philae/v4/prober"
	"github.com/Scalingo/go-philae/v4/sampleprobe"
	"github.com/Scalingo/link/v2/api"
	"github.com/Scalingo/link/v2/config"
	"github.com/Scalingo/link/v2/models"
)

func TestFromChecks(t *testing.T) {
	examples := []struct {
		Name           string
		Checks         []models.HealthCheck
		ExpectedChecks []string
	}{
		{
			Name:           "With no checks",
			Checks:         []models.HealthCheck{},
			ExpectedChecks: []string{},
		}, {
			Name: "With a tcp check",
			Checks: []models.HealthCheck{
				{
					Type: api.TCPHealthCheck,
				},
			},
			ExpectedChecks: []string{"tcp"},
		},
	}

	for _, example := range examples {
		t.Run(example.Name, func(t *testing.T) {
			ctx := context.Background()
			c := config.Config{
				HealthCheckTimeout: 10 * time.Millisecond,
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
		Name           string
		Probes         []prober.Probe
		ExpectedResult bool
		ExpectedError  string
	}{
		{
			Name:           "With no probe configured",
			Probes:         []prober.Probe{},
			ExpectedResult: true,
		}, {
			Name: "With only failing probes",
			Probes: []prober.Probe{
				sampleprobe.NewSampleProbe("test", false),
				sampleprobe.NewSampleProbe("test-2", false),
			},
			ExpectedResult: false,
			ExpectedError:  "error",
		}, {
			Name: "With failing and valid probes",
			Probes: []prober.Probe{
				sampleprobe.NewSampleProbe("test", false),
				sampleprobe.NewSampleProbe("test-2", true),
			},
			ExpectedResult: false,
			ExpectedError:  "error",
		}, {
			Name: "With only valid probes",
			Probes: []prober.Probe{
				sampleprobe.NewSampleProbe("test", true),
				sampleprobe.NewSampleProbe("test-2", true),
			},
			ExpectedResult: true,
		},
	}

	for _, example := range examples {
		t.Run(example.Name, func(t *testing.T) {
			prober := prober.NewProber()

			for _, probe := range example.Probes {
				prober.AddProbe(probe)
			}

			checker := HeathChecker{
				prober: prober,
			}

			healthy, err := checker.IsHealthy(context.Background())
			assert.Equal(t, example.ExpectedResult, healthy)
			if example.ExpectedError != "" {
				assert.Contains(t, err.Error(), example.ExpectedError)
			}
		})
	}
}
