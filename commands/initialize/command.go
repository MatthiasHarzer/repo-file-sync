package initialize

import (
	"fmt"
	"os"
	"private-ide-config-sync/fs"
	"private-ide-config-sync/persistance"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var databaseRepositoryURL string

func init() {
	Command.Flags().StringVar(&databaseRepositoryURL, "url", "", "Database repository URL")

}

var Command = &cobra.Command{
	Use:   "init",
	Short: "Initialize IDE config sync",
	Long:  "Initialize IDE config sync",
	RunE: func(c *cobra.Command, args []string) error {
		if databaseRepositoryURL == "" {
			return fmt.Errorf("you must provide a database repository URL")
		}

		dbRepoPath := persistance.DefaultDatabaseDir

		exists, err := fs.Exists(dbRepoPath)
		if err != nil {
			panic(err)
		}

		if exists {
			isEmpty, err := fs.IsDirectoryEmpty(dbRepoPath)
			if err != nil {
				panic(err)
			}
			if !isEmpty {
				return fmt.Errorf("database repository directory must be empty")
			}
		} else {
			err = os.Mkdir(dbRepoPath, os.ModePerm)
			if err != nil {
				panic(err)
			}
		}

		_, err = persistance.NewDatabaseFromURL(databaseRepositoryURL, dbRepoPath)
		if err != nil {
			panic(fmt.Errorf("failed to create database repository: %s", err))
		}

		color.Green("Database repository created at %s", dbRepoPath)
		return nil
	},
}
