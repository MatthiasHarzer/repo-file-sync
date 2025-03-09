package discover

import (
	"ide-config-sync/persistance"
	"ide-config-sync/repository"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var baseDir string

func init() {
	Command.Flags().StringVarP(&baseDir, "dir", "d", "", "The directory to search for repositories. Defaults to the current working directory.")
}

var Command = &cobra.Command{
	Use:   "discover",
	Short: "Discover IDE config files",
	Long:  "Discover IDE config files",
	RunE: func(c *cobra.Command, args []string) error {
		var err error
		if baseDir == "" {
			baseDir, err = os.Getwd()
			if err != nil {
				return err
			}
		}

		repos := repository.FindRepositories(baseDir, persistance.DefaultDatabaseDir)

		for repo := range repos {
			println(color.GreenString("+"), "Discovered ", color.GreenString(repo))

			ideConfigs := repository.GetIDEFolderPaths(repo)

			for path := range ideConfigs {
				println(color.BlueString("  +"), "IDE config found at", color.BlueString(path))
			}

		}

		return nil
	},
}
