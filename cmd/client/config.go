package main

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Server     string   `envconfig:"SERVER"`
	Guides     []string `envconfig:"GUIDES"`
	RetryCount int      `envconfig:"RETRY_COUNT"`

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
