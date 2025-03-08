package repository

import (
	"fmt"
	"net/url"
	"path/filepath"
	"private-ide-config-sync/persistance"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/otiai10/copy"
	"golang.org/x/exp/maps"
)

var ideFolders = []string{
	".idea",
	".vscode",
}

func getOrigins(repoDir string) ([]url.URL, error) {
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

type ConfigManager struct {
	LocalRepoDirectory string
	remotes            []url.URL
	db                 *persistance.DatabaseRepo
}

func NewRepoConfigManager(db *persistance.DatabaseRepo, localRepoDirectory string) (*ConfigManager, error) {
	if strings.HasSuffix(localRepoDirectory, ".git") {
		localRepoDirectory = strings.TrimSuffix(localRepoDirectory, ".git")
	}
	localRepoDirectory = filepath.ToSlash(localRepoDirectory)

	remoteURLs, err := getOrigins(localRepoDirectory)
	if err != nil {
		return nil, err
	}

	return &ConfigManager{
		LocalRepoDirectory: localRepoDirectory,
		remotes:            remoteURLs,
		db:                 db,
	}, err
}

func (c *ConfigManager) getIdeFolderPaths() ([]string, error) {
	return persistance.FindFolders(c.LocalRepoDirectory, ideFolders)
}

func (c *ConfigManager) SaveOrigin(originURL url.URL) error {
	idePaths, err := c.getIdeFolderPaths()
	if err != nil {
		return err
	}

	for _, ideFolder := range idePaths {
		err := c.db.Write(originURL, c.LocalRepoDirectory, ideFolder)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ConfigManager) Save() error {
	for _, remote := range c.remotes {
		err := c.SaveOrigin(remote)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ConfigManager) PullRemote(origin url.URL) error {
	dbPaths, err := c.db.Read(origin)
	if err != nil {
		return err
	}

	for _, path := range dbPaths {
		repoFolderPath := fmt.Sprintf("%s/%s", c.LocalRepoDirectory, path.RelativePath)
		err := copy.Copy(path.FsPath, repoFolderPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *ConfigManager) Pull() error {
	for _, remote := range c.remotes {
		err := c.PullRemote(remote)
		if err != nil {
			return err
		}
	}
	return nil
}
