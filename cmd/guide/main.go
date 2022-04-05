package main

import (
	"time"
	"word-of-wisom/internal/guide"
	"word-of-wisom/pkg/gtp"

	"github.com/sirupsen/logrus"
)

func main() {
	conf, err := loadConfig()
	if err != nil {
		logrus.WithError(err).Fatal("cannot init guide")
	}

	if conf.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	logrus.Debugf("%+v", conf)

	g := guide.NewGuide(conf.Host, conf.Port, conf.Secret, gtp.NewGTP(time.Now))
	if err := g.Run(); err != nil {
		logrus.WithError(err).Fatal("cannot run guide")
	}
}
