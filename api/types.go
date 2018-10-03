package api

import "github.com/Scalingo/link/models"

const (
	ACTIVATED = "ACTIVATED"
	STANDBY   = "STANDBY"
	FAILING   = "FAILING"
)

type IP struct {
	models.IP
	Status string `json:"status,omitempty"`
}
