package controllers

import (
	"testing"

	"github.com/kevholditch/vigilant/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestHeaderController(t *testing.T) {
	t.Run("should_correctly_count_worker_nodes", func(t *testing.T) {
		scenario := NewHeaderControllerScenario(t)
		defer scenario.Cleanup()

		scenario.Given().
			ConfigureCluster(func(builder *ClusterBuilder) {
				builder.WithWorkerNodes(3)
			}).
			When().
			the_header_model_is_built().
			Then().
			the_header_model_should_be(func(m *models.HeaderModel) {
				assert.Equal(t, 3, m.WorkerNodes)
				assert.Equal(t, 0, m.ControlPlaneNodes)
				assert.NotEmpty(t, m.KubernetesVersion)
				assert.Empty(t, m.ClusterName)
			})
	})

	t.Run("should_correctly_count_control_plane_nodes", func(t *testing.T) {
		scenario := NewHeaderControllerScenario(t)
		defer scenario.Cleanup()

		scenario.Given().
			ConfigureCluster(func(builder *ClusterBuilder) {
				builder.WithControlPlaneNodes(2)
			}).
			When().
			the_header_model_is_built().
			Then().
			the_header_model_should_be(func(m *models.HeaderModel) {
				assert.Equal(t, 0, m.WorkerNodes)
				assert.Equal(t, 2, m.ControlPlaneNodes)
				assert.NotEmpty(t, m.KubernetesVersion)
				assert.Empty(t, m.ClusterName)
			})
	})

	t.Run("should_correctly_count_mixed_nodes", func(t *testing.T) {
		scenario := NewHeaderControllerScenario(t)
		defer scenario.Cleanup()

		scenario.Given().
			ConfigureCluster(func(builder *ClusterBuilder) {
				builder.WithControlPlaneNodes(2).WithWorkerNodes(3)
			}).
			When().
			the_header_model_is_built().
			Then().
			the_header_model_should_be(func(m *models.HeaderModel) {
				assert.Equal(t, 3, m.WorkerNodes)
				assert.Equal(t, 2, m.ControlPlaneNodes)
				assert.NotEmpty(t, m.KubernetesVersion)
				assert.Empty(t, m.ClusterName)
			})
	})

	t.Run("should_handle_legacy_master_nodes", func(t *testing.T) {
		scenario := NewHeaderControllerScenario(t)
		defer scenario.Cleanup()

		scenario.Given().
			ConfigureCluster(func(builder *ClusterBuilder) {
				builder.WithMasterNodes(1)
			}).
			When().
			the_header_model_is_built().
			Then().
			the_header_model_should_be(func(m *models.HeaderModel) {
				assert.Equal(t, 0, m.WorkerNodes)
				assert.Equal(t, 1, m.ControlPlaneNodes)
				assert.NotEmpty(t, m.KubernetesVersion)
				assert.Empty(t, m.ClusterName)
			})
	})

	t.Run("should_handle_empty_cluster", func(t *testing.T) {
		scenario := NewHeaderControllerScenario(t)
		defer scenario.Cleanup()

		scenario.Given().
			ConfigureCluster(func(builder *ClusterBuilder) {
				// No nodes added - empty cluster
			}).
			When().
			the_header_model_is_built().
			Then().
			the_header_model_should_be(func(m *models.HeaderModel) {
				assert.Equal(t, 0, m.WorkerNodes)
				assert.Equal(t, 0, m.ControlPlaneNodes)
				assert.NotEmpty(t, m.KubernetesVersion)
				assert.Empty(t, m.ClusterName)
			})
	})

	t.Run("should_handle_complex_cluster_with_fluent_builder", func(t *testing.T) {
		scenario := NewHeaderControllerScenario(t)
		defer scenario.Cleanup()

		scenario.Given().
			ConfigureCluster(func(builder *ClusterBuilder) {
				builder.WithControlPlaneNodes(3).
					WithWorkerNodes(5).
					WithMasterNodes(1)
			}).
			When().
			the_header_model_is_built().
			Then().
			the_header_model_should_be(func(m *models.HeaderModel) {
				assert.Equal(t, 4, m.ControlPlaneNodes) // 3 control plane + 1 master
				assert.Equal(t, 5, m.WorkerNodes)
				assert.NotEmpty(t, m.KubernetesVersion)
				assert.Empty(t, m.ClusterName)
			})
	})
}
