package tx

import (
	"piprim.net/gbcl/app"
)

func New(from, to app.Account, value uint, data string) app.Tx {
	return app.Tx{
		From:  from,
		To:    to,
		Value: value,
		Data:  data,
	}
}
