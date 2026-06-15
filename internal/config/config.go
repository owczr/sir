package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	GitHubHost     string   `json:"github_host"`
	GitHubUsername string   `json:"github_username"`
	GitHubToken    string   `json:"github_token"`
	Repositories   []string `json:"repositories"`
}

func getConfigPath() string {
	return filepath.Join(os.Getenv("HOME"), ".config", "sir", "config.json")
}

func Load() (*Config, error) {
	configPath := getConfigPath()
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{
				GitHubHost:   "github.com",
				Repositories: []string{},
			}, nil
		}
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Set default host if not specified
	if cfg.GitHubHost == "" {
		cfg.GitHubHost = "github.com"
	}

	return &cfg, nil
}

func Save(cfg *Config) error {
	configPath := getConfigPath()
	
	// Ensure directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
