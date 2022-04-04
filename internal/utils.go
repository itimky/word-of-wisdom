package internal

import (
	"runtime/debug"

	"github.com/sirupsen/logrus"
)

func Recovery() {
	if rMsg := recover(); rMsg != nil {
		logrus.WithField("stacktrace", string(debug.Stack())).Error(rMsg)
	}
}
