package models

import (
	"context"
	"net"
	"strconv"

	"github.com/Scalingo/go-utils/errors/v2"
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

func (h HealthCheck) Validate(_ context.Context) error {
	validation := errors.NewValidationErrorsBuilder()
	if h.Type == "" {
		validation.Set("type", "Health check type is required")
	}
	if h.Type != api.TCPHealthCheck {
		validation.Set("type", "Health check type is not supported")
	}
	if h.Host == "" {
		validation.Set("host", "Host is required")
	}
	if h.Port <= 0 || h.Port > 65535 {
		validation.Set("port", "Port must be between 1 and 65535")
	}

	validationErr := validation.Build()
	if validationErr != nil {
		return validationErr
	}
	return nil
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

func (h HealthChecks) Validate(ctx context.Context) error {
	validation := errors.NewValidationErrorsBuilder()
	for i, check := range h {
		if err := check.Validate(ctx); err != nil {
			validation.Set(strconv.Itoa(i), err.Error())
		}
	}
	validationErr := validation.Build()
	if validationErr != nil {
		return validationErr
	}
	return nil
}
