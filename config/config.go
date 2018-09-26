package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	Interface             string        `envconfig:"INTERFACE"`
	Hostname              string        `envconfig:"HOSTNAME"`
	Port                  int           `envconfig:"PORT" default:"1313"`
	KeepAliveInterval     time.Duration `envconfig:"KEEPALIVE_INTERVAL" default:"3s"`
	HealthcheckInterval   time.Duration `envconfig:"HEALTH_CHECK_INTERVAL" default:"5s"`
	HealthcheckTimeout    time.Duration `envconfig:"HEALTH_CHECK_TIMEOUT" default:"5s"`
	ARPGratuitousInterval time.Duration `envconfig:"ARP_GRATUITOUS_INTERVAL" default:"10s"`
}

func (c Config) LeaseTime() time.Duration {
	return 2 * c.KeepAliveInterval
}

func Build() (Config, error) {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		return config, errors.Wrap(err, "fail to parse environment")
	}

	return config, nil
}
