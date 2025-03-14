package commands

import "github.com/fatih/color"

func FormatRepositoryDiscovered(repo string, remotes []string) string {
	remotesString := ""
	for i, remote := range remotes {
		if i > 0 {
			remotesString += ", "
		}
		remotesString += color.YellowString(remote)
	}
	if remotesString == "" {
		remotesString = color.YellowString("no remotes, ignoring")
	}

	return color.GreenString("+") + " Discovered " + color.GreenString(repo) + " (" + remotesString + ")"
}

func FormatFailedToReadRemotes(repo string, err error) string {
	return color.RedString("- Failed to read remotes for ") + color.RedString(repo) + ": " + color.RedString(err.Error())
}
