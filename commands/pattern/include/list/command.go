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
	Command.Flags().BoolVarP(&isGlobalPattern, "global", "g", false, "Whether to list global pattern or repository patterns")
}

var Command = &cobra.Command{
	Use:   "list",
	Short: "List all include patterns",
	Long:  "List all include patterns",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, _, _, _, globalDiscoveryOptions, err := commands.Setup(baseDir)
		if err != nil {
			panic(err)
		}

		_, err = git.PlainOpen(baseDir)
		isRepo := err == nil

		if isGlobalPattern {
			println(color.YellowString("Global include patterns:"))
			for pattern := range globalDiscoveryOptions.IncludePatterns {
				println(color.GreenString("  +"), pattern)
			}
		} else if !isRepo {
			println(color.RedString("Include patterns can only be listed from withing repositories or from global pattern. Please enter a git repository directory first or use the `--global` flag."))
			return nil
		} else {
			println(color.YellowString("Repository include patterns:"))
			options, err := db.ReadRepoDiscoveryOptions(baseDir)
			if err != nil {
				panic(err)
			}

			for pattern := range options.IncludePatterns {
				println(color.GreenString("  +"), pattern)
			}
		}

		return nil
	},
}
