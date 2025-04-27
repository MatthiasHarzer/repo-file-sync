package config

import (
	"fmt"
	"path/filepath"

	"github.com/MatthiasHarzer/repo-file-sync/util/fsutil"
	"github.com/MatthiasHarzer/repo-file-sync/util/units"
)

var (
	StoragePath             string
	File                    string
	DefaultDatabaseRepoPath string
	MaxFileSize             = 5 * units.MiB
)

func init() {
	StoragePath = filepath.ToSlash(fmt.Sprintf("%s/.repo-file-sync", fsutil.HomeDir()))
	File = fmt.Sprintf("%s/config.json", StoragePath)
	DefaultDatabaseRepoPath = fmt.Sprintf("%s/database", StoragePath)
}
