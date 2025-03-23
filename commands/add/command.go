package add

import (
	"ide-config-sync/commands"
	"os"

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

		currentIncludes, err := db.ReadRepoIncludes(baseDir)
		if err != nil {
			panic(err)
		}

		newIncludes := append(currentIncludes, args...)

		err = db.WriteRepoIncludes(baseDir, newIncludes)
		if err != nil {
			panic(err)
		}

		if !cfg.LocalOnly {
			err := db.Push()
			if err != nil {
				panic(err)
			}
		}

		return nil
	},
}
