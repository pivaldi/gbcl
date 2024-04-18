package dbblock

import (
	"crypto/sha256"
	"encoding/json"

	"github.com/pkg/errors"
	tx "piprim.net/gbcl/app/tx"
	apptype "piprim.net/gbcl/app/type"
)

type Header struct {
	Parent apptype.Hash // parent block reference
	Time   uint64
}

type Block struct {
	Header Header
	TXs    []tx.Tx // new transactions only (payload)
}

func (b *Block) Hash() (apptype.Hash, error) {
	blockJSON, err := json.Marshal(b)

	return sha256.Sum256(blockJSON), errors.Wrap(err, "")
}

func New(parent apptype.Hash, time uint64, txs []tx.Tx) Block {
	return Block{
		Header: Header{Parent: parent, Time: time},
		TXs:    txs,
	}
}

type BlockFS struct {
	Key   apptype.Hash `json:"hash"`
	Value Block        `json:"block"`
}
