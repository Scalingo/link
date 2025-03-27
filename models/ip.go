package models

import (
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/Scalingo/link/v2/api"
)

// Endpoint stores the configuration of a virtual IP for one host
type Endpoint struct {
	ID string `json:"id"` // ID of this Virtual IP (strarting with vip-)

	// Full IP with mask (i.e. 10.0.0.1/32)
	IP string `json:"ip"`

	Status              string       `json:"status,omitempty"`     // Status of this VIP
	Checks              HealthChecks `json:"checks,omitempty"`     // Healthcheck configured with this VIP
	HealthCheckInterval int          `json:"healthcheck_interval"` // HealthcheckIntevals for this VIP
}

// ToLogrusFields returns a Logrus representation of an IP
func (i Endpoint) ToLogrusFields() logrus.Fields {
	return logrus.Fields{
		"ip":    i.IP,
		"ip_id": i.ID,
	}
}

// StorableIP transforms the IP to a string that is compatible with ETCD key naming rules
func (i Endpoint) StorableIP() string {
	return strings.Replace(i.IP, "/", "_", -1)
}

func (i Endpoint) ToAPIType() api.Endpoint {
	return api.Endpoint{
		ID:                  i.ID,
		IP:                  i.IP,
		Status:              i.Status,
		Checks:              i.Checks.ToAPIType(),
		HealthCheckInterval: i.HealthCheckInterval,
	}
}

type Endpoints []Endpoint

func (e Endpoints) ToAPIType() []api.Endpoint {
	endpoints := make([]api.Endpoint, 0, len(e))
	for _, endpoint := range e {
		endpoints = append(endpoints, endpoint.ToAPIType())
	}
	return endpoints
}

// IPLink is the structure stored when an IP is linked to an Host
type IPLink struct {
	UpdatedAt time.Time `json:"updated_at"`
}
