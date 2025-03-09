package fsutil

import (
	"errors"
	"io/fs"
	"os"
)

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func HomeDir() string {
	home, _ := os.UserHomeDir()
	return home
}

func IsDirectory(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fi.IsDir(), nil
}

func IsDirectoryEmpty(path string) (bool, error) {
	isDir, err := IsDirectory(path)
	if err != nil {
		return false, err
	}
	if !isDir {
		return false, errors.New("not a directory")
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}

	return len(entries) == 0, nil
}
