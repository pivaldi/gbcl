package config

import (
	"embed"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	libfile "piprim.net/gbcl/lib/file"
)

//go:embed etc/db/genesis.json
var FS embed.FS

const (
	databaseDirName  = "database"
	databaseFileName = "block.db"
	genesisFilename  = "genesis.json"
)

var config *Config
var isInit bool

type Config struct {
	dataDirectory   string
	genesisFilePath string
	dbFilePath      string
	port            uint16
}

func (c *Config) GetPort() uint16 {
	return c.port
}

func (c *Config) SetPort(port uint16) {
	c.port = port
}

func (c *Config) GetGenesisFilePath() string {
	return c.genesisFilePath
}

func (c *Config) GetDBFilePath() string {
	return c.dbFilePath
}

func Init(c *Config) {
	if isInit {
		log.Debug().Msg("App is already initialized…")
		return
	}

	config = c

	isInit = true
}

func Get() *Config {
	if !isInit {
		panic(errors.New("app config is not initialized"))
	}

	return config
}

func (c *Config) SetDataDirectory(path string) error {
	if isInit {
		return errors.New("config allready initialized")
	}

	if path == "" {
		return errors.New("empty data directory is not allowed")
	}

	c.dataDirectory = path

	err := initDataDirs(path)
	if err != nil {
		return err
	}

	c.dbFilePath = filepath.Join(c.dataDirectory, databaseDirName, databaseFileName)
	c.genesisFilePath = filepath.Join(c.dataDirectory, databaseDirName, genesisFilename)

	return nil
}

func initDataDirs(dataDir string) error {
	err := libfile.CreateDirIfNotExists(dataDir)
	if err != nil {
		return errors.Wrap(err, "")
	}

	for _, path := range []string{databaseDirName} {
		err := libfile.CreateDirIfNotExists(filepath.Join(dataDir, path))
		if err != nil {
			return errors.Wrap(err, "")
		}
	}

	return nil
}
