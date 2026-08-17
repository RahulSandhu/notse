package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/RahulSandhu/notse/internal/config"
	"github.com/RahulSandhu/notse/internal/display"
	"github.com/RahulSandhu/notse/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

// Entry point of the application
func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// Run initializes and starts the TUI application
func run() error {
	// Get config directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".config", "notse")
	configPath := filepath.Join(configDir, "config.json")

	// Load or create config
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Load theme
	theme, err := config.LoadTheme(cfg.ThemeFile)
	if err != nil {
		return fmt.Errorf("failed to load theme: %w", err)
	}
	display.SetTheme(theme)

	// Initialize storage
	store, err := storage.NewStorage(cfg.NotesFile)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	// Create TUI model
	model := display.NewModel(store)

	// Start the program
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run program: %w", err)
	}

	return nil
}
