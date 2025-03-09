package commands

import (
	"fmt"
	"ide-config-sync/config"
	"ide-config-sync/ide"
	"ide-config-sync/persistance"
	"os"
	"path/filepath"
)

func Setup(baseDir string) (*persistance.DatabaseRepo, <-chan string, *config.Config, error) {
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

	db, err := persistance.NewDatabase(cfg.DatabasePath)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not open database repo\n\ndid you run `ide-config-sync init`?")
	}

	if !cfg.LocalOnly {
		err = db.Pull()
		if err != nil {
			return nil, nil, nil, err
		}
	}

	reposCh := ide.FindRepositories(baseDir, cfg.DatabasePath)

	return db, reposCh, cfg, nil
}
