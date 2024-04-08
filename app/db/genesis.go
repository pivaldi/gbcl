package appdb

import (
	"encoding/json"

	"github.com/pkg/errors"
	"piprim.net/gbcl/app"
	apptype "piprim.net/gbcl/app/type"
)

type genesis struct {
	Balances map[apptype.Account]uint `json:"balances"`
}

func loadGenesis() (*genesis, error) {
	content, err := app.FS.ReadFile("etc/db/genesis.json")
	if err != nil {
		return nil, errors.Wrap(err, "error loadind genesis file")
	}

	loadedGenesis := new(genesis)
	err = json.Unmarshal(content, loadedGenesis)
	if err != nil {
		return nil, errors.Wrap(err, "error loadind genesis file")
	}

	return loadedGenesis, nil
}
