package server

import (
	"fmt"

	"github.com/tinylib/msgp/msgp"

	"github.com/itimky/word-of-wisom/api"
	srvapi "github.com/itimky/word-of-wisom/api/server"
	"github.com/itimky/word-of-wisom/internal/service/shield"
	"github.com/itimky/word-of-wisom/pkg/gtp"
)

func convertPuzzleSolution(solution msgp.Raw) (*shield.PuzzleSolution, error) {
	if solution == nil {
		return nil, nil
	}

	var requestPayload srvapi.PuzzleSolution
	if _, err := requestPayload.UnmarshalMsg(solution); err != nil {
		return nil, fmt.Errorf("unmarshal tour complete request payload: %w", err)
	}

	return &shield.PuzzleSolution{
		InitialHash: gtp.Hash(requestPayload.InitialHash),
		LastHash:    gtp.Hash(requestPayload.LastHash),
	}, nil
}

func newRestrictedResponse(checkResult shield.PuzzleCheckResult) (*srvapi.ResponseMsg, error) {
	var responseType srvapi.ResponseType
	var payload msgp.Raw
	switch checkResult.Type {
	case shield.Restricted:
		responseType = srvapi.Restricted
		puzzleResponse := &srvapi.PuzzleResponse{
			InitialHash: api.Hash(checkResult.Puzzle.InitialHash),
			TourLength:  byte(checkResult.Puzzle.TourLength),
		}
		var err error
		payload, err = puzzleResponse.MarshalMsg(nil)
		if err != nil {
			return nil, fmt.Errorf("marshal puzzle payload: %w", err)
		}
	case shield.WrongSolution:
		responseType = srvapi.WrongPuzzle
	case shield.Ok:
		return nil, fmt.Errorf("unexpected result type: %v", checkResult.Type)
	}

	return &srvapi.ResponseMsg{
		Type:    responseType,
		Payload: payload,
	}, nil
}

func newQuoteResponse(quote string) (*srvapi.ResponseMsg, error) {
	payload := srvapi.QuoteResponse{Quote: quote}
	responsePayload, err := payload.MarshalMsg(nil)
	if err != nil {
		return nil, fmt.Errorf("marshal service granted payload: %w", err)
	}
	return &srvapi.ResponseMsg{
		Type:    srvapi.Granted,
		Payload: responsePayload,
	}, nil
}
