package server

import (
	"fmt"

	"github.com/tinylib/msgp/msgp"

	"github.com/itimky/word-of-wisom/api"
	srvapi "github.com/itimky/word-of-wisom/api/server"
	"github.com/itimky/word-of-wisom/internal/gtp/server"
	"github.com/itimky/word-of-wisom/pkg/gtp"
)

func convertPuzzleSolution(solution *srvapi.PuzzleSolution) *gtp.PuzzleSolution {
	if solution == nil {
		return nil
	}

	return &gtp.PuzzleSolution{
		InitialHash: gtp.Hash(solution.InitialHash),
		LastHash:    gtp.Hash(solution.LastHash),
	}
}

func newRestrictedResponse(checkResult server.PuzzleCheckResult) (*srvapi.ResponseMsg, error) {
	var responseType srvapi.ResponseType
	var payload msgp.Raw
	switch checkResult.Type {
	case server.Restricted:
		responseType = srvapi.Restricted
		puzzleResponse := &srvapi.PuzzleResponse{
			InitialHash: api.Hash(checkResult.Puzzle.InitialHash),
			TourLength:  byte(checkResult.Puzzle.TourLength),
		}
		var err error
		payload, err = puzzleResponse.MarshalMsg(nil)
		if err != nil {
			return nil, fmt.Errorf("marshal msg: %w", err)
		}
	case server.WrongSolution:
		responseType = srvapi.WrongPuzzle
	case server.Ok:
		return nil, fmt.Errorf("unexpected result type (%v)", checkResult.Type)
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
		return nil, fmt.Errorf("marshal msg: %w", err)
	}
	return &srvapi.ResponseMsg{
		Type:    srvapi.Granted,
		Payload: responsePayload,
	}, nil
}
