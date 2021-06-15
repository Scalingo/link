package models

// Host store the host configuration. This is used by an host to retrieve his configuration after a restart.
type Host struct {
	Hostname string `json:"hostname"` // Hostname of this host
	LeaseID  int64  `json:"lease_id"` // LeaseID current lease ID of this host
}
