package server

import "github.com/itimky/word-of-wisom/internal/service/shield"

//go:generate mockery --all --exported=true --with-expecter=true

type shieldService interface {
	CheckPuzzle(clientIP string, solution *shield.PuzzleSolution) shield.PuzzleCheckResult
}

type quoteService interface {
	Get() string
}
