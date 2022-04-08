package server

import (
	"fmt"
	"time"

	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/pool/goroutine"
	"github.com/sirupsen/logrus"
	"github.com/tinylib/msgp/msgp"

	srvapi "github.com/itimky/word-of-wisom/api/server"
	gtpserver "github.com/itimky/word-of-wisom/internal/gtp/server"
	"github.com/itimky/word-of-wisom/internal/tcp"
)

type Server struct {
	gnet.BuiltinEventEngine
	pool *goroutine.Pool

	addr      string
	multicore bool
	timeout   time.Duration

	gtpServer       gtpServer
	quoteRepository quoteRepository
}

func NewServer(
	pool *goroutine.Pool,
	addr string,
	multicore bool,
	timeout time.Duration,
	gtpServer gtpServer,
	quoteSvc quoteRepository,
) *Server {
	return &Server{
		pool:            pool,
		addr:            addr,
		multicore:       multicore,
		timeout:         timeout,
		gtpServer:       gtpServer,
		quoteRepository: quoteSvc,
	}
}

func (s *Server) Run() error {
	err := gnet.Run(s, s.addr, gnet.WithOptions(gnet.Options{Multicore: s.multicore, TCPKeepAlive: s.timeout}))
	if err != nil {
		return fmt.Errorf("gnet run: %w", err)
	}
	return nil
}

func (s *Server) OnBoot(eng gnet.Engine) gnet.Action {
	logrus.Infof("server (multi-core=%t timeout=%v) is listening on %s\n", s.multicore, s.timeout, s.addr)
	return gnet.None
}

func (s *Server) OnTraffic(conn gnet.Conn) gnet.Action {
	err := s.pool.Submit(func() {
		if err := s.handleConnection(conn); err != nil {
			logrus.WithError(err).Error("handle connection")
		}
		if err := conn.Close(nil); err != nil {
			logrus.WithError(err).Error("close connection")
		}
	})

	if err != nil {
		logrus.WithError(err).Error("submit")
		return gnet.Close
	}

	return gnet.None
}

func (s *Server) handleConnection(conn gnet.Conn) error {
	requestMsg := srvapi.RequestMsg{}
	if err := requestMsg.DecodeMsg(msgp.NewReader(conn)); err != nil {
		return fmt.Errorf("decode message: %w", err)
	}

	clientIP := tcp.GetClientIP(conn)

	responseMsg, err := s.handleRequest(clientIP, requestMsg)
	if err != nil {
		return fmt.Errorf("handle request: %w", err)
	}

	data, err := responseMsg.MarshalMsg(nil)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	err = conn.AsyncWrite(data, func(c gnet.Conn) error {
		return nil
	})
	if err != nil {
		return fmt.Errorf("write data: %w", err)
	}

	return nil
}

func (s *Server) handleRequest(clientIP string, requestMsg srvapi.RequestMsg) (*srvapi.ResponseMsg, error) {
	restrictedResponse, err := s.checkPuzzle(clientIP, requestMsg)
	if err != nil {
		return nil, fmt.Errorf("check puzzle: %w", err)
	}
	if restrictedResponse != nil {
		return restrictedResponse, nil
	}

	switch requestMsg.Type {
	case srvapi.Quote:
		quote := s.quoteRepository.Get()
		responseMsg, err := newQuoteResponse(quote)
		if err != nil {
			return responseMsg, fmt.Errorf("new quote response")
		}
		return responseMsg, nil
	default:
		return &srvapi.ResponseMsg{
			Type: srvapi.Unsupported,
		}, nil
	}
}

func (s *Server) checkPuzzle(clientIP string, requestMsg srvapi.RequestMsg) (*srvapi.ResponseMsg, error) {
	puzzleSolution := convertPuzzleSolution(requestMsg.PuzzleSolution)

	checkResult := s.gtpServer.CheckPuzzle(clientIP, puzzleSolution)
	if checkResult.Type != gtpserver.Ok {
		restrictedResponse, err := newRestrictedResponse(checkResult)
		if err != nil {
			return nil, fmt.Errorf("new restricted response")
		}
		return restrictedResponse, nil
	}

	return nil, nil
}
