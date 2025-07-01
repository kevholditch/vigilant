package controllers

import (
	"context"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kevholditch/vigilant/internal/models"
	"github.com/kevholditch/vigilant/internal/theme"
	"github.com/kevholditch/vigilant/internal/utils"
	"github.com/kevholditch/vigilant/internal/views"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

// UpdateMsg is sent when a controller needs to trigger a re-render
type UpdateMsg struct{}

var debugLogger *log.Logger

func init() {
	// Create debug log file
	debugFile, err := os.OpenFile("vigilant-debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		// Fallback to stderr if file creation fails
		debugLogger = log.New(os.Stderr, "[DEBUG] ", log.LstdFlags)
	} else {
		debugLogger = log.New(debugFile, "[DEBUG] ", log.LstdFlags)
	}
}

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

	// Watch-related fields
	pods            *utils.OrderedMap[models.Pod] // ordered collection of pods
	watchStarted    bool
	resourceVersion string // <--- store resource version here
	needsUpdate     bool   // Flag to indicate if view needs updating

	// Message channel for updates
	updateChan chan tea.Msg

	ctx    context.Context
	cancel context.CancelFunc
}

// NewPodListController creates a new pod list controller
func NewPodListController(clientset *kubernetes.Clientset, theme *theme.Theme, clusterName string, onDescribePod func(*views.PodListView) tea.Cmd, onOpenLogs func(*views.PodListView) tea.Cmd) *PodListController {
	ctx, cancel := context.WithCancel(context.Background())
	controller := &PodListController{
		onDescribePod: onDescribePod,
		onOpenLogs:    onOpenLogs,
		clientset:     clientset,
		theme:         theme,
		clusterName:   clusterName,
		pods:          utils.NewOrderedMap[models.Pod](),
		updateChan:    make(chan tea.Msg),
		ctx:           ctx,
		cancel:        cancel,
	}

	// Initialize with initial pod list
	controller.initializePods()

	// Create the view with initial pods
	podView := views.NewPodListView(controller.getPodsList(), theme, clusterName)
	controller.podView = podView

	// Start watching for changes
	controller.startWatch()

	return controller
}

// initializePods fetches initial pods and populates the map
func (c *PodListController) initializePods() {
	podList, err := c.clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		debugLogger.Printf("error getting initial pods: %v", err)
		return
	}

	// Clear existing data
	c.pods.Clear()

	// Add pods in a consistent order (sorted by namespace, then name)
	for _, k8sPod := range podList.Items {
		key := k8sPod.Namespace + "/" + k8sPod.Name
		c.pods.Set(key, models.ToPodModel(k8sPod))
	}

	c.resourceVersion = podList.ResourceVersion // <--- store resource version
}

// startWatch starts watching for pod changes
func (c *PodListController) startWatch() {
	if c.watchStarted {
		return
	}

	go func() {
		c.watchPods()
	}()

	c.watchStarted = true
}

// watchPods watches for pod changes and updates the local state
func (c *PodListController) watchPods() {
	watcher, err := c.clientset.CoreV1().Pods("").Watch(c.ctx, metav1.ListOptions{
		ResourceVersion: c.resourceVersion,
	})
	if err != nil {
		debugLogger.Printf("error starting pod watch: %v", err)
		return
	}
	defer watcher.Stop()

	debugLogger.Printf("Started watching pods from resource version: %s", c.resourceVersion)

	for {
		select {
		case <-c.ctx.Done():
			debugLogger.Printf("Pod watch stopped by context cancellation")
			return
		case event, ok := <-watcher.ResultChan():
			if !ok {
				debugLogger.Printf("Pod watch channel closed")
				return
			}
			pod, ok := event.Object.(*corev1.Pod)
			if !ok {
				debugLogger.Printf("unexpected object type in watch event: %T", event.Object)
				continue
			}

			key := pod.Namespace + "/" + pod.Name

			switch event.Type {
			case watch.Added:
				c.pods.Set(key, models.ToPodModel(*pod))
				debugLogger.Printf("Pod added: %s", key)
			case watch.Modified:
				c.pods.Set(key, models.ToPodModel(*pod))
				debugLogger.Printf("Pod modified: %s", key)
			case watch.Deleted:
				c.pods.Delete(key)
				debugLogger.Printf("Pod deleted: %s", key)
			}
			c.needsUpdate = true

			// Send update message to trigger re-render
			SendUpdate(c.updateChan)
			debugLogger.Printf("Sent UpdateMsg for %s event", event.Type)
		}
	}
}

// updateView updates the pod list view with current pods
func (c *PodListController) updateView() {
	pods := c.getPodsList()
	debugLogger.Printf("[updateView] Controller pod map count: %d", c.pods.Len())
	if len(pods) > 0 {
		debugLogger.Printf("[updateView] First 3 pods in controller: %v", podNamesPreview(pods, 3))
	}
	c.podView.UpdatePods(pods)
	debugLogger.Printf("[updateView] View pod count after update: %d", len(c.podView.Pods()))
	if len(c.podView.Pods()) > 0 {
		debugLogger.Printf("[updateView] First 3 pods in view: %v", podNamesPreview(c.podView.Pods(), 3))
	}
}

// getPodsList returns the current pods as a slice in consistent order
func (c *PodListController) getPodsList() []models.Pod {
	return c.pods.Values()
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

	debugLogger.Printf("[Render] Called. needsUpdate: %v", c.needsUpdate)

	// Always update the view if needed, regardless of needsUpdate flag
	if c.needsUpdate {
		debugLogger.Printf("[Render] needsUpdate is true, calling updateView")
		c.updateView()
		c.needsUpdate = false
	}

	return c.podView.Render()
}

// refreshPods refreshes the pods data by reinitializing and restarting watch
func (c *PodListController) refreshPods() tea.Cmd {
	return func() tea.Msg {
		debugLogger.Printf("Refreshing pods")

		// Clear current pods
		c.pods.Clear()

		// Reinitialize pods
		c.initializePods()

		// Mark that we need to update the view
		c.needsUpdate = true

		return nil
	}
}

// GetPods returns the current pods (for testing)
func (c *PodListController) GetPods() []models.Pod {
	return c.getPodsList()
}

// SendUpdate sends an UpdateMsg through the given channel
func SendUpdate(updateChan chan<- tea.Msg) {
	select {
	case updateChan <- UpdateMsg{}:
		// Message sent successfully
	default:
		// Channel is full, skip sending message
	}
}

// GetUpdateChannel returns the channel for pod update messages
func (c *PodListController) GetUpdateChannel() <-chan tea.Msg {
	return c.updateChan
}

// Stop stops the watch goroutine
func (c *PodListController) Stop() {
	if c.cancel != nil {
		c.cancel()
	}
}

// podNamesPreview returns a preview of pod names for logging
func podNamesPreview(pods []models.Pod, max int) []string {
	var names []string
	for i, pod := range pods {
		if i >= max {
			break
		}
		names = append(names, pod.Namespace+"/"+pod.Name)
	}
	return names
}
