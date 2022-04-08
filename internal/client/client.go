package client

import (
	"fmt"

	srvapi "github.com/itimky/word-of-wisom/api/server"
)

type Client struct {
	gtpClient gtpClient
}

func NewClient(client gtpClient) *Client {
	return &Client{
		gtpClient: client,
	}
}

func (s *Client) RequestQuote() (*QuoteResponse, error) {
	response, err := s.gtpClient.MakeRequest(srvapi.Quote)
	if err != nil {
		return nil, fmt.Errorf("make request: %w", err)
	}

	quoteResponse, err := convertPayloadToQuoteResponse(response)
	if err != nil {
		return nil, fmt.Errorf("convert payload to quote response")
	}

	return quoteResponse, nil
}
