package gtp

import (
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func now() time.Time {
	return time.Date(2022, 4, 3, 13, 3, 29, 0, time.UTC)
}

func hexHash(hexStr string) Hash {
	if len(hexStr) != 2*HashSize {
		panic(fmt.Sprintf("wrong hash size (%v)", len(hexStr)))
	}

	hashSlice, err := hex.DecodeString(hexStr)
	if err != nil {
		panic(fmt.Sprintf("hex decode string %s", hexStr))
	}
	var hash Hash
	copy(hash[:], hashSlice)
	return hash
}

type GTPSuite struct {
	suite.Suite
	gtp *Calc
}

func (s *GTPSuite) SetupSuite() {
	s.gtp = NewCalc(now)
}

func (s *GTPSuite) TestGTP_CalcInitialHash() {
	clientIP := "127.0.0.1"
	tourLength := 5
	secret := "test-secret-key"

	hash := s.gtp.CalcInitialHash(clientIP, tourLength, secret)

	s.Equal(hexHash("820888B1A040503A82AFA97EB0AE59E8214866C2D74F3DBC705A002FB17C86E9"), hash)
}

func (s *GTPSuite) TestGTP_CalcGuideHash() {
	clientIP := "127.0.0.1"
	prevHash := hexHash("820888B1A040503A82AFA97EB0AE59E8214866C2D74F3DBC705A002FB17C86E9")
	tourNumber := 1
	tourLength := 5
	secret := "test-secret-key-2"

	hash := s.gtp.CalcGuideHash(prevHash, tourNumber, tourLength, clientIP, secret)
	logrus.Info(hash)
	s.Equal(hexHash("801B39EFB47CAEE5EA342A0F6FC885E1A2A32C15BEF7DA0B58DD12A35A031CC7"), hash)
}

func (s *GTPSuite) TestGTP_VerifyHash() {
	clientIP := "127.0.0.1"
	tourLength := 5
	secret := "test-secret-key"
	guideSecrets := []string{"test-guide-secret-key-1", "test-guide-secret-key-2"}

	args := []struct {
		initialHash Hash
		lastHash    Hash
		result      bool
	}{
		{
			initialHash: hexHash("820888B1A040503A82AFA97EB0AE59E8214866C2D74F3DBC705A002FB17C86E9"),
			lastHash:    hexHash("4DDD8EB388374180706E41B24A19AA29B5E58A7281E6CBEEA4C8BD223D3A4B67"),
			result:      true,
		},
		{
			initialHash: hexHash("820888B1A040503A82AFA97EB0AE59E8214866C2D74F3DBC705A002FB17C86E9"),
			lastHash:    hexHash("5DDD8EB388374180706E41B24A19AA29B5E58A7281E6CBEEA4C8BD223D3A4B67"),
			result:      false,
		},
	}

	for _, tt := range args {
		s.Run("verify-hash", func() {
			match := s.gtp.VerifyHash(tt.initialHash, tt.lastHash, tourLength, clientIP, secret, guideSecrets)
			s.Equal(tt.result, match)
		})
	}
}

func TestGTP_Timestamp(t *testing.T) {
	args := []struct {
		name   string
		now    func() time.Time
		result string
	}{
		{
			name: "round down",
			now: func() time.Time {
				return time.Date(2022, 4, 3, 13, 3, 29, 0, time.UTC)
			},
			result: time.Date(2022, 4, 3, 13, 3, 0, 0, time.UTC).String(),
		},
		{
			name: "round up",
			now: func() time.Time {
				return time.Date(2022, 4, 3, 13, 3, 30, 0, time.UTC)
			},
			result: time.Date(2022, 4, 3, 13, 4, 0, 0, time.UTC).String(),
		},
	}
	for _, tt := range args {
		t.Run(tt.name, func(t *testing.T) {
			gtp := NewCalc(tt.now)
			assert.Equal(t, tt.result, gtp.timestamp())
		})
	}
}

func TestGTPSuite(t *testing.T) {
	suite.Run(t, new(GTPSuite))
}

func TestGuideIndex(t *testing.T) {
	args := []struct {
		hash       Hash
		guideCount int
		result     int
	}{
		{
			hash:       hexHash("0000000000000000000000000000000000000000000000000000000000000000"),
			guideCount: 3,
			result:     0,
		},
		{
			hash:       hexHash("0000000000000000000000000000000000000000000000000000000000000001"),
			guideCount: 3,
			result:     1,
		},
		{
			hash:       hexHash("0000000000000000000000000000000000000000000000000000000000000002"),
			guideCount: 3,
			result:     2,
		},
		{
			hash:       hexHash("0000000000000000000000000000000000000000000000000000000000000003"),
			guideCount: 3,
			result:     0,
		},
	}
	for _, tt := range args {
		t.Run("test-guide-index", func(t *testing.T) {
			assert.Equal(t, tt.result, GuideIndex(tt.hash, tt.guideCount))
		})
	}
}

func BenchmarkGTP_CalcInitialHash(b *testing.B) {
	clientIP := "127.0.0.1"
	tourLength := 5
	secret := "test-secret-key"
	gtp := NewCalc(time.Now)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		gtp.CalcInitialHash(clientIP, tourLength, secret)
	}
}

func BenchmarkGTP_CalcGuideHash(b *testing.B) {
	clientIP := "127.0.0.1"
	prevHash := hexHash("820888B1A040503A82AFA97EB0AE59E8214866C2D74F3DBC705A002FB17C86E9")
	tourNumber := 1
	tourLength := 5
	secret := "test-secret-key-2"
	gtp := NewCalc(time.Now)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		gtp.CalcGuideHash(prevHash, tourNumber, tourLength, clientIP, secret)
	}
}

func BenchmarkGTP_VerifyHash(b *testing.B) {
	clientIP := "127.0.0.1"
	tourLength := 5
	secret := "test-secret-key"
	guideSecrets := []string{"test-guide-secret-key-1", "test-guide-secret-key-2"}
	initialHash := hexHash("820888B1A040503A82AFA97EB0AE59E8214866C2D74F3DBC705A002FB17C86E9")
	lastHash := hexHash("4DDD8EB388374180706E41B24A19AA29B5E58A7281E6CBEEA4C8BD223D3A4B67")
	gtp := NewCalc(time.Now)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		gtp.VerifyHash(initialHash, lastHash, tourLength, clientIP, secret, guideSecrets)
	}
}
