package server

import (
	"math/rand"
	"net"
	"testing"
	"time"
	srvcontracts "word-of-wisom/internal/contracts/server"
	"word-of-wisom/internal/gtp"
	utilmocks "word-of-wisom/internal/testutils/mocks"
	"word-of-wisom/server/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ServerSuite struct {
	suite.Suite
	server       *Server
	connMock     *utilmocks.Conn
	hashCalcMock *mocks.HashCalc
	clientIP     *net.TCPAddr
}

func (s *ServerSuite) SetupSuite() {
	host := "host"
	port := "port"
	secretLength := 10
	secretUpdateIntervalMinutes := time.Duration(2)
	tourLength := 5
	guideSecrets := []string{"s1", "s2"}

	// init in SetupTest
	var rnd *rand.Rand
	var hashCalcMock *mocks.HashCalc

	s.server = NewServer(
		host,
		port,
		secretLength,
		secretUpdateIntervalMinutes,
		tourLength,
		guideSecrets,
		hashCalcMock,
		rnd,
	)
	s.clientIP = &net.TCPAddr{IP: []byte{127, 0, 0, 1}}
}

func (s *ServerSuite) SetupTest() {
	s.connMock = &utilmocks.Conn{}
	s.hashCalcMock = &mocks.HashCalc{}
	s.server.hashCalc = s.hashCalcMock
	s.server.rand = rand.New(rand.NewSource(1))
}

func (s *ServerSuite) TestServer_initialRequestHandler() {
	expectedHash := gtp.Hash{188, 50, 238, 132, 222, 117, 223, 120, 12, 44, 45, 67, 206, 160, 197, 63, 165, 211, 117, 233}

	s.connMock.EXPECT().RemoteAddr().Return(s.clientIP)
	s.hashCalcMock.EXPECT().CalcInitialHash(s.clientIP.IP.String(), s.server.tourLength, s.server.secret).Return(expectedHash)

	result, err := s.server.initialRequestHandler(s.connMock)
	s.NoError(err)
	s.IsType(&srvcontracts.ResponseMsg{}, result)
	resultMsg := result.(*srvcontracts.ResponseMsg)
	s.Equal(byte(srvcontracts.ServiceRestricted), resultMsg.Type)
	var resultPayload srvcontracts.ServiceRestrictedPayload
	_, err = resultPayload.UnmarshalMsg(resultMsg.Payload)
	s.NoError(err)
	s.Equal(byte(5), resultPayload.TourLength)
	s.Equal(expectedHash, gtp.Hash(resultPayload.InitialHash))
}

func (s *ServerSuite) TestServer_tourCompleteRequestHandler() {
	initialHash := gtp.Hash{188, 50, 238, 132, 222, 117, 223, 120, 12, 44, 45, 67, 206, 160, 197, 63, 165, 211, 117, 233}
	lastHash := gtp.Hash{189, 50, 238, 132, 222, 117, 223, 120, 12, 44, 45, 67, 206, 160, 197, 63, 165, 211, 117, 233}

	tourCompletePayload := srvcontracts.TourCompletePayload{InitialHash: initialHash, LastHash: lastHash}
	payload, err := tourCompletePayload.MarshalMsg(nil)
	s.NoError(err)

	requestMsg := srvcontracts.RequestMsg{Type: byte(srvcontracts.TourCompleteRequest), Payload: payload}

	s.connMock.EXPECT().RemoteAddr().Return(s.clientIP)
	s.hashCalcMock.EXPECT().VerifyHash(initialHash, lastHash, s.server.tourLength, s.clientIP.IP.String(),
		s.server.secret, s.server.guideSecrets).Return(true)

	result, err := s.server.tourCompleteRequestHandler(s.connMock, requestMsg)
	s.NoError(err)
	s.IsType(&srvcontracts.ResponseMsg{}, result)
	resultMsg := result.(*srvcontracts.ResponseMsg)
	s.Equal(byte(srvcontracts.ServiceGranted), resultMsg.Type)
	var resultPayload srvcontracts.ServiceGrantedPayload
	_, err = resultPayload.UnmarshalMsg(resultMsg.Payload)
	s.NoError(err)
	s.Equal("Carpe Diem", resultPayload.Quote)
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
