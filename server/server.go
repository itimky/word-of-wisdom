package server

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tinylib/msgp/msgp"
	"math/rand"
	"net"
	"time"
	"word-of-wisom/internal"
	srvcontracts "word-of-wisom/internal/contracts/server"
)

type Server struct {
	host                        string
	port                        string
	secretLength                int
	secretUpdateIntervalMinutes time.Duration
	secret                      string
	tourLength                  int
	guideSecrets                []string
	hashCalc                    HashCalc

	rand *rand.Rand
}

func NewServer(
	host string,
	port string,
	secretLength int,
	secretUpdateIntervalMinutes time.Duration,
	tourLength int,
	guideSecrets []string,
	hashCalc HashCalc,
	rand *rand.Rand,
) *Server {
	return &Server{
		host:                        host,
		port:                        port,
		secretLength:                secretLength,
		secretUpdateIntervalMinutes: secretUpdateIntervalMinutes,
		tourLength:                  tourLength,
		guideSecrets:                guideSecrets,
		hashCalc:                    hashCalc,

		rand: rand,
	}
}

func (s *Server) Run() error {
	err := s.updateSecret()
	if err != nil {
		return fmt.Errorf("run server: %w", err)
	}

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

	go s.periodicSecretUpdate()

	for {
		conn, err := l.Accept()
		if err != nil {
			logrus.WithError(err).Error("accept connection")
			continue
		}
		go s.handleRequest(conn)
	}
}

func (s *Server) handleRequest(conn net.Conn) {
	defer internal.Recovery()
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			logrus.WithError(err).Error("close connection")
		}
	}(conn)
	logrus.Debug("new connection")

	request := srvcontracts.RequestMsg{}
	if err := request.DecodeMsg(msgp.NewReader(conn)); err != nil {
		logrus.WithError(err).Error("decode message")
		return
	}

	logrus.Debug(request)

	var data msgp.Encodable
	var err error
	switch srvcontracts.RequestType(request.Type) {
	case srvcontracts.InitialRequest:
		data, err = s.initialRequestHandler(conn)
	case srvcontracts.TourCompleteRequest:
		data, err = s.tourCompleteRequestHandler(conn, request)
	default:
		logrus.WithField("type", request.Type).Error("unsupported request type")
		data, err = s.unsupportedRequestHandler(conn)
	}

	if err != nil {
		logrus.WithError(err).Error("handle response")
		return
	}

	writer := msgp.NewWriter(conn)
	if err := data.EncodeMsg(writer); err != nil {
		logrus.WithError(err).WithField("data", data).Error("encode msg")
		return
	}
	if err = writer.Flush(); err != nil {
		logrus.WithError(err).Error("flush response")
	}
}

func (s *Server) hostPort() string {
	return s.host + ":" + s.port
}

func (s *Server) getClientIP(conn net.Conn) string {
	return conn.RemoteAddr().(*net.TCPAddr).IP.String()
}

func (s *Server) updateSecret() error {
	secret, err := randomSecret(s.secretLength)
	if err != nil {
		return fmt.Errorf("update secret: %w", err)
	}
	s.secret = secret
	return nil
}

func (s *Server) periodicSecretUpdate() {
	ticker := time.NewTicker(s.secretUpdateIntervalMinutes * time.Minute)
	for range ticker.C {
		err := s.updateSecret()
		if err != nil {
			logrus.WithError(err).Error("periodic secret update")
		}
	}
}
