package api

import "github.com/Scalingo/link/v2/models"

const (
	Activated = "ACTIVATED"
	Standby   = "STANDBY"
	Failing   = "FAILING"
)

type IP struct {
	models.IP
	Status string `json:"status,omitempty"`
}
