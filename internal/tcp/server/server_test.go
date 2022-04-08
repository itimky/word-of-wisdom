package server

import (
	"testing"

	"github.com/itimky/word-of-wisom/pkg/gtp"

	"github.com/stretchr/testify/suite"
	"github.com/tinylib/msgp/msgp"

	"github.com/itimky/word-of-wisom/api"
	srvapi "github.com/itimky/word-of-wisom/api/server"
	"github.com/itimky/word-of-wisom/internal/gtp/server"
	"github.com/itimky/word-of-wisom/internal/tcp/server/mocks"
	"github.com/itimky/word-of-wisom/pkg/testutils"
)

type ServerSuite struct {
	suite.Suite
	shieldMock       *mocks.ShieldService
	quoteServiceMock *mocks.QuoteService
	srv              *Server
	clientIP         string
}

func (s *ServerSuite) SetupSuite() {
	// Init shieldMock & quoteServiceMock in SetupTest
	s.srv = &Server{}
	s.clientIP = "127.0.0.1"
}

func (s *ServerSuite) SetupTest() {
	s.shieldMock = &mocks.ShieldService{}
	s.quoteServiceMock = &mocks.QuoteService{}
	s.srv.gtpServer = s.shieldMock
	s.srv.quoteRepository = s.quoteServiceMock
}

func (s *ServerSuite) TestServer_handleRequest__Restricted() {
	request := srvapi.RequestMsg{Type: srvapi.Quote}

	hash := testutils.HexHash("820888B1A040503A82AFA97EB0AE59E8214866C2D74F3DBC705A002FB17C86E9")
	tourLength := 7
	checkResult := server.PuzzleCheckResult{
		Type: server.Restricted,
		Puzzle: &gtp.Puzzle{
			InitialHash: hash,
			TourLength:  tourLength,
		},
	}

	payload := srvapi.PuzzleResponse{
		InitialHash: api.Hash(hash),
		TourLength:  byte(tourLength),
	}
	rawPayload, err := payload.MarshalMsg(nil)
	s.NoError(err)

	s.shieldMock.EXPECT().CheckPuzzle(s.clientIP, (*gtp.PuzzleSolution)(nil)).Return(checkResult)

	response, err := s.srv.handleRequest(s.clientIP, request)
	s.NoError(err)
	s.Equal(srvapi.Restricted, response.Type)
	s.Equal(msgp.Raw(rawPayload), response.Payload)
}

func (s *ServerSuite) TestWrongSolution() {
	initialHash := testutils.HexHash("820888B1A040503A82AFA97EB0AE59E8214866C2D74F3DBC705A002FB17C86E9")
	latHash := testutils.HexHash("820888B1A040503A82AFA97EB0AE59E8214866C2D74F3DBC705A002FB17C86E5")
	puzzleSolution := gtp.PuzzleSolution{
		InitialHash: initialHash,
		LastHash:    latHash,
	}
	checkResult := server.PuzzleCheckResult{
		Type: server.WrongSolution,
	}
	apiPuzzleSolution := &srvapi.PuzzleSolution{
		InitialHash: api.Hash(initialHash),
		LastHash:    api.Hash(latHash),
	}

	request := srvapi.RequestMsg{Type: srvapi.Quote, PuzzleSolution: apiPuzzleSolution}

	s.shieldMock.EXPECT().CheckPuzzle(s.clientIP, &puzzleSolution).Return(checkResult)

	response, err := s.srv.handleRequest(s.clientIP, request)
	s.NoError(err)
	s.Equal(srvapi.WrongPuzzle, response.Type)
	s.Nil(response.Payload)
}

func (s *ServerSuite) TestGranted() {
	initialHash := testutils.HexHash("820888B1A040503A82AFA97EB0AE59E8214866C2D74F3DBC705A002FB17C86E9")
	latHash := testutils.HexHash("820888B1A040503A82AFA97EB0AE59E8214866C2D74F3DBC705A002FB17C86E5")
	quote := "Some quote"
	puzzleSolution := &gtp.PuzzleSolution{
		InitialHash: initialHash,
		LastHash:    latHash,
	}
	checkResult := server.PuzzleCheckResult{
		Type: server.Ok,
	}
	apiPuzzleSolution := &srvapi.PuzzleSolution{
		InitialHash: api.Hash(initialHash),
		LastHash:    api.Hash(latHash),
	}
	request := srvapi.RequestMsg{Type: srvapi.Quote, PuzzleSolution: apiPuzzleSolution}

	s.shieldMock.EXPECT().CheckPuzzle(s.clientIP, puzzleSolution).Return(checkResult)
	s.quoteServiceMock.EXPECT().Get().Return(quote)

	responsePayload := srvapi.QuoteResponse{Quote: quote}
	rawResponsePayload, err := responsePayload.MarshalMsg(nil)
	s.NoError(err)

	response, err := s.srv.handleRequest(s.clientIP, request)
	s.NoError(err)
	s.Equal(srvapi.Granted, response.Type)
	s.Equal(msgp.Raw(rawResponsePayload), response.Payload)
}

func (s *ServerSuite) TestUnsupported() {
	request := srvapi.RequestMsg{Type: srvapi.RequestType(10)}
	checkResult := server.PuzzleCheckResult{
		Type: server.Ok,
	}
	s.shieldMock.EXPECT().CheckPuzzle(s.clientIP, (*gtp.PuzzleSolution)(nil)).Return(checkResult)
	response, err := s.srv.handleRequest(s.clientIP, request)
	s.NoError(err)
	s.Equal(srvapi.Unsupported, response.Type)
	s.Nil(response.Payload)
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerSuite))
}
