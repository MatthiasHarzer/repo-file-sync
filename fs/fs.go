package fs

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
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

func shouldSkipDir(path string, excludePatterns []string) bool {
	for _, pattern := range excludePatterns {
		matched, _ := regexp.MatchString(pattern, filepath.Base(path))
		if matched {
			return true
		}
	}
	return false
}

func shouldIncludeDir(path string, includePatterns []string) bool {
	if len(includePatterns) == 0 {
		return true
	}

	for _, pattern := range includePatterns {
		matched, _ := regexp.MatchString(pattern, filepath.Base(path))
		if matched {
			return true
		}
	}
	return false
}

func FindFolders(root string, includePatterns []string, excludePatterns []string) <-chan string {
	queue := []string{root}

	folders := make(chan string)

	go func() {
		defer close(folders)

		for len(queue) > 0 {
			dir := queue[0]
			queue = queue[1:]

			entries, err := os.ReadDir(dir)
			if err != nil {
				continue
			}

			for _, entry := range entries {
				fullPath := filepath.Join(dir, entry.Name())
				if !entry.IsDir() {
					continue
				}

				if shouldSkipDir(fullPath, excludePatterns) {
					continue
				}
				if shouldIncludeDir(fullPath, includePatterns) {
					relPath, err := filepath.Rel(root, fullPath)
					if err != nil {
						continue
					}
					folders <- relPath
					continue
				}

				queue = append(queue, fullPath)
			}
		}

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
