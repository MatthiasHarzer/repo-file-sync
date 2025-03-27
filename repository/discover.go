package repository

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/go-git/go-git/v5"
)

var repoSearchIgnoreFolders = []string{
	"node_modules",
	"venv",
	".venv",
}

func shouldSkipPath(path string, excludedFolders []string) bool {
	for _, pattern := range excludedFolders {
		if strings.EqualFold(pattern, filepath.Base(path)) {
			return true
		}
	}
	return false
}

func DiscoverRepositories(base, ignoredRepo string) <-chan string {
	queue := []string{base}
	repos := make(chan string)

	go func() {
		defer close(repos)
		for len(queue) > 0 {
			dir := queue[0]
			queue = queue[1:]

			entries, err := os.ReadDir(dir)
			if err != nil {
				continue
			}

			for _, entry := range entries {
				if !entry.IsDir() {
					continue
				}
				fullPath := filepath.Join(dir, entry.Name())

				if shouldSkipPath(fullPath, repoSearchIgnoreFolders) {
					continue
				}

				if filepath.ToSlash(fullPath) == filepath.ToSlash(ignoredRepo) {
					continue
				}

				if entry.Name() == ".git" {
					_, err := git.PlainOpen(dir)
					if err != nil {
						continue
					}
					repos <- dir
					continue
				}

				queue = append(queue, fullPath)
			}
		}
	}()

	return repos
}

func DiscoverRepositoryFiles(repo string, config DiscoveryOptions) <-chan File {
	files := make(chan File)
	knownFiles := make(map[string]bool)

	go func() {
		defer close(files)

		for _, pattern := range config.IncludePatterns.Slice() {
			matching, err := doublestar.Glob(os.DirFS(repo), pattern, doublestar.WithFilesOnly())
			if err != nil {
				continue
			}

			for _, file := range matching {
				absolutePath := filepath.Join(repo, file)

				if _, ok := knownFiles[absolutePath]; ok {
					continue
				}

				files <- File{
					AbsolutePath:     absolutePath,
					PathFromRepoRoot: file,
				}
			}
		}

	}()

	return files
}
