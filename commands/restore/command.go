package restore

import (
	"github.com/MatthiasHarzer/repo-file-sync/commands"

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
	Short: "Restore repository files from the database",
	Long:  "Restore repository files from the database",
	RunE: func(c *cobra.Command, args []string) error {
		db, dir, repos, _, _, err := commands.Setup(baseDir)
		if err != nil {
			panic(err)
		}

		println("Discovering repositories in", color.GreenString(dir))
		for repo := range repos {
			println(commands.RepositoryDiscovered(repo))

			files, err := db.ReadRepoFiles(repo)
			if err != nil {
				panic(err)
			}

			for file := range files {
				err := db.RestoreRepoFile(repo, file)
				if err != nil {
					color.Red("Failed to restore files: %s", err)
					continue
				}

				println(color.BlueString("  +"), "File restored:", color.BlueString(file.PathFromRepoRoot))
			}
		}

		return nil
	},
}
