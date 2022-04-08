package client

import (
	"testing"

	srvapi "github.com/itimky/word-of-wisom/api/server"
	"github.com/itimky/word-of-wisom/internal/client/mocks"
	"github.com/stretchr/testify/assert"
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
