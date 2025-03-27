package addpattern

import (
	"os"

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
	Command.Flags().BoolVarP(&isGlobalPattern, "global", "g", false, "Whether to add the patterns as global pattern")
}

func addPatterns(options *repository.DiscoveryOptions, args []string) {
	for _, arg := range args {
		options.IncludePatterns.Add(arg)
		println(color.GreenString("  +"), "Added include pattern", color.GreenString(arg))
	}
}

var Command = &cobra.Command{
	Use:   "addpattern",
	Short: "Adds a custom pattern to include",
	Long:  "Adds a custom pattern to include",
	RunE: func(c *cobra.Command, args []string) error {
		var err error
		if baseDir == "" {
			baseDir, err = os.Getwd()
			if err != nil {
				panic(err)
			}
		}

		db, _, cfg, err := commands.Setup(baseDir)
		if err != nil {
			panic(err)
		}

		_, err = git.PlainOpen(baseDir)
		isRepo := err == nil

		if isGlobalPattern {
			println(color.YellowString("Adding global patterns:"))
			options, err := db.ReadGlobalDiscoveryOptions()
			if err != nil {
				panic(err)
			}
			addPatterns(&options, args)

			err = db.WriteGlobalDiscoveryOptions(options)
			if err != nil {
				panic(err)
			}
		} else if !isRepo {
			println(color.RedString("Custom ignores can only be added to repositories or as a global pattern. Please enter a git repository directory first or use the `--global` flag."))
			return nil
		} else {
			println(color.YellowString("Adding patterns to repository:"))
			options, err := db.ReadRepoDiscoveryOptions(baseDir)
			if err != nil {
				panic(err)
			}

			addPatterns(&options, args)

			err = db.WriteRepoDiscoveryOptions(baseDir, options)
			if err != nil {
				panic(err)
			}
		}

		return commands.Push(cfg, db)
	},
}
