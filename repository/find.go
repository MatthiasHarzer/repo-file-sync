package repository

import (
	"fmt"
	"path/filepath"
	"private-ide-config-sync/persistance"
	"strings"
)

func FindRepositories(base string, ignore string) ([]string, error) {
	repositories, err := persistance.FindFolders(base, []string{".git"})
	if err != nil {
		return nil, err
	}

	var repositoryPaths []string
	for _, repository := range repositories {
		correctSlash := filepath.ToSlash(repository)
		noDotGit := strings.TrimSuffix(correctSlash, ".git")

		f, err := filepath.Abs(noDotGit)
		if err != nil {
			return nil, err
		}
		i, err := filepath.Abs(ignore)
		if err != nil {
			return nil, err
		}

		if f == i {
			continue
		}
		repositoryPaths = append(repositoryPaths, fmt.Sprintf("%s/%s", base, noDotGit))
	}

	return repositoryPaths, nil
}
