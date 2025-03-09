package repository

import (
	"fmt"
	"ide-config-sync/fs"
	"path/filepath"
	"strings"
)

func FindRepositories(base string, ignore string) <-chan string {
	repositories := fs.FindFolders(base, []string{".git$"}, ignorePatterns)

	repos := make(chan string)

	go func() {
		defer close(repos)

		known := make(map[string]bool)

		for repository := range repositories {
			correctSlash := filepath.ToSlash(repository)
			noDotGit := strings.TrimSuffix(correctSlash, ".git")

			f, _ := filepath.Abs(noDotGit)
			i, _ := filepath.Abs(ignore)

			if f == i {
				continue
			}

			repoPath := strings.TrimSuffix(filepath.ToSlash(fmt.Sprintf("%s/%s", base, noDotGit)), "/")
			for folder := range known {
				if strings.HasPrefix(repoPath, folder) {
					continue
				}
			}

			known[repoPath] = true

			repos <- repoPath
		}
	}()

	return repos
}
