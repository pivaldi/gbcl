package account

import (
	"os"

	"github.com/pkg/errors"
	liberrors "piprim.net/gbcl/lib/errors"
)

type Account string

func New(value string) Account {
	if value == "" {
		liberrors.HandleError(errors.New("empty account name is not allowed"))
		os.Exit(1)
	}

	return Account(value)
}
