package restore

import (
	"path/filepath"

	"repo-file-sync/commands"
	"repo-file-sync/util/fsutil"

	"github.com/fatih/color"
	"github.com/otiai10/copy"
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
		db, repos, _, err := commands.Setup(baseDir)
		if err != nil {
			panic(err)
		}

		for repo := range repos {
			println(commands.RepositoryDiscovered(repo))

			files, err := db.ReadRepoFiles(repo)
			if err != nil {
				panic(err)
			}

			for _, file := range files {
				err := copy.Copy(file.AbsolutePath, repo)
				if err != nil {
					color.Red("Failed to restore files: %s", err)
					continue
				}

				restoredFiles, err := fsutil.ListFiles(file.AbsolutePath)
				if err != nil {
					color.Red("Failed to list files: %s", err)
					continue
				}

				for _, restoredFile := range restoredFiles {
					relPath, err := filepath.Rel(file.AbsolutePath, restoredFile)
					if err != nil {
						color.Red("Failed to get relative path: %s", err)
						continue
					}

					println(color.BlueString("  +"), "File restored:", color.BlueString(relPath))
				}
			}
		}

		return nil
	},
}
