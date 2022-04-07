package guide

import (
	"testing"

	"github.com/itimky/word-of-wisom/api"
	guideapi "github.com/itimky/word-of-wisom/api/guide"
	"github.com/itimky/word-of-wisom/internal/tcp/guide/mocks"

	"github.com/itimky/word-of-wisom/pkg/testutils"

	"github.com/stretchr/testify/assert"
)

func TestGuide_handleRequest(t *testing.T) {
	clientIP := "127.0.0.1"
	secret := "secret1"
	tourNumber := 2
	tourLength := 5
	hash := testutils.HexHash("820888B1A040503A82AFA97EB0AE59E8214866C2D74F3DBC705A002FB17C86E9")

	hashCalcMock := &mocks.HashCalc{}
	hashCalcMock.EXPECT().CalcGuideHash(hash, tourNumber, tourLength, clientIP, secret).Return(hash)

	guide := &Guide{secret: secret, hashCalc: hashCalcMock}
	request := guideapi.RequestMsg{
		PreviousHash: api.Hash(hash),
		TourNumber:   byte(tourNumber),
		TourLength:   byte(tourLength),
	}
	response := guide.handleRequest(clientIP, request)
	assert.Equal(t, api.Hash(hash), response.Hash)
}
