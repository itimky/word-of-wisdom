package client

import (
	guideapi "github.com/itimky/word-of-wisom/api/guide"
	srvapi "github.com/itimky/word-of-wisom/api/server"
)

//go:generate mockery --all --exported=true --with-expecter=true

type tcpClient interface {
	RequestServer(request srvapi.RequestMsg) (*srvapi.ResponseMsg, error)
	RequestGuideByIdx(idx int, request guideapi.RequestMsg) (*guideapi.ResponseMsg, error)
}
