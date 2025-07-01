package controllers

import (
	"testing"
	"time"

	"github.com/kevholditch/vigilant/internal/models"
	"github.com/kevholditch/vigilant/internal/theme"
)

type DeploymentListControllerScenario struct {
	t              *testing.T
	builder        *ClusterBuilder
	controller     *DeploymentListController
	deploymentView *models.Deployment
}

func NewDeploymentListControllerScenario(t *testing.T) *DeploymentListControllerScenario {
	builder := NewClusterBuilder(t)
	return &DeploymentListControllerScenario{
		t:       t,
		builder: builder,
	}
}

func (s *DeploymentListControllerScenario) Given() *DeploymentListControllerScenario { return s }
func (s *DeploymentListControllerScenario) When() *DeploymentListControllerScenario  { return s }
func (s *DeploymentListControllerScenario) Then() *DeploymentListControllerScenario  { return s }
func (s *DeploymentListControllerScenario) and() *DeploymentListControllerScenario   { return s }

func (s *DeploymentListControllerScenario) ConfigureCluster(configFn func(*ClusterBuilder)) *DeploymentListControllerScenario {
	configFn(s.builder)
	return s
}

func (s *DeploymentListControllerScenario) the_deployment_list_controller_is_instantiated() *DeploymentListControllerScenario {
	theme := theme.NewDefaultTheme()
	s.controller = NewDeploymentListController(s.builder.GetClientset(), theme, "test-cluster", nil)
	return s
}

func (s *DeploymentListControllerScenario) the_deployment_list_view_is_built() *DeploymentListControllerScenario {
	// No-op for now, as controller instantiation builds the view
	return s
}

func (s *DeploymentListControllerScenario) select_next_deployment() *DeploymentListControllerScenario {
	s.controller.deploymentView.SelectNext()
	return s
}

func (s *DeploymentListControllerScenario) select_prev_deployment() *DeploymentListControllerScenario {
	s.controller.deploymentView.SelectPrev()
	return s
}

func (s *DeploymentListControllerScenario) the_deployment_list_should_be(assertFn func([]models.Deployment)) *DeploymentListControllerScenario {
	deployments := s.controller.GetDeployments()
	assertFn(deployments)
	return s
}

func (s *DeploymentListControllerScenario) the_selected_deployment_should_be(assertFn func(*models.Deployment)) *DeploymentListControllerScenario {
	assertFn(s.controller.deploymentView.GetSelected())
	return s
}

func (s *DeploymentListControllerScenario) refresh_deployments() *DeploymentListControllerScenario {
	cmd := s.controller.refreshDeployments()
	if cmd != nil {
		cmd() // simulate running the command
	}
	return s
}

func (s *DeploymentListControllerScenario) a_new_deployment_is_added_to_cluster(name, namespace string) *DeploymentListControllerScenario {
	// Add the deployment to the cluster
	s.builder.WithDeployment(name, namespace)

	// Wait for the watch to detect the change by polling until the deployment appears
	// This is deterministic because we're waiting for a specific condition
	maxAttempts := 10
	for i := 0; i < maxAttempts; i++ {
		deployments := s.controller.GetDeployments()
		for _, deployment := range deployments {
			if deployment.Name == name && deployment.Namespace == namespace {
				return s // Deployment found, test can continue
			}
		}
		// Small delay between attempts
		time.Sleep(50 * time.Millisecond)
	}

	// If we get here, the deployment wasn't detected - this will cause the test to fail
	s.t.Errorf("Deployment %s in namespace %s was not detected by watch after %d attempts", name, namespace, maxAttempts)
	return s
}

func (s *DeploymentListControllerScenario) Cleanup() {
	if s.controller != nil {
		s.controller.Stop()
	}
	if s.builder != nil {
		s.builder.Cleanup()
	}
}
