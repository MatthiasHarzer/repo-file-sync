package database

import "github.com/MatthiasHarzer/repo-file-sync/repository"

type Database interface {
	WriteRepoFile(repo string, file repository.File) error
	WriteRepoDiscoveryOptions(repo string, repoConfiguration repository.DiscoveryOptions) error
	ReadRepoFiles(repo string) (<-chan repository.File, error)
	ReadRepoDiscoveryOptions(repo string) (repository.DiscoveryOptions, error)
	WriteGlobalDiscoveryOptions(globalConfig repository.DiscoveryOptions) error
	ReadGlobalDiscoveryOptions() (repository.DiscoveryOptions, error)
	Push() error
	Pull() error
}
