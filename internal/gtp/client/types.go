package client

import (
	"github.com/tinylib/msgp/msgp"

	"github.com/itimky/word-of-wisom/pkg/gtp"
)

type serverResponse struct {
	Granted bool
	Puzzle  *gtp.Puzzle
	Payload msgp.Raw
}
