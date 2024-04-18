package apptype

import (
	"encoding/hex"

	"github.com/pkg/errors"
)

type Hash [32]byte

func (h Hash) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(h[:])), nil
}

func (h *Hash) UnmarshalText(data []byte) error {
	_, err := hex.Decode(h[:], data)

	return errors.Wrap(err, "unmarshaling text error")
}
