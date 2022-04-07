package shield

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_randomSecret_Length(t *testing.T) {
	args := []struct {
		length int
	}{
		{length: 5},
		{length: 10},
		{length: 15},
		{length: 20},
	}
	for _, tt := range args {
		t.Run("test-random-secret-length", func(t *testing.T) {
			secret, err := randomSecret(tt.length)
			assert.NoError(t, err)
			assert.Equal(t, tt.length, len(secret))
		})
	}
}

func Test_randomSecret_Unique(t *testing.T) {
	secret1, err := randomSecret(10)
	assert.NoError(t, err)
	secret2, err := randomSecret(10)
	assert.NoError(t, err)
	assert.NotEqual(t, secret1, secret2)
}
