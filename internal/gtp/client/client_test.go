package client

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/tinylib/msgp/msgp"

	"github.com/itimky/word-of-wisom/api"
	guideapi "github.com/itimky/word-of-wisom/api/guide"
	srvapi "github.com/itimky/word-of-wisom/api/server"
	"github.com/itimky/word-of-wisom/internal/gtp/client/mocks"
	"github.com/itimky/word-of-wisom/pkg/gtp"
	"github.com/itimky/word-of-wisom/pkg/testutils"
)

type ServiceSuite struct {
	suite.Suite
	svc           *Client
	tcpClientMock *mocks.TcpClient
}

func (s *ServiceSuite) SetupSuite() {
	s.svc = &Client{
		guideCount: 2,
	}
}

func (s *ServiceSuite) SetupTest() {
	s.tcpClientMock = &mocks.TcpClient{}
	s.svc.tcpClient = s.tcpClientMock
}

func (s *ServiceSuite) TestService_makeRequest___Restricted_NotGranted() {
	tourLength := 5
	hash := testutils.HexHash("820888B1A040503A82AFA97EB0AE59E8214866C2D74F3DBC705A002FB17C86E9")
	request := srvapi.RequestMsg{Type: srvapi.Quote, PuzzleSolution: nil}
	responsePayload := srvapi.PuzzleResponse{InitialHash: api.Hash(hash), TourLength: byte(tourLength)}
	payloadRaw, err := responsePayload.MarshalMsg(nil)
	s.NoError(err)
	responseMsg := &srvapi.ResponseMsg{
		Type:    srvapi.Restricted,
		Payload: payloadRaw,
	}

	s.tcpClientMock.EXPECT().RequestServer(request).Return(responseMsg, nil)

	response, err := s.svc.makeRequest(request.Type, nil)
	s.NoError(err)
	s.False(response.Granted)
	s.NotNil(response.Puzzle)
	s.Equal(tourLength, response.Puzzle.TourLength)
	s.Equal(hash, response.Puzzle.InitialHash)
	s.Nil(response.Payload)
}

func (s *ServiceSuite) TestService_makeRequest___WrongPuzzle_NotGranted() {
	tourLength := 5
	initialHash := testutils.HexHash("820888B1A040503A82AFA97EB0AE59E8214866C2D74F3DBC705A002FB17C86E9")
	lastHash := testutils.HexHash("920888B1A040503A82AFA97EB0AE59E8214866C2D74F3DBC705A002FB17C86E9")
	puzzleSolution := &gtp.PuzzleSolution{
		InitialHash: initialHash,
		LastHash:    lastHash,
	}
	apiPuzzleSolution := &srvapi.PuzzleSolution{
		InitialHash: api.Hash(initialHash),
		LastHash:    api.Hash(lastHash),
	}

	request := srvapi.RequestMsg{Type: srvapi.Quote, PuzzleSolution: apiPuzzleSolution}
	responsePayload := srvapi.PuzzleResponse{InitialHash: api.Hash(initialHash), TourLength: byte(tourLength)}
	payloadRaw, err := responsePayload.MarshalMsg(nil)
	s.NoError(err)
	responseMsg := &srvapi.ResponseMsg{
		Type:    srvapi.WrongPuzzle,
		Payload: payloadRaw,
	}

	s.tcpClientMock.EXPECT().RequestServer(request).Return(responseMsg, nil)

	response, err := s.svc.makeRequest(request.Type, puzzleSolution)
	s.NoError(err)
	s.False(response.Granted)
	s.NotNil(response.Puzzle)
	s.Equal(tourLength, response.Puzzle.TourLength)
	s.Equal(initialHash, response.Puzzle.InitialHash)
	s.Nil(response.Payload)
}

func (s *ServiceSuite) TestService_solvePuzzle() {
	tourLength := 5
	guideCount := 2
	hashes := []gtp.Hash{
		// Server initial hash
		testutils.HexHash("020888B1A040503A82AFA97EB0AE59E8214866C2D74F3DBC705A002FB17C86E9"),

		// Guide hashes
		testutils.HexHash("120888B1A040503A82AFA97EB0AE59E8214866C2D74F3DBC705A002FB17C86E9"),
		testutils.HexHash("220888B1A040503A82AFA97EB0AE59E8214866C2D74F3DBC705A002FB17C86E9"),
		testutils.HexHash("320888B1A040503A82AFA97EB0AE59E8214866C2D74F3DBC705A002FB17C86E9"),
		testutils.HexHash("420888B1A040503A82AFA97EB0AE59E8214866C2D74F3DBC705A002FB17C86E9"),
		testutils.HexHash("520888B1A040503A82AFA97EB0AE59E8214866C2D74F3DBC705A002FB17C86E9"),
	}

	for i := 0; i < len(hashes)-1; i++ {
		request := guideapi.RequestMsg{
			PreviousHash: api.Hash(hashes[i]),
			TourNumber:   byte(i + 1),
			TourLength:   byte(tourLength),
		}
		response := &guideapi.ResponseMsg{
			Hash: api.Hash(hashes[i+1]),
		}

		guideIdx := gtp.GuideIndex(hashes[i], guideCount)
		s.tcpClientMock.EXPECT().RequestGuideByIdx(guideIdx, request).Return(response, nil)
	}

	puzzle := gtp.Puzzle{
		InitialHash: hashes[0],
		TourLength:  tourLength,
	}

	solution, err := s.svc.solvePuzzle(puzzle)
	s.NoError(err)
	s.Equal(hashes[0], solution.InitialHash)
	s.Equal(hashes[len(hashes)-1], solution.LastHash)
}

func (s *ServiceSuite) TestService_makeRequest___Granted_Granted() {
	quote := "some quote"
	request := srvapi.RequestMsg{Type: srvapi.Quote, PuzzleSolution: nil}
	responsePayload := srvapi.QuoteResponse{Quote: quote}
	payloadRaw, err := responsePayload.MarshalMsg(nil)
	s.NoError(err)
	responseMsg := &srvapi.ResponseMsg{
		Type:    srvapi.Granted,
		Payload: payloadRaw,
	}

	s.tcpClientMock.EXPECT().RequestServer(request).Return(responseMsg, nil)

	response, err := s.svc.makeRequest(request.Type, nil)
	s.NoError(err)
	s.True(response.Granted)
	s.Nil(response.Puzzle)
	s.NotNil(response.Payload)
	s.Equal(msgp.Raw(payloadRaw), response.Payload)
}

func (s *ServiceSuite) TestService_MakeRequest() {
	quote := "some quote"
	request := srvapi.RequestMsg{Type: srvapi.Quote, PuzzleSolution: nil}
	responsePayload := srvapi.QuoteResponse{Quote: quote}
	rawPayload, err := responsePayload.MarshalMsg(nil)
	s.NoError(err)
	responseMsg := &srvapi.ResponseMsg{
		Type:    srvapi.Granted,
		Payload: rawPayload,
	}

	s.tcpClientMock.EXPECT().RequestServer(request).Return(responseMsg, nil)

	response, err := s.svc.MakeRequest(request.Type)
	s.NoError(err)
	s.Equal(msgp.Raw(rawPayload), response)
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
