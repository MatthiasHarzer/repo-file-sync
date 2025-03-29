package add

import (
	"github.com/MatthiasHarzer/repo-file-sync/commands"
	"github.com/MatthiasHarzer/repo-file-sync/repository"

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
	Command.Flags().BoolVarP(&isGlobalPattern, "global", "g", false, "Whether to add the include patterns as global pattern")
}

func addIncludePatterns(options *repository.DiscoveryOptions, args []string) {
	for _, arg := range args {
		options.IncludePatterns.Add(arg)
		println(color.GreenString("  +"), "Added include pattern", color.GreenString(arg))
	}
}

var Command = &cobra.Command{
	Use:   "add",
	Short: "Add a custom glob-pattern to include",
	Long:  "Add a custom glob-pattern to include",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, _, _, cfg, globalDiscoveryOptions, err := commands.Setup(baseDir)
		if err != nil {
			panic(err)
		}

		_, err = git.PlainOpen(baseDir)
		isRepo := err == nil

		if isGlobalPattern {
			println(color.YellowString("Adding global patterns:"))
			addIncludePatterns(globalDiscoveryOptions, args)

			err = db.WriteGlobalDiscoveryOptions(*globalDiscoveryOptions)
			if err != nil {
				panic(err)
			}
		} else if !isRepo {
			println(color.RedString("Includes patterns can only be added to repositories or as a global pattern. Please enter a git repository directory first or use the `--global` flag."))
			return nil
		} else {
			println(color.YellowString("Adding include patterns to repository:"))
			options, err := db.ReadRepoDiscoveryOptions(baseDir)
			if err != nil {
				panic(err)
			}

			addIncludePatterns(&options, args)

			err = db.WriteRepoDiscoveryOptions(baseDir, options)
			if err != nil {
				panic(err)
			}
		}

		return commands.Push(cfg, db)
	},
}
