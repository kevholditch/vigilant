package controllers

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kevholditch/vigilant/internal/theme"
	"github.com/kevholditch/vigilant/internal/views"
	"k8s.io/client-go/kubernetes"
)

// CommandBarController handles the command bar functionality
type CommandBarController struct {
	commandBarView     *views.CommandBarView
	theme              *theme.Theme
	width              int
	height             int
	clientset          *kubernetes.Clientset
	clusterName        string
	availableResources []string
	// Callback for switching views
	onSwitchView func(string) tea.Cmd
}

// NewCommandBarController creates a new command bar controller
func NewCommandBarController(clientset *kubernetes.Clientset, theme *theme.Theme, clusterName string, availableResources []string, onSwitchView func(string) tea.Cmd) *CommandBarController {
	return &CommandBarController{
		commandBarView:     views.NewCommandBarView(theme, availableResources),
		theme:              theme,
		clientset:          clientset,
		clusterName:        clusterName,
		availableResources: availableResources,
		onSwitchView:       onSwitchView,
	}
}

// HandleKey handles key press events for the command bar
func (cbc *CommandBarController) HandleKey(msg tea.KeyMsg) tea.Cmd {
	if !cbc.commandBarView.IsActive() {
		return nil
	}

	switch msg.String() {
	case "esc":
		// Deactivate command bar
		cbc.commandBarView.Deactivate()
		return nil
	case "enter":
		// Execute command
		return cbc.executeCommand()
	case "tab":
		// Tab completion
		cbc.commandBarView.TabComplete()
		return nil
	case "up":
		// Select previous suggestion
		cbc.commandBarView.SelectPrevSuggestion()
		return nil
	case "down":
		// Select next suggestion
		cbc.commandBarView.SelectNextSuggestion()
		return nil
	case "left":
		// Move cursor left
		cbc.commandBarView.MoveCursorLeft()
		return nil
	case "right":
		// Move cursor right
		cbc.commandBarView.MoveCursorRight()
		return nil
	case "home":
		// Move cursor to start
		cbc.commandBarView.MoveCursorToStart()
		return nil
	case "end":
		// Move cursor to end
		cbc.commandBarView.MoveCursorToEnd()
		return nil
	case "backspace":
		// Delete character
		cbc.commandBarView.DeleteChar()
		return nil
	default:
		// Handle regular character input
		if len(msg.Runes) > 0 {
			cbc.commandBarView.AddChar(msg.Runes[0])
		}
		return nil
	}
}

// Activate activates the command bar
func (cbc *CommandBarController) Activate() {
	cbc.commandBarView.Activate()
}

// IsActive returns whether the command bar is active
func (cbc *CommandBarController) IsActive() bool {
	return cbc.commandBarView.IsActive()
}

// ActionText returns the action text (not used for command bar)
func (cbc *CommandBarController) ActionText() string {
	return "Command Bar"
}

// Render returns the rendered command bar view
func (cbc *CommandBarController) Render(width, height int) string {
	cbc.width = width
	cbc.height = height
	cbc.commandBarView.SetSize(width, height)
	return cbc.commandBarView.Render()
}

// executeCommand executes the current command
func (cbc *CommandBarController) executeCommand() tea.Cmd {
	input := cbc.commandBarView.GetInput()

	// Deactivate command bar
	cbc.commandBarView.Deactivate()

	// Execute the view switch command
	if cbc.onSwitchView != nil {
		return cbc.onSwitchView(input)
	}

	return nil
}
