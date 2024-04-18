package db

import (
	"os"

	"github.com/pkg/errors"
	"piprim.net/gbcl/app/config"
	libfile "piprim.net/gbcl/lib/file"
)

func Init() error {
	if libfile.Exists(config.Get().GetGenesisFilePath()) {
		return nil
	}

	if err := writeGenesisToDisk(); err != nil {
		return err
	}

	if err := writeEmptyBlocksDBToDisk(); err != nil {
		return err
	}

	return nil
}

func writeEmptyBlocksDBToDisk() error {
	conf := config.Get()

	return errors.Wrap(os.WriteFile(conf.GetDBFilePath(), []byte(""), libfile.GetDefaultFileMode()), "")
}
