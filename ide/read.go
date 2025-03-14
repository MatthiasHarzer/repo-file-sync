package ide

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
)

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

func findFolders(root string, includePatterns []string, excludePatterns []string) <-chan string {
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

func ReadIDEFolderPaths(repo string) <-chan string {
	return findFolders(repo, ideFolders, ignorePatterns)
}

func FindRepositories(base string, ignoredRepo string) <-chan string {
	repoIgnorePatterns := append([]string{regexp.QuoteMeta(ignoredRepo)}, ignorePatterns...)
	repositories := findFolders(base, []string{"^.git$"}, repoIgnorePatterns)

	repos := make(chan string)

	go func() {
		defer close(repos)

		known := make(map[string]bool)

		for repository := range repositories {
			correctSlash := filepath.ToSlash(repository)
			noDotGit := strings.TrimSuffix(correctSlash, ".git")
			repoPath := strings.TrimSuffix(filepath.ToSlash(fmt.Sprintf("%s/%s", base, noDotGit)), "/")
			for folder := range known {
				if strings.HasPrefix(repoPath, folder) {
					continue
				}
			}

			_, err := git.PlainOpen(repoPath)
			if err != nil {
				continue
			}

			known[repoPath] = true
			repos <- repoPath
		}
	}()

	return repos
}
