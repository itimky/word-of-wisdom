package guide

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tinylib/msgp/msgp"
	"net"
	guidecontracts "word-of-wisom/internal/contracts/guide"
)

type Guide struct {
	host     string
	port     string
	secret   string
	hashCalc HashCalc
}

func NewGuide(host, port, secret string, hashCalc HashCalc) *Guide {
	return &Guide{
		host:     host,
		port:     port,
		secret:   secret,
		hashCalc: hashCalc,
	}
}

func (g *Guide) Run() error {
	l, err := net.Listen("tcp", g.hostPort())
	if err != nil {
		return fmt.Errorf("run server listener: %w", err)
	}
	defer func(l net.Listener) {
		err := l.Close()
		if err != nil {
			logrus.WithError(err).Error("close listener")
		}
	}(l)
	logrus.Infof("Listening on %v", g.hostPort())

	for {
		conn, err := l.Accept()
		if err != nil {
			logrus.WithError(err).Error("accept connection")
			continue
		}
		go g.handleRequest(conn)
	}
}

func (g *Guide) handleRequest(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			logrus.WithError(err).Error("close connection")
		}
	}(conn)

	request := guidecontracts.RequestMsg{}
	if err := request.DecodeMsg(msgp.NewReader(conn)); err != nil {
		logrus.WithError(err).Error("decode message")
		return
	}

	logrus.Debug(request)

	response := g.tourGuideHandler(conn, request)

	writer := msgp.NewWriter(conn)
	if err := response.EncodeMsg(writer); err != nil {
		logrus.WithError(err).WithField("data", response).Error("encode msg")
		return
	}
	if err := writer.Flush(); err != nil {
		logrus.WithError(err).Error("flush response")
	}
}

func (g *Guide) getClientIP(conn net.Conn) string {
	return conn.RemoteAddr().(*net.TCPAddr).IP.String()
}

func (g *Guide) hostPort() string {
	return g.host + ":" + g.port
}
