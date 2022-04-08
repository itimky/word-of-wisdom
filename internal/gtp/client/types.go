package client

import (
	"github.com/itimky/word-of-wisom/pkg/gtp"
	"github.com/tinylib/msgp/msgp"
)

type serverResponse struct {
	Granted bool
	Puzzle  *gtp.Puzzle
	Payload msgp.Raw
}
