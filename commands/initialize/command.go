package initialize

import (
	"bufio"
	"fmt"
	"ide-config-sync/fs"
	"ide-config-sync/persistance"
	"net/url"
	"os"

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
	Run: func(c *cobra.Command, args []string) {
		if databaseRepositoryURL == "" {
			scanner := bufio.NewScanner(os.Stdin)
			fmt.Print("Enter the database repository URL: ")
			scanned := scanner.Scan()
			if !scanned {
				color.Red("failed to read input")
				return
			}
			databaseRepositoryURL = scanner.Text()
		}
		_, err := url.Parse(databaseRepositoryURL)
		if err != nil {
			color.Red("invalid URL: %s", err)
			return
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
				color.Red("database repository directory must be empty")
				return
			}
		} else {
			err = os.Mkdir(dbRepoPath, os.ModePerm)
			if err != nil {
				panic(err)
			}
		}

		_, err = persistance.InitializeFromURL(databaseRepositoryURL, dbRepoPath)
		if err != nil {
			color.Red("failed to create database repository: %s", err)
			return
		}

		color.Green("Database repository created at %s", dbRepoPath)
	},
}
