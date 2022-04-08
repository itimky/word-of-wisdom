package client

import (
	"errors"
	"fmt"

	"github.com/tinylib/msgp/msgp"

	"github.com/itimky/word-of-wisom/api"
	guideapi "github.com/itimky/word-of-wisom/api/guide"
	srvapi "github.com/itimky/word-of-wisom/api/server"
	"github.com/itimky/word-of-wisom/pkg/gtp"
)

var errRetriesOver = errors.New("retries over")

type Client struct {
	tcpClient  tcpClient
	guideCount int
	retryCount int
}

func NewClient(client tcpClient, guideCount, retryCount int) *Client {
	return &Client{
		tcpClient:  client,
		guideCount: guideCount,
		retryCount: retryCount,
	}
}

func (s *Client) MakeRequest(reqType srvapi.RequestType) (msgp.Raw, error) {
	response, err := s.makeRequest(reqType, nil)
	if err != nil {
		return nil, fmt.Errorf("make request: %w", err)
	}

	if response.Granted {
		return response.Payload, nil
	}

	puzzleSolution, err := s.solvePuzzle(*response.Puzzle)
	if err != nil {
		return nil, fmt.Errorf("solve puzzle")
	}

	for i := 0; i < s.retryCount; i++ {
		response, err = s.makeRequest(reqType, puzzleSolution)
		if err != nil {
			return nil, fmt.Errorf("quote request: %w", err)
		}

		if response.Granted {
			return response.Payload, nil
		}
	}

	return nil, errRetriesOver
}

func (s *Client) makeRequest(reqType srvapi.RequestType, solution *gtp.PuzzleSolution) (*serverResponse, error) {
	apiSolution := convertPuzzleSolutionToSrvPuzzleSolution(solution)
	request := srvapi.RequestMsg{Type: reqType, PuzzleSolution: apiSolution}

	responseMsg, err := s.tcpClient.RequestServer(request)
	if err != nil {
		return nil, fmt.Errorf("request server: %w", err)
	}

	response, err := convertResponseMsgToServerResponse(responseMsg)
	if err != nil {
		return nil, fmt.Errorf("convert response msg to server response")
	}

	return response, nil
}

func (s *Client) solvePuzzle(puzzle gtp.Puzzle) (*gtp.PuzzleSolution, error) {
	hash := puzzle.InitialHash
	for i := 1; i < puzzle.TourLength+1; i++ {
		request := guideapi.RequestMsg{
			PreviousHash: api.Hash(hash),
			TourNumber:   byte(i),
			TourLength:   byte(puzzle.TourLength),
		}
		guideIdx := gtp.GuideIndex(hash, s.guideCount)
		response, err := s.tcpClient.RequestGuideByIdx(guideIdx, request)
		if err != nil {
			return nil, fmt.Errorf("request guide by index: %w", err)
		}
		hash = gtp.Hash(response.Hash)
	}

	return &gtp.PuzzleSolution{
		InitialHash: puzzle.InitialHash,
		LastHash:    hash,
	}, nil
}
