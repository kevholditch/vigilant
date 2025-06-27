package controllers

import (
	"testing"

	"github.com/kevholditch/vigilant/internal/models"
	"github.com/kevholditch/vigilant/internal/theme"
)

type PodListControllerScenario struct {
	t          *testing.T
	builder    *ClusterBuilder
	controller *PodListController
	podView    *models.Pod
}

func NewPodListControllerScenario(t *testing.T) *PodListControllerScenario {
	builder := NewClusterBuilder(t)
	return &PodListControllerScenario{
		t:       t,
		builder: builder,
	}
}

func (s *PodListControllerScenario) Given() *PodListControllerScenario { return s }
func (s *PodListControllerScenario) When() *PodListControllerScenario  { return s }
func (s *PodListControllerScenario) Then() *PodListControllerScenario  { return s }
func (s *PodListControllerScenario) and() *PodListControllerScenario   { return s }

func (s *PodListControllerScenario) ConfigureCluster(configFn func(*ClusterBuilder)) *PodListControllerScenario {
	configFn(s.builder)
	return s
}

func (s *PodListControllerScenario) the_pod_list_controller_is_instantiated() *PodListControllerScenario {
	theme := theme.NewDefaultTheme()
	s.controller = NewPodListController(s.builder.GetClientset(), theme, "test-cluster", nil, nil)
	return s
}

func (s *PodListControllerScenario) the_pod_list_view_is_built() *PodListControllerScenario {
	// No-op for now, as controller instantiation builds the view
	return s
}

func (s *PodListControllerScenario) select_next_pod() *PodListControllerScenario {
	s.controller.podView.SelectNext()
	return s
}

func (s *PodListControllerScenario) select_prev_pod() *PodListControllerScenario {
	s.controller.podView.SelectPrev()
	return s
}

func (s *PodListControllerScenario) the_pod_list_should_be(assertFn func([]models.Pod)) *PodListControllerScenario {
	pods := s.controller.podView.Pods()
	assertFn(pods)
	return s
}

func (s *PodListControllerScenario) the_selected_pod_should_be(assertFn func(*models.Pod)) *PodListControllerScenario {
	assertFn(s.controller.podView.GetSelected())
	return s
}

func (s *PodListControllerScenario) refresh_pods() *PodListControllerScenario {
	cmd := s.controller.refreshPods()
	if cmd != nil {
		cmd() // simulate running the command
	}
	return s
}

func (s *PodListControllerScenario) Cleanup() {
	if s.builder != nil {
		s.builder.Cleanup()
	}
}
