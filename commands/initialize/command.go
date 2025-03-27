package initialize

import (
	"bufio"
	"fmt"
	"net/url"
	"os"

	"repo-file-sync/config"
	"repo-file-sync/database"
	"repo-file-sync/repository"
	"repo-file-sync/util/commandutil"
	"repo-file-sync/util/fsutil"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var defaultIncludePatterns = []string{
	"**/.idea/**",
	"**/.vscode/**",
}

var defaultExcludePatterns = []string{
	"**/node_modules/**",
	"**/*venv/**",
}

func readUseLocalOnly() (bool, error) {
	return commandutil.BooleanPrompt("Do you want to use local only mode?", false)
}

func readDatabasePath() string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Where should the configurations be stored? [%s]: ", config.DefaultDatabaseRepoPath)
	scanned := scanner.Scan()
	if !scanned {
		color.Red("failed to read input")
		return ""
	}
	text := scanner.Text()
	if text == "" {
		return config.DefaultDatabaseRepoPath
	}

	return text
}

func readDatabaseRepositoryURL(c *config.Config) string {
	scanner := bufio.NewScanner(os.Stdin)
	var text string
	for {
		if c.DatabaseRepoURL != "" {
			fmt.Printf("Enter the database repository URL [%s]: ", c.DatabaseRepoURL)
		} else {
			fmt.Print("Enter the database repository URL: ")
		}
		scanned := scanner.Scan()
		if !scanned {
			color.Red("failed to read input")
			return ""
		}
		text = scanner.Text()
		_, err := url.Parse(text)
		if err != nil {
			color.Red("invalid URL: %s", err)
			continue
		}
		break
	}
	return text
}

func readGlobalDiscoveryOptions() (repository.DiscoveryOptions, error) {
	options := repository.NewDiscoveryOptions()

	color.Set(color.Bold).Println("Include patterns:")

	for _, pattern := range defaultIncludePatterns {
		fmt.Println(" ", color.GreenString(pattern))
	}
	accept, err := commandutil.BooleanPrompt("Do you want to add these default include patterns?", true)
	if err != nil {
		return repository.DiscoveryOptions{}, err
	}

	if accept {
		options.IncludePatterns.Add(defaultIncludePatterns...)
	}

	color.Set(color.Bold).Println("Exclude patterns:")
	for _, pattern := range defaultExcludePatterns {
		fmt.Println(" ", color.GreenString(pattern))
	}
	accept, err = commandutil.BooleanPrompt("Do you want to add these default exclude patterns?", true)
	if err != nil {
		return repository.DiscoveryOptions{}, err
	}

	if accept {
		options.ExcludePatterns.Add(defaultExcludePatterns...)
	}

	return options, nil
}

var Command = &cobra.Command{
	Use:   "init",
	Short: "Initialize external file sync",
	Long:  "Initialize external file sync",
	Run: func(c *cobra.Command, args []string) {
		exists, err := fsutil.Exists(config.StoragePath)
		if err != nil {
			panic(err)
		}

		if !exists {
			err = os.MkdirAll(config.StoragePath, 0755)
			if err != nil {
				color.Red("failed to create storage directory: %s", err)
			}
		}

		cfg, err := config.Load()
		if err != nil {
			color.Red("failed to load config: %s", err)
			return
		}

		cfg.DatabasePath = readDatabasePath()
		localOnly, err := readUseLocalOnly()
		if err != nil {
			color.Red("failed to read input: %s", err)
			return
		}

		cfg.LocalOnly = localOnly
		if !cfg.LocalOnly {
			cfg.DatabaseRepoURL = readDatabaseRepositoryURL(cfg)
		}

		dbPathExists, _ := fsutil.Exists(cfg.DatabasePath)
		if !dbPathExists {
			err = os.MkdirAll(cfg.DatabasePath, 0755)
			if err != nil {
				color.Red("failed to create database path: %s", err)
				return
			}
		} else {
			dbPathEmpty, _ := fsutil.IsDirectoryEmpty(cfg.DatabasePath)
			if !dbPathEmpty {
				color.Red("database path must be an empty directory")
				return
			}
		}

		var db database.Database
		if cfg.LocalOnly {
			db, err = database.InitializeRepoDatabaseFromPath(cfg.DatabasePath)
		} else {
			db, err = database.InitializeRepoDatabaseFromURL(cfg.DatabaseRepoURL, cfg.DatabasePath)
		}
		if err != nil {
			color.Red("failed to create database repository: %s", err)
			return
		}
		color.Green("Database repository created at %s", cfg.DatabasePath)

		globalDiscoveryOptions, err := readGlobalDiscoveryOptions()
		if err != nil {
			color.Red("failed to read discovery options: %s", err)
			return
		}

		err = db.WriteGlobalDiscoveryOptions(globalDiscoveryOptions)
		if err != nil {
			color.Red("failed to write global discovery options: %s", err)
			return
		}

		err = config.Save(cfg)
		if err != nil {
			color.Red("failed to save config: %s", err)
		}

	},
}
