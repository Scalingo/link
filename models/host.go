package models

type Host struct {
	Hostname string `json:"hostname"`
	LeaseID  int64  `json:"lease_id"`
}
