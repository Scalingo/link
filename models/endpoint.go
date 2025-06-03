package models

import (
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/Scalingo/link/v3/api"
)

// Endpoint stores the configuration of an endpoint for one host
type Endpoint struct {
	ID string `json:"id"` // ID of this endpoint (starting with vip-)

	// Full IP with mask (i.e. 10.0.0.1/32)
	// Deprecated: The IP is now stored in the ARP Plugin config. This field is kept for backward compatibility
	IP string `json:"ip"`

	Checks              HealthChecks `json:"checks,omitempty"`     // Health check configured with this Endpoint
	HealthCheckInterval int          `json:"healthcheck_interval"` // Health check Intervals for this Endpoint

	Plugin       string          `json:"plugin,omitempty"`        // Plugin to use for this Endpoint
	PluginConfig json.RawMessage `json:"plugin_config,omitempty"` // Plugin configuration
}

// ToLogrusFields returns a Logrus representation of an IP
func (i Endpoint) LogFields() logrus.Fields {
	return logrus.Fields{
		"id":     i.ID,
		"plugin": i.Plugin,
	}
}

func (i Endpoint) ToAPIType() api.Endpoint {
	// Retro Compatibility: Can be removed when all LinK switched to the v3 API and storage.
	plugin := i.Plugin
	if plugin == "" {
		plugin = "arp"
	}

	return api.Endpoint{
		ID:                  i.ID,
		Checks:              i.Checks.ToAPIType(),
		HealthCheckInterval: i.HealthCheckInterval,
		Plugin:              plugin,
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

// EndpointLink is the structure stored when an IP is linked to an Host
type EndpointLink struct {
	UpdatedAt time.Time `json:"updated_at"`
}
