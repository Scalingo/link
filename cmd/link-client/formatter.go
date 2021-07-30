package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/Scalingo/link/v2/api"
	"github.com/logrusorgru/aurora"
	"github.com/olekukonko/tablewriter"
)

func formatIPs(ips []api.IP) {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "IP", "Status", "CHECKS"})

	for _, ip := range ips {
		status := formatStatus(ip)

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

func formatIP(ip api.IP) {
	fmt.Printf("ID:\t%s\n", ip.ID)
	fmt.Printf("Status:\t%s\n", formatStatus(ip))
	if len(ip.Checks) == 0 {
		fmt.Printf("Checks:\tNone\n")
	} else {
		fmt.Println("Checks:")
		for _, check := range ip.Checks {
			fmt.Printf(" - Type: %s, Host: %s, Port: %v\n", check.Type, check.Host, check.Port)
		}
	}
}

func formatStatus(ip api.IP) string {
	switch ip.Status {
	case api.Activated:
		return aurora.Green("ACTIVATED").String()
	case api.Standby:
		return aurora.Brown("STANDBY").String()
	case api.Failing:
		return aurora.Red("FAILING").String()
	default:
		return ip.Status
	}
}
