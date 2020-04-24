package models

import (
	"strings"

	"github.com/sirupsen/logrus"
)

// IP stores the configuration of a virtual IP for one host
type IP struct {
	ID                  string        `json:"id"`                   // ID of this Virtual IP (strarting with vip-)
	IP                  string        `json:"ip"`                   // Full IP with mask (i.e. 10.0.0.1/32)
	Status              string        `json:"status,omitempty"`     // Status of this VIP
	Checks              []Healthcheck `json:"checks,omitempty"`     // Healthcheck configured with this VIP
	KeepaliveInterval   int           `json:"keepalive_interval"`   // Deprecated
	HealthcheckInterval int           `json:"healthcheck_interval"` // HealthcheckIntevals for this VIP
}

// ToLogrusFields this method is used to add important vip fields to our logger
func (i IP) ToLogrusFields() logrus.Fields {
	return logrus.Fields{
		"ip":    i.IP,
		"ip_id": i.ID,
	}
}

func (ip IP) StorableIP() string {
	return strings.Replace(ip.IP, "/", "_", -1)
}
