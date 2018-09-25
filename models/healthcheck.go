package models

import (
	"net"
	"strconv"
)

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
	return net.JoinHostPort(h.Host, strconv.Itoa(h.Port))
}
