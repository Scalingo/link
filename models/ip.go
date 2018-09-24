package models

type IP struct {
	ID     string `json:"id"`
	IP     string `json:"ip"`
	Status string `json:"status,omitempty"`
}
