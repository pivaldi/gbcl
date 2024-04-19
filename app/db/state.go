package db

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	appaccount "piprim.net/gbcl/app/account"
	"piprim.net/gbcl/app/config"
	dbblock "piprim.net/gbcl/app/db/block"
	tx "piprim.net/gbcl/app/tx"
	apptype "piprim.net/gbcl/app/type"
	libfile "piprim.net/gbcl/lib/file"
)

type State struct {
	Balances        map[appaccount.Account]uint
	txMempool       []tx.Tx
	dbFile          *os.File
	latestBlockHash apptype.Hash
}

func (s *State) LatestBlockHash() apptype.Hash {
	return s.latestBlockHash
}

func (s *State) apply(txv tx.Tx) error {
	if txv.IsReward() {
		s.Balances[txv.To] += txv.Value
		return nil
	}

	if txv.Value > s.Balances[txv.From] {
		return errors.New("insufficient balance")
	}

	s.Balances[txv.From] -= txv.Value
	s.Balances[txv.To] += txv.Value

	return nil
}

func (s *State) AddTx(t tx.Tx) error {
	if err := s.apply(t); err != nil {
		return err
	}
	s.txMempool = append(s.txMempool, t)

	return nil
}

func (s *State) AddBlock(b dbblock.Block) error {
	for _, tx := range b.TXs {
		if err := s.AddTx(tx); err != nil {
			return err
		}
	}

	return nil
}

func (s *State) Persist() (apptype.Hash, error) {
	block := dbblock.New(
		s.latestBlockHash,
		uint64(time.Now().Unix()),
		s.txMempool,
	)

	blockHash, err := block.Hash()
	if err != nil {
		return apptype.Hash{}, errors.Wrap(err, "")
	}

	blockFs := dbblock.BlockFS{Key: blockHash, Value: block}

	blockFsJSON, err := json.Marshal(blockFs)
	if err != nil {
		return apptype.Hash{}, errors.Wrap(err, "")
	}

	fmt.Printf("Persisting new Block to disk:\n")
	fmt.Printf("\t%s\n", blockFsJSON)

	if _, err = s.dbFile.Write(append(blockFsJSON, '\n')); err != nil {
		return apptype.Hash{}, errors.Wrap(err, "")
	}
	s.latestBlockHash = blockHash

	s.txMempool = []tx.Tx{}

	return s.latestBlockHash, nil
}

func (s *State) Close() {
	if s != nil {
		_ = s.dbFile.Close()
	}
}

func NewStateFromDisk() (*State, error) {
	conf := config.Get()

	gen, err := getGenesisFromMemory()
	if err != nil {
		return nil, err
	}

	balances := make(map[appaccount.Account]uint)
	for account, balance := range gen.Balances {
		balances[account] = balance
	}

	f, err := os.OpenFile(conf.GetDBFilePath(), os.O_CREATE|os.O_APPEND|os.O_RDWR, libfile.GetDefaultFileMode())
	if err != nil {
		return nil, errors.Wrap(err, "error reading database")
	}

	scanner := bufio.NewScanner(f)
	state := &State{
		Balances:        balances,
		txMempool:       make([]tx.Tx, 0),
		dbFile:          f,
		latestBlockHash: apptype.Hash{},
	}

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, errors.Wrap(err, "error reading database")
		}

		blockFsJSON := scanner.Bytes()
		var blockFs dbblock.BlockFS
		err = json.Unmarshal(blockFsJSON, &blockFs)
		if err != nil {
			return nil, errors.Wrap(err, "error unmarshaling block")
		}

		err = state.applyBlock(blockFs.Value)
		if err != nil {
			return nil, err
		}

		state.latestBlockHash = blockFs.Key
	}

	return state, nil
}

func (s *State) applyBlock(b dbblock.Block) error {
	for _, tx := range b.TXs {
		if err := s.apply(tx); err != nil {
			return err
		}
	}

	return nil
}
