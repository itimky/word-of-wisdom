package client

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/tinylib/msgp/msgp"

	guideapi "github.com/itimky/word-of-wisom/api/guide"
	srvapi "github.com/itimky/word-of-wisom/api/server"
)

type Client struct {
	server string
	guides []string
}

func NewClient(server string, guides []string) *Client {
	return &Client{
		server: server,
		guides: guides,
	}
}

func (c *Client) RequestServer(request srvapi.RequestMsg) (*srvapi.ResponseMsg, error) {
	conn, err := net.Dial("tcp", c.server)
	if err != nil {
		return nil, fmt.Errorf("connect to guide: %w", err)
	}

	writer := msgp.NewWriter(conn)

	err = request.EncodeMsg(writer)
	if err != nil {
		return nil, fmt.Errorf("encode msg: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return nil, fmt.Errorf("flush writer: %w", err)
	}

	logrus.Debugf("Server request sent: {Type: %v, PuzzleSolution: %+v}", request.Type, request.PuzzleSolution)

	response := srvapi.ResponseMsg{}
	if err := response.DecodeMsg(msgp.NewReader(conn)); err != nil {
		return nil, fmt.Errorf("decode server response: %w", err)
	}

	return &response, nil
}

func (c *Client) RequestGuideByIdx(idx int, request guideapi.RequestMsg) (*guideapi.ResponseMsg, error) {
	guideAddr := c.guides[idx]

	conn, err := net.Dial("tcp", guideAddr)
	if err != nil {
		return nil, fmt.Errorf("connect to guide: %w", err)
	}

	writer := msgp.NewWriter(conn)

	err = request.EncodeMsg(writer)
	if err != nil {
		return nil, fmt.Errorf("encode msg: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return nil, fmt.Errorf("flush writer: %w", err)
	}

	logrus.Debugf("Guide request sent: %+v", request)

	response := guideapi.ResponseMsg{}
	if err := response.DecodeMsg(msgp.NewReader(conn)); err != nil {
		return nil, fmt.Errorf("decode guide response: %w", err)
	}

	return &response, nil
}
