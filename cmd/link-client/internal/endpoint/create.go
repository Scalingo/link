package endpoint

import (
	"context"
	"fmt"

	"github.com/logrusorgru/aurora/v3"
	"github.com/urfave/cli/v3"

	"github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/link/v2/api"
	"github.com/Scalingo/link/v2/cmd/link-client/internal/utils"
	"github.com/Scalingo/link/v2/plugin/arp"
)

func Create(ctx context.Context, c *cli.Command) error {
	params := api.AddEndpointParams{
		HealthCheckInterval: int(c.Int("health-check-interval")),
		Plugin:              c.String("plugin"),
	}

	var pluginConfig any
	var err error

	if params.Plugin != arp.Name {
		return cli.Exit(fmt.Sprintf("Plugin %s does not exist", params.Plugin), 1)
	}
	pluginConfig, err = getArpPluginConfig(ctx, c)
	if err != nil {
		return cli.Exit(errors.Wrap(ctx, err, "invalid plugin config").Error(), 1)
	}
	params.PluginConfig = pluginConfig

	checks, err := parseHealthChecks(ctx, c)
	if err != nil {
		return errors.Wrap(ctx, err, "parse health checks")
	}
	params.Checks = checks

	client := utils.GetClient(c)
	newEndpoint, err := client.AddEndpoint(ctx, params)
	if err != nil {
		return errors.Wrap(ctx, err, "create endpoint")
	}

	fmt.Println(aurora.Green(fmt.Sprintf("Endpoint %s (%s) successfully added", newEndpoint.ID, newEndpoint.ElectionKey)))

	return nil
}

func getArpPluginConfig(ctx context.Context, c *cli.Command) (arp.PluginConfig, error) {
	value := c.String("ip")
	if value == "" {
		return arp.PluginConfig{}, errors.New(ctx, "ip is required for arp plugin")
	}

	return arp.PluginConfig{
		IP: value,
	}, nil
}
