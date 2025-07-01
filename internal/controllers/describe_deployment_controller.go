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

// DescribeDeploymentController handles input for the describe deployment view
type DescribeDeploymentController struct {
	describeDeploymentView *views.DescribeDeploymentView
	onBack                 func() tea.Cmd
	clientset              *kubernetes.Clientset
	theme                  *theme.Theme
	deploymentName         string
	namespace              string
	width                  int
	height                 int
}

// NewDescribeDeploymentController creates a new describe deployment controller
func NewDescribeDeploymentController(clientset *kubernetes.Clientset, theme *theme.Theme, deploymentName, namespace string, onBack func() tea.Cmd) *DescribeDeploymentController {
	// Fetch deployment details
	deployment, err := models.GetDeployment(clientset, namespace, deploymentName)
	if err != nil {
		log.Printf("error getting deployment details: %v", err)
		// Create a placeholder deployment for error case
		deployment = &models.Deployment{
			Name:      deploymentName,
			Namespace: namespace,
			Status:    "Error",
		}
	}

	describeDeploymentView := views.NewDescribeDeploymentView(deployment, theme)

	return &DescribeDeploymentController{
		describeDeploymentView: describeDeploymentView,
		onBack:                 onBack,
		clientset:              clientset,
		theme:                  theme,
		deploymentName:         deploymentName,
		namespace:              namespace,
	}
}

// HandleKey handles key press events for the describe deployment view
func (c *DescribeDeploymentController) HandleKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "up", "k":
		c.describeDeploymentView.ScrollUp()
		return nil
	case "down", "j":
		c.describeDeploymentView.ScrollDown()
		return nil
	case "pgup", "ctrl+u":
		c.describeDeploymentView.ScrollPageUp()
		return nil
	case "pgdown", "ctrl+d":
		c.describeDeploymentView.ScrollPageDown()
		return nil
	case "g":
		c.describeDeploymentView.ScrollToTop()
		return nil
	case "G":
		c.describeDeploymentView.ScrollToBottom()
		return nil
	case "esc":
		return c.onBack()
	case "r":
		// Refresh deployment details
		return c.refreshDeployment()
	default:
		return nil
	}
}

// ActionText returns the text to describe the action the controller is performing for the header bar
func (c *DescribeDeploymentController) ActionText() string {
	return fmt.Sprintf("Describing deployment %s", c.deploymentName)
}

// Render returns the rendered describe deployment view
func (c *DescribeDeploymentController) Render(width, height int) string {
	c.width = width
	c.height = height
	c.describeDeploymentView.SetSize(width, height)
	return c.describeDeploymentView.Render()
}

// refreshDeployment refreshes the deployment details
func (c *DescribeDeploymentController) refreshDeployment() tea.Cmd {
	return func() tea.Msg {
		deployment, err := models.GetDeployment(c.clientset, c.namespace, c.deploymentName)
		if err != nil {
			log.Printf("error refreshing deployment details: %v", err)
			return nil
		}

		// Update the describe deployment view with new data
		c.describeDeploymentView.UpdateDeployment(deployment)
		return nil
	}
}
