package controllers

import (
	"context"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kevholditch/vigilant/internal/models"
	"github.com/kevholditch/vigilant/internal/theme"
	"github.com/kevholditch/vigilant/internal/views"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// PodListController handles input for the pod list view
type PodListController struct {
	podView       *views.PodListView
	onDescribePod func(*views.PodListView) tea.Cmd
	onOpenLogs    func(*views.PodListView) tea.Cmd
	clientset     *kubernetes.Clientset
	theme         *theme.Theme
	clusterName   string
	width         int
	height        int
}

// NewPodListController creates a new pod list controller
func NewPodListController(clientset *kubernetes.Clientset, theme *theme.Theme, clusterName string, onDescribePod func(*views.PodListView) tea.Cmd, onOpenLogs func(*views.PodListView) tea.Cmd) *PodListController {
	// Fetch pods data (logic moved from models)
	pods := fetchPods(clientset)

	podView := views.NewPodListView(pods, theme, clusterName)

	return &PodListController{
		podView:       podView,
		onDescribePod: onDescribePod,
		onOpenLogs:    onOpenLogs,
		clientset:     clientset,
		theme:         theme,
		clusterName:   clusterName,
	}
}

// fetchPods fetches pods from the Kubernetes cluster
func fetchPods(clientset *kubernetes.Clientset) []models.Pod {
	podList, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	var pods []models.Pod
	if err != nil {
		log.Printf("error getting pods: %v", err)
		return []models.Pod{} // Use empty slice on error
	}
	for _, k8sPod := range podList.Items {
		pods = append(pods, models.ToPodModel(k8sPod))
	}
	return pods
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
	case "l":
		return c.onOpenLogs(c.podView)
	case "r":
		// Refresh pods data
		return c.refreshPods()
	default:
		return nil
	}
}

// ActionText returns the text to describe the action the controller is performing for the header bar
func (c *PodListController) ActionText() string {
	return "Viewing pods"
}

// Render returns the rendered pod list view
func (c *PodListController) Render(width, height int) string {
	c.width = width
	c.height = height
	c.podView.SetSize(width, height)
	return c.podView.Render()
}

// refreshPods refreshes the pods data
func (c *PodListController) refreshPods() tea.Cmd {
	return func() tea.Msg {
		pods := fetchPods(c.clientset)
		// Update the pod view with new data
		c.podView.UpdatePods(pods)
		return nil
	}
}
