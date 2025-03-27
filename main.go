package main

import (
	"repo-file-sync/commands/discover"
	"repo-file-sync/commands/initialize"
	"repo-file-sync/commands/patterns"
	"repo-file-sync/commands/pull"
	"repo-file-sync/commands/removepattern"
	"repo-file-sync/commands/restore"
	"repo-file-sync/commands/save"

	"github.com/spf13/cobra"
)

func init() {
	command.AddCommand(patterns.Command)
	command.AddCommand(removepattern.Command)
	command.AddCommand(save.Command)
	command.AddCommand(restore.Command)
	command.AddCommand(initialize.Command)
	command.AddCommand(discover.Command)
	command.AddCommand(pull.Command)
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
