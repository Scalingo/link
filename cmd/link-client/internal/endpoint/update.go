package endpoint

import (
	"context"
	"fmt"

	"github.com/logrusorgru/aurora/v3"
	"github.com/urfave/cli/v3"

	"github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/link/v2/api"
	"github.com/Scalingo/link/v2/cmd/link-client/internal/utils"
)

func UpdateChecks(ctx context.Context, c *cli.Command) error {
	client := utils.GetClient(c)
	endpointID := c.String("endpoint-id")
	if endpointID == "" {
		return cli.Exit("endpoint-id is required", 1)
	}
	checks, err := parseHealthChecks(ctx, c)
	if err != nil {
		return errors.Wrap(ctx, err, "parse health checks")
	}
	endpoint, err := client.UpdateEndpoint(ctx, endpointID, api.UpdateEndpointParams{
		HealthChecks: checks,
	})
	if err != nil {
		return errors.Wrap(ctx, err, "update endpoint")
	}

	fmt.Println(aurora.Green(fmt.Sprintf("Health checks of the Endpoint %s (%s) successfully updated", endpoint.ID, endpoint.ElectionKey)))

	return nil
}
