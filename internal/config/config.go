package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config holds application configuration
type Config struct {
	NotesFile   string `json:"notes_file"`
	MaxHistory  int    `json:"max_history"`
	ThemeFile   string `json:"theme_file"`
	DefaultTags []string `json:"default_tags,omitempty"`
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".config", "notse")
	
	return &Config{
		NotesFile:   filepath.Join(configDir, "notes_history.json"),
		MaxHistory:  100,
		ThemeFile:   filepath.Join(configDir, "theme.json"),
		DefaultTags: []string{},
	}
}

// Load reads config from file, creates default if not exists
func Load(configPath string) (*Config, error) {
	// If config doesn't exist, create default
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		cfg := DefaultConfig()
		if err := cfg.Save(configPath); err != nil {
			return nil, err
		}
		return cfg, nil
	}

	// Read existing config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Save writes config to file
func (c *Config) Save(configPath string) error {
	// Ensure directory exists
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}
