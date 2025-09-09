package config

import (
	"context"
	"math/rand"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/go-utils/logger"
)

type Config struct {
	Hostname            string        `envconfig:"HOSTNAME"`
	User                string        `envconfig:"USER"`
	Password            string        `envconfig:"PASSWORD"`
	Port                int           `envconfig:"PORT" default:"1313"`
	KeepAliveInterval   time.Duration `envconfig:"KEEPALIVE_INTERVAL" default:"3s"`
	KeepAliveRetry      int           `envconfig:"KEEPALIVE_RETRY" default:"5"`
	HealthCheckInterval time.Duration `envconfig:"HEALTH_CHECK_INTERVAL" default:"5s"`
	HealthCheckTimeout  time.Duration `envconfig:"HEALTH_CHECK_TIMEOUT" default:"5s"`

	PluginEnsureInterval           time.Duration `envconfig:"PLUGIN_ENSURE_INTERVAL" default:"1s"`
	PluginEnsureMaxBackoffInterval time.Duration `envconfig:"PLUGIN_ENSURE_MAX_BACKOFF_INTERVAL" default:"10m"`

	ARPGratuitousInterval   time.Duration `envconfig:"ARP_GRATUITOUS_INTERVAL" default:"1s"` // Deprecated: Use PluginEnsureInterval
	FailCountBeforeFailover int           `envconfig:"FAIL_COUNT_BEFORE_FAILOVER" default:"3"`

	SecretStorageEncryptionKey string   `envconfig:"SECRET_STORAGE_ENCRYPTION_KEY" default:""`
	SecretStorageAlternateKeys []string `envconfig:"SECRET_STORAGE_ALTERNATE_KEYS" default:""`

	MaxNumberOfEndpoints int `envconfig:"MAX_NUMBER_OF_ENDPOINTS" default:"1000"`
}

// LeaseTime is 5 * the global keepalive interval
func (c Config) LeaseTime() time.Duration {
	return 5 * c.KeepAliveInterval
}

// RandomDurationAround returns a duration (1.0 - percent) * duration < n < (1.0 + percent) * duration
func RandomDurationAround(duration time.Duration, scatteringPercentage float64) time.Duration {
	delta := int64(float64(duration) * 2 * scatteringPercentage)
	return time.Duration(int64(duration) + rand.Int63n(delta) - (delta / 2))
}

func Build(ctx context.Context) (Config, error) {
	log := logger.Get(ctx)
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		return config, errors.Wrap(ctx, err, "fail to parse environment")
	}

	if _, ok := os.LookupEnv("ARP_GRATUITOUS_INTERVAL"); ok {
		log.Error("ARP_GRATUITOUS_INTERVAL is deprecated, please use PLUGIN_ENSURE_INTERVAL instead")
		config.PluginEnsureInterval = config.ARPGratuitousInterval
	}

	if config.PluginEnsureInterval >= config.PluginEnsureMaxBackoffInterval {
		return config, errors.New(ctx, "PLUGIN_ENSURE_MAX_BACKOFF_INTERVAL must be greater than PLUGIN_ENSURE_INTERVAL")
	}

	return config, nil
}
