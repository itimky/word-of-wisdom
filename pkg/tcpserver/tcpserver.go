package tcpserver

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/tinylib/msgp/msgp"
)

type RequestHandler interface {
	HandleRequest(conn net.Conn) (msgp.Encodable, error)
}

type TCPServer struct {
	host    string
	port    string
	handler RequestHandler
}

func NewTCPServer(host, port string, handler RequestHandler) *TCPServer {
	return &TCPServer{host: host, port: port, handler: handler}
}

func (s *TCPServer) Run() error {
	l, err := net.Listen("tcp", s.hostPort())
	if err != nil {
		return fmt.Errorf("run server listener: %w", err)
	}
	defer func(l net.Listener) {
		err := l.Close()
		if err != nil {
			logrus.WithError(err).Error("close listener")
		}
	}(l)

	logrus.Infof("Listening on %v", s.hostPort())

	for {
		conn, err := l.Accept()
		if err != nil {
			logrus.WithError(err).Error("accept connection")
			continue
		}
		go s.handleRequest(conn)
	}
}

func (s *TCPServer) handleRequest(conn net.Conn) {
	defer recovery()

	data, err := s.handler.HandleRequest(conn)
	if err != nil {
		logrus.WithError(err).Error("handle request")
		return
	}

	writer := msgp.NewWriter(conn)
	if err = data.EncodeMsg(writer); err != nil {
		logrus.WithError(err).WithField("data", data).Error("encode msg")
		return
	}
	if err = writer.Flush(); err != nil {
		logrus.WithError(err).Error("flush response")
		return
	}
}

func (s *TCPServer) hostPort() string {
	return fmt.Sprintf("%v:%v", s.host, s.port)
}
