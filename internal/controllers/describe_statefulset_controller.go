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

// DescribeStatefulSetController handles input for the describe StatefulSet view
type DescribeStatefulSetController struct {
	describeStatefulSetView *views.DescribeStatefulSetView
	onBack                  func() tea.Cmd
	clientset               *kubernetes.Clientset
	theme                   *theme.Theme
	statefulsetName         string
	namespace               string
	width                   int
	height                  int
}

// NewDescribeStatefulSetController creates a new describe StatefulSet controller
func NewDescribeStatefulSetController(clientset *kubernetes.Clientset, theme *theme.Theme, statefulsetName, namespace string, onBack func() tea.Cmd) *DescribeStatefulSetController {
	statefulset, err := models.GetStatefulSet(clientset, namespace, statefulsetName)
	if err != nil {
		log.Printf("error getting statefulset details: %v", err)
		statefulset = &models.StatefulSet{
			Name:      statefulsetName,
			Namespace: namespace,
			Status:    "Error",
		}
	}
	describeStatefulSetView := views.NewDescribeStatefulSetView(statefulset, theme)
	return &DescribeStatefulSetController{
		describeStatefulSetView: describeStatefulSetView,
		onBack:                  onBack,
		clientset:               clientset,
		theme:                   theme,
		statefulsetName:         statefulsetName,
		namespace:               namespace,
	}
}

// HandleKey handles key press events for the describe StatefulSet view
func (c *DescribeStatefulSetController) HandleKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "up", "k":
		c.describeStatefulSetView.ScrollUp()
		return nil
	case "down", "j":
		c.describeStatefulSetView.ScrollDown()
		return nil
	case "pgup", "ctrl+u":
		c.describeStatefulSetView.ScrollPageUp()
		return nil
	case "pgdown", "ctrl+d":
		c.describeStatefulSetView.ScrollPageDown()
		return nil
	case "g":
		c.describeStatefulSetView.ScrollToTop()
		return nil
	case "G":
		c.describeStatefulSetView.ScrollToBottom()
		return nil
	case "esc":
		return c.onBack()
	case "r":
		return c.refreshStatefulSet()
	default:
		return nil
	}
}

// ActionText returns the text to describe the action the controller is performing for the header bar
func (c *DescribeStatefulSetController) ActionText() string {
	return fmt.Sprintf("Describing StatefulSet %s", c.statefulsetName)
}

// Render returns the rendered describe StatefulSet view
func (c *DescribeStatefulSetController) Render(width, height int) string {
	c.width = width
	c.height = height
	c.describeStatefulSetView.SetSize(width, height)
	return c.describeStatefulSetView.Render()
}

// refreshStatefulSet refreshes the StatefulSet details
func (c *DescribeStatefulSetController) refreshStatefulSet() tea.Cmd {
	return func() tea.Msg {
		statefulset, err := models.GetStatefulSet(c.clientset, c.namespace, c.statefulsetName)
		if err != nil {
			log.Printf("error refreshing statefulset details: %v", err)
			return nil
		}
		c.describeStatefulSetView.UpdateStatefulSet(statefulset)
		return nil
	}
}
