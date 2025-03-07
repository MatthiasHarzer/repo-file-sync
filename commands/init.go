package commands

import (
	"fmt"
	"os"
	"private-ide-config-sync/persistance"
	"private-ide-config-sync/repository"
)

func Init(baseDir, dbDir string) (*persistance.DatabaseRepo, []*repository.RepoConfigManager, error) {
	var err error
	if baseDir == "" {
		baseDir, err = os.Getwd()
		if err != nil {
			return nil, nil, err
		}
	}

	if dbDir == "" {
		dbDir = persistance.DefaultDatabaseDir
	}

	fmt.Printf("Using database dir %s\n", dbDir)
	fmt.Printf("Discovering repositories in %s\n", baseDir)

	db, err := persistance.NewDatabase(dbDir)
	if err != nil {
		return nil, nil, fmt.Errorf("could not open database repo: %e", err)
	}

	err = db.Pull()
	if err != nil {
		return nil, nil, err
	}

	repos, err := repository.FindRepositories(baseDir, dbDir)
	if err != nil {
		return nil, nil, err
	}
	var repoConfigs []*repository.RepoConfigManager
	for _, repo := range repos {
		r, err := repository.NewRepoConfigManager(db, repo)
		if err != nil {
			return nil, nil, err
		}
		repoConfigs = append(repoConfigs, r)
	}

	return db, repoConfigs, nil
}
