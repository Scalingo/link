package endpoint

import (
	"context"
	"fmt"

	"github.com/logrusorgru/aurora/v3"
	"github.com/urfave/cli/v3"

	"github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/link/v3/api"
	"github.com/Scalingo/link/v3/cmd/link-client/internal/utils"
	"github.com/Scalingo/link/v3/plugin/arp"
	outscalepublicip "github.com/Scalingo/link/v3/plugin/outscale_public_ip"
)

func Create(ctx context.Context, c *cli.Command) error {
	params := api.AddEndpointParams{
		HealthCheckInterval: c.Int("health-check-interval"),
		Plugin:              c.String("plugin"),
	}

	var pluginConfig any
	var err error

	switch params.Plugin {
	case arp.Name:
		pluginConfig, err = getArpPluginConfig(ctx, c)
	case outscalepublicip.Name:
		pluginConfig, err = getOutscalePublicIPPluginConfig(ctx, c)
	default:
		err = fmt.Errorf("plugin %s not supported", params.Plugin)
	}
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

func getOutscalePublicIPPluginConfig(ctx context.Context, c *cli.Command) (outscalepublicip.PluginConfig, error) {
	publicIPID := c.String("public-ip-id")
	if publicIPID == "" {
		return outscalepublicip.PluginConfig{}, errors.New(ctx, "public-ip-id is required for outscale public ip plugin")
	}

	nicID := c.String("nic-id")
	if nicID == "" {
		return outscalepublicip.PluginConfig{}, errors.New(ctx, "nic-id is required for outscale public ip plugin")
	}
	region := c.String("region")
	if region == "" {
		return outscalepublicip.PluginConfig{}, errors.New(ctx, "region is required for outscale public ip plugin")
	}
	accessKey := c.String("access-key")
	if accessKey == "" {
		return outscalepublicip.PluginConfig{}, errors.New(ctx, "access-key is required for outscale public ip plugin")
	}
	secretKey := c.String("secret-key")
	if secretKey == "" {
		return outscalepublicip.PluginConfig{}, errors.New(ctx, "secret-key is required for outscale public ip plugin")
	}

	return outscalepublicip.PluginConfig{
		PublicIPID: publicIPID,
		NICID:      nicID,
		Region:     region,
		AccessKey:  accessKey,
		SecretKey:  secretKey,
	}, nil
}
