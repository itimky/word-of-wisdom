package main

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Host string `envconfig:"HOST"`
	Port string `envconfig:"PORT"`

	SecretLength                int           `envconfig:"SECRET_LENGTH"`
	SecretUpdateIntervalSeconds time.Duration `envconfig:"SECRET_UPDATE_INTERVAL_SECONDS"`
	TourLength                  int           `envconfig:"TOUR_LENGTH"`
	GuideSecrets                []string      `envconfig:"GUIDE_SECRETS"`

	Debug bool `envconfig:"RETRY_COUNT"`
}

func loadConfig() (config, error) {
	var conf config

	err := envconfig.Process("wow", &conf)

	if err != nil {
		return conf, fmt.Errorf("load config: %w", err)
	}

	return conf, nil
}
