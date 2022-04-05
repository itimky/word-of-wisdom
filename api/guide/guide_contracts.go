package guide

//go:generate msgp

type RequestMsg struct {
	PreviousHash [20]byte
	TourNumber   byte
	TourLength   byte
}

type ResponseMsg struct {
	Hash [20]byte
}
