package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/Scalingo/link/api"
	"github.com/logrusorgru/aurora"
	"github.com/olekukonko/tablewriter"
)

func formatIPs(ips []api.IP) {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "IP", "Status", "CHECKS"})

	for _, ip := range ips {
		var status string

		switch ip.Status {
		case api.ACTIVATED:
			status = aurora.Green("ACTIVATED").String()
		case api.STANDBY:
			status = aurora.Brown("STANDBY").String()
		case api.FAILING:
			status = aurora.Red("FAILING").String()
		default:
			status = ip.Status
		}

		checks := "None"
		if len(ip.Checks) > 0 {
			var c []string
			for _, check := range ip.IP.Checks {
				c = append(c, fmt.Sprintf("%s - %s", check.Type, net.JoinHostPort(check.Host, strconv.Itoa(check.Port))))
			}

			checks = strings.Join(c, ",")
		}
		table.Append([]string{
			ip.ID,
			ip.IP.IP,
			status,
			checks,
		})
	}
	table.Render()
}
