package server

import (
	"github.com/itimky/word-of-wisom/internal/gtp/server"
	"github.com/itimky/word-of-wisom/pkg/gtp"
)

//go:generate mockery --all --exported=true --with-expecter=true

type gtpServer interface {
	CheckPuzzle(clientIP string, solution *gtp.PuzzleSolution) server.PuzzleCheckResult
}

type quoteRepository interface {
	Get() string
}
