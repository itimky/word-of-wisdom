package server

import (
	"crypto/rand"
	"github.com/pkg/errors"
	"math/big"
)

func randomSecret(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", errors.WithMessage(err, "random secret")
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}
