package guide

import "github.com/itimky/word-of-wisom/api"

//go:generate msgp

type RequestMsg struct {
	PreviousHash api.Hash `msg:"previous_hash,extension"`
	TourNumber   byte     `msg:"tour_number"`
	TourLength   byte     `msg:"tour_length"`
}

type ResponseMsg struct {
	Hash api.Hash `msg:"hash,extension"`
}
