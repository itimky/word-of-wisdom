package main

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Host string `envconfig:"HOST"`
	Port string `envconfig:"PORT"`

	SecretLength                string   `envconfig:"SECRET_LENGTH"`
	SecretUpdateIntervalMinutes string   `envconfig:"SECRET_UPDATE_INTERVAL_MINUTES"`
	TourLength                  string   `envconfig:"TOUR_LENGTH"`
	GuideSecrets                []string `envconfig:"GUIDE_SECRETS"`

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
