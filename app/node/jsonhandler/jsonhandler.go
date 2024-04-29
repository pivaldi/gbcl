package jsonhandler

import (
	"github.com/pkg/errors"
	"piprim.net/gbcl/app/account"
	"piprim.net/gbcl/app/db"
	"piprim.net/gbcl/app/tx"
	apptype "piprim.net/gbcl/app/type"
)

// State is set when starting/running node
var State *db.State

type TxAddReq struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value uint   `json:"value"`
	Data  string `json:"data"`
}

type TxAddResp struct {
	Hash apptype.Hash `json:"blockHash"`
}

type BalancesResp struct {
	Hash     apptype.Hash             `json:"blockHash"`
	Balances map[account.Account]uint `json:"balances"`
}

func TxAdd(txr *TxAddReq) (*TxAddResp, error) {
	ttx := tx.New(account.New(txr.From), account.New(txr.To), txr.Value, txr.Data)

	err := State.AddTx(ttx)
	if err != nil {
		return nil, errors.Wrap(err, "error adding tx to state")
	}

	hash, err := State.Persist()
	if err != nil {
		return nil, errors.Wrap(err, "error persisting state")
	}

	return &TxAddResp{hash}, nil
}

func ListBalances(_ *struct{}) (*BalancesResp, error) {
	return &BalancesResp{
		Hash:     State.LatestBlockHash(),
		Balances: State.Balances,
	}, nil
}
