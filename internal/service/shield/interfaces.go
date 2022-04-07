package shield

import (
	"github.com/itimky/word-of-wisom/pkg/gtp"
)

//go:generate mockery --all --exported=true --with-expecter=true

type hashCalc interface {
	CalcInitialHash(clientIP string, tourLength int, secret string) gtp.Hash
	VerifyHash(initialHash, lastHash gtp.Hash, tourLength int, clientIP string, secret string, guideSecrets []string) bool
}

type quoteService interface {
	Get() string
}
