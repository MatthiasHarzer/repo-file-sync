package list

import (
	"os"

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
	Short: "Add a custom glob-pattern to include",
	Long:  "Add a custom glob-pattern to include",
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		if baseDir == "" {
			baseDir, err = os.Getwd()
			if err != nil {
				panic(err)
			}
		}

		db, _, _, globalDiscoveryOptions, err := commands.Setup(baseDir)
		if err != nil {
			panic(err)
		}

		_, err = git.PlainOpen(baseDir)
		isRepo := err == nil

		if isGlobalPattern {
			println(color.YellowString("Global patterns:"))
			for pattern := range globalDiscoveryOptions.IncludePatterns {
				println(color.GreenString("  +"), pattern)
			}
		} else if !isRepo {
			println(color.RedString("Ignores can only be listed from repositories or from global pattern. Please enter a git repository directory first or use the `--global` flag."))
			return nil
		} else {
			println(color.YellowString("Repository patterns:"))
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
