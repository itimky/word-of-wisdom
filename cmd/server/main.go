package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/panjf2000/gnet/v2/pkg/pool/goroutine"

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
	quoteRepo := quote.NewRepository(rand.New(rand.NewSource(time.Now().Unix())))

	gtpServer := gtpserver.NewServer(
		gtpserver.Config{
			TourLength:           conf.TourLength,
			SecretLength:         conf.SecretLength,
			SecretUpdateInterval: conf.SecretUpdateInterval,
			GuideSecrets:         conf.GuideSecrets,
		},
		hashCalc,
	)
	if err = gtpServer.Init(); err != nil {
		logrus.WithError(err).Fatal("shield service init")
	}

	pool := goroutine.Default()
	defer pool.Release()

	srv := tcpserver.NewServer(
		pool,
		fmt.Sprintf("%v:%v", conf.Host, conf.Port),
		conf.Multicore,
		conf.Timeout,
		gtpServer,
		quoteRepo,
	)
	if err = srv.Run(); err != nil {
		logrus.WithError(err).Fatal("server run")
	}
}
