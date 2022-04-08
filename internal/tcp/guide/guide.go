package guide

import (
	"fmt"
	"time"

	"github.com/panjf2000/gnet/v2/pkg/pool/goroutine"

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
	pool *goroutine.Pool

	addr      string
	multicore bool
	timeout   time.Duration

	secret   string
	hashCalc hashCalc
}

func NewGuide(
	pool *goroutine.Pool,
	addr string,
	multicore bool,
	timeout time.Duration,
	secret string,
	hashCalc hashCalc,
) *Guide {
	return &Guide{
		pool:      pool,
		addr:      addr,
		multicore: multicore,
		timeout:   timeout,
		secret:    secret,
		hashCalc:  hashCalc,
	}
}

func (g *Guide) Run() error {
	err := gnet.Run(g, g.addr, gnet.WithOptions(gnet.Options{Multicore: g.multicore, TCPKeepAlive: g.timeout}))
	if err != nil {
		return fmt.Errorf("gnet run: %w", err)
	}
	return nil
}

func (g *Guide) OnBoot(eng gnet.Engine) gnet.Action {
	logrus.Infof("server (multi-core=%t timeout=%v) is listening on %s\n", g.multicore, g.timeout, g.addr)
	return gnet.None
}

func (g *Guide) OnTraffic(conn gnet.Conn) gnet.Action {
	err := g.pool.Submit(func() {
		if err := g.handleConnection(conn); err != nil {
			logrus.WithError(err).Error("handle connection")
		}
		if err := conn.Close(nil); err != nil {
			logrus.WithError(err).Error("close connection")
		}
	})

	if err != nil {
		logrus.WithError(err).Error("submit")
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

	err = conn.AsyncWrite(data, nil)
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
