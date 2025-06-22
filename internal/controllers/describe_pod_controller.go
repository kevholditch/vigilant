package controllers

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kevholditch/vigilant/internal/models"
	"github.com/kevholditch/vigilant/internal/theme"
	"github.com/kevholditch/vigilant/internal/views"
	"k8s.io/client-go/kubernetes"
)

// DescribePodController handles input for the describe pod view
type DescribePodController struct {
	describePodView *views.DescribePodView
	onBack          func() tea.Cmd
	clientset       *kubernetes.Clientset
	theme           *theme.Theme
	podName         string
	namespace       string
	width           int
	height          int
}

// NewDescribePodController creates a new describe pod controller
func NewDescribePodController(clientset *kubernetes.Clientset, theme *theme.Theme, podName, namespace string, onBack func() tea.Cmd) *DescribePodController {
	// Fetch pod details
	pod, err := models.GetPod(clientset, namespace, podName)
	if err != nil {
		log.Printf("error getting pod details: %v", err)
		// Create a placeholder pod for error case
		pod = &models.Pod{
			Name:      podName,
			Namespace: namespace,
			Status:    "Error",
		}
	}

	describePodView := views.NewDescribePodView(pod, theme)

	return &DescribePodController{
		describePodView: describePodView,
		onBack:          onBack,
		clientset:       clientset,
		theme:           theme,
		podName:         podName,
		namespace:       namespace,
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
	case "r":
		// Refresh pod details
		return c.refreshPod()
	default:
		return nil
	}
}

// ActionText returns the text to describe the action the controller is performing for the header bar
func (c *DescribePodController) ActionText() string {
	return fmt.Sprintf("Describing pod %s", c.podName)
}

// Render returns the rendered describe pod view
func (c *DescribePodController) Render(width, height int) string {
	c.width = width
	c.height = height
	c.describePodView.SetSize(width, height)
	return c.describePodView.Render()
}

// refreshPod refreshes the pod details
func (c *DescribePodController) refreshPod() tea.Cmd {
	return func() tea.Msg {
		pod, err := models.GetPod(c.clientset, c.namespace, c.podName)
		if err != nil {
			log.Printf("error refreshing pod details: %v", err)
			return nil
		}

		// Update the describe pod view with new data
		c.describePodView.UpdatePod(pod)
		return nil
	}
}
