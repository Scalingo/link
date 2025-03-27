package models

import (
	"net"
	"strconv"

	"github.com/Scalingo/link/v2/api"
)

type HealthChecks []HealthCheck

func (h HealthChecks) ToAPIType() []api.HealthCheck {
	checks := make([]api.HealthCheck, 0, len(h))
	for _, check := range h {
		checks = append(checks, check.ToAPIType())
	}
	return checks
}

type HealthCheck struct {
	Type api.HealthCheckType `json:"type"`
	Host string              `json:"host"`
	Port int                 `json:"port"`
}

func (h HealthCheck) Addr() string {
	return net.JoinHostPort(h.Host, strconv.Itoa(h.Port))
}

func (h HealthCheck) ToAPIType() api.HealthCheck {
	return api.HealthCheck{
		Type: h.Type,
		Host: h.Host,
		Port: h.Port,
	}
}

func HealthCheckFromAPIType(h api.HealthCheck) HealthCheck {
	return HealthCheck{
		Type: h.Type,
		Host: h.Host,
		Port: h.Port,
	}
}

func HealthChecksFromAPIType(hs []api.HealthCheck) HealthChecks {
	checks := make(HealthChecks, 0, len(hs))
	for _, check := range hs {
		checks = append(checks, HealthCheckFromAPIType(check))
	}
	return checks
}
