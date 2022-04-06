package guide

import (
	"net"
	"testing"
	"word-of-wisom/internal/guide/mocks"
	"word-of-wisom/pkg/testutils"

	"github.com/stretchr/testify/assert"
)

func TestGuide_tourGuideHandler(t *testing.T) {
	clientIP := &net.TCPAddr{IP: []byte{127, 0, 0, 1}}
	secret := "secret1"
	tourNumber := 2
	tourLength := 5
	hash := testutils.HexHash("820888B1A040503A82AFA97EB0AE59E8214866C2D74F3DBC705A002FB17C86E9")

	hashCalcMock := &mocks.HashCalc{}
	hashCalcMock.EXPECT().CalcGuideHash(hash, tourNumber, tourLength, clientIP.IP.String(), secret).Return(hash)

	guide := NewGuide(secret, hashCalcMock)
	request := Request{PreviousHash: hash, TourNumber: tourNumber, TourLength: tourLength}
	response := guide.tourGuideHandler(request, clientIP.IP.String())
	assert.Equal(t, hash, response.Hash)
}
