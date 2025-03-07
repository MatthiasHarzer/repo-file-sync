package pull

import (
	"fmt"
	"private-ide-config-sync/commands"

	"github.com/spf13/cobra"
)

var (
	baseDir string
	dbDir   string
)

func init() {
	Command.Flags().StringVarP(&baseDir, "dir", "d", "", "Base directory")
	Command.Flags().StringVar(&dbDir, "db-dir", "", "DB directory")
}

var Command = &cobra.Command{
	Use:   "pull",
	Short: "Pull IDE config files",
	Long:  "Pull IDE config files",
	RunE: func(c *cobra.Command, args []string) error {
		_, repos, err := commands.Init(baseDir, dbDir)
		if err != nil {
			panic(err)
		}

		for _, repo := range repos {
			err = repo.Pull()
			if err != nil {
				panic(err)
			}

			fmt.Printf("Pulled configs to %s\n", repo.LocalRepoDirectory)
		}

		return nil
	},
}
