package main

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Host      string        `envconfig:"HOST"`
	Port      string        `envconfig:"PORT"`
	Multicore bool          `envconfig:"MULTICORE"`
	Timeout   time.Duration `envconfig:"TIMEOUT"`

	SecretLength         int           `envconfig:"SECRET_LENGTH"`
	SecretUpdateInterval time.Duration `envconfig:"SECRET_UPDATE_INTERVAL"`
	TourLength           int           `envconfig:"TOUR_LENGTH"`
	GuideSecrets         []string      `envconfig:"GUIDE_SECRETS"`

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
