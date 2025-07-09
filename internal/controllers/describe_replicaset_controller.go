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

// DescribeReplicaSetController handles input for the describe ReplicaSet view
type DescribeReplicaSetController struct {
	describeReplicaSetView *views.DescribeReplicaSetView
	onBack                 func() tea.Cmd
	clientset              *kubernetes.Clientset
	theme                  *theme.Theme
	replicasetName         string
	namespace              string
	width                  int
	height                 int
}

// NewDescribeReplicaSetController creates a new describe ReplicaSet controller
func NewDescribeReplicaSetController(clientset *kubernetes.Clientset, theme *theme.Theme, replicasetName, namespace string, onBack func() tea.Cmd) *DescribeReplicaSetController {
	replicaset, err := models.GetReplicaSet(clientset, namespace, replicasetName)
	if err != nil {
		log.Printf("error getting replicaset details: %v", err)
		replicaset = &models.ReplicaSet{
			Name:      replicasetName,
			Namespace: namespace,
			Status:    "Error",
		}
	}
	describeReplicaSetView := views.NewDescribeReplicaSetView(replicaset, theme)
	return &DescribeReplicaSetController{
		describeReplicaSetView: describeReplicaSetView,
		onBack:                 onBack,
		clientset:              clientset,
		theme:                  theme,
		replicasetName:         replicasetName,
		namespace:              namespace,
	}
}

// HandleKey handles key press events for the describe ReplicaSet view
func (c *DescribeReplicaSetController) HandleKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "up", "k":
		c.describeReplicaSetView.ScrollUp()
		return nil
	case "down", "j":
		c.describeReplicaSetView.ScrollDown()
		return nil
	case "pgup", "ctrl+u":
		c.describeReplicaSetView.ScrollPageUp()
		return nil
	case "pgdown", "ctrl+d":
		c.describeReplicaSetView.ScrollPageDown()
		return nil
	case "g":
		c.describeReplicaSetView.ScrollToTop()
		return nil
	case "G":
		c.describeReplicaSetView.ScrollToBottom()
		return nil
	case "esc":
		return c.onBack()
	case "r":
		return c.refreshReplicaSet()
	default:
		return nil
	}
}

// ActionText returns the text to describe the action the controller is performing for the header bar
func (c *DescribeReplicaSetController) ActionText() string {
	return fmt.Sprintf("Describing ReplicaSet %s", c.replicasetName)
}

// Render returns the rendered describe ReplicaSet view
func (c *DescribeReplicaSetController) Render(width, height int) string {
	c.width = width
	c.height = height
	c.describeReplicaSetView.SetSize(width, height)
	return c.describeReplicaSetView.Render()
}

// refreshReplicaSet refreshes the ReplicaSet details
func (c *DescribeReplicaSetController) refreshReplicaSet() tea.Cmd {
	return func() tea.Msg {
		replicaset, err := models.GetReplicaSet(c.clientset, c.namespace, c.replicasetName)
		if err != nil {
			log.Printf("error refreshing replicaset details: %v", err)
			return nil
		}
		c.describeReplicaSetView.UpdateReplicaSet(replicaset)
		return nil
	}
}
