package appdb

import (
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"piprim.net/gbcl/app"
	apptype "piprim.net/gbcl/app/type"
)

const dbFileMode = 0600

type Snapshot [32]byte

type State struct {
	Balances  map[apptype.Account]uint
	txMempool []apptype.Tx
	dbFile    *os.File
	snapshot  Snapshot
}

func (s *State) apply(tx apptype.Tx) error {
	if tx.IsReward() {
		s.Balances[tx.To] += tx.Value
		return nil
	}

	if tx.Value > s.Balances[tx.From] {
		return errors.New("insufficient balance")
	}

	s.Balances[tx.From] -= tx.Value
	s.Balances[tx.To] += tx.Value

	return nil
}

func (s *State) Add(tx apptype.Tx) error {
	if err := s.apply(tx); err != nil {
		return err
	}
	s.txMempool = append(s.txMempool, tx)

	return nil
}

func (s *State) doSnapshot() error {
	// Re-read the whole file from the first byte
	_, err := s.dbFile.Seek(0, 0)
	if err != nil {
		return errors.Wrap(err, "")
	}

	txsData, err := io.ReadAll(s.dbFile)
	if err != nil {
		return errors.Wrap(err, "")
	}

	s.snapshot = sha256.Sum256(txsData)

	return nil
}

func (s *State) Persist() (Snapshot, error) {
	mempool := make([]apptype.Tx, len(s.txMempool))

	copy(mempool, s.txMempool)
	for i := 0; i < len(mempool); i++ {
		txJSON, err := json.Marshal(mempool[i])
		if err != nil {
			return Snapshot{}, errors.Wrap(err, "")
		}

		log.Debug().Str("tx", string(txJSON)).Msg("Persisting new TX to disk")

		if _, err = s.dbFile.Write(append(txJSON, '\n')); err != nil {
			return Snapshot{}, errors.Wrap(err, "")
		}

		err = s.doSnapshot()
		if err != nil {
			return Snapshot{}, err
		}

		log.Debug().Msg(fmt.Sprintf("New DB Snapshot: %x\n", s.snapshot))

		s.txMempool = s.txMempool[1:]
	}

	return s.snapshot, nil
}

func (s *State) Close() {
	if s != nil {
		_ = s.dbFile.Close()
	}
}

func NewStateFromDisk() (*State, error) {
	config := app.GetConfig()

	gen, err := loadGenesis()
	if err != nil {
		return nil, err
	}

	balances := make(map[apptype.Account]uint)
	for account, balance := range gen.Balances {
		balances[account] = balance
	}

	f, err := os.OpenFile(config.DBFilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, dbFileMode)
	if err != nil {
		return nil, errors.Wrap(err, "error reading database")
	}

	scanner := bufio.NewScanner(f)
	state := &State{
		Balances:  balances,
		txMempool: make([]apptype.Tx, 0),
		dbFile:    f,
	}

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, errors.Wrap(err, "error reading database")
		}

		var tx apptype.Tx
		err = json.Unmarshal(scanner.Bytes(), &tx)
		if err != nil {
			return nil, errors.Wrap(err, "error reading database")
		}

		if err := state.apply(tx); err != nil {
			return nil, err
		}
	}

	return state, nil
}
