package controllers

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Controller defines the interface for handling view-specific input
type Controller interface {
	// HandleKey handles key press events and returns a command
	HandleKey(msg tea.KeyMsg) tea.Cmd

	// GetViewType returns the view type this controller manages
	GetViewType() string
}
