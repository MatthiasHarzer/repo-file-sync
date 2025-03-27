package discover

import (
	"repo-file-sync/commands"
	"repo-file-sync/repository"

	"github.com/fatih/color"
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
		db, dir, repos, _, globalDiscoveryOptions, err := commands.Setup(baseDir)
		if err != nil {
			panic(err)
		}

		println("Discovering repositories in", color.GreenString(dir))
		for repo := range repos {
			println(commands.RepositoryDiscovered(repo))

			repoDiscoverOptions, err := db.ReadRepoDiscoveryOptions(repo)
			if err != nil {
				panic(err)
			}

			mergedOptions := globalDiscoveryOptions.Merge(repoDiscoverOptions)
			files := repository.DiscoverRepositoryFiles(repo, mergedOptions)

			for file := range files {
				println(color.BlueString("  +"), "File discovered:", color.BlueString(file.PathFromRepoRoot))
			}
		}

		return nil
	},
}
