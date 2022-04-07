package main

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Host      string `envconfig:"HOST"`
	Port      string `envconfig:"PORT"`
	Multicore bool   `envconfig:"MULTICORE"`

	Secret string `envconfig:"SECRET"`

	Debug bool `envconfig:"DEBUG"`
}

func loadConfig() (config, error) {
	var conf config

	err := envconfig.Process("wow", &conf)

	if err != nil {
		return conf, fmt.Errorf("load config: %w", err)
	}

	return conf, nil
}
