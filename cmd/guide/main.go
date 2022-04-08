package main

import (
	"fmt"
	"time"

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

	g := guide.NewGuide(
		fmt.Sprintf("%v:%v", conf.Host, conf.Port),
		conf.Multicore,
		conf.Secret,
		gtp.NewCalc(time.Now),
	)
	if err = g.Run(); err != nil {
		logrus.WithError(err).Fatal("guide run")
	}
}
