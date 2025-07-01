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

// DeploymentListController handles input for the deployment list view
type DeploymentListController struct {
	deploymentView       *views.DeploymentListView
	onDescribeDeployment func(*views.DeploymentListView) tea.Cmd
	clientset            *kubernetes.Clientset
	theme                *theme.Theme
	clusterName          string
	width                int
	height               int

	// Watch-related fields
	deployments     *utils.OrderedMap[models.Deployment] // ordered collection of deployments
	watchStarted    bool
	resourceVersion string // store resource version here
	needsUpdate     bool   // Flag to indicate if view needs updating

	// Message channel for updates
	updateChan chan tea.Msg

	ctx    context.Context
	cancel context.CancelFunc
}

// NewDeploymentListController creates a new deployment list controller
func NewDeploymentListController(clientset *kubernetes.Clientset, theme *theme.Theme, clusterName string, onDescribeDeployment func(*views.DeploymentListView) tea.Cmd) *DeploymentListController {
	ctx, cancel := context.WithCancel(context.Background())
	controller := &DeploymentListController{
		onDescribeDeployment: onDescribeDeployment,
		clientset:            clientset,
		theme:                theme,
		clusterName:          clusterName,
		deployments:          utils.NewOrderedMap[models.Deployment](),
		updateChan:           make(chan tea.Msg),
		ctx:                  ctx,
		cancel:               cancel,
	}

	// Initialize with initial deployment list
	controller.initializeDeployments()

	// Create the view with initial deployments
	deploymentView := views.NewDeploymentListView(controller.getDeploymentsList(), theme, clusterName)
	controller.deploymentView = deploymentView

	// Start watching for changes
	controller.startWatch()

	return controller
}

// initializeDeployments fetches initial deployments and populates the map
func (c *DeploymentListController) initializeDeployments() {
	deploymentList, err := c.clientset.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		debugLogger.Printf("error getting initial deployments: %v", err)
		return
	}

	// Clear existing data
	c.deployments.Clear()

	// Add deployments in a consistent order (sorted by namespace, then name)
	for _, k8sDeployment := range deploymentList.Items {
		key := k8sDeployment.Namespace + "/" + k8sDeployment.Name
		c.deployments.Set(key, models.ToDeploymentModel(k8sDeployment))
	}

	c.resourceVersion = deploymentList.ResourceVersion
}

// startWatch starts watching for deployment changes
func (c *DeploymentListController) startWatch() {
	if c.watchStarted {
		return
	}

	go func() {
		c.watchDeployments()
	}()

	c.watchStarted = true
}

// watchDeployments watches for deployment changes and updates the local state
func (c *DeploymentListController) watchDeployments() {
	watcher, err := c.clientset.AppsV1().Deployments("").Watch(c.ctx, metav1.ListOptions{
		ResourceVersion: c.resourceVersion,
	})
	if err != nil {
		debugLogger.Printf("error starting deployment watch: %v", err)
		return
	}
	defer watcher.Stop()

	debugLogger.Printf("Started watching deployments from resource version: %s", c.resourceVersion)

	for {
		select {
		case <-c.ctx.Done():
			debugLogger.Printf("Deployment watch stopped by context cancellation")
			return
		case event, ok := <-watcher.ResultChan():
			if !ok {
				debugLogger.Printf("Deployment watch channel closed")
				return
			}
			deployment, ok := event.Object.(*appsv1.Deployment)
			if !ok {
				debugLogger.Printf("unexpected object type in watch event: %T", event.Object)
				continue
			}

			key := deployment.Namespace + "/" + deployment.Name

			switch event.Type {
			case watch.Added:
				c.deployments.Set(key, models.ToDeploymentModel(*deployment))
				debugLogger.Printf("Deployment added: %s", key)
			case watch.Modified:
				c.deployments.Set(key, models.ToDeploymentModel(*deployment))
				debugLogger.Printf("Deployment modified: %s", key)
			case watch.Deleted:
				c.deployments.Delete(key)
				debugLogger.Printf("Deployment deleted: %s", key)
			}
			c.needsUpdate = true

			// Send update message to trigger re-render
			SendUpdate(c.updateChan)
			debugLogger.Printf("Sent UpdateMsg for %s event", event.Type)
		}
	}
}

// updateView updates the deployment list view with current deployments
func (c *DeploymentListController) updateView() {
	deployments := c.getDeploymentsList()
	debugLogger.Printf("[updateView] Controller deployment map count: %d", c.deployments.Len())
	if len(deployments) > 0 {
		debugLogger.Printf("[updateView] First 3 deployments in controller: %v", deploymentNamesPreview(deployments, 3))
	}
	c.deploymentView.UpdateDeployments(deployments)
	debugLogger.Printf("[updateView] View deployment count after update: %d", len(c.deploymentView.Deployments()))
	if len(c.deploymentView.Deployments()) > 0 {
		debugLogger.Printf("[updateView] First 3 deployments in view: %v", deploymentNamesPreview(c.deploymentView.Deployments(), 3))
	}
}

// getDeploymentsList returns the current deployments as a slice in consistent order
func (c *DeploymentListController) getDeploymentsList() []models.Deployment {
	return c.deployments.Values()
}

// HandleKey handles key press events for the deployment list view
func (c *DeploymentListController) HandleKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "up", "k":
		c.deploymentView.SelectPrev()
		return nil
	case "down", "j":
		c.deploymentView.SelectNext()
		return nil
	case "d":
		return c.onDescribeDeployment(c.deploymentView)
	case "r":
		// Refresh deployments
		return c.refreshDeployments()
	default:
		return nil
	}
}

// ActionText returns the text to describe the action the controller is performing for the header bar
func (c *DeploymentListController) ActionText() string {
	return "Listing deployments"
}

// Render returns the rendered deployment list view
func (c *DeploymentListController) Render(width, height int) string {
	c.width = width
	c.height = height
	c.deploymentView.SetSize(width, height)

	// Update view if needed
	if c.needsUpdate {
		c.updateView()
		c.needsUpdate = false
	}

	return c.deploymentView.Render()
}

// refreshDeployments refreshes the deployment list
func (c *DeploymentListController) refreshDeployments() tea.Cmd {
	return func() tea.Msg {
		c.initializeDeployments()
		c.updateView()
		return nil
	}
}

// GetDeployments returns the current list of deployments
func (c *DeploymentListController) GetDeployments() []models.Deployment {
	return c.getDeploymentsList()
}

// GetUpdateChannel returns the update channel
func (c *DeploymentListController) GetUpdateChannel() <-chan tea.Msg {
	return c.updateChan
}

// Stop stops the controller and cleans up resources
func (c *DeploymentListController) Stop() {
	if c.cancel != nil {
		c.cancel()
	}
}

// deploymentNamesPreview returns a preview of deployment names for debugging
func deploymentNamesPreview(deployments []models.Deployment, max int) []string {
	var names []string
	for i, deployment := range deployments {
		if i >= max {
			break
		}
		names = append(names, deployment.Name)
	}
	return names
}
