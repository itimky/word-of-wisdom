package main

import (
	"math/rand"
	"time"
	"word-of-wisom/internal/gtp"
	"word-of-wisom/server"

	"github.com/sirupsen/logrus"
)

func main() {
	conf, err := loadConfig()
	if err != nil {
		logrus.WithError(err).Fatal("cannot init server")
	}

	if conf.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	logrus.Debugf("%+v", conf)

	srv := server.NewServer(
		conf.Host,
		conf.Port,
		conf.SecretLength,
		time.Duration(conf.SecretUpdateIntervalSeconds),
		conf.TourLength,
		conf.GuideSecrets,
		gtp.NewGTP(time.Now),
		rand.New(rand.NewSource(time.Now().Unix())), //nolint:gosec
	)
	if err := srv.Run(); err != nil {
		logrus.WithError(err).Fatal("cannot run server")
	}
}
