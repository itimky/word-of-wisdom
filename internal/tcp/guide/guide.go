package guide

import (
	"fmt"
	"log"

	"github.com/panjf2000/gnet/v2"
	"github.com/sirupsen/logrus"
	"github.com/tinylib/msgp/msgp"

	"github.com/itimky/word-of-wisom/api"
	guideapi "github.com/itimky/word-of-wisom/api/guide"
	"github.com/itimky/word-of-wisom/internal/tcp"
	"github.com/itimky/word-of-wisom/pkg/gtp"
)

type Guide struct {
	gnet.BuiltinEventEngine
	//pool *goroutine.Pool
	//eng  gnet.Engine

	addr      string
	multicore bool

	secret   string
	hashCalc hashCalc
}

func NewGuide(addr string, multicore bool, secret string, hashCalc hashCalc) *Guide {
	return &Guide{
		addr:      addr,
		multicore: multicore,
		secret:    secret,
		hashCalc:  hashCalc,
	}
}

func (g *Guide) Run() error {
	err := gnet.Run(g, g.addr, gnet.WithMulticore(g.multicore))
	if err != nil {
		return fmt.Errorf("gnet run: %w", err)
	}
	return nil
}

func (g *Guide) OnBoot(eng gnet.Engine) gnet.Action {
	log.Printf("server with multi-core=%t is listening on %s\n", g.multicore, g.addr)
	return gnet.None
}

func (g *Guide) OnTraffic(conn gnet.Conn) gnet.Action {
	if err := g.handleConnection(conn); err != nil {
		logrus.WithError(err).Error("handle connection")
		return gnet.Close
	}

	return gnet.None
}

func (g *Guide) handleConnection(conn gnet.Conn) error {
	requestMsg := guideapi.RequestMsg{}
	if err := requestMsg.DecodeMsg(msgp.NewReader(conn)); err != nil {
		return fmt.Errorf("decode message: %w", err)
	}

	responseMsg := g.handleRequest(tcp.GetClientIP(conn), requestMsg)

	data, err := responseMsg.MarshalMsg(nil)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	_, err = conn.Write(data)
	if err != nil {
		return fmt.Errorf("write data: %w", err)
	}

	return nil
}

func (g *Guide) handleRequest(clientAPI string, requestMsg guideapi.RequestMsg) guideapi.ResponseMsg {
	hash := g.hashCalc.CalcGuideHash(
		gtp.Hash(requestMsg.PreviousHash),
		int(requestMsg.TourNumber),
		int(requestMsg.TourLength),
		clientAPI,
		g.secret,
	)
	return guideapi.ResponseMsg{Hash: api.Hash(hash)}
}
