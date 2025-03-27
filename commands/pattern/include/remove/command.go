package remove

import (
	"repo-file-sync/commands"
	"repo-file-sync/repository"

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
	Command.Flags().BoolVarP(&isGlobalPattern, "global", "g", false, "Whether to remove the patterns from global patterns")
}

func removeIncludePatterns(options *repository.DiscoveryOptions, args []string) {
	for _, arg := range args {
		options.IncludePatterns.Remove(arg)
		println(color.GreenString("+"), "Removed include pattern", color.GreenString(arg))
	}
}

var Command = &cobra.Command{
	Use:   "remove",
	Short: "Remove a custom glob-pattern from includes",
	Long:  "Remove a custom glob-pattern from includes",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, _, cfg, globalDiscoveryOptions, err := commands.Setup(baseDir)
		if err != nil {
			panic(err)
		}

		_, err = git.PlainOpen(baseDir)
		isRepo := err == nil

		if isGlobalPattern {
			println(color.YellowString("Removing global include patterns:"))
			removeIncludePatterns(globalDiscoveryOptions, args)

			err = db.WriteGlobalDiscoveryOptions(*globalDiscoveryOptions)
			if err != nil {
				panic(err)
			}
		} else if !isRepo {
			println(color.RedString("Includes can only be removed from within repositories or from global pattern. Please enter a git repository directory first or use the `--global` flag."))
			return nil
		} else {
			println(color.YellowString("Removing include patterns to repository:"))
			options, err := db.ReadRepoDiscoveryOptions(baseDir)
			if err != nil {
				panic(err)
			}

			removeIncludePatterns(&options, args)

			err = db.WriteRepoDiscoveryOptions(baseDir, options)
			if err != nil {
				panic(err)
			}
		}

		return commands.Push(cfg, db)
	},
}
