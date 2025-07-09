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

// StatefulSetListController handles input for the StatefulSet list view
type StatefulSetListController struct {
	statefulsetView       *views.StatefulSetListView
	onDescribeStatefulSet func(*views.StatefulSetListView) tea.Cmd
	clientset             *kubernetes.Clientset
	theme                 *theme.Theme
	clusterName           string
	width                 int
	height                int

	// Watch-related fields
	statefulsets    *utils.OrderedMap[models.StatefulSet]
	watchStarted    bool
	resourceVersion string
	needsUpdate     bool

	// Message channel for updates
	updateChan chan tea.Msg

	ctx    context.Context
	cancel context.CancelFunc
}

// NewStatefulSetListController creates a new StatefulSet list controller
func NewStatefulSetListController(clientset *kubernetes.Clientset, theme *theme.Theme, clusterName string, onDescribeStatefulSet func(*views.StatefulSetListView) tea.Cmd) *StatefulSetListController {
	ctx, cancel := context.WithCancel(context.Background())
	controller := &StatefulSetListController{
		onDescribeStatefulSet: onDescribeStatefulSet,
		clientset:             clientset,
		theme:                 theme,
		clusterName:           clusterName,
		statefulsets:          utils.NewOrderedMap[models.StatefulSet](),
		updateChan:            make(chan tea.Msg),
		ctx:                   ctx,
		cancel:                cancel,
	}

	controller.initializeStatefulSets()
	statefulsetView := views.NewStatefulSetListView(controller.getStatefulSetsList(), theme, clusterName)
	controller.statefulsetView = statefulsetView
	controller.startWatch()
	return controller
}

// initializeStatefulSets fetches initial StatefulSets and populates the map
func (c *StatefulSetListController) initializeStatefulSets() {
	ssList, err := c.clientset.AppsV1().StatefulSets("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return
	}
	c.statefulsets.Clear()
	for _, k8sSS := range ssList.Items {
		key := k8sSS.Namespace + "/" + k8sSS.Name
		c.statefulsets.Set(key, models.ToStatefulSetModel(k8sSS))
	}
	c.resourceVersion = ssList.ResourceVersion
}

// startWatch starts watching for StatefulSet changes
func (c *StatefulSetListController) startWatch() {
	if c.watchStarted {
		return
	}
	go func() {
		c.watchStatefulSets()
	}()
	c.watchStarted = true
}

// watchStatefulSets watches for StatefulSet changes and updates the local state
func (c *StatefulSetListController) watchStatefulSets() {
	watcher, err := c.clientset.AppsV1().StatefulSets("").Watch(c.ctx, metav1.ListOptions{
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
			ss, ok := event.Object.(*appsv1.StatefulSet)
			if !ok {
				continue
			}
			key := ss.Namespace + "/" + ss.Name
			switch event.Type {
			case watch.Added, watch.Modified:
				c.statefulsets.Set(key, models.ToStatefulSetModel(*ss))
			case watch.Deleted:
				c.statefulsets.Delete(key)
			}
			c.needsUpdate = true
			SendUpdate(c.updateChan)
		}
	}
}

// updateView updates the StatefulSet list view with current StatefulSets
func (c *StatefulSetListController) updateView() {
	ss := c.getStatefulSetsList()
	c.statefulsetView.UpdateStatefulSets(ss)
}

// getStatefulSetsList returns the current StatefulSets as a slice in consistent order
func (c *StatefulSetListController) getStatefulSetsList() []models.StatefulSet {
	return c.statefulsets.Values()
}

// HandleKey handles key press events for the StatefulSet list view
func (c *StatefulSetListController) HandleKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "up", "k":
		c.statefulsetView.SelectPrev()
		return nil
	case "down", "j":
		c.statefulsetView.SelectNext()
		return nil
	case "d":
		return c.onDescribeStatefulSet(c.statefulsetView)
	case "r":
		return c.refreshStatefulSets()
	default:
		return nil
	}
}

// ActionText returns the text to describe the action the controller is performing for the header bar
func (c *StatefulSetListController) ActionText() string {
	return "Listing StatefulSets"
}

// Render returns the rendered StatefulSet list view
func (c *StatefulSetListController) Render(width, height int) string {
	c.width = width
	c.height = height
	c.statefulsetView.SetSize(width, height)
	if c.needsUpdate {
		c.updateView()
		c.needsUpdate = false
	}
	return c.statefulsetView.Render()
}

// refreshStatefulSets refreshes the StatefulSet list by reinitializing and restarting watch
func (c *StatefulSetListController) refreshStatefulSets() tea.Cmd {
	return func() tea.Msg {
		c.statefulsets.Clear()
		c.initializeStatefulSets()
		c.needsUpdate = true
		return nil
	}
}

// GetStatefulSets returns the current StatefulSets (for testing)
func (c *StatefulSetListController) GetStatefulSets() []models.StatefulSet {
	return c.getStatefulSetsList()
}

// GetUpdateChannel returns the channel for StatefulSet update messages
func (c *StatefulSetListController) GetUpdateChannel() <-chan tea.Msg {
	return c.updateChan
}

// Stop stops the controller and cancels the context
func (c *StatefulSetListController) Stop() {
	if c.cancel != nil {
		c.cancel()
	}
}
