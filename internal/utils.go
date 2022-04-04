package internal

import (
	"github.com/sirupsen/logrus"
	"runtime/debug"
)

func Recovery() {
	if rMsg := recover(); rMsg != nil {
		logrus.WithField("stacktrace", string(debug.Stack())).Error(rMsg)
	}
}
