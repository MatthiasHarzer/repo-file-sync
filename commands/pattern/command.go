package pattern

import (
	"errors"

	"github.com/MatthiasHarzer/repo-file-sync/commands/pattern/exclude"
	"github.com/MatthiasHarzer/repo-file-sync/commands/pattern/include"

	"github.com/spf13/cobra"
)

func init() {
	Command.AddCommand(include.Command)
	Command.AddCommand(exclude.Command)
}

var Command = &cobra.Command{
	Use:   "pattern",
	Short: "Add or remove include or exclude patterns",
	Long:  "Add or remove include or exclude patterns",
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("use the `include` or `exclude` subcommands")
	},
}
