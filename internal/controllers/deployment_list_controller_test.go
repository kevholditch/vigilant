package controllers

import (
	"testing"

	"github.com/kevholditch/vigilant/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestDeploymentListController(t *testing.T) {
	t.Run("should_show_no_deployments_when_cluster_is_empty", func(t *testing.T) {
		s := NewDeploymentListControllerScenario(t)
		defer s.Cleanup()
		s.Given().
			ConfigureCluster(func(builder *ClusterBuilder) {
				// No deployments
			}).
			When().
			the_deployment_list_controller_is_instantiated().
			Then().
			the_deployment_list_should_be(func(deployments []models.Deployment) {
				assert.Empty(t, deployments)
			})
	})

	t.Run("should_list_multiple_deployments_in_different_namespaces", func(t *testing.T) {
		s := NewDeploymentListControllerScenario(t)
		defer s.Cleanup()
		s.Given().
			ConfigureCluster(func(builder *ClusterBuilder) {
				builder.WithDeployment("deployment-a", "ns1").WithDeployment("deployment-b", "ns2")
			}).
			When().
			the_deployment_list_controller_is_instantiated().
			Then().
			the_deployment_list_should_be(func(deployments []models.Deployment) {
				assert.Len(t, deployments, 2)
				assert.ElementsMatch(t, []string{"deployment-a", "deployment-b"}, []string{deployments[0].Name, deployments[1].Name})
			})
	})

	t.Run("should_select_next_and_prev_deployment", func(t *testing.T) {
		s := NewDeploymentListControllerScenario(t)
		defer s.Cleanup()
		s.Given().
			ConfigureCluster(func(builder *ClusterBuilder) {
				builder.WithDeployment("deployment-a", "ns1").WithDeployment("deployment-b", "ns1")
			}).
			When().
			the_deployment_list_controller_is_instantiated().
			select_next_deployment().
			Then().
			the_selected_deployment_should_be(func(deployment *models.Deployment) {
				assert.Equal(t, "deployment-b", deployment.Name)
			}).
			and().
			select_prev_deployment().
			the_selected_deployment_should_be(func(deployment *models.Deployment) {
				assert.Equal(t, "deployment-a", deployment.Name)
			})
	})

	t.Run("should_refresh_deployments_when_new_deployment_is_added", func(t *testing.T) {
		s := NewDeploymentListControllerScenario(t)
		defer s.Cleanup()
		s.Given().
			ConfigureCluster(func(builder *ClusterBuilder) {
				builder.WithDeployment("deployment-a", "ns1")
			}).
			When().
			the_deployment_list_controller_is_instantiated().
			Then().
			the_deployment_list_should_be(func(deployments []models.Deployment) {
				assert.Len(t, deployments, 1)
				assert.Equal(t, "deployment-a", deployments[0].Name)
			}).
			and().
			ConfigureCluster(func(builder *ClusterBuilder) {
				builder.WithDeployment("deployment-b", "ns1")
			}).
			refresh_deployments().
			Then().
			the_deployment_list_should_be(func(deployments []models.Deployment) {
				assert.Len(t, deployments, 2)
				assert.ElementsMatch(t, []string{"deployment-a", "deployment-b"}, []string{deployments[0].Name, deployments[1].Name})
			})
	})

	t.Run("should_automatically_detect_new_deployments_via_watch", func(t *testing.T) {
		s := NewDeploymentListControllerScenario(t)
		defer s.Cleanup()
		s.Given().
			ConfigureCluster(func(builder *ClusterBuilder) {
				builder.WithDeployment("deployment-a", "ns1")
			}).
			the_deployment_list_controller_is_instantiated().
			Then().
			the_deployment_list_should_be(func(deployments []models.Deployment) {
				assert.Len(t, deployments, 1)
				assert.Equal(t, "deployment-a", deployments[0].Name)
			}).
			When().
			a_new_deployment_is_added_to_cluster("deployment-b", "ns1").
			Then().
			the_deployment_list_should_be(func(deployments []models.Deployment) {
				assert.Len(t, deployments, 2)
				assert.ElementsMatch(t, []string{"deployment-a", "deployment-b"}, []string{deployments[0].Name, deployments[1].Name})
			})
	})
}
