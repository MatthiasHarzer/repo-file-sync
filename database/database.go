package database

import "ide-config-sync/repository"

type Database interface {
	WriteRepoFile(repo string, file repository.File) error
	WriteRepoOptions(repo string, repoConfiguration repository.Options) error
	ReadRepoFiles(repo string) ([]repository.File, error)
	ReadRepoOptions(repo string) (repository.Options, error)
	WriteGlobalOptions(globalConfig repository.Options) error
	ReadGlobalOptions() (repository.Options, error)
	Push() error
	Pull() error
}
