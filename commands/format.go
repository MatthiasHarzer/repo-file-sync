package commands

import (
	"repo-file-sync/repository"

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

	if len(remotes) == 0 {
		remotesString = color.RGB(150, 150, 150).Sprintf("No remotes - no files included")
	}

	return color.GreenString("+") + " Discovered " + color.GreenString(repo) + " (" + remotesString + ")"
}
