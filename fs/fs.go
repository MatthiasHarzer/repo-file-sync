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

func ListDirs(path string) ([]string, error) {
	var dirs []string
	err := filepath.Walk(path, func(subPath string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && subPath != path {
			dirs = append(dirs, filepath.ToSlash(subPath))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return dirs, nil
}

func GetFolderName(path string) string {
	_, folderName := filepath.Split(path)
	return folderName
}

func FindFolders(root string, folderNames []string) ([]string, error) {
	var folderLookup = make(map[string]bool)
	for _, folderName := range folderNames {
		folderLookup[folderName] = true
	}

	var paths []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
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
			paths = append(paths, filepath.ToSlash(rel))
		}
		return nil
	})
	return paths, err
}

func HomeDir() string {
	home, _ := os.UserHomeDir()
	return home
}
