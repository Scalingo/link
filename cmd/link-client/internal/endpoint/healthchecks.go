package endpoint

import (
	"context"
	"net"
	"strconv"
	"strings"

	"github.com/urfave/cli/v3"

	"github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/link/v3/api"
)

func parseHealthChecks(ctx context.Context, c *cli.Command) ([]api.HealthCheck, error) {
	checks := c.StringSlice("health-check")
	result := make([]api.HealthCheck, 0, len(checks))
	for _, check := range checks {
		checkOpts := strings.Split(check, " ")
		if len(checkOpts) != 2 {
			return nil, errors.New(ctx, "invalid check format: "+check)
		}
		checkType := checkOpts[0]
		host, port, err := net.SplitHostPort(checkOpts[1])
		if err != nil {
			return nil, errors.Wrapf(ctx, err, "invalid host/port format: %s", checkOpts[1])
		}
		portI, err := strconv.Atoi(port)
		if err != nil {
			return nil, errors.Wrapf(ctx, err, "invalid port on check %s", checkOpts[1])
		}
		result = append(result, api.HealthCheck{
			Type: api.HealthCheckType(checkType),
			Host: host,
			Port: portI,
		})
	}
	return result, nil
}
