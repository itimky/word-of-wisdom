package server

import (
	"fmt"
	"log"

	srvapi "github.com/itimky/word-of-wisom/api/server"
	"github.com/itimky/word-of-wisom/internal/service/shield"
	"github.com/itimky/word-of-wisom/internal/tcp"
	"github.com/panjf2000/gnet/v2"
	"github.com/sirupsen/logrus"
	"github.com/tinylib/msgp/msgp"
)

type Server struct {
	gnet.BuiltinEventEngine
	//pool   *goroutine.Pool
	//eng    gnet.Engine
	shield *shield.Service

	addr      string
	multicore bool
}

func NewServer(shield *shield.Service, addr string, multicore bool) *Server {
	return &Server{
		shield:    shield,
		addr:      addr,
		multicore: multicore,
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
	//s.eng = eng
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

	responseMsg, err := s.handleRequest(conn, requestMsg)
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

func (s *Server) handleRequest(conn gnet.Conn, requestMsg srvapi.RequestMsg) (srvapi.ResponseMsg, error) {
	var responseMsg srvapi.ResponseMsg
	var err error
	switch requestMsg.Type {
	case srvapi.Initial:
		result := s.shield.HandleInitial(tcp.GetClientIP(conn))
		responseMsg, err = convertInitialResultToResponseMsg(result)
		if err != nil {
			return responseMsg, fmt.Errorf("convert initial result")
		}
	case srvapi.TourComplete:
		var request shield.TourCompleteRequest
		request, err = convertRequestMsgToTourCompleteRequest(requestMsg)
		if err != nil {
			return responseMsg, fmt.Errorf("convert tour complete request: %w", err)
		}
		result := s.shield.HandleTourComplete(request, tcp.GetClientIP(conn))
		responseMsg, err = convertTourCompleteResultToResponseMsg(result)
		if err != nil {
			return responseMsg, fmt.Errorf("convert tour complete result")
		}
	default:
		responseMsg.Type = srvapi.Unsupported
	}

	return responseMsg, nil
}
