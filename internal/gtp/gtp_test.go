package gtp

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func now() time.Time {
	return time.Date(2022, 4, 3, 13, 3, 29, 0, time.UTC)
}

func newMockedGTP(now func() time.Time) GTP {
	return GTP{now: now}
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
			gtp := newMockedGTP(tt.now)
			assert.Equal(t, tt.result, gtp.timestamp())
		})
	}
}

func TestGTP_NewInitialGTPMsg(t *testing.T) {
	clientIP := "127.0.0.1"
	tourLength := 5
	secret := "test-secret-key"
	gtp := newMockedGTP(now)

	gtpMsg := gtp.NewInitialGTPMsg(clientIP, tourLength, secret)
	assert.Equal(t, [20]byte{188, 50, 238, 132, 222, 117, 223, 120, 12, 44, 45, 67, 206, 160, 197, 63, 165, 211, 117, 233}, gtpMsg.Hash)
	assert.Equal(t, byte(5), gtpMsg.TourLength)
	assert.Equal(t, byte(0), gtpMsg.TourNumber)
}

func TestGTP_NewGTPMsg(t *testing.T) {
	clientIP := "127.0.0.1"
	prevHash := [20]byte{188, 50, 238, 132, 222, 117, 223, 120, 12, 44, 45, 67, 206, 160, 197, 63, 165, 211, 117, 233}
	tourNumber := 1
	tourLength := 5
	secret := "test-secret-key-2"
	gtp := newMockedGTP(now)

	gtpMsg := gtp.NewGTPMsg(prevHash, tourNumber, tourLength, clientIP, secret)
	assert.Equal(t, [20]byte{109, 110, 182, 130, 70, 255, 23, 151, 42, 59, 23, 94, 135, 57, 235, 196, 65, 41, 151, 178}, gtpMsg.Hash)
	assert.Equal(t, byte(5), gtpMsg.TourLength)
	assert.Equal(t, byte(1), gtpMsg.TourNumber)
}

func TestGuideIndex(t *testing.T) {
	args := []struct {
		hash       [20]byte
		guideCount byte
		result     int64
	}{
		{
			hash:       [20]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			guideCount: 3,
			result:     0,
		},
		{
			hash:       [20]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			guideCount: 3,
			result:     1,
		},
		{
			hash:       [20]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
			guideCount: 3,
			result:     2,
		},
		{
			hash:       [20]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3},
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
