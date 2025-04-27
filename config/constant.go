package config

import (
	"fmt"
	"path/filepath"

	"github.com/MatthiasHarzer/repo-file-sync/util/fsutil"
)

var (
	StoragePath             string
	File                    string
	DefaultDatabaseRepoPath string
	MaxFileSize             int64 = 5 * 1024 * 1024 // 5 MB
)

func init() {
	StoragePath = filepath.ToSlash(fmt.Sprintf("%s/.repo-file-sync", fsutil.HomeDir()))
	File = fmt.Sprintf("%s/config.json", StoragePath)
	DefaultDatabaseRepoPath = fmt.Sprintf("%s/database", StoragePath)
}
