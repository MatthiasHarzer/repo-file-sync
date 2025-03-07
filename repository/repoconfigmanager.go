package repository

import (
	"fmt"
	"net/url"
	"path/filepath"
	"private-ide-config-sync/fs"
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

type RepoConfigManager struct {
	LocalRepoDirectory string
	remotes            []url.URL
	db                 *persistance.DatabaseRepo
}

func NewRepoConfigManager(db *persistance.DatabaseRepo, localRepoDirectory string) (*RepoConfigManager, error) {
	if strings.HasSuffix(localRepoDirectory, ".git") {
		localRepoDirectory = strings.TrimSuffix(localRepoDirectory, ".git")
	}
	localRepoDirectory = filepath.ToSlash(localRepoDirectory)

	remoteURLs, err := getOrigins(localRepoDirectory)
	if err != nil {
		return nil, err
	}

	return &RepoConfigManager{
		LocalRepoDirectory: localRepoDirectory,
		remotes:            remoteURLs,
		db:                 db,
	}, err
}

func (g *RepoConfigManager) getIdeFolderPaths() ([]string, error) {
	return fs.FindFolders(g.LocalRepoDirectory, ideFolders)
}

func (g *RepoConfigManager) SaveOrigin(originURL url.URL) error {
	idePaths, err := g.getIdeFolderPaths()
	if err != nil {
		return err
	}

	for _, ideFolder := range idePaths {
		err := g.db.Write(originURL, g.LocalRepoDirectory, ideFolder)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *RepoConfigManager) Save() error {
	for _, remote := range g.remotes {
		err := g.SaveOrigin(remote)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *RepoConfigManager) PullRemote(origin url.URL) error {
	dbPaths, err := g.db.Read(origin)
	if err != nil {
		return err
	}

	for _, path := range dbPaths {
		repoFolderPath := fmt.Sprintf("%s/%s", g.LocalRepoDirectory, path.RelativePath)
		err := copy.Copy(path.FsPath, repoFolderPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *RepoConfigManager) Pull() error {
	for _, remote := range g.remotes {
		err := g.PullRemote(remote)
		if err != nil {
			return err
		}
	}
	return nil
}
