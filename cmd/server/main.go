package main

import (
	"github.com/sirupsen/logrus"
	"word-of-wisom/server"
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

	srv := server.NewServer(conf.Host, conf.Port, 10, 2, 5, conf.GuideSecrets)
	if err := srv.Run(); err != nil {
		logrus.WithError(err).Fatal("cannot run server")
	}
}
