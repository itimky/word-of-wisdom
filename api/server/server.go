package server

import (
	"github.com/itimky/word-of-wisom/api"
	"github.com/tinylib/msgp/msgp"
)

//go:generate msgp

type RequestType byte

const (
	Initial      RequestType = 0
	TourComplete RequestType = 1
)

type RequestMsg struct {
	Type    RequestType `msg:"type"`
	Payload msgp.Raw    `msg:"payload"` // nil, TourCompletePayload
}

type TourCompletePayload struct {
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
	Payload msgp.Raw     `msg:"payload"` // nil, ServiceRestrictedPayload, ServiceGrantedPayload
}

type ServiceRestrictedPayload struct {
	InitialHash api.Hash `msg:"initial_hash,extension"`
	TourLength  byte     `msg:"tour_length"`
}

type ServiceGrantedPayload struct {
	Quote string `msg:"quote"`
}
