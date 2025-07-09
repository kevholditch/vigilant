package controllers

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kevholditch/vigilant/internal/models"
	"github.com/kevholditch/vigilant/internal/theme"
	"github.com/kevholditch/vigilant/internal/utils"
	"github.com/kevholditch/vigilant/internal/views"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

// ReplicaSetListController handles input for the ReplicaSet list view
type ReplicaSetListController struct {
	replicasetView       *views.ReplicaSetListView
	onDescribeReplicaSet func(*views.ReplicaSetListView) tea.Cmd
	clientset            *kubernetes.Clientset
	theme                *theme.Theme
	clusterName          string
	width                int
	height               int

	// Watch-related fields
	replicasets     *utils.OrderedMap[models.ReplicaSet]
	watchStarted    bool
	resourceVersion string
	needsUpdate     bool

	// Message channel for updates
	updateChan chan tea.Msg

	ctx    context.Context
	cancel context.CancelFunc
}

// NewReplicaSetListController creates a new ReplicaSet list controller
func NewReplicaSetListController(clientset *kubernetes.Clientset, theme *theme.Theme, clusterName string, onDescribeReplicaSet func(*views.ReplicaSetListView) tea.Cmd) *ReplicaSetListController {
	ctx, cancel := context.WithCancel(context.Background())
	controller := &ReplicaSetListController{
		onDescribeReplicaSet: onDescribeReplicaSet,
		clientset:            clientset,
		theme:                theme,
		clusterName:          clusterName,
		replicasets:          utils.NewOrderedMap[models.ReplicaSet](),
		updateChan:           make(chan tea.Msg),
		ctx:                  ctx,
		cancel:               cancel,
	}

	controller.initializeReplicaSets()
	replicasetView := views.NewReplicaSetListView(controller.getReplicaSetsList(), theme, clusterName)
	controller.replicasetView = replicasetView
	controller.startWatch()
	return controller
}

// initializeReplicaSets fetches initial ReplicaSets and populates the map
func (c *ReplicaSetListController) initializeReplicaSets() {
	rsList, err := c.clientset.AppsV1().ReplicaSets("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return
	}
	c.replicasets.Clear()
	for _, k8sRS := range rsList.Items {
		key := k8sRS.Namespace + "/" + k8sRS.Name
		c.replicasets.Set(key, models.ToReplicaSetModel(k8sRS))
	}
	c.resourceVersion = rsList.ResourceVersion
}

// startWatch starts watching for ReplicaSet changes
func (c *ReplicaSetListController) startWatch() {
	if c.watchStarted {
		return
	}
	go func() {
		c.watchReplicaSets()
	}()
	c.watchStarted = true
}

// watchReplicaSets watches for ReplicaSet changes and updates the local state
func (c *ReplicaSetListController) watchReplicaSets() {
	watcher, err := c.clientset.AppsV1().ReplicaSets("").Watch(c.ctx, metav1.ListOptions{
		ResourceVersion: c.resourceVersion,
	})
	if err != nil {
		return
	}
	defer watcher.Stop()
	for {
		select {
		case <-c.ctx.Done():
			return
		case event, ok := <-watcher.ResultChan():
			if !ok {
				return
			}
			rs, ok := event.Object.(*appsv1.ReplicaSet)
			if !ok {
				continue
			}
			key := rs.Namespace + "/" + rs.Name
			switch event.Type {
			case watch.Added, watch.Modified:
				c.replicasets.Set(key, models.ToReplicaSetModel(*rs))
			case watch.Deleted:
				c.replicasets.Delete(key)
			}
			c.needsUpdate = true
			SendUpdate(c.updateChan)
		}
	}
}

// updateView updates the ReplicaSet list view with current ReplicaSets
func (c *ReplicaSetListController) updateView() {
	rs := c.getReplicaSetsList()
	c.replicasetView.UpdateReplicaSets(rs)
}

// getReplicaSetsList returns the current ReplicaSets as a slice in consistent order
func (c *ReplicaSetListController) getReplicaSetsList() []models.ReplicaSet {
	return c.replicasets.Values()
}

// HandleKey handles key press events for the ReplicaSet list view
func (c *ReplicaSetListController) HandleKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "up", "k":
		c.replicasetView.SelectPrev()
		return nil
	case "down", "j":
		c.replicasetView.SelectNext()
		return nil
	case "d":
		return c.onDescribeReplicaSet(c.replicasetView)
	case "r":
		return c.refreshReplicaSets()
	default:
		return nil
	}
}

// ActionText returns the text to describe the action the controller is performing for the header bar
func (c *ReplicaSetListController) ActionText() string {
	return "Listing ReplicaSets"
}

// Render returns the rendered ReplicaSet list view
func (c *ReplicaSetListController) Render(width, height int) string {
	c.width = width
	c.height = height
	c.replicasetView.SetSize(width, height)
	if c.needsUpdate {
		c.updateView()
		c.needsUpdate = false
	}
	return c.replicasetView.Render()
}

// refreshReplicaSets refreshes the ReplicaSet list by reinitializing and restarting watch
func (c *ReplicaSetListController) refreshReplicaSets() tea.Cmd {
	return func() tea.Msg {
		c.replicasets.Clear()
		c.initializeReplicaSets()
		c.needsUpdate = true
		return nil
	}
}

// GetReplicaSets returns the current ReplicaSets (for testing)
func (c *ReplicaSetListController) GetReplicaSets() []models.ReplicaSet {
	return c.getReplicaSetsList()
}

// GetUpdateChannel returns the channel for ReplicaSet update messages
func (c *ReplicaSetListController) GetUpdateChannel() <-chan tea.Msg {
	return c.updateChan
}

// Stop stops the controller and cancels the context
func (c *ReplicaSetListController) Stop() {
	if c.cancel != nil {
		c.cancel()
	}
}
