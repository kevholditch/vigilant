package controllers

import (
	"testing"
	"time"

	"github.com/kevholditch/vigilant/internal/models"
	"github.com/kevholditch/vigilant/internal/theme"
)

type StatefulSetListControllerScenario struct {
	t               *testing.T
	builder         *ClusterBuilder
	controller      *StatefulSetListController
	statefulsetView *models.StatefulSet
}

func NewStatefulSetListControllerScenario(t *testing.T) *StatefulSetListControllerScenario {
	builder := NewClusterBuilder(t)
	return &StatefulSetListControllerScenario{
		t:       t,
		builder: builder,
	}
}

func (s *StatefulSetListControllerScenario) Given() *StatefulSetListControllerScenario { return s }
func (s *StatefulSetListControllerScenario) When() *StatefulSetListControllerScenario  { return s }
func (s *StatefulSetListControllerScenario) Then() *StatefulSetListControllerScenario  { return s }
func (s *StatefulSetListControllerScenario) and() *StatefulSetListControllerScenario   { return s }

func (s *StatefulSetListControllerScenario) ConfigureCluster(configFn func(*ClusterBuilder)) *StatefulSetListControllerScenario {
	configFn(s.builder)
	return s
}

func (s *StatefulSetListControllerScenario) the_statefulset_list_controller_is_instantiated() *StatefulSetListControllerScenario {
	theme := theme.NewDefaultTheme()
	s.controller = NewStatefulSetListController(s.builder.GetClientset(), theme, "test-cluster", nil)
	return s
}

func (s *StatefulSetListControllerScenario) select_next_statefulset() *StatefulSetListControllerScenario {
	s.controller.statefulsetView.SelectNext()
	return s
}

func (s *StatefulSetListControllerScenario) select_prev_statefulset() *StatefulSetListControllerScenario {
	s.controller.statefulsetView.SelectPrev()
	return s
}

func (s *StatefulSetListControllerScenario) the_statefulset_list_should_be(assertFn func([]models.StatefulSet)) *StatefulSetListControllerScenario {
	ss := s.controller.GetStatefulSets()
	assertFn(ss)
	return s
}

func (s *StatefulSetListControllerScenario) the_selected_statefulset_should_be(assertFn func(*models.StatefulSet)) *StatefulSetListControllerScenario {
	assertFn(s.controller.statefulsetView.GetSelected())
	return s
}

func (s *StatefulSetListControllerScenario) refresh_statefulsets() *StatefulSetListControllerScenario {
	cmd := s.controller.refreshStatefulSets()
	if cmd != nil {
		cmd()
	}
	return s
}

func (s *StatefulSetListControllerScenario) a_new_statefulset_is_added_to_cluster(name, namespace string) *StatefulSetListControllerScenario {
	s.builder.WithStatefulSet(name, namespace)
	maxAttempts := 10
	for i := 0; i < maxAttempts; i++ {
		ss := s.controller.GetStatefulSets()
		for _, sset := range ss {
			if sset.Name == name && sset.Namespace == namespace {
				return s
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
	s.t.Errorf("StatefulSet %s in namespace %s was not detected by watch after %d attempts", name, namespace, maxAttempts)
	return s
}

func (s *StatefulSetListControllerScenario) Cleanup() {
	if s.controller != nil {
		s.controller.Stop()
	}
	if s.builder != nil {
		s.builder.Cleanup()
	}
}
