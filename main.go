package main

import (
	"repo-file-sync/commands/addpattern"
	"repo-file-sync/commands/discover"
	"repo-file-sync/commands/initialize"
	"repo-file-sync/commands/removepattern"
	"repo-file-sync/commands/restore"
	"repo-file-sync/commands/save"

	"github.com/spf13/cobra"
)

func init() {
	command.AddCommand(addpattern.Command)
	command.AddCommand(removepattern.Command)
	command.AddCommand(save.Command)
	command.AddCommand(restore.Command)
	command.AddCommand(initialize.Command)
	command.AddCommand(discover.Command)
}

var command = &cobra.Command{
	Use:   "repo-file-sync",
	Short: "Sync IDE config files",
	Long:  "Sync IDE config files",
}

func main() {
	err := command.Execute()
	if err != nil {
		panic(err)
	}
}
