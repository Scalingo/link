package models

import "fmt"

type HealthcheckType string

const (
	TCPHealthCheck HealthcheckType = "TCP"
)

type Healthcheck struct {
	Type HealthcheckType `json:"type"`
	Host string          `json:"host"`
	Port int             `json:"port"`
}

func (h Healthcheck) Addr() string {
	return fmt.Sprintf("%s:%v", h.Host, h.Port)
}
