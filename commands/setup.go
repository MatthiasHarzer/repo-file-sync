package commands

import (
	"fmt"
	"ide-config-sync/persistance"
	"ide-config-sync/repository"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

func Setup(baseDir, dbDir string) (*persistance.DatabaseRepo, <-chan string, error) {
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

	baseDir = filepath.ToSlash(baseDir)
	dbDir = filepath.ToSlash(dbDir)

	color.RGB(200, 200, 200).Print("Using database dir ")
	color.Green(dbDir)

	db, err := persistance.NewDatabase(dbDir)
	if err != nil {
		return nil, nil, fmt.Errorf("could not open database repo: %e", err)
	}

	err = db.Pull()
	if err != nil {
		return nil, nil, err
	}

	reposCh := repository.FindRepositories(baseDir, dbDir)
	return db, reposCh, nil

	//if len(repos) == 1 {
	//	color.RGB(200, 200, 200).Print("Discovered 1 repository in ")
	//} else {
	//	color.RGB(200, 200, 200).Print(fmt.Sprintf("Discovered %d repositories in ", len(repos)))
	//}
	//color.Green(baseDir)
	//
	//var repoConfigs []*repository.ConfigManager
	//for _, repo := range repos {
	//	r, err := repository.NewRepoConfigManager(db, repo)
	//	if err != nil {
	//		return nil, nil, err
	//	}
	//	repoConfigs = append(repoConfigs, r)
	//}
	//
	//println()
	//
	//return db, repoConfigs, nil
}
