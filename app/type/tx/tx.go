package tx

import apptype "piprim.net/gbcl/app/type"

func New(from, to apptype.Account, value uint, data string) apptype.Tx {
	return apptype.Tx{
		From:  from,
		To:    to,
		Value: value,
		Data:  data,
	}
}
