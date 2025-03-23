package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"ide-config-sync/config"
	"ide-config-sync/database"
	"ide-config-sync/repository"

	"github.com/fatih/color"
)

func Setup(baseDir string) (database.Database, <-chan string, *config.Config, error) {
	var err error
	if baseDir == "" {
		baseDir, err = os.Getwd()
		if err != nil {
			return nil, nil, nil, err
		}
	}

	baseDir = filepath.ToSlash(baseDir)

	cfg, err := config.Load()
	if err != nil {
		return nil, nil, nil, err
	}

	db, err := database.NewRepoDatabase(cfg.DatabasePath)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not open database repo\n\ndid you run `ide-config-sync init`?")
	}

	if !cfg.LocalOnly {
		err = db.Pull()
		if err != nil {
			return nil, nil, nil, err
		}
	}

	reposCh := repository.DiscoverRepositories(baseDir, cfg.DatabasePath)

	return db, reposCh, cfg, nil
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
