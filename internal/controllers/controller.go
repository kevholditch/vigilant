package controllers

import (
	tea "github.com/charmbracelet/bubbletea"
)

// KeyBinding represents a key binding with description
type KeyBinding struct {
	Key         string
	Description string
}

// Controller defines the interface for handling view-specific input and rendering
type Controller interface {
	// HandleKey handles key press events and returns a command
	HandleKey(msg tea.KeyMsg) tea.Cmd

	// Render returns the rendered view content
	Render(width, height int) string

	// ActionText returns the text to describe the action the controller is performing for the header bar
	ActionText() string

	// GetKeyBindings returns the key bindings for this controller
	GetKeyBindings() []KeyBinding
}

// UpdateableController extends Controller with update functionality
type UpdateableController interface {
	Controller

	// GetUpdateChannel returns a channel for receiving update messages
	GetUpdateChannel() <-chan tea.Msg
}
