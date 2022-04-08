package client

import (
	"github.com/tinylib/msgp/msgp"

	srvapi "github.com/itimky/word-of-wisom/api/server"
)

//go:generate mockery --all --exported=true --with-expecter=true

type gtpClient interface {
	MakeRequest(reqType srvapi.RequestType) (msgp.Raw, error)
}
