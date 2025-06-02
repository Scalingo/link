package endpoint

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v3"

	"github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/link/v3/cmd/link-client/internal/utils"
)

func List(ctx context.Context, c *cli.Command) error {
	client := utils.GetClient(c)
	endpoints, err := client.ListEndpoints(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "list endpoints")
	}

	if len(endpoints) == 0 {
		fmt.Println("No endpoints configured.")
		return nil
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Election Key", "Status", "Plugin", "CHECKS"})

	for _, endpoint := range endpoints {
		status := FormatStatus(endpoint)

		checks := "None"
		if len(endpoint.Checks) > 0 {
			var c []string
			for _, check := range endpoint.Checks {
				c = append(c, fmt.Sprintf("%s - %s", check.Type, net.JoinHostPort(check.Host, strconv.Itoa(check.Port))))
			}

			checks = strings.Join(c, ",")
		}
		table.Append([]string{
			endpoint.ID,
			endpoint.ElectionKey,
			status,
			endpoint.Plugin,
			checks,
		})
	}
	table.Render()

	return nil
}
