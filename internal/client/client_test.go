package client

import (
	"testing"

	"github.com/stretchr/testify/assert"

	srvapi "github.com/itimky/word-of-wisom/api/server"
	"github.com/itimky/word-of-wisom/internal/client/mocks"
)

func TestClient_RequestQuote(t *testing.T) {
	gtpClientMock := &mocks.GtpClient{}
	c := NewClient(gtpClientMock)
	quote := "some quote"

	response := srvapi.QuoteResponse{
		Quote: quote,
	}
	rawResponse, err := response.MarshalMsg(nil)
	assert.NoError(t, err)

	gtpClientMock.EXPECT().MakeRequest(srvapi.Quote).Return(rawResponse, nil)

	quoteResponse, err := c.RequestQuote()
	assert.NoError(t, err)
	assert.Equal(t, quote, quoteResponse.Quote)
}
