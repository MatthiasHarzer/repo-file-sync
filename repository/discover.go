package repository

import (
	"io/fs"
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

func isRepo(path string) bool {
	gitDir := filepath.Join(path, ".git")
	info, err := os.Stat(gitDir)
	if err != nil || !info.IsDir() {
		return false
	}
	_, err = git.PlainOpen(path)
	return err == nil
}

func discoverChildRepositories(base, ignoredRepo string) <-chan string {
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

func discoverParentRepository(base, ignoredRepo string) *string {
	currentDir := filepath.Dir(base)
	for {
		if filepath.ToSlash(currentDir) == filepath.ToSlash(ignoredRepo) {
			return nil
		}

		if isRepo(currentDir) {
			return &currentDir
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			break
		}
		currentDir = parentDir
	}
	return nil
}

func DiscoverRepositories(base, ignoredRepo string) <-chan string {
	repos := make(chan string)

	go func() {
		defer close(repos)

		parentRepo := discoverParentRepository(base, ignoredRepo)
		if parentRepo != nil {
			repos <- *parentRepo
		}

		childRepos := discoverChildRepositories(base, ignoredRepo)
		for repo := range childRepos {
			repos <- repo
		}
	}()

	return repos
}

func DiscoverRepositoryFiles(repo string, config DiscoveryOptions) <-chan File {
	files := make(chan File)

	go func() {
		defer close(files)

		for _, pattern := range config.IncludePatterns.Slice() {
			_ = doublestar.GlobWalk(os.DirFS(repo), pattern, func(path string, d fs.DirEntry) error {
				absolutePath := filepath.ToSlash(filepath.Join(repo, path))

				for _, excludePattern := range config.ExcludePatterns.Slice() {
					if matched, _ := doublestar.PathMatch(excludePattern, absolutePath); matched {
						return nil
					}
				}

				info, err := d.Info()
				if err != nil {
					return err
				}

				files <- File{
					AbsolutePath:     absolutePath,
					PathFromRepoRoot: path,
					Size:             info.Size(),
				}

				return nil
			}, doublestar.WithFilesOnly(), doublestar.WithNoFollow())
		}
	}()

	return files
}
