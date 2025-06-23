package controllers

import (
	"testing"

	"github.com/kevholditch/vigilant/internal/models"
	"github.com/kevholditch/vigilant/internal/theme"
	"github.com/stretchr/testify/assert"
)

// HeaderControllerScenario represents a test scenario with given-when-then structure
type HeaderControllerScenario struct {
	t           *testing.T
	builder     *ClusterBuilder
	controller  *HeaderController
	headerModel *models.HeaderModel
}

// NewHeaderControllerScenario creates a new test scenario with a new cluster builder
func NewHeaderControllerScenario(t *testing.T) *HeaderControllerScenario {
	builder := NewClusterBuilder(t)
	return &HeaderControllerScenario{
		t:       t,
		builder: builder,
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

// WithClusterBuilder returns the underlying cluster builder for custom cluster configurations
func (ts *HeaderControllerScenario) WithClusterBuilder() *ClusterBuilder {
	return ts.builder
}

// ConfigureCluster configures the cluster using the builder and returns the scenario for chaining
func (ts *HeaderControllerScenario) ConfigureCluster(configFn func(*ClusterBuilder)) *HeaderControllerScenario {
	configFn(ts.builder)
	return ts
}

// the_header_controller_is_instantiated creates a new header controller
func (ts *HeaderControllerScenario) the_header_controller_is_instantiated() *HeaderControllerScenario {
	theme := theme.NewDefaultTheme()
	ts.controller = NewHeaderController(theme, ts.builder.GetClientset())
	return ts
}

// the_header_model_is_built builds the header model
func (ts *HeaderControllerScenario) the_header_model_is_built() *HeaderControllerScenario {
	ts.headerModel = buildHeaderModel(ts.builder.GetClientset())
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
	if ts.builder != nil {
		ts.builder.Cleanup()
	}
}

// the_header_model_should_be allows fluent assertions on the header model
func (ts *HeaderControllerScenario) the_header_model_should_be(assertFn func(*models.HeaderModel)) *HeaderControllerScenario {
	assertFn(ts.headerModel)
	return ts
}
