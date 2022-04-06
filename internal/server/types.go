package server

import (
	"fmt"
	srvcontracts "word-of-wisom/api/server"
	"word-of-wisom/pkg/gtp"

	"github.com/tinylib/msgp/msgp"
)

type TourCompleteRequest struct {
	InitialHash gtp.Hash
	LastHash    gtp.Hash
}

func NewTourCompleteRequestFromMsg(msg srvcontracts.RequestMsg) (TourCompleteRequest, error) {
	var tourCompleteRequest TourCompleteRequest
	var requestPayload srvcontracts.TourCompletePayload
	if _, err := requestPayload.UnmarshalMsg(msg.Payload); err != nil {
		return tourCompleteRequest, fmt.Errorf("unmarshal tour complete request payload: %w", err)
	}

	tourCompleteRequest.InitialHash = requestPayload.InitialHash
	tourCompleteRequest.LastHash = requestPayload.LastHash

	return tourCompleteRequest, nil
}

type encodabler interface {
	Encodable() (msgp.Encodable, error)
}

type ServiceRestrictedResponse struct {
	InitialHash gtp.Hash
	TourLength  int
}

func (r ServiceRestrictedResponse) Encodable() (msgp.Encodable, error) {
	serviceRestrictedPayload := srvcontracts.ServiceRestrictedPayload{InitialHash: r.InitialHash, TourLength: byte(r.TourLength)}

	responsePayload, err := serviceRestrictedPayload.MarshalMsg(nil)
	if err != nil {
		return nil, fmt.Errorf("marshal restricted response payload: %w", err)
	}

	responseMsg := srvcontracts.ResponseMsg{Type: byte(srvcontracts.ServiceRestricted), Payload: responsePayload}

	return &responseMsg, nil
}

type ServiceGrantedResponse struct {
	Quote string
}

func (r ServiceGrantedResponse) Encodable() (msgp.Encodable, error) {
	serviceGrantedPayload, err := srvcontracts.ServiceGrantedPayload{Quote: r.Quote}.MarshalMsg(nil)
	if err != nil {
		return nil, fmt.Errorf("marshal service granted payload: %w", err)
	}

	responseMsg := srvcontracts.ResponseMsg{Type: byte(srvcontracts.ServiceGranted), Payload: serviceGrantedPayload}
	return &responseMsg, nil
}

type UnsupportedRequest struct{}

func (r UnsupportedRequest) Encodable() (msgp.Encodable, error) {
	responseMsg := srvcontracts.ResponseMsg{Type: byte(srvcontracts.UnsupportedRequest)}
	return &responseMsg, nil
}
