package server

import (
	"net"
	"word-of-wisom/internal/gtp"
)

//go:generate mockery --name Conn --name HashCalc --with-expecter=true

type Conn interface {
	net.Conn
}

type HashCalc interface {
	CalcInitialHash(clientIP string, tourLength int, secret string) gtp.Hash
	CalcGuideHash(prevHash gtp.Hash, tourNumber int, tourLength int, clientIP string, secret string) gtp.Hash
	VerifyHash(initialHash, lastHash gtp.Hash, tourLength int, clientIP string, secret string, guideSecrets []string) bool
}
