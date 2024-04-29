package db

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"reflect"

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
	latestBlock     *dbblock.Block
	latestBlockHash apptype.Hash
}

func (s *State) LatestBlockHash() apptype.Hash {
	return s.latestBlockHash
}

func (s *State) NextBlockNumber() uint64 {
	if s == nil || s.latestBlock == nil {
		return uint64(0)
	}

	return s.LatestBlock().Header.Number + 1
}

func (s *State) LatestBlock() *dbblock.Block {
	return s.latestBlock
}

func (s *State) AddTx(t tx.Tx) error {
	if err := s.applyTX(t); err != nil {
		return err
	}
	s.txMempool = append(s.txMempool, t)

	return nil
}

func (s *State) AddBlock(b dbblock.Block) (*apptype.Hash, error) {
	pendingState := s.copy()
	err := pendingState.applyBlock(b)
	if err != nil {
		return nil, err
	}

	blockHash, err := b.Hash()
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	blockFs := dbblock.BlockFS{Key: blockHash, Value: b}

	blockFsJSON, err := json.Marshal(blockFs)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	fmt.Printf("Persisting new Block to disk:\n")
	fmt.Printf("\t%s\n", blockFsJSON)

	if _, err = s.dbFile.Write(append(blockFsJSON, '\n')); err != nil {
		return nil, errors.Wrap(err, "")
	}

	s.Balances = pendingState.Balances
	s.latestBlockHash = blockHash
	s.latestBlock = &b

	return &blockHash, nil
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

	state := &State{
		Balances:        balances,
		txMempool:       make([]tx.Tx, 0),
		dbFile:          f,
		latestBlockHash: apptype.Hash{},
	}

	scanner := bufio.NewScanner(f)
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

		state.latestBlock = &blockFs.Value
		state.latestBlockHash = blockFs.Key
	}

	return state, nil
}

func (s *State) applyBlock(b dbblock.Block) error {
	nextExpectedBlockNumber := uint64(0)
	if s.latestBlock != nil {
		nextExpectedBlockNumber = s.latestBlock.Header.Number + 1
	}

	if b.Header.Number != nextExpectedBlockNumber {
		return fmt.Errorf("next expected block must be '%d' not '%d'", nextExpectedBlockNumber, b.Header.Number)
	}

	if nextExpectedBlockNumber > 0 && !reflect.DeepEqual(b.Header.Parent, s.latestBlockHash) {
		return fmt.Errorf("next block parent hash must be '%x' not '%x'", s.latestBlockHash, b.Header.Parent)
	}

	return s.applyTXs(b.TXs)
}

func (s *State) applyTX(txv tx.Tx) error {
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

func (s *State) applyTXs(txs []tx.Tx) error {
	for _, tx := range txs {
		if err := s.applyTX(tx); err != nil {
			return err
		}
	}

	return nil
}

func (s *State) copy() State {
	c := State{}
	c.latestBlock = s.latestBlock
	c.latestBlockHash = s.latestBlockHash
	c.txMempool = make([]tx.Tx, len(s.txMempool))
	c.Balances = make(map[appaccount.Account]uint)

	for acc, balance := range s.Balances {
		c.Balances[acc] = balance
	}

	c.txMempool = append(c.txMempool, s.txMempool...)

	return c
}
