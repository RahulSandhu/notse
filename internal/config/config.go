package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// Config holds application configuration
type Config struct {
	NotesFile   string   `json:"notes_file"`
	MaxHistory  int      `json:"max_history"`
	ThemeFile   string   `json:"theme_file"`
	DefaultTags []string `json:"default_tags,omitempty"`
}

// Theme holds customizable UI colors
type Theme struct {
	Accent           string `json:"accent"`
	Muted            string `json:"muted"`
	TextOnAccent     string `json:"text_on_accent"`
	Selected         string `json:"selected"`
	TitleInfo        string `json:"title_info"`
	NormalTitle      string `json:"normal_title"`
	PinIndicatorColor string `json:"pin_indicator_color"`
	HelpKey          string `json:"help_key"`
	PageActiveDot    string `json:"page_active_dot"`
}

// DefaultTheme returns the built-in color scheme
func DefaultTheme() *Theme {
	return &Theme{
		Accent:            "#b8d9ae",
		Muted:             "#6e6e70",
		TextOnAccent:      "#000000",
		Selected:          "#445c3d",
		TitleInfo:         "#b8d9ae",
		NormalTitle:       "#e0e0e0",
		PinIndicatorColor: "#b8d9ae",
		HelpKey:           "#b8d9ae",
		PageActiveDot:     "#b8d9ae",
	}
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

// expandPath replaces leading ~ and expands environment variables in a path
func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		homeDir, _ := os.UserHomeDir()
		path = filepath.Join(homeDir, path[2:])
	}
	return os.ExpandEnv(path)
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

	cfg.NotesFile = expandPath(cfg.NotesFile)
	cfg.ThemeFile = expandPath(cfg.ThemeFile)

	return &cfg, nil
}

// LoadTheme reads the theme JSON file, falling back to defaults
func LoadTheme(themePath string) (*Theme, error) {
	if _, err := os.Stat(themePath); os.IsNotExist(err) {
		return DefaultTheme(), nil
	}

	data, err := os.ReadFile(themePath)
	if err != nil {
		return nil, err
	}

	var theme Theme
	if err := json.Unmarshal(data, &theme); err != nil {
		return nil, err
	}

	// Fill in missing values with defaults
	defaults := DefaultTheme()
	if theme.Accent == "" {
		theme.Accent = defaults.Accent
	}
	if theme.Muted == "" {
		theme.Muted = defaults.Muted
	}
	if theme.TextOnAccent == "" {
		theme.TextOnAccent = defaults.TextOnAccent
	}
	if theme.Selected == "" {
		theme.Selected = defaults.Selected
	}
	if theme.TitleInfo == "" {
		theme.TitleInfo = defaults.TitleInfo
	}
	if theme.NormalTitle == "" {
		theme.NormalTitle = defaults.NormalTitle
	}
	if theme.PinIndicatorColor == "" {
		theme.PinIndicatorColor = defaults.PinIndicatorColor
	}
	if theme.HelpKey == "" {
		theme.HelpKey = defaults.HelpKey
	}
	if theme.PageActiveDot == "" {
		theme.PageActiveDot = defaults.PageActiveDot
	}

	return &theme, nil
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
