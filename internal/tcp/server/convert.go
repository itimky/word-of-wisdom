package server

import (
	"fmt"

	"github.com/itimky/word-of-wisom/api"
	srvapi "github.com/itimky/word-of-wisom/api/server"
	"github.com/itimky/word-of-wisom/internal/service/shield"
	"github.com/itimky/word-of-wisom/pkg/gtp"
)

func convertInitialResultToResponseMsg(result shield.InitialResult) (srvapi.ResponseMsg, error) {
	var responseMsg srvapi.ResponseMsg
	serviceRestrictedPayload := srvapi.ServiceRestrictedPayload{
		InitialHash: api.Hash(result.InitialHash),
		TourLength:  byte(result.TourLength),
	}

	responsePayload, err := serviceRestrictedPayload.MarshalMsg(nil)
	if err != nil {
		return responseMsg, fmt.Errorf("marshal restricted response payload: %w", err)
	}

	responseMsg.Type = srvapi.Restricted
	responseMsg.Payload = responsePayload

	return responseMsg, nil
}

func convertRequestMsgToTourCompleteRequest(msg srvapi.RequestMsg) (shield.TourCompleteRequest, error) {
	var tourCompleteRequest shield.TourCompleteRequest
	var requestPayload srvapi.TourCompletePayload
	if _, err := requestPayload.UnmarshalMsg(msg.Payload); err != nil {
		return tourCompleteRequest, fmt.Errorf("unmarshal tour complete request payload: %w", err)
	}

	tourCompleteRequest.InitialHash = gtp.Hash(requestPayload.InitialHash)
	tourCompleteRequest.LastHash = gtp.Hash(requestPayload.LastHash)

	return tourCompleteRequest, nil
}

func convertTourCompleteResultToResponseMsg(result shield.TourCompleteResult) (srvapi.ResponseMsg, error) {
	var responseMsg srvapi.ResponseMsg
	if result.Granted {
		responseMsg.Type = srvapi.Granted
		payload := srvapi.ServiceGrantedPayload{Quote: result.Quote}
		responsePayload, err := payload.MarshalMsg(nil)
		if err != nil {
			return responseMsg, fmt.Errorf("marshal service granted payload: %w", err)
		}
		responseMsg.Payload = responsePayload
	} else {
		responseMsg.Type = srvapi.WrongPuzzle
	}

	return responseMsg, nil
}
