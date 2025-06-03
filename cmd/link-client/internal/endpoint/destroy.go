package endpoint

import (
	"context"
	"fmt"

	"github.com/logrusorgru/aurora/v3"
	"github.com/urfave/cli/v3"

	"github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/link/v3/cmd/link-client/internal/utils"
)

func Destroy(ctx context.Context, c *cli.Command) error {
	client := utils.GetClient(c)
	endpointID := c.String("endpoint-id")
	if endpointID == "" {
		return cli.Exit("endpoint-id is required", 1)
	}

	err := client.RemoveEndpoint(ctx, endpointID)
	if err != nil {
		return errors.Wrap(ctx, err, "destroy endpoint")
	}

	fmt.Println(aurora.Green("Endpoint destroyed."))
	return nil
}
