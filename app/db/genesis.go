package db

import (
	"encoding/json"
	"os"
	"time"

	"github.com/pkg/errors"
	"piprim.net/gbcl/app/account"
	"piprim.net/gbcl/app/config"
	libfile "piprim.net/gbcl/lib/file"
)

type genesis struct {
	Balances    map[account.Account]uint `json:"balances"`
	GenesisTime time.Time                `json:"genesisTime"`
	ChainID     string                   `json:"chainId"`
}

func getGenesisFromMemory() (*genesis, error) {
	content, err := config.FS.ReadFile("etc/db/genesis.json")
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

func writeGenesisToDisk() error {
	g, err := getGenesisFromMemory()
	if err != nil {
		return err
	}

	j, err := json.Marshal(g)
	if err != nil {
		return errors.Wrap(err, "")
	}

	path := getGenesisJSONFilePath()
	err = os.WriteFile(path, j, libfile.GetDefaultFileMode())

	return errors.Wrap(err, "writeGenesisToDisk error on "+path)
}

func getGenesisJSONFilePath() string {
	conf := config.Get()

	return conf.GetGenesisFilePath()
}
