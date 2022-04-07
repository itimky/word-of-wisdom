package shield

import (
	"testing"

	"github.com/itimky/word-of-wisom/internal/service/shield/mocks"
	"github.com/itimky/word-of-wisom/pkg/testutils"
	"github.com/stretchr/testify/suite"
)

type ServerSuite struct {
	suite.Suite
	service          *Service
	hashCalcMock     *mocks.HashCalc
	quoteServiceMock *mocks.QuoteService
	clientIP         string
}

func (s *ServerSuite) SetupSuite() {
	// Init hashCalc & quoteService in SetupTest
	s.service = &Service{}
	s.clientIP = "127.0.0.1"
}

func (s *ServerSuite) SetupTest() {
	s.hashCalcMock = &mocks.HashCalc{}
	s.quoteServiceMock = &mocks.QuoteService{}
	s.service.hashCalc = s.hashCalcMock
	s.service.quoteService = s.quoteServiceMock
}

func (s *ServerSuite) TestServer_initialRequestHandler() {
	hash := testutils.HexHash("820888B1A040503A82AFA97EB0AE59E8214866C2D74F3DBC705A002FB17C86E9")

	s.hashCalcMock.EXPECT().CalcInitialHash(s.clientIP, s.service.cfg.TourLength, s.service.secret).Return(hash)

	result := s.service.HandleInitial(s.clientIP)
	s.Equal(s.service.cfg.TourLength, result.TourLength)
	s.Equal(hash, result.InitialHash)
}

func (s *ServerSuite) TestServer_tourCompleteRequestHandler() {
	hash := testutils.HexHash("820888B1A040503A82AFA97EB0AE59E8214866C2D74F3DBC705A002FB17C86E9")
	quote := "Some quote"

	s.hashCalcMock.EXPECT().VerifyHash(hash, hash, s.service.cfg.TourLength, s.clientIP,
		s.service.secret, s.service.cfg.GuideSecrets).Return(true)
	s.quoteServiceMock.EXPECT().Get().Return(quote)

	result := s.service.HandleTourComplete(TourCompleteRequest{InitialHash: hash, LastHash: hash}, s.clientIP)
	s.Equal(quote, result.Quote)
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerSuite))
}
