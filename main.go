package main

import (
	"ide-config-sync/commands/add"
	"ide-config-sync/commands/discover"
	"ide-config-sync/commands/initialize"
	"ide-config-sync/commands/restore"
	"ide-config-sync/commands/save"

	"github.com/spf13/cobra"
)

func init() {
	command.AddCommand(add.Command)
	command.AddCommand(save.Command)
	command.AddCommand(restore.Command)
	command.AddCommand(initialize.Command)
	command.AddCommand(discover.Command)
}

var command = &cobra.Command{
	Use:   "ide-config-sync",
	Short: "Sync IDE config files",
	Long:  "Sync IDE config files",
}

func main() {
	err := command.Execute()
	if err != nil {
		panic(err)
	}
}
