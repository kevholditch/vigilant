package controllers

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kevholditch/vigilant/internal/theme"
	"github.com/kevholditch/vigilant/internal/views"
	"k8s.io/client-go/kubernetes"
)

// ReplicaSetController manages both listing and describing ReplicaSets
type ReplicaSetController struct {
	clientset   *kubernetes.Clientset
	theme       *theme.Theme
	clusterName string

	// Current state
	isShowingList bool

	// Controllers
	listCtrl     *ReplicaSetListController
	describeCtrl *DescribeReplicaSetController
}

// NewReplicaSetController creates a new ReplicaSet controller that manages both list and describe views
func NewReplicaSetController(clientset *kubernetes.Clientset, theme *theme.Theme, clusterName string) *ReplicaSetController {
	rc := &ReplicaSetController{
		clientset:     clientset,
		theme:         theme,
		clusterName:   clusterName,
		isShowingList: true,
	}

	// Initialize the list controller with callback to switch to describe view
	rc.listCtrl = NewReplicaSetListController(clientset, theme, clusterName, rc.handleDescribeReplicaSet)

	return rc
}

// handleDescribeReplicaSet handles the transition to describe ReplicaSet view
func (rc *ReplicaSetController) handleDescribeReplicaSet(rsView *views.ReplicaSetListView) tea.Cmd {
	return func() tea.Msg {
		selectedRS := rsView.GetSelected()
		if selectedRS != nil {
			rc.isShowingList = false
			rc.describeCtrl = NewDescribeReplicaSetController(
				rc.clientset,
				rc.theme,
				selectedRS.Name,
				selectedRS.Namespace,
				rc.handleBackToList,
			)
		}
		return nil
	}
}

// handleBackToList handles the transition back to ReplicaSet list view
func (rc *ReplicaSetController) handleBackToList() tea.Cmd {
	return func() tea.Msg {
		rc.isShowingList = true
		rc.describeCtrl = nil
		return nil
	}
}

// HandleKey handles key press events and forwards them to the active controller
func (rc *ReplicaSetController) HandleKey(msg tea.KeyMsg) tea.Cmd {
	if rc.isShowingList {
		cmd := rc.listCtrl.HandleKey(msg)
		if cmd != nil {
			cmd()
		}
		return cmd
	} else if rc.describeCtrl != nil {
		return rc.describeCtrl.HandleKey(msg)
	}
	return nil
}

// Render returns the rendered view content from the active controller
func (rc *ReplicaSetController) Render(width, height int) string {
	if rc.isShowingList {
		return rc.listCtrl.Render(width, height)
	} else if rc.describeCtrl != nil {
		return rc.describeCtrl.Render(width, height)
	}
	return "No view available"
}

// GetListController returns the ReplicaSet list controller
func (rc *ReplicaSetController) GetListController() *ReplicaSetListController {
	return rc.listCtrl
}

// ActionText returns the action text from the active controller
func (rc *ReplicaSetController) ActionText() string {
	if rc.isShowingList {
		return rc.listCtrl.ActionText()
	} else if rc.describeCtrl != nil {
		return rc.describeCtrl.ActionText()
	}
	return "Unknown action"
}

// GetUpdateChannel returns the update channel from the active controller
func (rc *ReplicaSetController) GetUpdateChannel() <-chan tea.Msg {
	if rc.isShowingList && rc.listCtrl != nil {
		return rc.listCtrl.GetUpdateChannel()
	}
	return nil
}
