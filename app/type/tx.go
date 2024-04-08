package apptype

type Tx struct {
	From  Account `json:"from"`
	To    Account `json:"to"`
	Value uint    `json:"value"`
	Data  string  `json:"data"`
}

func (t *Tx) IsReward() bool {
	if t == nil {
		return false
	}

	return t.Data == "reward"
}
