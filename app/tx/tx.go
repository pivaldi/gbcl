package tx

import (
	appaccount "piprim.net/gbcl/app/account"
)

type Tx struct {
	From  appaccount.Account `json:"from"`
	To    appaccount.Account `json:"to"`
	Value uint               `json:"value"`
	Data  string             `json:"data"`
}

func (t *Tx) IsReward() bool {
	if t == nil {
		return false
	}

	return t.Data == "reward"
}

func New(from, to appaccount.Account, value uint, data string) Tx {
	return Tx{
		From:  from,
		To:    to,
		Value: value,
		Data:  data,
	}
}
