package guide

import (
	guidecontracts "word-of-wisom/api/guide"
	"word-of-wisom/pkg/gtp"

	"github.com/tinylib/msgp/msgp"
)

type Request struct {
	PreviousHash gtp.Hash
	TourNumber   int
	TourLength   int
}

func newRequestFromMsg(msg guidecontracts.RequestMsg) Request {
	return Request{
		PreviousHash: msg.PreviousHash,
		TourNumber:   int(msg.TourNumber),
		TourLength:   int(msg.TourLength),
	}
}

type Response struct {
	Hash gtp.Hash
}

func (r Response) Encodable() (msgp.Encodable, error) {
	responseMsg := guidecontracts.ResponseMsg{Hash: r.Hash}
	return &responseMsg, nil
}