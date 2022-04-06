package guide

import "github.com/itimky/word-of-wisom/api"

//go:generate msgp

type RequestMsg struct {
	PreviousHash api.Hash
	TourNumber   byte
	TourLength   byte
}

type ResponseMsg struct {
	PreviousHash api.Hash
	Hash         [32]byte
}
