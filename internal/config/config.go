package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config represents the structure of the config file
type Config struct {
	CurrentUserName string `json:"current_user_name"`
	DatabaseURL     string `json:"database_url"`
}

const configFileName = ".gatorconfig.json"

// Read reads the config file from ~/.gatorconfig.json and returns a Config struct
func Read() (Config, error) {
	configPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		// If the file doesn't exist, return a default config
		if os.IsNotExist(err) {
			return Config{
				DatabaseURL: "postgres://example",
			}, nil
		}
		return Config{}, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return Config{}, err
	}

	// Set default values if not present
	if config.DatabaseURL == "" {
		config.DatabaseURL = "postgres://example"
	}

	return config, nil
}

// SetUser sets the current_user_name field and writes the config to the file
func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	return write(*c)
}

// getConfigFilePath returns the path to the config file
func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, configFileName), nil
}

// write writes the config to the file
func write(cfg Config) error {
	configPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}
