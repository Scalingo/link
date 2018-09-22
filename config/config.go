package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	Interface string `envconfig:"INTERFACE"`
	Hostname  string `envconfig:"HOSTNAME"`
	Port      int    `envconfig:"PORT" default:"1313"`
}

func Build() (Config, error) {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		return config, errors.Wrap(err, "fail to parse environment")
	}

	return config, nil
}
