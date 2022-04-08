package main

import (
	"time"

	client "github.com/itimky/word-of-wisom/internal/client"

	"github.com/sirupsen/logrus"

	gtpclient "github.com/itimky/word-of-wisom/internal/gtp/client"
	tcpclient "github.com/itimky/word-of-wisom/internal/tcp/client"
	"github.com/itimky/word-of-wisom/pkg/utils"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	conf, err := loadConfig()
	if err != nil {
		logrus.WithError(err).Fatal("cannot init client")
	}

	utils.SetupLogrus(conf.Debug)

	logrus.Debugf("%+v", conf)

	tcpClient := tcpclient.NewClient(conf.Server, conf.Guides)
	gtpClient := gtpclient.NewClient(tcpClient, len(conf.Guides), conf.RetryCount)
	c := client.NewClient(gtpClient)

	ticker := time.NewTicker(conf.RequestInterval)
	for range ticker.C {
		response, err := c.RequestQuote()
		if err != nil {
			logrus.WithError(err).Error("request quote")
			continue
		}

		logrus.Infof("Quote: %v\n\n", response.Quote)
	}
}
