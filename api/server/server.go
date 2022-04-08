package server

import (
	"github.com/itimky/word-of-wisom/api"
	"github.com/tinylib/msgp/msgp"
)

//go:generate msgp

type RequestType byte

const (
	Quote RequestType = 0
)

type RequestMsg struct {
	Type           RequestType     `msg:"type"`
	PuzzleSolution *PuzzleSolution `msg:"puzzle_solution"` // nil, PuzzleSolution
}

type PuzzleSolution struct {
	InitialHash api.Hash `msg:"initial_hash,extension"`
	LastHash    api.Hash `msg:"last_hash,extension"`
}

// --------------------------

type ResponseType byte

const (
	Unsupported ResponseType = 0
	Restricted  ResponseType = 1
	WrongPuzzle ResponseType = 2
	Granted     ResponseType = 3
)

type ResponseMsg struct {
	Type    ResponseType `msg:"type"`
	Payload msgp.Raw     `msg:"payload"` // nil, PuzzleResponse, QuoteResponse
}

type PuzzleResponse struct {
	InitialHash api.Hash `msg:"initial_hash,extension"`
	TourLength  byte     `msg:"tour_length"`
}

type QuoteResponse struct {
	Quote string `msg:"quote"`
}
