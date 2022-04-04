package guide

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"word-of-wisom/guide/mocks"
	guidecontracts "word-of-wisom/internal/contracts/guide"
	"word-of-wisom/internal/gtp"
	utilmocks "word-of-wisom/internal/testutils/mocks"
)

func TestGuide_handleRequest(t *testing.T) {
	host := "host"
	port := "port"
	clientIP := &net.TCPAddr{IP: []byte{127, 0, 0, 1}}
	secret := "secret1"
	tourNumber := 2
	tourLength := 5
	prevHash := [20]byte{188, 50, 238, 132, 222, 117, 223, 120, 12, 44, 45, 67, 206, 160, 197, 63, 165, 211, 117, 233}
	newHash := [20]byte{180, 50, 238, 132, 222, 117, 223, 120, 12, 44, 45, 67, 206, 160, 197, 63, 165, 211, 117, 233}

	request := guidecontracts.RequestMsg{
		PreviousHash: prevHash,
		TourNumber:   byte(tourNumber),
		TourLength:   byte(tourLength),
	}

	connMock := &utilmocks.Conn{}
	connMock.EXPECT().RemoteAddr().Return(clientIP)

	hashCalcMock := &mocks.HashCalc{}
	hashCalcMock.EXPECT().CalcGuideHash(gtp.Hash(request.PreviousHash), int(request.TourNumber),
		int(request.TourLength), clientIP.IP.String(), secret).Return(newHash)

	guide := NewGuide(host, port, secret, hashCalcMock)
	response := guide.tourGuideHandler(connMock, request)
	assert.IsType(t, &guidecontracts.ResponseMsg{}, response)
	responseMsg := response.(*guidecontracts.ResponseMsg)
	assert.Equal(t, newHash, responseMsg.Hash)
}
