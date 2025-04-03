package api

const (
	Activated = "ACTIVATED"
	Standby   = "STANDBY"
	Failing   = "FAILING"
)

type Endpoint struct {
	ID string `json:"id"`

	Status              string        `json:"status,omitempty"`
	Checks              []HealthCheck `json:"checks,omitempty"`
	HealthCheckInterval int           `json:"healthcheck_interval"`
	Plugin              string        `json:"plugin,omitempty"`
	ElectionKey         string        `json:"election_key,omitempty"`
}

type HealthCheckType string

const (
	TCPHealthCheck HealthCheckType = "TCP"
)

type HealthCheck struct {
	Type HealthCheckType `json:"type"`
	Host string          `json:"host"`
	Port int             `json:"port"`
}

type EndpointGetResponse struct {
	Endpoint Endpoint `json:"endpoint"`
}

type EndpointListResponse struct {
	Endpoints []Endpoint `json:"endpoints"`
}

type UpdateEndpointParams struct {
	HealthChecks []HealthCheck `json:"healthchecks"`
}

type AddEndpointParams struct {
	HealthCheckInterval int           `json:"healthcheck_interval"`
	Checks              []HealthCheck `json:"checks"`
	Plugin              string        `json:"plugin"`
	PluginConfig        any           `json:"plugin_config,omitempty"`
}

type GetEndpointHostsResponse struct {
	Hosts []Host `json:"hosts"`
}

type Host struct {
	Hostname string `json:"hostname"`
}
