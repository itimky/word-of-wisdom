package server

import (
	"testing"

	"github.com/itimky/word-of-wisom/api"
	srvapi "github.com/itimky/word-of-wisom/api/server"
	"github.com/itimky/word-of-wisom/internal/service/shield"
	"github.com/itimky/word-of-wisom/internal/tcp/server/mocks"
	"github.com/itimky/word-of-wisom/pkg/testutils"
	"github.com/stretchr/testify/suite"
	"github.com/tinylib/msgp/msgp"
)

type ServerSuite struct {
	suite.Suite
	shieldMock *mocks.ShieldService
	srv        *Server
	clientIP   string
}

func (s *ServerSuite) SetupSuite() {
	// Init shieldMock in SetupTest
	s.srv = &Server{}
	s.clientIP = "127.0.0.1"
}

func (s *ServerSuite) SetupTest() {
	s.shieldMock = &mocks.ShieldService{}
	s.srv.shield = s.shieldMock
}

func (s *ServerSuite) TestInitialRequest() {
	request := srvapi.RequestMsg{Type: srvapi.Initial}

	hash := testutils.HexHash("820888B1A040503A82AFA97EB0AE59E8214866C2D74F3DBC705A002FB17C86E9")
	tourLength := 7
	shieldResult := shield.InitialResult{
		InitialHash: hash,
		TourLength:  tourLength,
	}
	payload := srvapi.ServiceRestrictedPayload{
		InitialHash: api.Hash(hash),
		TourLength:  byte(tourLength),
	}
	rawPayload, err := payload.MarshalMsg(nil)
	s.NoError(err)

	s.shieldMock.EXPECT().HandleInitial(s.clientIP).Return(shieldResult)

	response, err := s.srv.handleRequest(s.clientIP, request)
	s.NoError(err)
	s.Equal(srvapi.Restricted, response.Type)
	s.Equal(msgp.Raw(rawPayload), response.Payload)
}

func (s *ServerSuite) TestWrongPuzzle() {
	initialHash := testutils.HexHash("820888B1A040503A82AFA97EB0AE59E8214866C2D74F3DBC705A002FB17C86E9")
	latHash := testutils.HexHash("820888B1A040503A82AFA97EB0AE59E8214866C2D74F3DBC705A002FB17C86E5")
	tourCompleteRequest := shield.TourCompleteRequest{
		InitialHash: initialHash,
		LastHash:    latHash,
	}
	tourCompleteResult := shield.TourCompleteResult{
		Granted: false,
	}
	payload := srvapi.TourCompletePayload{
		InitialHash: api.Hash(initialHash),
		LastHash:    api.Hash(latHash),
	}
	rawPayload, err := payload.MarshalMsg(nil)
	s.NoError(err)
	request := srvapi.RequestMsg{Type: srvapi.TourComplete, Payload: rawPayload}

	s.shieldMock.EXPECT().HandleTourComplete(s.clientIP, tourCompleteRequest).Return(tourCompleteResult)

	response, err := s.srv.handleRequest(s.clientIP, request)
	s.NoError(err)
	s.Equal(srvapi.WrongPuzzle, response.Type)
	s.Nil(response.Payload)
}

func (s *ServerSuite) TestGranted() {
	initialHash := testutils.HexHash("820888B1A040503A82AFA97EB0AE59E8214866C2D74F3DBC705A002FB17C86E9")
	latHash := testutils.HexHash("820888B1A040503A82AFA97EB0AE59E8214866C2D74F3DBC705A002FB17C86E5")
	quote := "Some quote"
	tourCompleteRequest := shield.TourCompleteRequest{
		InitialHash: initialHash,
		LastHash:    latHash,
	}
	tourCompleteResult := shield.TourCompleteResult{
		Granted: true,
		Quote:   quote,
	}
	payload := srvapi.TourCompletePayload{
		InitialHash: api.Hash(initialHash),
		LastHash:    api.Hash(latHash),
	}
	rawPayload, err := payload.MarshalMsg(nil)
	s.NoError(err)
	request := srvapi.RequestMsg{Type: srvapi.TourComplete, Payload: rawPayload}

	s.shieldMock.EXPECT().HandleTourComplete(s.clientIP, tourCompleteRequest).Return(tourCompleteResult)

	responsePayload := srvapi.ServiceGrantedPayload{Quote: quote}
	rawResponsePayload, err := responsePayload.MarshalMsg(nil)
	s.NoError(err)

	response, err := s.srv.handleRequest(s.clientIP, request)
	s.NoError(err)
	s.Equal(srvapi.Granted, response.Type)
	s.Equal(msgp.Raw(rawResponsePayload), response.Payload)
}

func (s *ServerSuite) TestUnsupported() {
	request := srvapi.RequestMsg{Type: srvapi.RequestType(10)}
	response, err := s.srv.handleRequest(s.clientIP, request)
	s.NoError(err)
	s.Equal(srvapi.Unsupported, response.Type)
	s.Nil(response.Payload)
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerSuite))
}
