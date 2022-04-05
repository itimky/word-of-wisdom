package main

import (
	"time"
	"word-of-wisom/internal/client"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	conf, err := loadConfig()
	if err != nil {
		logrus.WithError(err).Fatal("cannot init client")
	}

	if conf.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	logrus.Debugf("%+v", conf)

	cl := client.NewClient(conf.Server, conf.Guides)

	ticker := time.NewTicker(conf.RequestInterval)
	for range ticker.C {
		quote, err := cl.RequestQuote(conf.RetryCount)
		if err != nil {
			logrus.WithError(err).Error("request quote")
			continue
		}

		logrus.Infof("Quote: %v", quote)
	}
}
