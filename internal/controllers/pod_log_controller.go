package controllers

import (
	"context"
	"fmt"
	"io"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kevholditch/vigilant/internal/theme"
	"github.com/kevholditch/vigilant/internal/views"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

// LogFetcher is a function type that fetches logs for a pod
type LogFetcher func(podName, namespace string) (string, error)

// PodLogController handles input for the pod log view
type PodLogController struct {
	podLogView *views.PodLogView
	onBack     func() tea.Cmd
	logFetcher LogFetcher
	theme      *theme.Theme
	podName    string
	namespace  string
	width      int
	height     int
}

// NewKubernetesLogFetcher creates a LogFetcher that uses the Kubernetes API
func NewKubernetesLogFetcher(clientset *kubernetes.Clientset) LogFetcher {
	return func(podName, namespace string) (string, error) {
		req := clientset.CoreV1().Pods(namespace).GetLogs(podName, &corev1.PodLogOptions{})
		readCloser, err := req.Stream(context.TODO())
		if err != nil {
			return "", fmt.Errorf("error getting pod logs: %v", err)
		}
		defer readCloser.Close()

		// Read all logs
		content, err := io.ReadAll(readCloser)
		if err != nil {
			return "", fmt.Errorf("error reading pod logs: %v", err)
		}

		return string(content), nil
	}
}

// NewPodLogController creates a new pod log controller
func NewPodLogController(logFetcher LogFetcher, theme *theme.Theme, podName, namespace string, onBack func() tea.Cmd) *PodLogController {
	podLogView := views.NewPodLogView(podName, namespace, theme)

	controller := &PodLogController{
		podLogView: podLogView,
		onBack:     onBack,
		logFetcher: logFetcher,
		theme:      theme,
		podName:    podName,
		namespace:  namespace,
	}

	// Load logs initially
	controller.loadLogs()

	return controller
}

// loadLogs fetches pod logs and updates the view
func (c *PodLogController) loadLogs() {
	content, err := c.logFetcher(c.podName, c.namespace)
	if err != nil {
		c.podLogView.UpdateContent(fmt.Sprintf("Error getting pod logs: %v", err))
		return
	}

	c.podLogView.UpdateContent(content)
}

// HandleKey handles key press events for the pod log view
func (c *PodLogController) HandleKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "up", "k":
		c.podLogView.ScrollUp()
		return nil
	case "down", "j":
		c.podLogView.ScrollDown()
		return nil
	case "pgup", "b", "ctrl+u":
		c.podLogView.PageUp()
		return nil
	case "pgdown", "f", "ctrl+d":
		c.podLogView.PageDown()
		return nil
	case "g", "home":
		c.podLogView.GoToStart()
		return nil
	case "G", "end":
		c.podLogView.GoToEnd()
		return nil
	case "esc":
		return c.onBack()
	case "r":
		// Refresh pod logs
		return c.refreshLogs()
	default:
		return nil
	}
}

// ActionText returns the text to describe the action the controller is performing for the header bar
func (c *PodLogController) ActionText() string {
	return fmt.Sprintf("Viewing logs for pod %s", c.podName)
}

// Render returns the rendered pod log view
func (c *PodLogController) Render(width, height int) string {
	c.width = width
	c.height = height
	c.podLogView.SetSize(width, height)
	return c.podLogView.Render()
}

// refreshLogs refreshes the pod logs
func (c *PodLogController) refreshLogs() tea.Cmd {
	return func() tea.Msg {
		log.Printf("Refreshing logs for pod %s in namespace %s", c.podName, c.namespace)
		c.loadLogs()
		return nil
	}
}
