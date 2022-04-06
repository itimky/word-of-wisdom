package server

import (
	"fmt"
	"net"

	srvapi "github.com/itimky/word-of-wisom/api/server"
	"github.com/itimky/word-of-wisom/pkg/tcpserver"

	"github.com/sirupsen/logrus"
	"github.com/tinylib/msgp/msgp"
)

type quoteGetter interface {
	Get() string
}

type Server struct {
	secretLength int
	secret       string
	tourLength   int
	guideSecrets []string
	hashCalc     hashCalc
	quoteGetter  quoteGetter
}

func NewServer(
	secretLength int,
	tourLength int,
	guideSecrets []string,
	hashCalc hashCalc,
	quoteGetter quoteGetter,
) *Server {
	return &Server{
		secretLength: secretLength,
		tourLength:   tourLength,
		guideSecrets: guideSecrets,
		hashCalc:     hashCalc,
		quoteGetter:  quoteGetter,
	}
}

func (s *Server) HandleRequest(conn net.Conn) (msgp.Encodable, error) {
	requestMsg := srvapi.RequestMsg{}
	if err := requestMsg.DecodeMsg(msgp.NewReader(conn)); err != nil {
		return nil, fmt.Errorf("decode message: %w", err)
	}

	logrus.Debug(requestMsg)

	var result encodabler
	switch srvapi.RequestType(requestMsg.Type) {
	case srvapi.InitialRequest:
		result = s.initialRequestHandler(tcpserver.GetClientIP(conn))
	case srvapi.TourCompleteRequest:
		request, err := NewTourCompleteRequestFromMsg(requestMsg)
		if err != nil {
			return nil, fmt.Errorf("new tour complete request from msg: %w", err)
		}
		result = s.tourCompleteRequestHandler(request, tcpserver.GetClientIP(conn))
	default:
		result = s.unsupportedRequestHandler()
	}

	responseMsg, err := result.Encodable()
	if err != nil {
		return nil, fmt.Errorf("encodable: %w", err)
	}

	return responseMsg, nil
}

func (s *Server) UpdateSecret() error {
	secret, err := randomSecret(s.secretLength)
	if err != nil {
		return fmt.Errorf("update secret: %w", err)
	}
	s.secret = secret
	return nil
}

func (s *Server) initialRequestHandler(clientIP string) ServiceRestrictedResponse {
	initialHash := s.hashCalc.CalcInitialHash(clientIP, s.tourLength, s.secret)
	return ServiceRestrictedResponse{InitialHash: initialHash, TourLength: s.tourLength}
}

func (s *Server) tourCompleteRequestHandler(request TourCompleteRequest, clientIP string) encodabler {
	var response encodabler
	if s.hashCalc.VerifyHash(request.InitialHash, request.LastHash, s.tourLength, clientIP, s.secret, s.guideSecrets) {
		response = ServiceGrantedResponse{Quote: s.quoteGetter.Get()}
	} else {
		response = s.initialRequestHandler(clientIP)
	}

	return response
}

func (s *Server) unsupportedRequestHandler() UnsupportedRequest {
	return UnsupportedRequest{}
}
