package server

import (
	"fmt"
	"log"

	"github.com/itimky/word-of-wisom/internal/gtp/server"

	srvapi "github.com/itimky/word-of-wisom/api/server"
	"github.com/itimky/word-of-wisom/internal/tcp"
	"github.com/panjf2000/gnet/v2"
	"github.com/sirupsen/logrus"
	"github.com/tinylib/msgp/msgp"
)

type Server struct {
	gnet.BuiltinEventEngine

	addr      string
	multicore bool

	gtpServer       gtpServer
	quoteRepository quoteRepository
}

func NewServer(addr string, multicore bool, gtpServer gtpServer, quoteSvc quoteRepository) *Server {
	return &Server{
		addr:            addr,
		multicore:       multicore,
		gtpServer:       gtpServer,
		quoteRepository: quoteSvc,
	}
}

func (s *Server) Run() error {
	err := gnet.Run(s, s.addr, gnet.WithMulticore(s.multicore))
	if err != nil {
		return fmt.Errorf("gnet run: %w", err)
	}
	return nil
}

func (s *Server) OnBoot(eng gnet.Engine) gnet.Action {
	log.Printf("server with multi-core=%t is listening on %s\n", s.multicore, s.addr)
	return gnet.None
}

func (s *Server) OnTraffic(conn gnet.Conn) gnet.Action {
	if err := s.handleConnection(conn); err != nil {
		logrus.WithError(err).Error("handle connection")
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

	_, err = conn.Write(data)
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
	if checkResult.Type != server.Ok {
		restrictedResponse, err := newRestrictedResponse(checkResult)
		if err != nil {
			return nil, fmt.Errorf("new restricted response")
		}
		return restrictedResponse, nil
	}

	return nil, nil
}
