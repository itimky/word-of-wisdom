package client

import (
	srvapi "github.com/itimky/word-of-wisom/api/server"
	"github.com/tinylib/msgp/msgp"
)

//go:generate mockery --all --exported=true --with-expecter=true

type gtpClient interface {
	MakeRequest(reqType srvapi.RequestType) (msgp.Raw, error)
}
