package save

import (
	"fmt"

	"github.com/MatthiasHarzer/repo-file-sync/commands"
	"github.com/MatthiasHarzer/repo-file-sync/config"
	"github.com/MatthiasHarzer/repo-file-sync/repository"
	"github.com/MatthiasHarzer/repo-file-sync/util/units"

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
	Short: "Save repository files to the database",
	Long:  "Save repository files to the database",
	RunE: func(c *cobra.Command, args []string) error {
		db, dir, repos, cfg, globalDiscoveryOptions, err := commands.Setup(baseDir)
		if err != nil {
			panic(err)
		}

		println("Discovering repositories in", color.GreenString(dir))
		for repo := range repos {
			println(commands.RepositoryDiscovered(repo))

			localDiscoveryOptions, err := db.ReadRepoDiscoveryOptions(repo)
			if err != nil {
				panic(err)
			}

			mergedOptions := globalDiscoveryOptions.Merge(localDiscoveryOptions)

			files := repository.DiscoverRepositoryFiles(repo, mergedOptions)

			var existingFiles []repository.File
			for file := range files {
				if file.Size > config.MaxFileSize {
					maxFileSizeStr := units.ConvertBytesToHumanReadable(config.MaxFileSize)
					println(color.YellowString("  -"), "Skipping file", color.YellowString(file.PathFromRepoRoot), fmt.Sprintf("as it exceeds the maximum file size of %s", maxFileSizeStr))
					continue
				}

				existingFiles = append(existingFiles, file)
				err := db.WriteRepoFile(repo, file)
				if err != nil {
					color.Red("Failed to write %s to database: %s", file, err)
					continue
				}

				println(color.BlueString("  +"), "File saved:", color.BlueString(file.PathFromRepoRoot))
			}
			err = db.RemoveNonExistingRepoFies(repo, existingFiles)
			if err != nil {
				println(color.RedString("Failed to remove non-existing files from database: %s", err))
				return err
			}
		}

		return commands.Push(cfg, db)
	},
}
