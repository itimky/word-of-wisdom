package main

import (
	"time"

	"github.com/itimky/word-of-wisom/pkg/utils"

	"github.com/itimky/word-of-wisom/internal/tcp/client"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	conf, err := loadConfig()
	if err != nil {
		logrus.WithError(err).Fatal("cannot init client")
	}

	utils.SetupLogrus(conf.Debug)

	logrus.Debugf("%+v", conf)

	cl := client.NewClient(conf.Server, conf.Guides)

	ticker := time.NewTicker(conf.RequestInterval)
	for range ticker.C {
		quote, err := cl.RequestQuote(conf.RetryCount)
		if err != nil {
			logrus.WithError(err).Error("request quote")
			continue
		}

		logrus.Infof("Quote: %v\n\n", quote)
	}
}
