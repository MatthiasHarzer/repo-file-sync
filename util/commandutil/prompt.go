package commandutil

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

func BooleanPrompt(prompt string, defaultValue bool) (bool, error) {
	scanner := bufio.NewScanner(os.Stdin)
	var defaultOptionString string
	if defaultValue {
		defaultOptionString = "Y/n"
	} else {
		defaultOptionString = "y/N"
	}
	fmt.Printf("%s [%s]: ", prompt, defaultOptionString)
	for {
		scanned := scanner.Scan()
		if !scanned {
			return false, fmt.Errorf("failed to read input")
		}

		text := scanner.Text()
		if text == "" {
			return false, nil
		}

		switch strings.ToLower(text) {
		case "y", "yes", "true", "1":
			return true, nil

		case "n", "no", "false", "0":
			return false, nil
		default:
			color.Red("Invalid option '%s'", text)
		}
	}
}
