package patterns

import (
	"errors"

	"repo-file-sync/commands/patterns/include"

	"github.com/spf13/cobra"
)

func init() {
	Command.AddCommand(include.Command)
}

var Command = &cobra.Command{
	Use:   "patterns",
	Short: "Add or remove include or exclude patterns",
	Long:  "Add or remove include or exclude patterns",
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("use the `include` or `exclude` subcommands")
	},
}
