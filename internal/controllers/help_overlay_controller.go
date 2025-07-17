package controllers

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kevholditch/vigilant/internal/theme"
	"github.com/kevholditch/vigilant/internal/views"
)

// HelpOverlayController handles the help overlay functionality
type HelpOverlayController struct {
	helpOverlayView   *views.HelpOverlayView
	theme             *theme.Theme
	width             int
	height            int
	isActive          bool
	currentController Controller
}

// NewHelpOverlayController creates a new help overlay controller
func NewHelpOverlayController(theme *theme.Theme) *HelpOverlayController {
	return &HelpOverlayController{
		helpOverlayView:   views.NewHelpOverlayView(theme),
		theme:             theme,
		isActive:          false,
		currentController: nil,
	}
}

// Activate activates the help overlay
func (hoc *HelpOverlayController) Activate() {
	hoc.isActive = true
}

// Deactivate deactivates the help overlay
func (hoc *HelpOverlayController) Deactivate() {
	hoc.isActive = false
}

// IsActive returns whether the help overlay is active
func (hoc *HelpOverlayController) IsActive() bool {
	return hoc.isActive
}

// HandleKey handles key press events for the help overlay
func (hoc *HelpOverlayController) HandleKey(msg tea.KeyMsg) tea.Cmd {
	if !hoc.isActive {
		return nil
	}

	switch msg.String() {
	case "esc":
		hoc.Deactivate()
		return nil
	}
	return nil
}

// SetCurrentController sets the current controller for context-aware help
func (hoc *HelpOverlayController) SetCurrentController(controller Controller) {
	hoc.currentController = controller
}

// Render renders the help overlay
func (hoc *HelpOverlayController) Render(width, height int) string {
	if !hoc.isActive {
		return ""
	}

	var keyBindings []views.KeyBinding
	if hoc.currentController != nil {
		// Convert controller key bindings to view key bindings
		controllerBindings := hoc.currentController.GetKeyBindings()
		for _, binding := range controllerBindings {
			keyBindings = append(keyBindings, views.KeyBinding{
				Key:         binding.Key,
				Description: binding.Description,
			})
		}
	}

	return hoc.helpOverlayView.Render(width, height, keyBindings)
}
