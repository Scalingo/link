package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/Scalingo/link/api"
	"github.com/Scalingo/link/models"
	"github.com/logrusorgru/aurora"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var Version = "dev"

func main() {
	app := cli.NewApp()

	app.Version = Version

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "host",
			Value: "127.0.0.1:1313",
			Usage: "Host to contact",
		},
		cli.StringFlag{
			Name:  "user, u",
			Value: "",
			Usage: "Username for basic auth",
		},
		cli.StringFlag{
			Name:  "password, p",
			Value: "",
			Usage: "Password for basic auth",
		},
	}

	app.Commands = []cli.Command{
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
					cli.ShowCommandHelp(c, c.Command.Name)
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
					cli.ShowCommandHelp(c, c.Command.Name)
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
			Name:      "try-get-lock",
			ArgsUsage: "ID",
			Action: func(c *cli.Context) error {
				if c.NArg() != 1 {
					cli.ShowCommandHelp(c, c.Command.Name)
					return nil
				}
				client := getClientFromCtx(c)
				err := client.TryGetLock(context.Background(), c.Args().First())
				if err != nil {
					return err
				}
				fmt.Println(aurora.Green("Request sent."))
				return nil
			},
		}, {
			Name:      "update",
			ArgsUsage: "ID [CHECK_TYPE CHECK_ENDPOINT]...",
			Action: func(c *cli.Context) error {
				if len(c.Args())%2 == 0 {
					// 1 For the IP
					// And 2 per Healthchecks
					// So NArgs % 2 must be == 1
					cli.ShowCommandHelp(c, c.Command.Name)
					return nil
				}

				var healthchecks []models.Healthcheck
				curArg := 1
				for curArg < c.NArg() {
					endpoint := c.Args().Get(curArg + 1)
					host, port, err := net.SplitHostPort(endpoint)
					if err != nil {
						return fmt.Errorf("Invalid endpoint: %s", endpoint)
					}
					portI, err := strconv.Atoi(port)
					if err != nil {
						return fmt.Errorf("Invalid endpoint: %s", endpoint)
					}
					healthchecks = append(healthchecks, models.Healthcheck{
						Type: models.HealthcheckType(c.Args().Get(curArg)),
						Host: host,
						Port: portI,
					})
					curArg += 2
				}

				id := c.Args().First()

				client := getClientFromCtx(c)
				ip, err := client.UpdateIP(context.Background(), id, api.UpdateIPParams{Healthchecks: healthchecks})
				if err != nil {
					return errors.Wrapf(err, "fail to update the IP '%s'", id)
				}

				fmt.Println(aurora.Green(fmt.Sprintf("IP %s (%s) successfully updated", ip.IP.IP, ip.ID)))
				return nil
			},
		}, {
			Name:      "add",
			ArgsUsage: "IP [CHECK_TYPE CHECK_ENDPOINT]...",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "keepalive-interval",
					Value: 0,
					Usage: "Duration between keepalive requests",
				},
				cli.IntFlag{
					Name:  "healthcheck-interval",
					Value: 0,
					Usage: "Duration between healthchecks",
				},
			},
			Action: func(c *cli.Context) error {
				if len(c.Args())%2 == 0 {
					// 1 For the IP
					// And 2 per Healthchecks
					// So NArgs % 2 must be == 1
					cli.ShowCommandHelp(c, c.Command.Name)
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
						return fmt.Errorf("Invalid endpoint: %s", endpoint)
					}
					portI, err := strconv.Atoi(port)
					if err != nil {
						return fmt.Errorf("Invalid endpoint: %s", endpoint)
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
					KeepaliveInterval:   c.Int("keepalive-interval"),
				}
				newIP, err := client.AddIP(context.Background(), ip, params)
				if err != nil {
					return err
				}

				fmt.Println(aurora.Green(fmt.Sprintf("IP %s (%s) successfully added", newIP.IP.IP, newIP.ID)))
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
	if c.GlobalString("host") != "" {
		opts = append(opts, api.WithURL(c.GlobalString("host")))
	}

	if c.GlobalString("user") != "" {
		opts = append(opts, api.WithUser(c.GlobalString("user")))
	}

	if c.GlobalString("password") != "" {
		opts = append(opts, api.WithPassword(c.GlobalString("password")))
	}

	return api.NewHTTPClient(opts...)
}
