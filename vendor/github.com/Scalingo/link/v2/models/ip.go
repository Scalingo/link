package models

import (
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// IP stores the configuration of a virtual IP for one host
type IP struct {
	ID                  string        `json:"id"`                   // ID of this Virtual IP (strarting with vip-)
	IP                  string        `json:"ip"`                   // Full IP with mask (i.e. 10.0.0.1/32)
	Status              string        `json:"status,omitempty"`     // Status of this VIP
	Checks              []Healthcheck `json:"checks,omitempty"`     // Healthcheck configured with this VIP
	HealthcheckInterval int           `json:"healthcheck_interval"` // HealthcheckIntevals for this VIP
}

// ToLogrusFields returns a Logrus representation of an IP
func (i IP) ToLogrusFields() logrus.Fields {
	return logrus.Fields{
		"ip":    i.IP,
		"ip_id": i.ID,
	}
}

// StorableIP transforms the IP to a string that is compatible with ETCD key naming rules
func (i IP) StorableIP() string {
	return strings.Replace(i.IP, "/", "_", -1)
}

// IPLink is the structure stored when an IP is linked to an Host
type IPLink struct {
	UpdatedAt time.Time `json:"updated_at"`
}
