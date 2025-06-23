package controllers

import (
	"context"
	"fmt"
	"testing"

	"github.com/kevholditch/vigilant/internal/models"
	"github.com/kevholditch/vigilant/internal/theme"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

// HeaderControllerScenario represents a test scenario with given-when-then structure
type HeaderControllerScenario struct {
	t           *testing.T
	env         *envtest.Environment
	clientset   *kubernetes.Clientset
	controller  *HeaderController
	headerModel *models.HeaderModel
}

// NewHeaderControllerScenario creates a new test scenario
func NewHeaderControllerScenario(t *testing.T) *HeaderControllerScenario {
	return &HeaderControllerScenario{
		t: t,
	}
}

// Given returns the HeaderControllerScenario for the given step
func (ts *HeaderControllerScenario) Given() *HeaderControllerScenario {
	return ts
}

// When returns the HeaderControllerScenario for the when step
func (ts *HeaderControllerScenario) When() *HeaderControllerScenario {
	return ts
}

// Then returns the HeaderControllerScenario for the then step
func (ts *HeaderControllerScenario) Then() *HeaderControllerScenario {
	return ts
}

// a_kubernetes_cluster_with_3_worker_nodes sets up a test environment with 3 worker nodes
func (ts *HeaderControllerScenario) a_kubernetes_cluster_with_3_worker_nodes() *HeaderControllerScenario {
	ts.env = &envtest.Environment{
		CRDDirectoryPaths:     []string{},
		ErrorIfCRDPathMissing: false,
	}
	cfg, err := ts.env.Start()
	require.NoError(ts.t, err)
	ts.clientset, err = kubernetes.NewForConfig(cfg)
	require.NoError(ts.t, err)
	for i := 1; i <= 3; i++ {
		workerNode := &corev1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name:   fmt.Sprintf("worker-node-%d", i),
				Labels: map[string]string{},
			},
		}
		_, err := ts.clientset.CoreV1().Nodes().Create(context.TODO(), workerNode, metav1.CreateOptions{})
		require.NoError(ts.t, err)
	}
	return ts
}

// a_kubernetes_cluster_with_2_control_plane_nodes sets up a test environment with 2 control plane nodes
func (ts *HeaderControllerScenario) a_kubernetes_cluster_with_2_control_plane_nodes() *HeaderControllerScenario {
	ts.env = &envtest.Environment{
		CRDDirectoryPaths:     []string{},
		ErrorIfCRDPathMissing: false,
	}
	cfg, err := ts.env.Start()
	require.NoError(ts.t, err)
	ts.clientset, err = kubernetes.NewForConfig(cfg)
	require.NoError(ts.t, err)
	for i := 1; i <= 2; i++ {
		controlPlaneNode := &corev1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("control-plane-node-%d", i),
				Labels: map[string]string{
					"node-role.kubernetes.io/control-plane": "",
				},
			},
		}
		_, err := ts.clientset.CoreV1().Nodes().Create(context.TODO(), controlPlaneNode, metav1.CreateOptions{})
		require.NoError(ts.t, err)
	}
	return ts
}

// a_kubernetes_cluster_with_mixed_nodes sets up a test environment with both control plane and worker nodes
func (ts *HeaderControllerScenario) a_kubernetes_cluster_with_mixed_nodes() *HeaderControllerScenario {
	ts.env = &envtest.Environment{
		CRDDirectoryPaths:     []string{},
		ErrorIfCRDPathMissing: false,
	}
	cfg, err := ts.env.Start()
	require.NoError(ts.t, err)
	ts.clientset, err = kubernetes.NewForConfig(cfg)
	require.NoError(ts.t, err)
	for i := 1; i <= 2; i++ {
		controlPlaneNode := &corev1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("control-plane-node-%d", i),
				Labels: map[string]string{
					"node-role.kubernetes.io/control-plane": "",
				},
			},
		}
		_, err := ts.clientset.CoreV1().Nodes().Create(context.TODO(), controlPlaneNode, metav1.CreateOptions{})
		require.NoError(ts.t, err)
	}
	for i := 1; i <= 3; i++ {
		workerNode := &corev1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name:   fmt.Sprintf("worker-node-%d", i),
				Labels: map[string]string{},
			},
		}
		_, err := ts.clientset.CoreV1().Nodes().Create(context.TODO(), workerNode, metav1.CreateOptions{})
		require.NoError(ts.t, err)
	}
	return ts
}

