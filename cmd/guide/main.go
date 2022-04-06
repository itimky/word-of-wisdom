package main

import (
	"time"

	"github.com/itimky/word-of-wisom/internal/guide"
	"github.com/itimky/word-of-wisom/pkg/gtp"
	"github.com/itimky/word-of-wisom/pkg/tcpserver"

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

	g := guide.NewGuide(conf.Secret, gtp.NewGTP(time.Now))
	tcpServer := tcpserver.NewTCPServer(
		conf.Host,
		conf.Port,
		g,
	)
	if err := tcpServer.Run(); err != nil {
		logrus.WithError(err).Fatal("cannot run guide")
	}
}
