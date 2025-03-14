package discover

import (
	"ide-config-sync/commands"
	"ide-config-sync/ide"
	"ide-config-sync/repository"

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
		_, repos, _, err := commands.Setup(baseDir)
		if err != nil {
			panic(err)
		}

		for repo := range repos {
			remotes, err := repository.ReadRemotes(repo)
			if err != nil {
				println(color.RedString("- Failed to read remotes for"), color.RedString(repo), ":", color.RedString(err.Error()))
				continue
			}

			remotesString := ""
			for i, remote := range remotes {
				if i > 0 {
					remotesString += ", "
				}
				remotesString += color.YellowString(remote)
			}

			println(color.GreenString("+"), "Discovered", color.GreenString(repo), "("+remotesString+")")

			ideConfigs := ide.ReadIDEFolderPaths(repo)

			for path := range ideConfigs {
				println(color.BlueString("  +"), "IDE config found at", color.BlueString(path))
			}
		}

		return nil
	},
}
