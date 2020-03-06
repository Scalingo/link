package config

import (
	"math/rand"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Config struct {
	Interface             string        `envconfig:"INTERFACE"`
	Hostname              string        `envconfig:"HOSTNAME"`
	User                  string        `envconfig:"USER"`
	Password              string        `envconfig:"PASSWORD"`
	Port                  int           `envconfig:"PORT" default:"1313"`
	KeepAliveInterval     time.Duration `envconfig:"KEEPALIVE_INTERVAL" default:"3s"`
	KeepAliveRetry        int           `envconfig:"KEEPALIVE_RETRY" default:"5"`
	HealthcheckInterval   time.Duration `envconfig:"HEALTH_CHECK_INTERVAL" default:"5s"`
	HealthcheckTimeout    time.Duration `envconfig:"HEALTH_CHECK_TIMEOUT" default:"5s"`
	ARPGratuitousInterval time.Duration `envconfig:"ARP_GRATUITOUS_INTERVAL" default:"1s"`
	// Number of gratuitous ARP (GARP) packets sent when the state becomes 'ACTIVATED'
	ARPGratuitousCount      int `envconfig:"ARP_GRATUITOUS_COUNT" default:"3"`
	FailCountBeforeFailover int `envconfig:"FAIL_COUNT_BEFORE_FAILOVER" default:"3"`
}

// LeaseTime is either 5* the global keepalive interval, or 5 times the one
// which is given as argument which can be speicific to the current IP.
func (c Config) LeaseTime(ipKeepalive int) time.Duration {
	duration := time.Duration(ipKeepalive) * time.Second
	if duration == 0 {
		duration = c.KeepAliveInterval
	}
	return 5 * duration
}

// RandomDurationAround returns a duration (1.0 - percent) * duration < n < (1.0 + percent) * duration
func RandomDurationAround(duration time.Duration, scatteringPercentage float64) time.Duration {
	delta := int64(float64(duration) * 2 * scatteringPercentage)
	return time.Duration(int64(duration) + rand.Int63n(delta) - (delta / 2))
}

func Build() (Config, error) {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		return config, errors.Wrap(err, "fail to parse environment")
	}

	return config, nil
}
