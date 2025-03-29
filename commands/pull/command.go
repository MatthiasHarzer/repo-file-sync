package pull

import (
	"errors"

	"github.com/MatthiasHarzer/repo-file-sync/config"
	"github.com/MatthiasHarzer/repo-file-sync/database"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	baseDir string
)

func init() {
	Command.Flags().StringVarP(&baseDir, "dir", "d", "", "The directory to search for repositories. Defaults to the current working directory.")
}

var Command = &cobra.Command{
	Use:   "pull",
	Short: "Pulls the latest changes from the remote repositories without replacing any local files",
	Long:  "Pulls the latest changes from the remote repositories without replacing any local files",
	RunE: func(c *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			panic(err)
		}

		if cfg.LocalOnly {
			return errors.New("pull command is not available in local only mode")
		}

		db, err := database.NewRepoDatabase(cfg.DatabasePath)
		if err != nil {
			return errors.New("could not open database repo\n\ndid you run `repo-file-sync init`?")
		}

		println("Pulling from", color.GreenString(cfg.DatabaseRepoURL))

		err = db.Pull()
		if err != nil {
			panic(err)
		}

		color.RGB(200, 200, 200).Print("Pulled changes")
		return nil
	},
}
