package save

import (
	"fmt"
	"ide-config-sync/commands"
	"ide-config-sync/persistance"
	"ide-config-sync/repository"

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
	Use:   "save",
	Short: "Save IDE configurations to the database",
	Long:  "Save IDE configurations to the database",
	RunE: func(c *cobra.Command, args []string) error {
		db, repos, err := commands.Setup(baseDir, persistance.DefaultDatabaseDir)
		if err != nil {
			panic(err)
		}

		for repo := range repos {
			println(color.GreenString("+"), "Discovered", color.GreenString(repo))

			origins, err := repository.GetOrigins(repo)
			if err != nil {
				panic(err)
			}

			ideFolders := repository.GetIDEFolderPaths(repo)
			for ideConfig := range ideFolders {
				err := db.Write(origins, repo, ideConfig)
				if err != nil {
					color.Red("Failed to write %s to database: %s", ideConfig, err)
					continue
				}

				println(color.BlueString("  +"), "IDE config saved:", color.BlueString(ideConfig))
			}
		}

		println()

		err = db.Push()
		if err != nil {
			fmt.Printf("unable to push to database: %s\n", err)
			panic(err)
		}

		color.RGB(200, 200, 200).Print("Pushed changes to database")

		return nil
	},
}
