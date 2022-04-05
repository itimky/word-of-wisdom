package guide

import "word-of-wisom/pkg/gtp"

//go:generate mockery --name HashCalc --with-expecter=true

type HashCalc interface {
	CalcGuideHash(prevHash gtp.Hash, tourNumber int, tourLength int, clientIP string, secret string) gtp.Hash
}
