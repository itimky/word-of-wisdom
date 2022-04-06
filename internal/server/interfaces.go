package server

import (
	"word-of-wisom/pkg/gtp"
)

//go:generate mockery --name hashCalc --exported=true --with-expecter=true

type hashCalc interface {
	CalcInitialHash(clientIP string, tourLength int, secret string) gtp.Hash
	VerifyHash(initialHash, lastHash gtp.Hash, tourLength int, clientIP string, secret string, guideSecrets []string) bool
}
