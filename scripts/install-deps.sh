#!/bin/bash

echo "Installing Go dependencies for Notse app..."

# Install BubbleTea and related libraries
go get github.com/charmbracelet/bubbletea@latest
go get github.com/charmbracelet/bubbles@latest
go get github.com/charmbracelet/lipgloss@latest

# Tidy up
go mod tidy

echo "Dependencies installed successfully!"
echo "Run 'make run' or 'go run cmd/notse/main.go' to start the app"
