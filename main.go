package main

import (
	"private-ide-config-sync/commands/initialize"
	"private-ide-config-sync/commands/pull"
	"private-ide-config-sync/commands/push"

	"github.com/spf13/cobra"
)

func init() {
	command.AddCommand(push.Command)
	command.AddCommand(pull.Command)
	command.AddCommand(initialize.Command)
}

var command = &cobra.Command{
	Use:   "private-ide-config-sync",
	Short: "Sync IDE config files",
	Long:  "Sync IDE config files",
}

func main() {
	err := command.Execute()
	if err != nil {
		panic(err)
	}
}