// a_kubernetes_cluster_with_legacy_master_nodes sets up a test environment with legacy master nodes
func (ts *HeaderControllerScenario) a_kubernetes_cluster_with_legacy_master_nodes() *HeaderControllerScenario {
	ts.env = &envtest.Environment{
		CRDDirectoryPaths:     []string{},
		ErrorIfCRDPathMissing: false,
	}
	cfg, err := ts.env.Start()
	require.NoError(ts.t, err)
	ts.clientset, err = kubernetes.NewForConfig(cfg)
	require.NoError(ts.t, err)
	masterNode := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "master-node",
			Labels: map[string]string{
				"node-role.kubernetes.io/master": "",
			},
		},
	}
	_, err = ts.clientset.CoreV1().Nodes().Create(context.TODO(), masterNode, metav1.CreateOptions{})
	require.NoError(ts.t, err)
	return ts
}

// an_empty_kubernetes_cluster sets up a test environment with no nodes
func (ts *HeaderControllerScenario) an_empty_kubernetes_cluster() *HeaderControllerScenario {
	ts.env = &envtest.Environment{
		CRDDirectoryPaths:     []string{},
		ErrorIfCRDPathMissing: false,
	}
	cfg, err := ts.env.Start()
	require.NoError(ts.t, err)
	ts.clientset, err = kubernetes.NewForConfig(cfg)
	require.NoError(ts.t, err)
	return ts
}

// the_header_controller_is_instantiated creates a new header controller
func (ts *HeaderControllerScenario) the_header_controller_is_instantiated() *HeaderControllerScenario {
	theme := theme.NewDefaultTheme()
	ts.controller = NewHeaderController(theme, ts.clientset)
	return ts
}

// the_header_model_is_built builds the header model
func (ts *HeaderControllerScenario) the_header_model_is_built() *HeaderControllerScenario {
	ts.headerModel = buildHeaderModel(ts.clientset)
	return ts
}

// the_header_model_worker_nodes_count_should_be asserts the worker nodes count
func (ts *HeaderControllerScenario) the_header_model_worker_nodes_count_should_be(expected int) *HeaderControllerScenario {
	assert.Equal(ts.t, expected, ts.headerModel.WorkerNodes)
	return ts
}

// the_header_model_control_plane_nodes_count_should_be asserts the control plane nodes count
func (ts *HeaderControllerScenario) the_header_model_control_plane_nodes_count_should_be(expected int) *HeaderControllerScenario {
	assert.Equal(ts.t, expected, ts.headerModel.ControlPlaneNodes)
	return ts
}

// the_header_model_kubernetes_version_should_not_be_empty asserts that the Kubernetes version is not empty
func (ts *HeaderControllerScenario) the_header_model_kubernetes_version_should_not_be_empty() *HeaderControllerScenario {
	assert.NotEmpty(ts.t, ts.headerModel.KubernetesVersion)
	return ts
}

// the_header_model_cluster_name_should_be_empty asserts that the cluster name is empty (as expected)
func (ts *HeaderControllerScenario) the_header_model_cluster_name_should_be_empty() *HeaderControllerScenario {
	assert.Empty(ts.t, ts.headerModel.ClusterName)
	return ts
}

// and returns the same HeaderControllerScenario for chaining
func (ts *HeaderControllerScenario) and() *HeaderControllerScenario {
	return ts
}

// Cleanup cleans up the test environment
func (ts *HeaderControllerScenario) Cleanup() {
	if ts.env != nil {
		ts.env.Stop()
	}
}
