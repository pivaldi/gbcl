package account

import (
	"os"

	"github.com/pkg/errors"
	apptype "piprim.net/gbcl/app/type"
	liberrors "piprim.net/gbcl/lib/errors"
)

func New(value string) apptype.Account {
	if value == "" {
		liberrors.HandleError(errors.New("empty account name is not allowed"))
		os.Exit(1)
	}

	return apptype.Account(value)
}
