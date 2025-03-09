package config

import (
	"fmt"
	"ide-config-sync/fsutil"
	"path/filepath"
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
