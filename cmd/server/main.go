package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/itimky/word-of-wisom/internal/repository/quote"

	"github.com/itimky/word-of-wisom/pkg/utils"

	gtpserver "github.com/itimky/word-of-wisom/internal/gtp/server"
	tcpserver "github.com/itimky/word-of-wisom/internal/tcp/server"

	"github.com/sirupsen/logrus"

	"github.com/itimky/word-of-wisom/pkg/gtp"
)

func main() {
	conf, err := loadConfig()
	if err != nil {
		logrus.WithError(err).Fatal("cannot init server")
	}

	utils.SetupLogrus(conf.Debug)

	logrus.Debugf("%+v", conf)

	hashCalc := gtp.NewCalc(time.Now)

	//nolint:gosec
	// Just random item, no security
	quoteSvc := quote.NewRepository(rand.New(rand.NewSource(time.Now().Unix())))

	shieldSvc := gtpserver.NewServer(
		gtpserver.Config{
			TourLength:           conf.TourLength,
			SecretLength:         conf.SecretLength,
			SecretUpdateInterval: conf.SecretUpdateInterval,
			GuideSecrets:         conf.GuideSecrets,
		},
		hashCalc,
	)
	if err = shieldSvc.Init(); err != nil {
		logrus.WithError(err).Fatal("shield service init")
	}

	srv := tcpserver.NewServer(
		fmt.Sprintf("%v:%v", conf.Host, conf.Port),
		conf.Multicore,
		shieldSvc,
		quoteSvc,
	)
	if err = srv.Run(); err != nil {
		logrus.WithError(err).Fatal("server run")
	}
}
