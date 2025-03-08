package repository

import (
	"fmt"
	"net/url"
	"private-ide-config-sync/fs"

	"github.com/go-git/go-git/v5"
	"github.com/otiai10/copy"
	"golang.org/x/exp/maps"
)

var ideFolders = []string{
	".idea",
	".vscode",
}

type IDEConfig struct {
	FsPath       string
	RelativePath string
}

func GetOrigins(repoDir string) ([]url.URL, error) {
	r, err := git.PlainOpen(repoDir)
	if err != nil {
		return nil, err
	}

	remotes, err := r.Remotes()
	if err != nil {
		return nil, err
	}

	remoteURLs := make(map[url.URL]bool)
	for _, remote := range remotes {
		for _, urlStr := range remote.Config().URLs {
			asURL, err := url.Parse(urlStr)
			if err != nil {
				return nil, err
			}
			remoteURLs[*asURL] = true
		}
	}

	return maps.Keys(remoteURLs), nil
}

func GetIDEFolderPaths(base string) <-chan string {
	return fs.FindFolders(base, ideFolders)
}

func OverwriteIDEFolder(repo string, config IDEConfig) error {
	repoFolderPath := fmt.Sprintf("%s/%s", repo, config.RelativePath)
	err := copy.Copy(config.FsPath, repoFolderPath)
	if err != nil {
		return err
	}
	return nil
}
