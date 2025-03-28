package main

import (
	"fmt"

	"repo-file-sync/commands/discover"
	"repo-file-sync/commands/initialize"
	"repo-file-sync/commands/pattern"
	"repo-file-sync/commands/pull"
	"repo-file-sync/commands/restore"
	"repo-file-sync/commands/save"

	"github.com/spf13/cobra"
)

var version = "unknown"

func init() {
	command.AddCommand(pattern.Command)
	command.AddCommand(save.Command)
	command.AddCommand(restore.Command)
	command.AddCommand(initialize.Command)
	command.AddCommand(discover.Command)
	command.AddCommand(pull.Command)
}

var command = &cobra.Command{
	Use:   "repo-file-sync",
	Short: "Sync any repo file",
	Long:  "Sync any repo file",
	Run: func(c *cobra.Command, args []string) {
		fmt.Println("repo-file-sync version", version)
	},
}

func main() {
	err := command.Execute()
	if err != nil {
		panic(err)
	}
}
