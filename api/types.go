package api

const (
	Activated = "ACTIVATED"
	Standby   = "STANDBY"
	Failing   = "FAILING"
)

type Endpoint struct {
	ID string `json:"id"`

	IP                  string        `json:"ip"`
	Status              string        `json:"status,omitempty"`
	Checks              []HealthCheck `json:"checks,omitempty"`
	HealthCheckInterval int           `json:"healthcheck_interval"`
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
	Endpoint Endpoint `json:"ip"`
}

type EndpointListResponse struct {
	Endpoints []Endpoint `json:"ips"`
}

type UpdateEndpointParams struct {
	HealthChecks []HealthCheck `json:"healthchecks"`
}
