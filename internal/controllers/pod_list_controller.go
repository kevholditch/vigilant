package controllers

import (
	"context"
	"log"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kevholditch/vigilant/internal/models"
	"github.com/kevholditch/vigilant/internal/theme"
	"github.com/kevholditch/vigilant/internal/views"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
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

	// Watch-related fields
	pods            map[string]models.Pod // key: namespace/name
	podsMutex       sync.RWMutex
	watchStarted    bool
	resourceVersion string // <--- store resource version here

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
		pods:          make(map[string]models.Pod),
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
		log.Printf("error getting initial pods: %v", err)
		return
	}

	c.podsMutex.Lock()
	defer c.podsMutex.Unlock()

	for _, k8sPod := range podList.Items {
		key := k8sPod.Namespace + "/" + k8sPod.Name
		c.pods[key] = models.ToPodModel(k8sPod)
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
		log.Printf("error starting pod watch: %v", err)
		return
	}
	defer watcher.Stop()

	log.Printf("Started watching pods from resource version: %s", c.resourceVersion)

	for {
		select {
		case <-c.ctx.Done():
			log.Printf("Pod watch stopped by context cancellation")
			return
		case event, ok := <-watcher.ResultChan():
			if !ok {
				log.Printf("Pod watch channel closed")
				return
			}
			pod, ok := event.Object.(*corev1.Pod)
			if !ok {
				log.Printf("unexpected object type in watch event: %T", event.Object)
				continue
			}

			key := pod.Namespace + "/" + pod.Name

			c.podsMutex.Lock()
			switch event.Type {
			case watch.Added:
				c.pods[key] = models.ToPodModel(*pod)
				log.Printf("Pod added: %s", key)
			case watch.Modified:
				c.pods[key] = models.ToPodModel(*pod)
				log.Printf("Pod modified: %s", key)
			case watch.Deleted:
				delete(c.pods, key)
				log.Printf("Pod deleted: %s", key)
			}
			c.podsMutex.Unlock()

			// Update the view with new pod list
			c.updateView()
		}
	}
}

// updateView updates the pod list view with current pods
func (c *PodListController) updateView() {
	c.podsMutex.RLock()
	pods := c.getPodsList()
	c.podsMutex.RUnlock()

	c.podView.UpdatePods(pods)
}

// getPodsList returns the current pods as a slice
func (c *PodListController) getPodsList() []models.Pod {
	var pods []models.Pod
	for _, pod := range c.pods {
		pods = append(pods, pod)
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

// refreshPods refreshes the pods data by reinitializing and restarting watch
func (c *PodListController) refreshPods() tea.Cmd {
	return func() tea.Msg {
		log.Printf("Refreshing pods")

		// Clear current pods
		c.podsMutex.Lock()
		c.pods = make(map[string]models.Pod)
		c.podsMutex.Unlock()

		// Reinitialize pods
		c.initializePods()

		// Update view
		c.updateView()

		return nil
	}
}

// GetPods returns the current pods (for testing)
func (c *PodListController) GetPods() []models.Pod {
	c.podsMutex.RLock()
	defer c.podsMutex.RUnlock()
	return c.getPodsList()
}

// Stop stops the watch goroutine
func (c *PodListController) Stop() {
	if c.cancel != nil {
		c.cancel()
	}
}
