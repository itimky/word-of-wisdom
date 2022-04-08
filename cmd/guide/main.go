package main

import (
	"fmt"
	"time"

	"github.com/panjf2000/gnet/v2/pkg/pool/goroutine"

	"github.com/itimky/word-of-wisom/pkg/utils"

	"github.com/itimky/word-of-wisom/internal/tcp/guide"

	"github.com/sirupsen/logrus"

	"github.com/itimky/word-of-wisom/pkg/gtp"
)

func main() {
	conf, err := loadConfig()
	if err != nil {
		logrus.WithError(err).Fatal("cannot init guide")
	}

	utils.SetupLogrus(conf.Debug)

	logrus.Debugf("%+v", conf)

	pool := goroutine.Default()
	defer pool.Release()

	g := guide.NewGuide(
		pool,
		fmt.Sprintf("%v:%v", conf.Host, conf.Port),
		conf.Multicore,
		conf.Timeout,
		conf.Secret,
		gtp.NewCalc(time.Now),
	)
	if err = g.Run(); err != nil {
		logrus.WithError(err).Fatal("guide run")
	}
}
