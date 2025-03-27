package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/logrusorgru/aurora/v3"
	"github.com/olekukonko/tablewriter"

	"github.com/Scalingo/link/v2/api"
)

func formatEndpoints(endpoints []api.Endpoint) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "IP", "Status", "CHECKS"})

	for _, endpoint := range endpoints {
		status := formatStatus(endpoint)

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
			endpoint.IP,
			status,
			checks,
		})
	}
	table.Render()
}

func formatEndpoint(ip api.Endpoint) {
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

func formatStatus(endpoint api.Endpoint) string {
	switch endpoint.Status {
	case api.Activated:
		return aurora.Green("ACTIVATED").String()
	case api.Standby:
		return aurora.Yellow("STANDBY").String()
	case api.Failing:
		return aurora.Red("FAILING").String()
	default:
		return endpoint.Status
	}
}
