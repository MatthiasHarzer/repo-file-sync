package save

import (
	"path/filepath"

	"ide-config-sync/commands"
	"ide-config-sync/repository"
	"ide-config-sync/util/fsutil"

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
		db, repos, cfg, err := commands.Setup(baseDir)
		if err != nil {
			panic(err)
		}

		_ = cfg

		globalRepoOptions, err := db.ReadGlobalOptions()
		if err != nil {
			panic(err)
		}

		for repo := range repos {
			println(color.GreenString("+"), "Discovered", color.GreenString(repo))

			localRepoOptions, err := db.ReadRepoOptions(repo)
			if err != nil {
				panic(err)
			}

			mergedOptions := globalRepoOptions.Merge(localRepoOptions)

			files := repository.DiscoverRepositoryFiles(repo, mergedOptions)

			for file := range files {
				err := db.WriteRepoFile(repo, file)
				if err != nil {
					color.Red("Failed to write %s to database: %s", file, err)
					continue
				}

				if file.IsFile() {
					println(color.BlueString("  +"), "File added:", color.BlueString(file.PathFromRepoRoot))
					continue
				}
				addedFiles, err := fsutil.ListFiles(file.AbsolutePath)
				if err != nil {
					color.Red("Failed to list files: %s", err)
					continue
				}

				for _, addedFile := range addedFiles {
					relPath, err := filepath.Rel(repo, addedFile)
					if err != nil {
						color.Red("Failed to get relative path: %s", err)
						continue
					}
					println(color.BlueString("  +"), "File added:", color.BlueString(relPath))
				}
			}
		}
		return commands.Push(cfg, db)
	},
}
