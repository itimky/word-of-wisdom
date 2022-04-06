package tcpserver

import (
	"net"
	"runtime/debug"

	"github.com/sirupsen/logrus"
)

func recovery() {
	if rMsg := recover(); rMsg != nil {
		logrus.WithField("stacktrace", string(debug.Stack())).Error(rMsg)
	}
}

func GetClientIP(conn net.Conn) string {
	return conn.RemoteAddr().(*net.TCPAddr).IP.String()
}
