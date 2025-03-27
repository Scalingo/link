package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/logrusorgru/aurora/v3"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/link/v2/api"
)

var Version = "dev"

func main() {
	app := cli.NewApp()

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
			Name: "list",
			Action: func(c *cli.Context) error {
				client := getClientFromCtx(c)
				ips, err := client.ListEndpoints(context.Background())
				if err != nil {
					return err
				}
				formatEndpoints(ips)
				return nil
			},
		}, {
			Name:      "destroy",
			ArgsUsage: "ID",
			Action: func(c *cli.Context) error {
				if c.NArg() != 1 {
					err := cli.ShowCommandHelp(c, c.Command.Name)
					if err != nil {
						return errors.Wrap(err, "show destroy command helper")
					}
					return nil
				}
				client := getClientFromCtx(c)
				err := client.RemoveEndpoint(context.Background(), c.Args().First())
				if err != nil {
					return err
				}
				fmt.Println(aurora.Green(fmt.Sprintf("Endpoint %v deleted.", c.Args().First())))
				return nil
			},
		}, {
			Name:      "get",
			ArgsUsage: "ID",
			Action: func(c *cli.Context) error {
				if c.NArg() != 1 {
					err := cli.ShowCommandHelp(c, c.Command.Name)
					if err != nil {
						return errors.Wrap(err, "show get command helper")
					}
					return nil
				}
				client := getClientFromCtx(c)
				endpoints, err := client.GetEndpoint(context.Background(), c.Args().First())
				if err != nil {
					return err
				}
				formatEndpoint(endpoints)
				return nil
			},
		}, {
			Name:      "failover",
			ArgsUsage: "ID",
			Action: func(c *cli.Context) error {
				if c.NArg() != 1 {
					err := cli.ShowCommandHelp(c, c.Command.Name)
					if err != nil {
						return errors.Wrap(err, "show failover command helper")
					}
					return nil
				}
				client := getClientFromCtx(c)
				err := client.Failover(context.Background(), c.Args().First())
				if err != nil {
					return err
				}
				fmt.Println(aurora.Green("Request sent."))
				return nil
			},
		}, {
			Name:      "add",
			ArgsUsage: "IP [CHECK_TYPE CHECK_ENDPOINT]...",
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:  "healthcheck-interval",
					Value: 0,
					Usage: "Duration between health checks",
				},
			},
			Action: func(c *cli.Context) error {
				if c.NArg()%2 == 0 {
					// 1 For the IP
					// And 2 per Health checks
					// So NArgs % 2 must be == 1
					err := cli.ShowCommandHelp(c, c.Command.Name)
					if err != nil {
						return errors.Wrap(err, "show add command helper")
					}
					return nil
				}
				client := getClientFromCtx(c)
				ip := c.Args().First()
				var checks []api.HealthCheck
				curArg := 1
				for curArg < c.NArg() {
					endpoint := c.Args().Get(curArg + 1)
					host, port, err := net.SplitHostPort(endpoint)
					if err != nil {
						return fmt.Errorf("invalid endpoint: %s", endpoint)
					}
					portI, err := strconv.Atoi(port)
					if err != nil {
						return fmt.Errorf("invalid endpoint: %s", endpoint)
					}
					checks = append(checks, api.HealthCheck{
						Type: api.HealthCheckType(c.Args().Get(curArg)),
						Host: host,
						Port: portI,
					})
					curArg += 2
				}

				params := api.AddEndpointParams{
					Checks:              checks,
					HealthCheckInterval: c.Int("health check-interval"),
					IP:                  ip,
				}
				newEndpoint, err := client.AddEndpoint(context.Background(), params)
				if err != nil {
					return err
				}

				fmt.Println(aurora.Green(fmt.Sprintf("Endpoint %s (%s) successfully added", newEndpoint.ID, newEndpoint.IP)))
				return nil
			},
		}, {
			Name:      "update-healthchecks",
			ArgsUsage: "ID [CHECK_TYPE CHECK_ENDPOINT]...",
			Action: func(c *cli.Context) error {
				if c.NArg()%2 == 0 {
					// 1 For the IP
					// And 2 per Health checks
					// So NArgs % 2 must be == 1
					err := cli.ShowCommandHelp(c, c.Command.Name)
					if err != nil {
						return errors.Wrap(err, "show update-healthchecks command helper")
					}
					return nil
				}

				var healthChecks []api.HealthCheck
				curArg := 1
				for curArg < c.NArg() {
					healthCheckType := c.Args().Get(curArg)
					endpoint := c.Args().Get(curArg + 1)
					host, port, err := net.SplitHostPort(endpoint)
					if err != nil {
						return fmt.Errorf("invalid health check endpoint: %s", endpoint)
					}
					portI, err := strconv.Atoi(port)
					if err != nil {
						return fmt.Errorf("invalid health check port: %s", port)
					}
					healthChecks = append(healthChecks, api.HealthCheck{
						Type: api.HealthCheckType(healthCheckType),
						Host: host,
						Port: portI,
					})
					curArg += 2
				}

				linkEndpointID := c.Args().First()

				client := getClientFromCtx(c)
				ip, err := client.UpdateEndpoint(context.Background(),
					linkEndpointID, api.UpdateEndpointParams{HealthChecks: healthChecks},
				)
				if err != nil {
					return errors.Wrapf(err, "update the endpoint health checks '%s'", linkEndpointID)
				}

				fmt.Println(aurora.Green(fmt.Sprintf("Health checks of the Endpoint %s (%s) successfully updated", ip.IP, ip.ID)))
				return nil
			},
		}, {
			Name: "version",
			Action: func(c *cli.Context) error {
				fmt.Printf("Client version: \t%s\n", app.Version)
				client := getClientFromCtx(c)
				version, err := client.Version(context.Background())
				if err != nil {
					version = aurora.Red(err.Error()).String()
				}
				fmt.Printf("Server version: \t%s\n", version)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(aurora.Red(fmt.Sprintf("Error: %s", err.Error())))
	}
}

func getClientFromCtx(c *cli.Context) api.HTTPClient {
	var opts []api.ClientOpt
	if c.String("host") != "" {
		opts = append(opts, api.WithURL(c.String("host")))
	}

	if c.String("user") != "" {
		opts = append(opts, api.WithUser(c.String("user")))
	}

	if c.String("password") != "" {
		opts = append(opts, api.WithPassword(c.String("password")))
	}

	return api.NewHTTPClient(opts...)
}
