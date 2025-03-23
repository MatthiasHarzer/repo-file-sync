package discover

import (
	"ide-config-sync/commands"
	"ide-config-sync/repository"

	"github.com/spf13/cobra"
)

var baseDir string

func init() {
	Command.Flags().StringVarP(&baseDir, "dir", "d", "", "The directory to search for repositories. Defaults to the current working directory.")
}

var Command = &cobra.Command{
	Use:   "discover",
	Short: "Discover repositories and files which will be included by the save / restore commands",
	Long:  "Discover repositories and files which will be included by the save / restore commands",
	RunE: func(c *cobra.Command, args []string) error {
		db, repos, _, err := commands.Setup(baseDir)
		if err != nil {
			panic(err)
		}

		globalDiscoverOptions, err := db.ReadGlobalDiscoveryOptions()
		if err != nil {
			panic(err)
		}

		for repo := range repos {
			println(commands.RepositoryDiscovered(repo))

			repoDiscoverOptions, err := db.ReadRepoDiscoveryOptions(repo)
			if err != nil {
				panic(err)
			}

			mergedOptions := globalDiscoverOptions.Merge(repoDiscoverOptions)
			files := repository.DiscoverRepositoryFiles(repo, mergedOptions)

			for file := range files {
				print(commands.FileProcessed(repo, file, "File discovered"))
			}
		}

		return nil
	},
}
