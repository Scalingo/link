package endpoint

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/link/v3/cmd/link-client/internal/utils"
)

func Show(ctx context.Context, c *cli.Command) error {
	endpointID := c.String("endpoint-id")
	if endpointID == "" {
		return cli.Exit("endpoint-id is required", 1)
	}

	client := utils.GetClient(c)
	endpoint, err := client.GetEndpoint(ctx, endpointID)
	if err != nil {
		return errors.Wrap(ctx, err, "get endpoint")
	}

	fmt.Printf("ID:\t\t%s\n", endpoint.ID)
	fmt.Printf("Status:\t\t%s\n", FormatStatus(endpoint))
	fmt.Printf("Election Key: \t%s\n", endpoint.ElectionKey)
	fmt.Printf("Plugin:\t\t%s\n", endpoint.Plugin)
	if len(endpoint.Checks) == 0 {
		fmt.Printf("Checks:\t\tNone\n")
	} else {
		fmt.Println("Checks:")
		for _, check := range endpoint.Checks {
			fmt.Printf(" - Type: %s, Host: %s, Port: %v\n", check.Type, check.Host, check.Port)
		}
	}

	hosts, err := client.GetEndpointHosts(ctx, endpointID)
	if err != nil {
		return errors.Wrap(ctx, err, "get endpoint hosts")
	}
	if len(hosts) == 0 {
		fmt.Println("This endpoint is configured on no hosts.")
		return nil
	}

	fmt.Println("Hosts:")
	for _, host := range hosts {
		fmt.Printf(" - %s\n", host.Hostname)
	}

	return nil
}
