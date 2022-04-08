package server

import "github.com/itimky/word-of-wisom/pkg/gtp"

type CheckResultType int

const (
	Ok            CheckResultType = 0
	Restricted    CheckResultType = 1
	WrongSolution CheckResultType = 2
)

type PuzzleCheckResult struct {
	Type   CheckResultType
	Puzzle *gtp.Puzzle
}
