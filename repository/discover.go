package repository

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"

	"github.com/go-git/go-git/v5"
)

var repoSearchIgnorePatterns = []string{
	"node_modules",
	"venv",
	".venv",
}

func shouldSkipPath(path string, excludePatterns []string) bool {
	for _, pattern := range excludePatterns {
		matched, _ := regexp.MatchString(pattern, filepath.Base(path))
		if matched {
			return true
		}
	}
	return false
}

func shouldIncludePath(path string, includePatterns []string) bool {
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

func find(root string, includePatterns []string, excludePatterns []string) <-chan File {
	queue := []string{root}

	folders := make(chan File)

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

				if shouldSkipPath(fullPath, excludePatterns) {
					continue
				}
				if shouldIncludePath(fullPath, includePatterns) {
					relPath, err := filepath.Rel(root, fullPath)
					if err != nil {
						continue
					}
					var fileType FileType
					if entry.IsDir() {
						fileType = FileTypeDir
					} else {
						fileType = FileTypeFile
					}

					folders <- File{
						Type:             fileType,
						AbsolutePath:     fullPath,
						PathFromRepoRoot: relPath,
					}
					continue
				}

				if entry.IsDir() {
					queue = append(queue, fullPath)
				}
			}
		}

	}()

	return folders
}

func DiscoverRepositories(base, ignoredRepo string) <-chan string {
	repositories := find(base, []string{"^.git$"}, repoSearchIgnorePatterns)
	repos := make(chan string)

	go func() {
		defer close(repos)

		for repository := range repositories {
			if repository.IsFile() {
				continue
			}

			repoDir := filepath.Dir(repository.AbsolutePath)
			if filepath.ToSlash(repoDir) == filepath.ToSlash(ignoredRepo) {
				continue
			}

			_, err := git.PlainOpen(repoDir)
			if err != nil {
				continue
			}

			repos <- repoDir
		}
	}()

	return repos
}

func DiscoverRepositoryFiles(repo string, config DiscoveryOptions) <-chan File {
	files := make(chan File)
	knownFiles := make(map[string]bool)

	go func() {
		defer close(files)

		foundFiles := find(repo, config.IncludePatterns.Slice(), config.ExcludePatterns.Slice())

		for file := range foundFiles {
			files <- file
			knownFiles[file.AbsolutePath] = true
		}

		// Interpret includes as globs as well
		for _, pattern := range config.IncludePatterns.Slice() {
			matching, err := fs.Glob(os.DirFS(repo), pattern)
			if err != nil {
				continue
			}

			for _, file := range matching {
				absolutePath := filepath.Join(repo, file)

				if _, ok := knownFiles[absolutePath]; ok {
					continue
				}

				files <- File{
					Type:             FileTypeFile,
					AbsolutePath:     absolutePath,
					PathFromRepoRoot: file,
				}
			}
		}

	}()

	return files
}
