package repository

import (
	"net/url"

	"github.com/go-git/go-git/v5"
	"golang.org/x/exp/maps"
)

func ReadRemotes(repo string) ([]url.URL, error) {
	r, err := git.PlainOpen(repo)
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
