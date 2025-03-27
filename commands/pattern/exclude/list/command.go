package list

import (
	"repo-file-sync/commands"

	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
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
		db, _, _, globalDiscoveryOptions, err := commands.Setup(baseDir)
		if err != nil {
			panic(err)
		}

		_, err = git.PlainOpen(baseDir)
		isRepo := err == nil

		if isGlobalPattern {
			println(color.YellowString("Global exclude patterns:"))
			for pattern := range globalDiscoveryOptions.ExcludePatterns {
				println(color.GreenString("  +"), pattern)
			}
		} else if !isRepo {
			println(color.RedString("Exclude patterns can only be listed from within repositories or from global patterns. Please enter a git repository directory first or use the `--global` flag."))
			return nil
		} else {
			println(color.YellowString("Repository exclude patterns:"))
			options, err := db.ReadRepoDiscoveryOptions(baseDir)
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
