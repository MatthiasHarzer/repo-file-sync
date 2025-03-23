package config

import (
	"fmt"
	"path/filepath"

	"repo-file-sync/util/fsutil"
)

var (
	StoragePath             string
	File                    string
	DefaultDatabaseRepoPath string
)

func init() {
	StoragePath = filepath.ToSlash(fmt.Sprintf("%s/.repo-file-sync", fsutil.HomeDir()))
	File = fmt.Sprintf("%s/config.json", StoragePath)
	DefaultDatabaseRepoPath = fmt.Sprintf("%s/database", StoragePath)
}
