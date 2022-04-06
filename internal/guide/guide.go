package guide

import (
	"fmt"
	"net"

	guidecontracts "github.com/itimky/word-of-wisom/api/guide"
	"github.com/itimky/word-of-wisom/pkg/tcpserver"

	"github.com/sirupsen/logrus"
	"github.com/tinylib/msgp/msgp"
)

type Guide struct {
	secret   string
	hashCalc hashCalc
}

func NewGuide(secret string, hashCalc hashCalc) *Guide {
	return &Guide{
		secret:   secret,
		hashCalc: hashCalc,
	}
}

func (g *Guide) HandleRequest(conn net.Conn) (msgp.Encodable, error) {
	requestMsg := guidecontracts.RequestMsg{}
	if err := requestMsg.DecodeMsg(msgp.NewReader(conn)); err != nil {
		return nil, fmt.Errorf("decode message: %w", err)
	}

	request := newRequestFromMsg(requestMsg)

	logrus.Debugf("%v+", request)

	response := g.tourGuideHandler(request, tcpserver.GetClientIP(conn))

	responseMsg, err := response.Encodable()
	if err != nil {
		return nil, fmt.Errorf("encodable: %w", err)
	}

	return responseMsg, nil
}

func (g *Guide) tourGuideHandler(request Request, clientIP string) Response {
	hash := g.hashCalc.CalcGuideHash(request.PreviousHash, request.TourNumber, request.TourLength, clientIP, g.secret)
	response := Response{Hash: hash}

	return response
}
