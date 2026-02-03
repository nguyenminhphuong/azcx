package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
)

// Config stores azcx-specific configuration
type Config struct {
	PreviousSubscription string            `json:"previousSubscription,omitempty"`
	Aliases              map[string]string `json:"aliases,omitempty"`
}

// GetConfigPath returns the path to the azcx config file
// Windows: %APPDATA%\azcx\config.json
// macOS/Linux: ~/.config/azcx/config.json
func GetConfigPath() (string, error) {
	var configDir string

	if runtime.GOOS == "windows" {
		// Use %APPDATA% on Windows
		appData := os.Getenv("APPDATA")
		if appData == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			appData = filepath.Join(home, "AppData", "Roaming")
		}
		configDir = filepath.Join(appData, "azcx")
	} else {
		// Use ~/.config on macOS/Linux
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configDir = filepath.Join(home, ".config", "azcx")
	}

	return filepath.Join(configDir, "config.json"), nil
}

// Load loads the azcx config from disk
func Load() (*Config, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{Aliases: make(map[string]string)}, nil
		}
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.Aliases == nil {
		cfg.Aliases = make(map[string]string)
	}

	return &cfg, nil
}

// Save saves the azcx config to disk
func Save(cfg *Config) error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}
