package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/itimky/word-of-wisom/internal/service/quote"
	"github.com/itimky/word-of-wisom/internal/service/shield"
	"github.com/itimky/word-of-wisom/internal/tcp/server"

	"github.com/itimky/word-of-wisom/pkg/gtp"
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

	hashCalc := gtp.NewGTP(time.Now)

	//nolint:gosec
	// Just random item, no security
	quoteSvc := quote.NewService(rand.New(rand.NewSource(time.Now().Unix())))

	shieldSvc := shield.NewService(
		shield.Config{
			TourLength:           conf.TourLength,
			SecretLength:         conf.SecretLength,
			SecretUpdateInterval: conf.SecretUpdateInterval,
			GuideSecrets:         conf.GuideSecrets,
		},
		hashCalc,
		quoteSvc,
	)
	if err = shieldSvc.Init(); err != nil {
		logrus.WithError(err).Fatal("shield service init")
	}

	srv := server.NewServer(
		fmt.Sprintf("%v:%v", conf.Host, conf.Port),
		conf.Multicore,
		shieldSvc,
	)
	if err = srv.Run(); err != nil {
		logrus.WithError(err).Fatal("server run")
	}
}
