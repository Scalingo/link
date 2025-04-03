package endpoint

import (
	"github.com/logrusorgru/aurora/v3"

	"github.com/Scalingo/link/v2/api"
)

func FormatStatus(endpoint api.Endpoint) string {
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
