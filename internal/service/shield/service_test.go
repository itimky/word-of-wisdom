package shield

import (
	"testing"

	"github.com/itimky/word-of-wisom/internal/service/shield/mocks"
	"github.com/itimky/word-of-wisom/pkg/testutils"
	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite
	service      *Service
	hashCalcMock *mocks.HashCalc
	clientIP     string
}

func (s *ServiceSuite) SetupSuite() {
	// Init hashCalc in SetupTest
	s.service = &Service{}
	s.clientIP = "127.0.0.1"
}

func (s *ServiceSuite) SetupTest() {
	s.hashCalcMock = &mocks.HashCalc{}
	s.service.hashCalc = s.hashCalcMock
}

func (s *ServiceSuite) TestService_Check__Restricted() {
	hash := testutils.HexHash("820888B1A040503A82AFA97EB0AE59E8214866C2D74F3DBC705A002FB17C86E9")

	s.hashCalcMock.EXPECT().CalcInitialHash(s.clientIP, s.service.cfg.TourLength, s.service.secret).Return(hash)

	result := s.service.CheckPuzzle(s.clientIP, nil)
	s.Equal(Restricted, result.Type)
	s.Equal(s.service.cfg.TourLength, result.Puzzle.TourLength)
	s.Equal(hash, result.Puzzle.InitialHash)
}

func (s *ServiceSuite) TestService_Check__Ok() {
	hash := testutils.HexHash("820888B1A040503A82AFA97EB0AE59E8214866C2D74F3DBC705A002FB17C86E9")

	s.hashCalcMock.EXPECT().VerifyHash(hash, hash, s.service.cfg.TourLength, s.clientIP,
		s.service.secret, s.service.cfg.GuideSecrets).Return(true)

	result := s.service.CheckPuzzle(s.clientIP, &PuzzleSolution{InitialHash: hash, LastHash: hash})
	s.Equal(Ok, result.Type)
	s.Nil(result.Puzzle)
}

func (s *ServiceSuite) TestService_Check__WrongSolution() {
	hash := testutils.HexHash("820888B1A040503A82AFA97EB0AE59E8214866C2D74F3DBC705A002FB17C86E9")

	s.hashCalcMock.EXPECT().VerifyHash(hash, hash, s.service.cfg.TourLength, s.clientIP,
		s.service.secret, s.service.cfg.GuideSecrets).Return(false)

	result := s.service.CheckPuzzle(s.clientIP, &PuzzleSolution{InitialHash: hash, LastHash: hash})
	s.Equal(WrongSolution, result.Type)
	s.Nil(result.Puzzle)
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
