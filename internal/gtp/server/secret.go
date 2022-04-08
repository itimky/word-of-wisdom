package server

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func randomSecret(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", fmt.Errorf("random secret: %w", err)
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}
