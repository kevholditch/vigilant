package controllers

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kevholditch/vigilant/internal/theme"
	"github.com/kevholditch/vigilant/internal/views"
	"k8s.io/client-go/kubernetes"
)

// PodController manages both listing and describing pods
type PodController struct {
	clientset   *kubernetes.Clientset
	theme       *theme.Theme
	clusterName string

	// Current state
	isShowingList bool
	isShowingLogs bool

	// Controllers
	listCtrl     *PodListController
	describeCtrl *DescribePodController
	logCtrl      *PodLogController
}

// NewPodController creates a new pod controller that manages both list and describe views
func NewPodController(clientset *kubernetes.Clientset, theme *theme.Theme, clusterName string) *PodController {
	pc := &PodController{
		clientset:     clientset,
		theme:         theme,
		clusterName:   clusterName,
		isShowingList: true,
		isShowingLogs: false,
	}

	// Initialize the list controller with callbacks to switch to describe view and logs view
	pc.listCtrl = NewPodListController(clientset, theme, clusterName, pc.handleDescribePod, pc.handleOpenLogs)

	return pc
}

// handleDescribePod handles the transition to describe pod view
func (pc *PodController) handleDescribePod(podView *views.PodListView) tea.Cmd {
	return func() tea.Msg {
		selectedPod := podView.GetSelected()
		if selectedPod != nil {
			pc.isShowingList = false
			pc.isShowingLogs = false
			pc.describeCtrl = NewDescribePodController(
				pc.clientset,
				pc.theme,
				selectedPod.Name,
				selectedPod.Namespace,
				pc.handleBackToList,
			)
		}
		return nil
	}
}

// handleOpenLogs handles the transition to pod logs view
func (pc *PodController) handleOpenLogs(podView *views.PodListView) tea.Cmd {
	return func() tea.Msg {
		selectedPod := podView.GetSelected()
		if selectedPod != nil {
			pc.isShowingList = false
			pc.isShowingLogs = true
			logFetcher := NewKubernetesLogFetcher(pc.clientset)
			pc.logCtrl = NewPodLogController(
				logFetcher,
				pc.theme,
				selectedPod.Name,
				selectedPod.Namespace,
				pc.handleBackToList,
			)
		}
		return nil
	}
}

// handleBackToList handles the transition back to pod list view
func (pc *PodController) handleBackToList() tea.Cmd {
	return func() tea.Msg {
		pc.isShowingList = true
		pc.isShowingLogs = false
		pc.describeCtrl = nil
		pc.logCtrl = nil
		return nil
	}
}

// HandleKey handles key press events and forwards them to the active controller
func (pc *PodController) HandleKey(msg tea.KeyMsg) tea.Cmd {
	if pc.isShowingList {
		cmd := pc.listCtrl.HandleKey(msg)
		// If this is a command that might change our state, execute it immediately
		if cmd != nil {
			// Execute the command to update our internal state
			cmd()
		}
		return cmd
	} else if pc.describeCtrl != nil {
		return pc.describeCtrl.HandleKey(msg)
	} else if pc.logCtrl != nil {
		return pc.logCtrl.HandleKey(msg)
	}
	return nil
}

// Render returns the rendered view content from the active controller
func (pc *PodController) Render(width, height int) string {
	if pc.isShowingList {
		return pc.listCtrl.Render(width, height)
	} else if pc.describeCtrl != nil {
		return pc.describeCtrl.Render(width, height)
	} else if pc.logCtrl != nil {
		return pc.logCtrl.Render(width, height)
	}
	return "No view available"
}

// ActionText returns the action text from the active controller
func (pc *PodController) ActionText() string {
	if pc.isShowingList {
		return pc.listCtrl.ActionText()
	} else if pc.describeCtrl != nil {
		return pc.describeCtrl.ActionText()
	} else if pc.logCtrl != nil {
		return pc.logCtrl.ActionText()
	}
	return "Unknown action"
}
