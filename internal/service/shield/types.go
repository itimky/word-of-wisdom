package shield

import (
	"github.com/itimky/word-of-wisom/pkg/gtp"
)

type CheckResultType int

const (
	Ok            CheckResultType = 0
	Restricted    CheckResultType = 1
	WrongSolution CheckResultType = 2
)

type Puzzle struct {
	InitialHash gtp.Hash
	TourLength  int
}

type PuzzleSolution struct {
	InitialHash gtp.Hash
	LastHash    gtp.Hash
}

type PuzzleCheckResult struct {
	Type   CheckResultType
	Puzzle *Puzzle
}
