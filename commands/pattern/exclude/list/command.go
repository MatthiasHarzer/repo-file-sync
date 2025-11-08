package list

import (
	"github.com/MatthiasHarzer/repo-file-sync/commands"
	"github.com/MatthiasHarzer/repo-file-sync/repository"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	baseDir         string
	isGlobalPattern bool
)

func init() {
	Command.Flags().StringVarP(&baseDir, "dir", "d", "", "The directory to search for repositories. Defaults to the current working directory.")
	Command.Flags().BoolVarP(&isGlobalPattern, "global", "g", false, "Whether to list global exclude pattern or repository exclude patterns")
}

var Command = &cobra.Command{
	Use:   "list",
	Short: "List all patterns to exclude",
	Long:  "List all patterns to exclude",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, usedBaseDir, _, _, globalDiscoveryOptions, err := commands.Setup(baseDir)
		if err != nil {
			panic(err)
		}

		repo, hasRepo := repository.FindRepositoryRoot(usedBaseDir)

		if isGlobalPattern {
			println(color.YellowString("Global exclude patterns:"))
			for pattern := range globalDiscoveryOptions.ExcludePatterns {
				println(color.GreenString("  +"), pattern)
			}
		} else if !hasRepo {
			println(color.RedString("Exclude patterns can only be listed from within repositories or from global patterns. Please enter a git repository directory first or use the `--global` flag."))
			return nil
		} else {
			println(color.YellowString("Repository exclude patterns:"))
			options, err := db.ReadRepoDiscoveryOptions(repo)
			if err != nil {
				panic(err)
			}

			for pattern := range options.ExcludePatterns {
				println(color.GreenString("  +"), pattern)
			}
		}

		return nil
	},
}
