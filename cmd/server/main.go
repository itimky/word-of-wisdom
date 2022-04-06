package main

import (
	"math/rand"
	"time"
	"word-of-wisom/internal/server"
	"word-of-wisom/pkg/gtp"
	"word-of-wisom/pkg/quotes"
	"word-of-wisom/pkg/tcpserver"

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
		conf.SecretLength,
		conf.TourLength,
		conf.GuideSecrets,
		gtp.NewGTP(time.Now),
		quotes.NewQuoteRandomizer(rand.New(rand.NewSource(time.Now().Unix()))), //nolint:gosec
	)
	err = srv.UpdateSecret()
	if err != nil {
		logrus.WithError(err).Fatal("update secret")
	}

	go func() {
		ticker := time.NewTicker(conf.SecretUpdateInterval)
		for range ticker.C {
			err := srv.UpdateSecret()
			if err != nil {
				logrus.WithError(err).Error("periodic secret update")
			}
		}
	}()

	tcpServer := tcpserver.NewTCPServer(
		conf.Host,
		conf.Port,
		srv,
	)
	if err := tcpServer.Run(); err != nil {
		logrus.WithError(err).Fatal("cannot run server")
	}
}
