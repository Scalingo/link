package models

// Host stores the host configuration. This is used by an host to retrieve his configuration after a restart.
type Host struct {
	Hostname string `json:"hostname"` // Hostname of this host
	LeaseID  int64  `json:"lease_id"` // LeaseID current lease ID of this host
	// DataVersion is the version number of how the data is stored in etcd for this host. This field is introduced in 2.0.0. But the data format version v0 correspond to the format before v1.9.0.
	DataVersion int `json:"data_version,omitempty"`
}
