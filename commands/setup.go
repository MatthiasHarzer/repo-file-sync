package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"repo-file-sync/config"
	"repo-file-sync/database"
	"repo-file-sync/repository"

	"github.com/fatih/color"
)

func uninitializedError(message string) error {
	return fmt.Errorf("%s\n\ndid you run `repo-file-sync init`?`", message)
}

func Setup(baseDir string) (db database.Database, repos <-chan string, cfg *config.Config, globalDiscoveryOptions *repository.DiscoveryOptions, err error) {
	if baseDir == "" {
		baseDir, err = os.Getwd()
		if err != nil {
			return nil, nil, nil, nil, err
		}
	}

	baseDir = filepath.ToSlash(baseDir)

	cfg, err = config.Load()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	db, err = database.NewRepoDatabase(cfg.DatabasePath)
	if err != nil {
		return nil, nil, nil, nil, uninitializedError("could not open database repo")
	}

	if !cfg.LocalOnly {
		println("Pulling changes from", color.GreenString(cfg.DatabaseRepoURL))
		err = db.Pull()
		if err != nil {
			return nil, nil, nil, nil, err
		}
	}

	options, err := db.ReadGlobalDiscoveryOptions()
	if err != nil {
		return nil, nil, nil, nil, uninitializedError("could not read global discovery options")
	}

	reposCh := repository.DiscoverRepositories(baseDir, cfg.DatabasePath)

	return db, reposCh, cfg, &options, nil
}

func Push(cfg *config.Config, db database.Database) error {
	if cfg.LocalOnly {
		return nil
	}

	println("Pushing changes to", color.GreenString(cfg.DatabaseRepoURL))
	err := db.Push()
	if err != nil {
		fmt.Printf("unable to push changes: %s\n", err)
		panic(err)
	}

	color.RGB(200, 200, 200).Print("Pushed changes")

	return nil
}
