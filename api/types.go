package api

import (
	"fmt"

	"github.com/tinylib/msgp/msgp"

	"github.com/itimky/word-of-wisom/pkg/gtp"
)

const HashType int8 = 99

func init() {
	msgp.RegisterExtension(HashType, func() msgp.Extension { return new(Hash) })
}

type Hash gtp.Hash

func (h Hash) String() string {
	return fmt.Sprintf("%X", h[:])
}

func (h *Hash) ExtensionType() int8 { return HashType }

func (h *Hash) Len() int { return gtp.HashSize }

func (h *Hash) MarshalBinaryTo(b []byte) error {
	copy(b, (*h)[:])
	return nil
}

func (h *Hash) UnmarshalBinary(b []byte) error {
	copy((*h)[:], b)
	return nil
}
