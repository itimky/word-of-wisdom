package server

import (
	"word-of-wisom/internal/gtp"
)

//go:generate mockery --name HashCalc --with-expecter=true

type HashCalc interface {
	CalcInitialHash(clientIP string, tourLength int, secret string) gtp.Hash
	VerifyHash(initialHash, lastHash gtp.Hash, tourLength int, clientIP string, secret string, guideSecrets []string) bool
}
