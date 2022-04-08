package gtp

import (
	"crypto/sha256"
	"fmt"
)

const HashSize = sha256.Size

type Hash [HashSize]byte

func (h Hash) String() string {
	return fmt.Sprintf("%X", h[:])
}

type Puzzle struct {
	InitialHash Hash
	TourLength  int
}

type PuzzleSolution struct {
	InitialHash Hash
	LastHash    Hash
}
