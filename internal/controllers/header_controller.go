package controllers

import (
	"context"

	"github.com/charmbracelet/lipgloss"
	"github.com/kevholditch/vigilant/internal/models"
	"github.com/kevholditch/vigilant/internal/theme"
	"github.com/kevholditch/vigilant/internal/views"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// HeaderController manages the header view
// It builds and holds a HeaderModel using the kubeclient
// and passes it to the HeaderView for rendering
type HeaderController struct {
	headerView  *views.HeaderView
	theme       *theme.Theme
	headerModel *models.HeaderModel
}

// NewHeaderController creates a new header controller and builds the header model
func NewHeaderController(theme *theme.Theme, clientset *kubernetes.Clientset) *HeaderController {
	headerView := views.NewHeaderView(theme)
	headerModel := buildHeaderModel(clientset)
	return &HeaderController{
		headerView:  headerView,
		theme:       theme,
		headerModel: headerModel,
	}
}

// buildHeaderModel fetches cluster info and builds a HeaderModel
func buildHeaderModel(clientset *kubernetes.Clientset) *models.HeaderModel {
	clusterName := ""
	k8sVersion := ""
	controlPlaneNodes := 0
	workerNodes := 0

	// Get cluster name from current context (if possible)
	// This is best-effort; if not available, leave blank
	// (Cluster name is not directly available from clientset)
	// You may want to pass it in if you want to guarantee it

	// Get Kubernetes version
	version, err := clientset.Discovery().ServerVersion()
	if err == nil && version != nil {
		k8sVersion = version.String()
	}

	// List all nodes and count control plane and worker nodes
	nodeList, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err == nil {
		for _, node := range nodeList.Items {
			labels := node.GetLabels()
			isControlPlane := false
			if _, ok := labels["node-role.kubernetes.io/control-plane"]; ok {
				isControlPlane = true
			}
			if _, ok := labels["node-role.kubernetes.io/master"]; ok {
				isControlPlane = true
			}
			if isControlPlane {
				controlPlaneNodes++
			} else {
				workerNodes++
			}
		}
	}

	return &models.HeaderModel{
		ClusterName:       clusterName,
		KubernetesVersion: k8sVersion,
		ControlPlaneNodes: controlPlaneNodes,
		WorkerNodes:       workerNodes,
	}
}

// Render renders the header with the given view text
func (hc *HeaderController) Render(width int, viewText string) string {
	hc.headerView.SetSize(width)
	return hc.headerView.Render(hc.headerModel, viewText)
}

// GetHeight returns the height of the header
func (hc *HeaderController) GetHeight() int {
	sampleHeader := hc.headerView.Render(hc.headerModel, "Sample")
	return lipgloss.Height(sampleHeader)
}
