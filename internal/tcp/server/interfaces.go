package server

import "github.com/itimky/word-of-wisom/internal/service/shield"

//go:generate mockery --all --exported=true --with-expecter=true

type shieldService interface {
	HandleInitial(clientIP string) shield.InitialResult
	HandleTourComplete(clientIP string, request shield.TourCompleteRequest) shield.TourCompleteResult
}
