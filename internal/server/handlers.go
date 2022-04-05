package server

import (
	"fmt"
	"net"
	srvcontracts "word-of-wisom/api/server"

	"github.com/tinylib/msgp/msgp"
)

func (s *Server) initialRequestHandler(conn net.Conn) (msgp.Encodable, error) {
	initialHash := s.hashCalc.CalcInitialHash(s.getClientIP(conn), s.tourLength, s.secret)
	serviceRestrictedPayload := srvcontracts.ServiceRestrictedPayload{InitialHash: initialHash, TourLength: byte(s.tourLength)}

	responsePayload, err := serviceRestrictedPayload.MarshalMsg(nil)
	if err != nil {
		return nil, fmt.Errorf("marshal restricted response payload: %w", err)
	}

	response := srvcontracts.ResponseMsg{Type: byte(srvcontracts.ServiceRestricted), Payload: responsePayload}

	return &response, nil
}

func (s *Server) tourCompleteRequestHandler(conn net.Conn, requestMsg srvcontracts.RequestMsg) (msgp.Encodable, error) {
	var requestPayload srvcontracts.TourCompletePayload
	if _, err := requestPayload.UnmarshalMsg(requestMsg.Payload); err != nil {
		return nil, fmt.Errorf("unmarshal tour complete request payload: %w", err)
	}

	var response msgp.Encodable
	if s.hashCalc.VerifyHash(requestPayload.InitialHash, requestPayload.LastHash, s.tourLength, s.getClientIP(conn), s.secret, s.guideSecrets) {
		serviceGrantedPayload, err := srvcontracts.ServiceGrantedPayload{Quote: quotes[s.rand.Intn(len(quotes))]}.MarshalMsg(nil)
		if err != nil {
			return nil, fmt.Errorf("marshal service granted payload: %w", err)
		}
		response = &srvcontracts.ResponseMsg{Type: byte(srvcontracts.ServiceGranted), Payload: serviceGrantedPayload}
	} else {
		var err error
		response, err = s.initialRequestHandler(conn)
		if err != nil {
			return nil, fmt.Errorf("repeated service restricted: %w", err)
		}
	}

	return response, nil
}

func (s *Server) unsupportedRequestHandler() msgp.Encodable {
	response := srvcontracts.ResponseMsg{Type: byte(srvcontracts.UnsupportedRequest)}
	return &response
}
