package testutils

import (
	"encoding/hex"
	"fmt"

	"github.com/itimky/word-of-wisom/pkg/gtp"
)

func HexHash(hexStr string) gtp.Hash {
	if len(hexStr) != 2*gtp.HashSize {
		panic(fmt.Sprintf("wrong hash size (%v)", len(hexStr)))
	}

	hashSlice, err := hex.DecodeString(hexStr)
	if err != nil {
		panic(fmt.Sprintf("hex decode string %s", hexStr))
	}
	var hash gtp.Hash
	copy(hash[:], hashSlice)
	return hash
}
