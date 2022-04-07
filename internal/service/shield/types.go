package shield

import (
	"github.com/itimky/word-of-wisom/pkg/gtp"
)

type InitialResult struct {
	InitialHash gtp.Hash
	TourLength  int
}

type TourCompleteRequest struct {
	InitialHash gtp.Hash
	LastHash    gtp.Hash
}

type TourCompleteResult struct {
	Granted bool
	Quote   string
}
