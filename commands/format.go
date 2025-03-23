package commands

import (
	"path/filepath"

	"ide-config-sync/repository"

	"github.com/fatih/color"
)

func RepositoryDiscovered(repo string) string {
	remotes, err := repository.ReadRemotes(repo)
	if err != nil {
		return color.RedString("- Failed to read remotes for") + color.RedString(repo) + ":" + color.RedString(err.Error())
	}

	remotesString := ""
	for i, remote := range remotes {
		if i > 0 {
			remotesString += ", "
		}
		remotesString += color.YellowString(remote)
	}
	return color.GreenString("+") + " Discovered " + color.GreenString(repo) + " (" + remotesString + ")"
}

func FileProcessed(repo string, file repository.File, messagePrefix string) string {
	files, err := file.ListFiles()
	if err != nil {
		return color.RedString("- Failed to list files")
	}

	var output string
	for _, addedFile := range files {
		relPath, err := filepath.Rel(repo, addedFile)
		if err != nil {
			return color.RedString("- Failed to get relative path: %s", err)
		}
		if relPath == "" {
			continue
		}
		output += color.BlueString("  + ") + messagePrefix + ": " + color.BlueString(relPath) + "\n"
	}
	return output
}
