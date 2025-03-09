package commands

import (
	"fmt"
	"ide-config-sync/ide"
	"ide-config-sync/persistance"
	"os"
	"path/filepath"
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

	db, err := persistance.NewDatabase(dbDir)
	if err != nil {
		return nil, nil, fmt.Errorf("could not open database repo\n\ndid you run `ide-config-sync init`?")
	}

	err = db.Pull()
	if err != nil {
		return nil, nil, err
	}

	reposCh := ide.FindRepositories(baseDir, dbDir)
	return db, reposCh, nil
}
