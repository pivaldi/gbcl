package libfile

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// Exists returns whether a file or directory exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// IdDir determines if a file represented
// by `path` is a directory or not.
func IsDir(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, errors.Wrap(err, "")
	}

	return fileInfo.IsDir(), err
}

func ensureNotEmpty(s string) (err error) {
	if s == "" {
		err = errors.New("empty string given")
	}

	return
}

// CreateTargetDirIfNotExists ensures that the target directory
// of the file/directory exists, creating it if it does not exist.
// If "path" ends by a slash it's
// considering as a directory, else it's a file
func CreateTargetDirIfNotExists(path string) (err error) {
	if err = ensureNotEmpty(path); err != nil {
		return
	}

	if path[len(path)-1] != '/' {
		path = filepath.Dir(path)
	}

	err = CreateDirIfNotExists(path)

	return
}

var (
	defaultPermsNum = 0775
	defaultPerms    = fs.FileMode(defaultPermsNum)
)

// CreateDirIfNotExists ensures that the target directory
// of the directory exists, creating it if it does not exist.
func CreateDirIfNotExists(path string) (err error) {
	if err = ensureNotEmpty(path); err != nil {
		return
	}

	if !Exists(path) {
		if err = os.MkdirAll(path, defaultPerms); err != nil {
			return
		}
	}

	return
}
