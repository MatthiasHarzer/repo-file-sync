package config

import (
	"encoding/json"
	"errors"
	"ide-config-sync/util/fsutil"
	"os"
)

type Config struct {
	DatabasePath    string `json:"databasePath"`
	LocalOnly       bool   `json:"localOnly"`
	DatabaseRepoURL string `json:"databaseRepoURL"`
}

func DefaultConfig() *Config {
	return &Config{
		DatabasePath: DefaultDatabaseRepoPath,
		LocalOnly:    false,
	}
}

func (c *Config) Validate() error {
	if c.DatabasePath == "" {
		return errors.New("database path can not be empty")
	}

	if !c.LocalOnly {
		if c.DatabaseRepoURL == "" {
			return errors.New("database repo URL can not be empty if local only is disabled")
		}
	}

	return nil
}

func Save(config *Config) error {
	bytes, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		return err
	}

	err = os.WriteFile(File, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func Load() (*Config, error) {
	exists, _ := fsutil.Exists(File)
	if !exists {
		return DefaultConfig(), nil
	}

	bytes, err := os.ReadFile(File)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = json.Unmarshal(bytes, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
