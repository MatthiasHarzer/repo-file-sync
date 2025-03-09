package restore

import (
	"ide-config-sync/commands"
	"ide-config-sync/ide"
	"ide-config-sync/persistance"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	baseDir string
)

func init() {
	Command.Flags().StringVarP(&baseDir, "dir", "d", "", "The directory to search for repositories. Defaults to the current working directory.")
}

var Command = &cobra.Command{
	Use:   "restore",
	Short: "Restore IDE configs from the database",
	Long:  "Restore IDE configs from the database",
	RunE: func(c *cobra.Command, args []string) error {
		db, repos, err := commands.Setup(baseDir, persistance.DefaultDatabaseDir)
		if err != nil {
			panic(err)
		}

		for repo := range repos {
			println(color.GreenString("+"), "Discovered", color.GreenString(repo))

			knownConfigs, err := db.Read(repo)
			if err != nil {
				panic(err)
			}

			for _, config := range knownConfigs {
				err = ide.WriteIDEFolder(repo, config)
				if err != nil {
					color.Red("Failed to restore %s: %s", config.RelativePath, err)
					continue
				}

				println(color.BlueString("  +"), "IDE config restored:", color.BlueString(config.RelativePath))
			}
		}

		return nil
	},
}
