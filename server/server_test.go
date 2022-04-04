package server

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"net"
	"testing"
	srvcontracts "word-of-wisom/internal/contracts/server"
	"word-of-wisom/internal/gtp"
	"word-of-wisom/server/mocks"
)

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

func TestServer_initialRequestHandler(t *testing.T) {
	expectedHash := gtp.Hash{188, 50, 238, 132, 222, 117, 223, 120, 12, 44, 45, 67, 206, 160, 197, 63, 165, 211, 117, 233}
	tourLength := 5

	connMock := &mocks.Conn{}
	connMock.EXPECT().RemoteAddr().Return(&net.TCPAddr{IP: []byte{127, 0, 0, 1}})

	hashCalcMock := &mocks.HashCalc{}
	hashCalcMock.EXPECT().CalcInitialHash("127.0.0.1", tourLength, "").Return(expectedHash)

	server := NewServer(
		"host",
		"port",
		10,
		5,
		tourLength,
		[]string{"s1", "s2"},
		hashCalcMock,
		rand.New(rand.NewSource(1)),
	)

	result, err := server.initialRequestHandler(connMock)
	assert.NoError(t, err)
	assert.IsType(t, &srvcontracts.ResponseMsg{}, result)
	resultMsg := result.(*srvcontracts.ResponseMsg)
	assert.Equal(t, byte(srvcontracts.ServiceRestricted), resultMsg.Type)
	var resultPayload srvcontracts.ServiceRestrictedPayload
	_, err = resultPayload.UnmarshalMsg(resultMsg.Payload)
	assert.NoError(t, err)
	assert.Equal(t, byte(5), resultPayload.TourLength)
	assert.Equal(t, expectedHash, gtp.Hash(resultPayload.InitialHash))
}

func TestServer_tourCompleteRequestHandler(t *testing.T) {
	initialHash := gtp.Hash{188, 50, 238, 132, 222, 117, 223, 120, 12, 44, 45, 67, 206, 160, 197, 63, 165, 211, 117, 233}
	lastHash := gtp.Hash{189, 50, 238, 132, 222, 117, 223, 120, 12, 44, 45, 67, 206, 160, 197, 63, 165, 211, 117, 233}
	tourLength := 5

	tourCompletePayload := srvcontracts.TourCompletePayload{InitialHash: initialHash, LastHash: lastHash}
	payload, err := tourCompletePayload.MarshalMsg(nil)
	assert.NoError(t, err)

	requestMsg := srvcontracts.RequestMsg{Type: byte(srvcontracts.TourCompleteRequest), Payload: payload}

	connMock := &mocks.Conn{}
	connMock.EXPECT().RemoteAddr().Return(&net.TCPAddr{IP: []byte{127, 0, 0, 1}})

	hashCalcMock := &mocks.HashCalc{}
	hashCalcMock.EXPECT().VerifyHash(initialHash, lastHash, 5, "127.0.0.1", "", []string{"s1", "s2"}).Return(true)

	server := NewServer(
		"host",
		"port",
		10,
		5,
		tourLength,
		[]string{"s1", "s2"}, hashCalcMock,
		rand.New(rand.NewSource(1)),
	)
	result, err := server.tourCompleteRequestHandler(connMock, requestMsg)
	assert.NoError(t, err)
	assert.IsType(t, &srvcontracts.ResponseMsg{}, result)
	resultMsg := result.(*srvcontracts.ResponseMsg)
	assert.Equal(t, byte(srvcontracts.ServiceGranted), resultMsg.Type)
	var resultPayload srvcontracts.ServiceGrantedPayload
	_, err = resultPayload.UnmarshalMsg(resultMsg.Payload)
	assert.NoError(t, err)
	assert.Equal(t, "Carpe Diem", resultPayload.Quote)
}
