package utils

import "github.com/sirupsen/logrus"

func SetupLogrus(debug bool) {
	formatter := &logrus.TextFormatter{}
	formatter.DisableQuote = true
	logrus.SetFormatter(formatter)

	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
}