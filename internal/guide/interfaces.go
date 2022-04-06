package guide

import "word-of-wisom/pkg/gtp"

//go:generate mockery --name hashCalc --exported=true --with-expecter=true

type hashCalc interface {
	CalcGuideHash(prevHash gtp.Hash, tourNumber int, tourLength int, clientIP string, secret string) gtp.Hash
}
