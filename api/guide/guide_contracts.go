package guide

//go:generate msgp

type RequestMsg struct {
	PreviousHash [32]byte
	TourNumber   byte
	TourLength   byte
}

type ResponseMsg struct {
	Hash [32]byte
}
