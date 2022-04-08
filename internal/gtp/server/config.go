package server

import "time"

type Config struct {
	TourLength           int
	SecretLength         int
	SecretUpdateInterval time.Duration
	GuideSecrets         []string
}
