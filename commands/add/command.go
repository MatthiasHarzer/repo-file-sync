package add

import (
	"os"

	"ide-config-sync/commands"

	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

var baseDir string

func init() {
	Command.Flags().StringVarP(&baseDir, "dir", "d", "", "The directory to search for repositories. Defaults to the current working directory.")
}

var Command = &cobra.Command{
	Use:   "add",
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

		if !isRepo {
			println(color.RedString("Custom ignores can only be added to repositories. Please enter a git repository directory first."))
			return nil
		}

		options, err := db.ReadRepoOptions(baseDir)
		if err != nil {
			panic(err)
		}

		for _, arg := range args {
			options.IncludePatterns.Add(arg)
			println(color.GreenString("+"), "Added include pattern", color.GreenString(arg))
		}

		err = db.WriteRepoOptions(baseDir, options)
		if err != nil {
			panic(err)
		}

		if !cfg.LocalOnly {
			println()
			err := db.Push()
			if err != nil {
				panic(err)
			}
		}

		return nil
	},
}
