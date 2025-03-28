package exclude

import (
	"errors"

	"github.com/MatthiasHarzer/repo-file-sync/commands/pattern/exclude/add"
	"github.com/MatthiasHarzer/repo-file-sync/commands/pattern/exclude/list"
	"github.com/MatthiasHarzer/repo-file-sync/commands/pattern/exclude/remove"

	"github.com/spf13/cobra"
)

func init() {
	Command.AddCommand(add.Command)
	Command.AddCommand(remove.Command)
	Command.AddCommand(list.Command)
}

var Command = &cobra.Command{
	Use:   "exclude",
	Short: "Add, remove or list exclude patterns",
	Long:  "Add, remove or list exclude patterns",
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("use the `add`, `remove` or `list` subcommands")
	},
}
