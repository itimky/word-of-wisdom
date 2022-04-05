package guide

import (
	"net"
	guidecontracts "word-of-wisom/api/guide"

	"github.com/tinylib/msgp/msgp"
)

func (g *Guide) tourGuideHandler(conn net.Conn, request guidecontracts.RequestMsg) msgp.Encodable {
	hash := g.hashCalc.CalcGuideHash(request.PreviousHash, int(request.TourNumber), int(request.TourLength), g.getClientIP(conn), g.secret)
	response := guidecontracts.ResponseMsg{Hash: hash}

	return &response
}
