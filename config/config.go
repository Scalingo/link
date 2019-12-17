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

func (c Config) LeaseTime() time.Duration {
	return 5 * c.KeepAliveInterval
}

type Delta float64

// DefaultDelta is the default value to use with PlusMinusDelta
const DefaultDelta Delta = 0.25

// PlusMinusDelta returns a number [value - value * delta ≤ value ≤ value + value * delta]
func PlusMinusDelta(value float64, delta Delta) float64 {
	lower := float64(value) * (1.0 - float64(delta))
	higher := float64(value) * (1.0 + float64(delta))
	diff := higher - lower
	rate := rand.Float64()
	return lower + rate*diff
}

// PlusMinusDeltaDuration is a helper to manipulate time.Duration with PlusMinusDelta
func PlusMinusDeltaDuration(value time.Duration, delta Delta) time.Duration {
	res := PlusMinusDelta(float64(value), delta)
	return time.Duration(int64(res))
}

func Build() (Config, error) {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		return config, errors.Wrap(err, "fail to parse environment")
	}

	return config, nil
}
