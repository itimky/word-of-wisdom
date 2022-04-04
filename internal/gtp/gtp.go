package gtp

import (
	"crypto/sha1" //nolint:gosec
	"math/big"
	"strconv"
	"time"
)

type Hash [20]byte

type GTP struct {
	now func() time.Time
}

func NewGTP(now func() time.Time) *GTP {
	return &GTP{now: now}
}

func (gtp *GTP) timestamp() string {
	// TODO: move precision to config
	return gtp.now().UTC().Round(time.Minute).String()
}

func (gtp *GTP) CalcInitialHash(clientIP string, tourLength int, secret string) Hash {
	//nolint:gosec
	return sha1.Sum([]byte(clientIP + strconv.Itoa(tourLength) + gtp.timestamp() + secret))
}

func (gtp *GTP) CalcGuideHash(prevHash Hash, tourNumber, tourLength int, clientIP, secret string) Hash {
	//nolint:gosec
	return sha1.Sum([]byte(string(prevHash[:]) + strconv.Itoa(tourNumber) + strconv.Itoa(tourLength) + clientIP + gtp.timestamp() + secret))
}

func (gtp *GTP) VerifyHash(initialHash, lastHash Hash, tourLength int, clientIP, secret string, guideSecrets []string) bool {
	if initialHash != gtp.CalcInitialHash(clientIP, tourLength, secret) {
		return false
	}

	hash := initialHash
	for i := 1; i < tourLength+1; i++ {
		hash = gtp.CalcGuideHash(hash, i, tourLength, clientIP, guideSecrets[GuideIndex(hash, len(guideSecrets))])
	}

	return hash == lastHash
}

//nolint:wsl
func GuideIndex(hash Hash, guideCount int) int64 {
	hashNumber := big.NewInt(0)
	hashNumber.SetBytes(hash[:])
	count := big.NewInt(int64(guideCount))
	index := big.NewInt(0)
	return index.Mod(hashNumber, count).Int64()
}
