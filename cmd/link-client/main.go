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
	"github.com/Scalingo/link/v2/models"
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
				ips, err := client.ListIPs(context.Background())
				if err != nil {
					return err
				}
				formatIPs(ips)
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
				err := client.RemoveIP(context.Background(), c.Args().First())
				if err != nil {
					return err
				}
				fmt.Println(aurora.Green(fmt.Sprintf("IP %v deleted.", c.Args().First())))
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
				ip, err := client.GetIP(context.Background(), c.Args().First())
				if err != nil {
					return err
				}
				formatIP(ip)
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
					Usage: "Duration between healthchecks",
				},
			},
			Action: func(c *cli.Context) error {
				if c.NArg()%2 == 0 {
					// 1 For the IP
					// And 2 per Healthchecks
					// So NArgs % 2 must be == 1
					err := cli.ShowCommandHelp(c, c.Command.Name)
					if err != nil {
						return errors.Wrap(err, "show add command helper")
					}
					return nil
				}
				client := getClientFromCtx(c)
				ip := c.Args().First()
				var checks []models.Healthcheck
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
					checks = append(checks, models.Healthcheck{
						Type: models.HealthcheckType(c.Args().Get(curArg)),
						Host: host,
						Port: portI,
					})
					curArg += 2
				}

				params := api.AddIPParams{
					Checks:              checks,
					HealthcheckInterval: c.Int("healthcheck-interval"),
				}
				newIP, err := client.AddIP(context.Background(), ip, params)
				if err != nil {
					return err
				}

				fmt.Println(aurora.Green(fmt.Sprintf("IP %s (%s) successfully added", newIP.IP.IP, newIP.ID)))
				return nil
			},
		}, {
			Name:      "update-healthchecks",
			ArgsUsage: "ID [CHECK_TYPE CHECK_ENDPOINT]...",
			Action: func(c *cli.Context) error {
				if c.NArg()%2 == 0 {
					// 1 For the IP
					// And 2 per Healthchecks
					// So NArgs % 2 must be == 1
					err := cli.ShowCommandHelp(c, c.Command.Name)
					if err != nil {
						return errors.Wrap(err, "show update-healthchecks command helper")
					}
					return nil
				}

				var healthchecks []models.Healthcheck
				curArg := 1
				for curArg < c.NArg() {
					healthcheckType := c.Args().Get(curArg)
					endpoint := c.Args().Get(curArg + 1)
					host, port, err := net.SplitHostPort(endpoint)
					if err != nil {
						return fmt.Errorf("invalid healthcheck endpoint: %s", endpoint)
					}
					portI, err := strconv.Atoi(port)
					if err != nil {
						return fmt.Errorf("invalid healthcheck port: %s", port)
					}
					healthchecks = append(healthchecks, models.Healthcheck{
						Type: models.HealthcheckType(healthcheckType),
						Host: host,
						Port: portI,
					})
					curArg += 2
				}

				linkIPId := c.Args().First()

				client := getClientFromCtx(c)
				ip, err := client.UpdateIP(context.Background(),
					linkIPId, api.UpdateIPParams{Healthchecks: healthchecks},
				)
				if err != nil {
					return errors.Wrapf(err, "fail to update the IP healthchecks '%s'", linkIPId)
				}

				fmt.Println(aurora.Green(fmt.Sprintf("Healthchecks of the IP %s (%s) successfully updated", ip.IP.IP, ip.ID)))
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
