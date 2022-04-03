package server

import "github.com/tinylib/msgp/msgp"

//go:generate msgp

type RequestType byte

const (
	InitialRequest      RequestType = 0
	TourCompleteRequest RequestType = 1
)

type RequestMsg struct {
	Type    byte
	Payload msgp.Raw // TourCompletePayload
}

type TourCompletePayload struct {
	InitialHash [20]byte
	LastHash    [20]byte
}

// --------------------------

type ResponseType byte

const (
	ServiceGranted     ResponseType = 0
	ServiceRestricted  ResponseType = 1
	UnsupportedRequest ResponseType = 2
)

type ResponseMsg struct {
	Type    byte
	Payload msgp.Raw // ServiceRestrictedPayload, ServiceGrantedPayload
}

type ServiceRestrictedPayload struct {
	InitialHash [20]byte
	TourLength  byte
}

type ServiceGrantedPayload struct {
	Quote string
}
