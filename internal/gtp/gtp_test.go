package gtp

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

func now() time.Time {
	return time.Date(2022, 4, 3, 13, 3, 29, 0, time.UTC)
}

type GTPSuite struct {
	suite.Suite
	gtp GTP
}

func (s *GTPSuite) SetupSuite() {
	s.gtp = NewGTP(now)
}

func (s *GTPSuite) TestGTP_NewInitialGTPMsg() {
	clientIP := "127.0.0.1"
	tourLength := 5
	secret := "test-secret-key"

	hash := s.gtp.CalcInitialHash(clientIP, tourLength, secret)
	s.Equal(Hash{188, 50, 238, 132, 222, 117, 223, 120, 12, 44, 45, 67, 206, 160, 197, 63, 165, 211, 117, 233}, hash)
}

func (s *GTPSuite) TestGTP_NewGTPMsg() {
	clientIP := "127.0.0.1"
	prevHash := Hash{188, 50, 238, 132, 222, 117, 223, 120, 12, 44, 45, 67, 206, 160, 197, 63, 165, 211, 117, 233}
	tourNumber := 1
	tourLength := 5
	secret := "test-secret-key-2"

	hash := s.gtp.CalcGuideHash(prevHash, tourNumber, tourLength, clientIP, secret)
	s.Equal(Hash{109, 110, 182, 130, 70, 255, 23, 151, 42, 59, 23, 94, 135, 57, 235, 196, 65, 41, 151, 178}, hash)
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
			initialHash: Hash{188, 50, 238, 132, 222, 117, 223, 120, 12, 44, 45, 67, 206, 160, 197, 63, 165, 211, 117, 233},
			lastHash:    Hash{197, 171, 194, 168, 171, 131, 188, 213, 250, 233, 86, 175, 183, 149, 123, 254, 75, 7, 98, 70},
			result:      true,
		},
		{
			initialHash: Hash{188, 50, 238, 132, 222, 117, 223, 120, 12, 44, 45, 67, 206, 160, 197, 63, 165, 211, 117, 233},
			lastHash:    Hash{198, 171, 194, 168, 171, 131, 188, 213, 250, 233, 86, 175, 183, 149, 123, 254, 75, 7, 98, 70},
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
			gtp := NewGTP(tt.now)
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
		result     int64
	}{
		{
			hash:       Hash{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			guideCount: 3,
			result:     0,
		},
		{
			hash:       Hash{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			guideCount: 3,
			result:     1,
		},
		{
			hash:       Hash{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
			guideCount: 3,
			result:     2,
		},
		{
			hash:       Hash{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3},
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
