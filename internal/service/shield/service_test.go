package shield

import (
	"testing"

	"github.com/itimky/word-of-wisom/internal/service/shield/mocks"
	"github.com/itimky/word-of-wisom/pkg/testutils"
	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite
	service          *Service
	hashCalcMock     *mocks.HashCalc
	quoteServiceMock *mocks.QuoteService
	clientIP         string
}

func (s *ServiceSuite) SetupSuite() {
	// Init hashCalc & quoteService in SetupTest
	s.service = &Service{}
	s.clientIP = "127.0.0.1"
}

func (s *ServiceSuite) SetupTest() {
	s.hashCalcMock = &mocks.HashCalc{}
	s.quoteServiceMock = &mocks.QuoteService{}
	s.service.hashCalc = s.hashCalcMock
	s.service.quoteService = s.quoteServiceMock
}

func (s *ServiceSuite) TestServer_initialRequestHandler() {
	hash := testutils.HexHash("820888B1A040503A82AFA97EB0AE59E8214866C2D74F3DBC705A002FB17C86E9")

	s.hashCalcMock.EXPECT().CalcInitialHash(s.clientIP, s.service.cfg.TourLength, s.service.secret).Return(hash)

	result := s.service.HandleInitial(s.clientIP)
	s.Equal(s.service.cfg.TourLength, result.TourLength)
	s.Equal(hash, result.InitialHash)
}

func (s *ServiceSuite) TestServer_tourCompleteRequestHandler_Granted() {
	hash := testutils.HexHash("820888B1A040503A82AFA97EB0AE59E8214866C2D74F3DBC705A002FB17C86E9")
	quote := "Some quote"

	s.hashCalcMock.EXPECT().VerifyHash(hash, hash, s.service.cfg.TourLength, s.clientIP,
		s.service.secret, s.service.cfg.GuideSecrets).Return(true)
	s.quoteServiceMock.EXPECT().Get().Return(quote)

	result := s.service.HandleTourComplete(s.clientIP, TourCompleteRequest{InitialHash: hash, LastHash: hash})
	s.True(result.Granted)
	s.Equal(quote, result.Quote)
}

func (s *ServiceSuite) TestServer_tourCompleteRequestHandler_NotGranted() {
	hash := testutils.HexHash("820888B1A040503A82AFA97EB0AE59E8214866C2D74F3DBC705A002FB17C86E9")
	quote := "Some quote"

	s.hashCalcMock.EXPECT().VerifyHash(hash, hash, s.service.cfg.TourLength, s.clientIP,
		s.service.secret, s.service.cfg.GuideSecrets).Return(false)
	s.quoteServiceMock.EXPECT().Get().Return(quote)

	result := s.service.HandleTourComplete(s.clientIP, TourCompleteRequest{InitialHash: hash, LastHash: hash})
	s.False(result.Granted)
	s.Empty(result.Quote)
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
