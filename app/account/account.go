package account

import (
	"os"

	"github.com/pkg/errors"
	"piprim.net/gbcl/app"
	liberrors "piprim.net/gbcl/lib/errors"
)

func New(value string) app.Account {
	if value == "" {
		liberrors.HandleError(errors.New("empty account name is not allowed"))
		os.Exit(1)
	}

	return app.Account(value)
}
