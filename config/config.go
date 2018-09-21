package config

import "github.com/kelseyhightower/envconfig"

var C Config

type Config struct {
	Interface string `envconfig:"INTERFACE"`
	Hostname  string `envconfig:"HOSTNAME"`
	Port      int    `envconfig:"PORT" default:"1313"`
}

func Init() {
	err := envconfig.Process("", &C)
	if err != nil {
		panic(err)
	}
}
