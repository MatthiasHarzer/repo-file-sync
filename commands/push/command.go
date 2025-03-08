package push

import (
	"fmt"
	"private-ide-config-sync/commands"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	baseDir string
	dbDir   string
)

func init() {
	Command.Flags().StringVarP(&baseDir, "dir", "d", "", "The directory to search for repositories")
	Command.Flags().StringVar(&dbDir, "database-repository-dir", "", "The directory to use for the database repository")
}

var Command = &cobra.Command{
	Use:   "push",
	Short: "Push IDE config files",
	Long:  "Push IDE config files",
	RunE: func(c *cobra.Command, args []string) error {
		db, repos, err := commands.Setup(baseDir, dbDir)
		if err != nil {
			panic(err)
		}

		for _, repo := range repos {
			err = repo.Save()
			if err != nil {
				panic(err)
			}

			println(color.BlueString("+"), color.WhiteString(" extracted "), color.BlueString(repo.LocalRepoDirectory))
		}

		err = db.Push()
		if err != nil {
			fmt.Printf("unable to push: %s\n", err)
			panic(err)
		}

		println()
		color.RGB(200, 200, 200).Print("Pushed changes to database")

		return nil
	},
}
