package client

import (
	"errors"
	"fmt"

	"github.com/itimky/word-of-wisom/api"
	srvapi "github.com/itimky/word-of-wisom/api/server"
	"github.com/itimky/word-of-wisom/pkg/gtp"
	"github.com/tinylib/msgp/msgp"
)

var errUnexpectedResponse = errors.New("unexpected response")

func convertResponseMsgToServerResponse(responseMsg *srvapi.ResponseMsg) (*serverResponse, error) {
	var response *serverResponse
	switch responseMsg.Type {
	case srvapi.Granted:
		response = &serverResponse{
			Granted: true,
			Puzzle:  nil,
			Payload: responseMsg.Payload,
		}
	case srvapi.Restricted, srvapi.WrongPuzzle:
		puzzle, err := convertResponseMsgPayloadToPuzzle(responseMsg.Payload)
		if err != nil {
			return nil, fmt.Errorf("convert response msg payload to puzzle")
		}
		response = &serverResponse{
			Granted: false,
			Puzzle:  puzzle,
			Payload: nil,
		}
	case srvapi.Unsupported:
	default:
		return nil, errUnexpectedResponse
	}

	return response, nil
}

func convertResponseMsgPayloadToPuzzle(payload msgp.Raw) (*gtp.Puzzle, error) {
	response := srvapi.PuzzleResponse{}
	_, err := response.UnmarshalMsg(payload)
	if err != nil {
		return nil, fmt.Errorf("unmarshal msg: %w", err)
	}

	return &gtp.Puzzle{
		InitialHash: gtp.Hash(response.InitialHash),
		TourLength:  int(response.TourLength),
	}, nil
}

func convertPuzzleSolutionToSrvPuzzleSolution(solution *gtp.PuzzleSolution) *srvapi.PuzzleSolution {
	if solution == nil {
		return nil
	}

	return &srvapi.PuzzleSolution{
		InitialHash: api.Hash(solution.InitialHash),
		LastHash:    api.Hash(solution.LastHash),
	}
}
