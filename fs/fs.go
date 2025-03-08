package fs

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
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

func FindFolders(root string, folderNames []string) <-chan string {
	var folderLookup = make(map[string]bool)
	for _, folderName := range folderNames {
		folderLookup[folderName] = true
	}

	folders := make(chan string)

	go func() {
		defer close(folders)

		_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				return nil
			}

			name := filepath.Base(path)

			if folderLookup[name] {
				rel, err := filepath.Rel(root, path)
				if err != nil {
					return err
				}
				folders <- filepath.ToSlash(rel)
			}
			return nil
		})
	}()
	return folders
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
