package config

import (
	"fmt"
	"path/filepath"

	"ide-config-sync/util/fsutil"
)

var (
	StoragePath             string
	File                    string
	DefaultDatabaseRepoPath string
)

func init() {
	StoragePath = filepath.ToSlash(fmt.Sprintf("%s/.ide-config-sync", fsutil.HomeDir()))
	File = fmt.Sprintf("%s/config.json", StoragePath)
	DefaultDatabaseRepoPath = fmt.Sprintf("%s/database", StoragePath)
}
