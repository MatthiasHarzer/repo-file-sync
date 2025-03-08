package push

import (
	"fmt"
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
	Command.Flags().StringVarP(&baseDir, "dir", "d", "", "The directory to search for repositories")
	Command.Flags().StringVar(&dbDir, "database-repository-dir", "", "The directory to use for the database repository")
}

var Command = &cobra.Command{
	Use:   "push",
	Short: "Push IDE config files",
	Long:  "Push IDE config files",
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
			fmt.Printf("unable to push: %s\n", err)
			panic(err)
		}

		color.RGB(200, 200, 200).Print("Pushed changes to database")

		return nil
	},
}
