package include

import (
	"errors"

	"repo-file-sync/commands/pattern/include/add"
	"repo-file-sync/commands/pattern/include/list"
	"repo-file-sync/commands/pattern/include/remove"

	"github.com/spf13/cobra"
)

func init() {
	Command.AddCommand(add.Command)
	Command.AddCommand(remove.Command)
	Command.AddCommand(list.Command)
}

var Command = &cobra.Command{
	Use:   "include",
	Short: "Add, remove or list include patterns",
	Long:  "Add, remove or list include patterns",
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("use the `add`, `remove` or `list` subcommands")
	},
}
