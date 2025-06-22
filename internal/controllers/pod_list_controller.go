package controllers

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kevholditch/vigilant/internal/views"
)

// PodListController handles input for the pod list view
type PodListController struct {
	podView       *views.PodView
	onDescribePod func(*views.PodView) tea.Cmd
}

// NewPodListController creates a new pod list controller
func NewPodListController(podView *views.PodView, onDescribePod func(*views.PodView) tea.Cmd) *PodListController {
	return &PodListController{
		podView:       podView,
		onDescribePod: onDescribePod,
	}
}

// HandleKey handles key press events for the pod list view
func (c *PodListController) HandleKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "up", "k":
		c.podView.SelectPrev()
		return nil
	case "down", "j":
		c.podView.SelectNext()
		return nil
	case "d":
		return c.onDescribePod(c.podView)
	default:
		return nil
	}
}

// GetViewType returns the view type this controller manages
func (c *PodListController) GetViewType() string {
	return "pod_list"
}
