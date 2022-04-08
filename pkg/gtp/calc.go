package gtp

import (
	"crypto/sha256"
	"math/big"
	"strconv"
	"time"
)

type Calc struct {
	now func() time.Time
}

func NewCalc(now func() time.Time) *Calc {
	return &Calc{now: now}
}

func (c *Calc) timestamp() string {
	// TODO: move precision to config
	return c.now().UTC().Round(time.Minute).String()
}

func (c *Calc) CalcInitialHash(clientIP string, tourLength int, secret string) Hash {
	return sha256.Sum256([]byte(clientIP + strconv.Itoa(tourLength) + c.timestamp() + secret))
}

func (c *Calc) CalcGuideHash(prevHash Hash, tourNumber, tourLength int, clientIP, secret string) Hash {
	return sha256.Sum256([]byte(string(prevHash[:]) + strconv.Itoa(tourNumber) + strconv.Itoa(tourLength) + clientIP + c.timestamp() + secret))
}

func (c *Calc) VerifyHash(initialHash, lastHash Hash, tourLength int, clientIP, secret string, guideSecrets []string) bool {
	if initialHash != c.CalcInitialHash(clientIP, tourLength, secret) {
		return false
	}

	hash := initialHash
	for i := 1; i < tourLength+1; i++ {
		hash = c.CalcGuideHash(hash, i, tourLength, clientIP, guideSecrets[GuideIndex(hash, len(guideSecrets))])
	}

	return hash == lastHash
}

//nolint:wsl
func GuideIndex(hash Hash, guideCount int) int {
	hashNumber := big.NewInt(0)
	hashNumber.SetBytes(hash[:])
	count := big.NewInt(int64(guideCount))
	index := big.NewInt(0)
	return int(index.Mod(hashNumber, count).Int64())
}
