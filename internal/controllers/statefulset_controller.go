package controllers

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kevholditch/vigilant/internal/theme"
	"github.com/kevholditch/vigilant/internal/views"
	"k8s.io/client-go/kubernetes"
)

// StatefulSetController manages both listing and describing StatefulSets
type StatefulSetController struct {
	clientset   *kubernetes.Clientset
	theme       *theme.Theme
	clusterName string

	// Current state
	isShowingList bool

	// Controllers
	listCtrl     *StatefulSetListController
	describeCtrl *DescribeStatefulSetController
}

// NewStatefulSetController creates a new StatefulSet controller that manages both list and describe views
func NewStatefulSetController(clientset *kubernetes.Clientset, theme *theme.Theme, clusterName string) *StatefulSetController {
	sc := &StatefulSetController{
		clientset:     clientset,
		theme:         theme,
		clusterName:   clusterName,
		isShowingList: true,
	}

	// Initialize the list controller with callback to switch to describe view
	sc.listCtrl = NewStatefulSetListController(clientset, theme, clusterName, sc.handleDescribeStatefulSet)

	return sc
}

// handleDescribeStatefulSet handles the transition to describe StatefulSet view
func (sc *StatefulSetController) handleDescribeStatefulSet(ssView *views.StatefulSetListView) tea.Cmd {
	return func() tea.Msg {
		selectedSS := ssView.GetSelected()
		if selectedSS != nil {
			sc.isShowingList = false
			sc.describeCtrl = NewDescribeStatefulSetController(
				sc.clientset,
				sc.theme,
				selectedSS.Name,
				selectedSS.Namespace,
				sc.handleBackToList,
			)
		}
		return nil
	}
}

// handleBackToList handles the transition back to StatefulSet list view
func (sc *StatefulSetController) handleBackToList() tea.Cmd {
	return func() tea.Msg {
		sc.isShowingList = true
		sc.describeCtrl = nil
		return nil
	}
}

// HandleKey handles key press events and forwards them to the active controller
func (sc *StatefulSetController) HandleKey(msg tea.KeyMsg) tea.Cmd {
	if sc.isShowingList {
		cmd := sc.listCtrl.HandleKey(msg)
		if cmd != nil {
			cmd()
		}
		return cmd
	} else if sc.describeCtrl != nil {
		return sc.describeCtrl.HandleKey(msg)
	}
	return nil
}

// Render returns the rendered view content from the active controller
func (sc *StatefulSetController) Render(width, height int) string {
	if sc.isShowingList {
		return sc.listCtrl.Render(width, height)
	} else if sc.describeCtrl != nil {
		return sc.describeCtrl.Render(width, height)
	}
	return "No view available"
}

// GetListController returns the StatefulSet list controller
func (sc *StatefulSetController) GetListController() *StatefulSetListController {
	return sc.listCtrl
}

// ActionText returns the action text from the active controller
func (sc *StatefulSetController) ActionText() string {
	if sc.isShowingList {
		return sc.listCtrl.ActionText()
	} else if sc.describeCtrl != nil {
		return sc.describeCtrl.ActionText()
	}
	return "Unknown action"
}

// GetUpdateChannel returns the update channel from the active controller
func (sc *StatefulSetController) GetUpdateChannel() <-chan tea.Msg {
	if sc.isShowingList && sc.listCtrl != nil {
		return sc.listCtrl.GetUpdateChannel()
	}
	return nil
}
