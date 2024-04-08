package app

import (
	"fmt"
	"os"
	"path/filepath"

	"embed"

	"github.com/pkg/errors"

	"github.com/rs/zerolog/log"
	liberrors "piprim.net/gbcl/lib/errors"
	libfile "piprim.net/gbcl/lib/file"
)

const Major = "0"
const Minor = "1"
const Fix = "0"
const Verbal = "TX Add && Balances List"

var isInit bool

//go:embed etc/db/genesis.json
var FS embed.FS

func GetVersion() string {
	return fmt.Sprintf("%s version : %s.%s.%s-beta %s", config.Name, Major, Minor, Fix, Verbal)
}

func Message(msg string) {
	fmt.Println(msg)
}

func Init() error {
	if isInit {
		log.Debug().Msg("App is already initialized…")
		return nil
	}

	liberrors.InitLog()
	config.Name = "gbcl"
	config.ShortDescription = "The Blockchain learning CLI"

	if config.RootDirectory == "" {
		dirname, err := os.UserHomeDir()
		if err != nil {
			return errors.Wrap(err, "app initialisation error")
		}

		config.RootDirectory = filepath.Join(dirname, "."+config.Name)
	}

	err := createRootDir()

	if config.DBFilePath == "" {
		config.DBFilePath = filepath.Join(config.RootDirectory, "db", "db.txt")
	}

	isInit = true

	return err
}

func createRootDir() error {
	err := libfile.CreateDirIfNotExists(config.RootDirectory)
	if err != nil {
		return errors.Wrap(err, "")
	}

	for _, path := range []string{"db"} {
		err := libfile.CreateDirIfNotExists(filepath.Join(config.RootDirectory, path))
		if err != nil {
			return errors.Wrap(err, "")
		}
	}

	return nil
}
