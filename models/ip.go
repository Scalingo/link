package models

type IPStatus string

const (
	IPStatusStandBy   IPStatus = "STANDBY"
	IPStatusActivated IPStatus = "ACTIVATED"
)

type IP struct {
	ID string `json:"id"`
	IP string `json:"ip"`
}
