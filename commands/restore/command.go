package restore

import (
	"private-ide-config-sync/commands"
	"private-ide-config-sync/repository"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	baseDir string
	dbDir   string
)

func init() {
	Command.Flags().StringVarP(&baseDir, "dir", "d", "", "Base directory")
	Command.Flags().StringVar(&dbDir, "db-dir", "", "DB directory")
}

var Command = &cobra.Command{
	Use:   "restore",
	Short: "Restore IDE configs from the database",
	Long:  "Restore IDE configs from the database",
	RunE: func(c *cobra.Command, args []string) error {
		db, repos, err := commands.Setup(baseDir, dbDir)
		if err != nil {
			panic(err)
		}

		for repo := range repos {
			println(color.GreenString("+"), "Discovered", color.GreenString(repo))

			origins, err := repository.GetOrigins(repo)
			if err != nil {
				panic(err)
			}

			knownConfigs, err := db.Read(origins)
			if err != nil {
				panic(err)
			}

			for _, config := range knownConfigs {
				err = repository.OverwriteIDEFolder(repo, config)
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
