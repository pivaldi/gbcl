package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"piprim.net/gbcl/app/config"
	"piprim.net/gbcl/app/db"
	"piprim.net/gbcl/app/node"
)

const (
	Major            = "0"
	Minor            = "1"
	Fix              = "0"
	Verbal           = "TX Add && Balances List"
	Name             = "gbcl"
	ShortDescription = "The Blockchain learning CLI"
)

var isInit bool

func GetVersion() string {
	return fmt.Sprintf("%s version : %s.%s.%s-beta %s", Name, Major, Minor, Fix, Verbal)
}

func Message(msg string) {
	fmt.Println(msg)
}

func Init(dataDir string) error {
	if isInit {
		return nil
	}

	conf := new(config.Config)

	if dataDir == "" {
		dirname, err := os.UserHomeDir()

		if err != nil {
			return errors.Wrap(err, "app initialisation error")
		}

		dataDir = filepath.Join(dirname, "."+Name)
	}

	err := conf.SetDataDirectory(dataDir)
	if err != nil {
		return errors.Wrap(err, "app error on init")
	}

	conf.SetPort(node.DefaultPort)
	config.Init(conf)

	err = db.Init()
	if err != nil {
		return err
	}

	isInit = true

	return nil
}
