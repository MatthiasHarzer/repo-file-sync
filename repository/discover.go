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
	_, err := git.PlainOpen(path)
	return err == nil
}

func discoverChildRepositories(base, ignoredRepo string) <-chan string {
	queue := []string{base}
	repos := make(chan string)

	go func() {
		defer close(repos)

		if isRepo(base) {
			repos <- base
			return
		}

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
				subDir := filepath.Join(dir, entry.Name())

				if shouldSkipPath(subDir, repoSearchIgnoreFolders) {
					continue
				}

				if filepath.ToSlash(subDir) == filepath.ToSlash(ignoredRepo) {
					continue
				}

				if isRepo(subDir) {
					repos <- subDir
					continue
				}

				queue = append(queue, subDir)
			}
		}
	}()

	return repos
}

// DiscoverParentRepository searches for a git repository starting from the given path
// and traversing up the directory tree. Returns the repository root path and true
// if a repository is found, otherwise returns empty string and false.
// The ignoredRepo parameter can be used to exclude a specific repository path from the search.
func DiscoverParentRepository(base, ignoredRepo string) (string, bool) {
	currentDir := base
	for {
		if filepath.ToSlash(currentDir) == filepath.ToSlash(ignoredRepo) {
			return "", false
		}

		if isRepo(currentDir) {
			return currentDir, true
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			break
		}
		currentDir = parentDir
	}
	return "", false
}

func DiscoverRepositories(base, ignoredRepo string) <-chan string {
	repos := make(chan string)

	go func() {
		defer close(repos)

		emitted := make(map[string]struct{})

		parentRepo, foundParentRepo := DiscoverParentRepository(base, ignoredRepo)
		if foundParentRepo {
			normalized := filepath.ToSlash(parentRepo)
			emitted[normalized] = struct{}{}
			repos <- parentRepo
		}

		childRepos := discoverChildRepositories(base, ignoredRepo)
		for repo := range childRepos {
			normalized := filepath.ToSlash(repo)
			if _, exists := emitted[normalized]; exists {
				continue
			}
			emitted[normalized] = struct{}{}
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
