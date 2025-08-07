package main

import (
	"context"
	"fmt"
	"os"

	"github.com/logrusorgru/aurora/v3"
	"github.com/urfave/cli/v3"

	"github.com/Scalingo/link/v3/cmd/link-client/internal/endpoint"
	"github.com/Scalingo/link/v3/cmd/link-client/internal/utils"
)

var Version = "dev"

func main() {
	app := cli.Command{}
	app.Version = Version

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "host",
			Value: "127.0.0.1:1313",
			Usage: "Host to contact",
		},
		&cli.StringFlag{
			Name:    "user",
			Aliases: []string{"u"},
			Value:   "",
			Usage:   "Username for basic auth",
		},
		&cli.StringFlag{
			Name:    "password",
			Aliases: []string{"p"},
			Value:   "",
			Usage:   "Password for basic auth",
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:   "list",
			Action: endpoint.List,
		}, {
			Name:    "destroy",
			Aliases: []string{"delete"},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "endpoint-id",
					Usage:    "ID of the endpoint to destroy",
					Aliases:  []string{"id", "endpoint"},
					Required: true,
				},
			},
			Action: endpoint.Destroy,
		}, {
			Name:    "show",
			Aliases: []string{"get"},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "endpoint-id",
					Aliases:  []string{"id", "endpoint"},
					Usage:    "ID of the endpoint to show",
					Required: true,
				},
			},
			Action: endpoint.Show,
		}, {
			Name: "failover",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "endpoint-id",
					Aliases:  []string{"id", "endpoint"},
					Usage:    "ID of the endpoint to failover",
					Required: true,
				},
			},
			Action: endpoint.Failover,
		}, {
			Name:    "create",
			Aliases: []string{"add"},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "plugin",
					Usage:    "Name of the plugin",
					Required: true,
				},
				// ARP Plugin
				&cli.StringFlag{
					Name:  "ip",
					Usage: "For ARP Plugin: IP to add",
				},
				// Outscale Public IP Plugin
				&cli.StringFlag{
					Name:  "public-ip-id",
					Usage: "For Outscale Public IP Plugin: ID of the public IP to add",
				},
				&cli.StringFlag{
					Name:  "nic-id",
					Usage: "For Outscale Public IP Plugin: ID of the NIC to add the public IP to",
				},
				&cli.StringFlag{
					Name:  "region",
					Usage: "For Outscale Public IP Plugin: Region of the public IP",
				},
				&cli.StringFlag{
					Name:  "access-key",
					Usage: "For Outscale Public IP Plugin: Access key for the Outscale API",
				},
				&cli.StringFlag{
					Name:  "secret-key",
					Usage: "For Outscale Public IP Plugin: Secret key for the Outscale API",
				},
				&cli.IntFlag{
					Name:  "health-check-interval",
					Value: 0,
					Usage: "Duration between health checks",
				},
				&cli.StringSliceFlag{
					Name:  "health-check",
					Usage: "Health checks to add format: [TYPE HOST:PORT]",
				},
			},
			Action: endpoint.Create,
		}, {
			Name: "set-health-checks",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "endpoint-id",
					Aliases:  []string{"id", "endpoint"},
					Usage:    "ID of the endpoint to update",
					Required: true,
				},
				&cli.StringSliceFlag{
					Name:  "health-check",
					Usage: "Health checks to add format: [TYPE HOST:PORT]",
				},
			},
			Action: endpoint.UpdateChecks,
		}, {
			Name:   "rotate-encryption-key",
			Action: RotateEncryptionKey,
		}, {
			Name: "version",
			Action: func(ctx context.Context, c *cli.Command) error {
				fmt.Printf("Client version: \t%s\n", app.Version)
				client := utils.GetClient(c)
				version, err := client.Version(ctx)
				if err != nil {
					version = aurora.Red(err.Error()).String()
				}
				fmt.Printf("Server version: \t%s\n", version)
				return nil
			},
		},
	}

	err := app.Run(context.Background(), os.Args)
	if err != nil {
		fmt.Println(aurora.Red("Error: " + err.Error()))
	}
}

func RotateEncryptionKey(ctx context.Context, c *cli.Command) error {
	client := utils.GetClient(c)

	fmt.Println("Rotating encryption key...")
	err := client.RotateEncryptionKey(ctx)
	if err != nil {
		return fmt.Errorf("rotate encryption key: %w", err)
	}

	fmt.Println(aurora.Green("Encryption key rotated successfully"))
	return nil
}
