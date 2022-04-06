package server

import (
	"math/rand"
	"testing"

	"github.com/itimky/word-of-wisom/internal/server/mocks"
	"github.com/itimky/word-of-wisom/pkg/quotes"
	"github.com/itimky/word-of-wisom/pkg/testutils"
	utilmocks "github.com/itimky/word-of-wisom/pkg/testutils/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ServerSuite struct {
	suite.Suite
	server       *Server
	connMock     *utilmocks.Conn
	hashCalcMock *mocks.HashCalc
	clientIP     string
}

func (s *ServerSuite) SetupSuite() {
	secretLength := 10
	tourLength := 5
	guideSecrets := []string{"s1", "s2"}

	// init in SetupTest
	var hashCalcMock *mocks.HashCalc

	s.server = NewServer(
		secretLength,
		tourLength,
		guideSecrets,
		hashCalcMock,
		nil,
	)
	s.clientIP = "127:0:0:1"
}

func (s *ServerSuite) SetupTest() {
	s.connMock = &utilmocks.Conn{}
	s.hashCalcMock = &mocks.HashCalc{}
	s.server.hashCalc = s.hashCalcMock
	s.server.quoteGetter = quotes.NewQuoteRandomizer(rand.New(rand.NewSource(1)))
}

func (s *ServerSuite) TestServer_initialRequestHandler() {
	hash := testutils.HexHash("820888B1A040503A82AFA97EB0AE59E8214866C2D74F3DBC705A002FB17C86E9")

	s.hashCalcMock.EXPECT().CalcInitialHash(s.clientIP, s.server.tourLength, s.server.secret).Return(hash)

	result := s.server.initialRequestHandler(s.clientIP)
	s.Equal(5, result.TourLength)
	s.Equal(hash, result.InitialHash)
}

func (s *ServerSuite) TestServer_tourCompleteRequestHandler() {
	hash := testutils.HexHash("820888B1A040503A82AFA97EB0AE59E8214866C2D74F3DBC705A002FB17C86E9")

	s.hashCalcMock.EXPECT().VerifyHash(hash, hash, s.server.tourLength, s.clientIP,
		s.server.secret, s.server.guideSecrets).Return(true)

	result := s.server.tourCompleteRequestHandler(TourCompleteRequest{InitialHash: hash, LastHash: hash}, s.clientIP)
	s.IsType(ServiceGrantedResponse{}, result)

	resultMsg := result.(ServiceGrantedResponse)
	s.Equal("Carpe Diem", resultMsg.Quote)
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerSuite))
}

func Test_randomSecret_Length(t *testing.T) {
	args := []struct {
		length int
	}{
		{length: 5},
		{length: 10},
		{length: 15},
		{length: 20},
	}
	for _, tt := range args {
		t.Run("test-random-secret-length", func(t *testing.T) {
			secret, err := randomSecret(tt.length)
			assert.NoError(t, err)
			assert.Equal(t, tt.length, len(secret))
		})
	}
}

func Test_randomSecret_Unique(t *testing.T) {
	secret1, err := randomSecret(10)
	assert.NoError(t, err)
	secret2, err := randomSecret(10)
	assert.NoError(t, err)
	assert.NotEqual(t, secret1, secret2)
}
