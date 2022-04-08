package client

import (
	"fmt"

	srvapi "github.com/itimky/word-of-wisom/api/server"
	"github.com/tinylib/msgp/msgp"
)

func convertPayloadToQuoteResponse(payload msgp.Raw) (*QuoteResponse, error) {
	response := srvapi.QuoteResponse{}
	_, err := response.UnmarshalMsg(payload)
	if err != nil {
		return nil, fmt.Errorf("unmarshal msg: %w", err)
	}

	return &QuoteResponse{
		Quote: response.Quote,
	}, nil
}
