package controllers

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kevholditch/vigilant/internal/theme"
	"github.com/kevholditch/vigilant/internal/views"
	"k8s.io/client-go/kubernetes"
)

// DeploymentController manages both listing and describing deployments
type DeploymentController struct {
	clientset   *kubernetes.Clientset
	theme       *theme.Theme
	clusterName string

	// Current state
	isShowingList bool

	// Controllers
	listCtrl     *DeploymentListController
	describeCtrl *DescribeDeploymentController
}

// NewDeploymentController creates a new deployment controller that manages both list and describe views
func NewDeploymentController(clientset *kubernetes.Clientset, theme *theme.Theme, clusterName string) *DeploymentController {
	dc := &DeploymentController{
		clientset:     clientset,
		theme:         theme,
		clusterName:   clusterName,
		isShowingList: true,
	}

	// Initialize the list controller with callback to switch to describe view
	dc.listCtrl = NewDeploymentListController(clientset, theme, clusterName, dc.handleDescribeDeployment)

	return dc
}

// handleDescribeDeployment handles the transition to describe deployment view
func (dc *DeploymentController) handleDescribeDeployment(deploymentView *views.DeploymentListView) tea.Cmd {
	return func() tea.Msg {
		selectedDeployment := deploymentView.GetSelected()
		if selectedDeployment != nil {
			dc.isShowingList = false
			dc.describeCtrl = NewDescribeDeploymentController(
				dc.clientset,
				dc.theme,
				selectedDeployment.Name,
				selectedDeployment.Namespace,
				dc.handleBackToList,
			)
		}
		return nil
	}
}

// handleBackToList handles the transition back to deployment list view
func (dc *DeploymentController) handleBackToList() tea.Cmd {
	return func() tea.Msg {
		dc.isShowingList = true
		dc.describeCtrl = nil
		return nil
	}
}

// HandleKey handles key press events and forwards them to the active controller
func (dc *DeploymentController) HandleKey(msg tea.KeyMsg) tea.Cmd {
	if dc.isShowingList {
		cmd := dc.listCtrl.HandleKey(msg)
		// If this is a command that might change our state, execute it immediately
		if cmd != nil {
			// Execute the command to update our internal state
			cmd()
		}
		return cmd
	} else if dc.describeCtrl != nil {
		return dc.describeCtrl.HandleKey(msg)
	}
	return nil
}

// Render returns the rendered view content from the active controller
func (dc *DeploymentController) Render(width, height int) string {
	if dc.isShowingList {
		return dc.listCtrl.Render(width, height)
	} else if dc.describeCtrl != nil {
		return dc.describeCtrl.Render(width, height)
	}
	return "No view available"
}

// GetListController returns the deployment list controller
func (dc *DeploymentController) GetListController() *DeploymentListController {
	return dc.listCtrl
}

// ActionText returns the action text from the active controller
func (dc *DeploymentController) ActionText() string {
	if dc.isShowingList {
		return dc.listCtrl.ActionText()
	} else if dc.describeCtrl != nil {
		return dc.describeCtrl.ActionText()
	}
	return "Unknown action"
}

// GetUpdateChannel returns the update channel from the active controller
func (dc *DeploymentController) GetUpdateChannel() <-chan tea.Msg {
	if dc.isShowingList && dc.listCtrl != nil {
		return dc.listCtrl.GetUpdateChannel()
	}
	// Return nil channel if no updateable controller is active
	return nil
}
