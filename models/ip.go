package models

import "github.com/sirupsen/logrus"

type IP struct {
	ID                  string        `json:"id"`
	IP                  string        `json:"ip"`
	Status              string        `json:"status,omitempty"`
	Checks              []Healthcheck `json:"checks,omitempty"`
	KeepaliveInterval   int           `json:"keepalive_interval"`
	HealthcheckInterval int           `json:"healthcheck_interval"`
}

func (i IP) ToLogrusFields() logrus.Fields {
	return logrus.Fields{
		"ip":    i.IP,
		"ip_id": i.ID,
	}
}
