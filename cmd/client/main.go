package main

import (
	"github.com/sirupsen/logrus"
	"time"
	"word-of-wisom/client"
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

	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		quote, err := cl.RequestQuote(conf.RetryCount)
		if err != nil {
			logrus.WithError(err).Error("request quote")
			continue
		}
		logrus.Debugf("Quote: %v", quote)
	}
}
