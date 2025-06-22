package controllers

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kevholditch/vigilant/internal/views"
)

// DescribePodController handles input for the describe pod view
type DescribePodController struct {
	describePodView *views.DescribePodView
	onBack          func() tea.Cmd
}

// NewDescribePodController creates a new describe pod controller
func NewDescribePodController(describePodView *views.DescribePodView, onBack func() tea.Cmd) *DescribePodController {
	return &DescribePodController{
		describePodView: describePodView,
		onBack:          onBack,
	}
}

// HandleKey handles key press events for the describe pod view
func (c *DescribePodController) HandleKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "up", "k":
		c.describePodView.ScrollUp()
		return nil
	case "down", "j":
		c.describePodView.ScrollDown()
		return nil
	case "esc":
		return c.onBack()
	default:
		return nil
	}
}

// GetViewType returns the view type this controller manages
func (c *DescribePodController) GetViewType() string {
	return "describe_pod"
}
