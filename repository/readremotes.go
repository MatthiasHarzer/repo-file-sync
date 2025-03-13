package repository

import (
	"github.com/go-git/go-git/v5"
	"golang.org/x/exp/maps"
)

func ReadRemotes(repo string) ([]string, error) {
	r, err := git.PlainOpen(repo)
	if err != nil {
		return nil, err
	}

	remotes, err := r.Remotes()
	if err != nil {
		return nil, err
	}

	remoteURLs := make(map[string]bool)
	for _, remote := range remotes {
		for _, urlStr := range remote.Config().URLs {
			remoteURLs[urlStr] = true
		}
	}

	return maps.Keys(remoteURLs), nil
}
