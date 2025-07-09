package controllers

import (
	"testing"
	"time"

	"github.com/kevholditch/vigilant/internal/models"
	"github.com/kevholditch/vigilant/internal/theme"
)

type ReplicaSetListControllerScenario struct {
	t              *testing.T
	builder        *ClusterBuilder
	controller     *ReplicaSetListController
	replicasetView *models.ReplicaSet
}

func NewReplicaSetListControllerScenario(t *testing.T) *ReplicaSetListControllerScenario {
	builder := NewClusterBuilder(t)
	return &ReplicaSetListControllerScenario{
		t:       t,
		builder: builder,
	}
}

func (s *ReplicaSetListControllerScenario) Given() *ReplicaSetListControllerScenario { return s }
func (s *ReplicaSetListControllerScenario) When() *ReplicaSetListControllerScenario  { return s }
func (s *ReplicaSetListControllerScenario) Then() *ReplicaSetListControllerScenario  { return s }
func (s *ReplicaSetListControllerScenario) and() *ReplicaSetListControllerScenario   { return s }

func (s *ReplicaSetListControllerScenario) ConfigureCluster(configFn func(*ClusterBuilder)) *ReplicaSetListControllerScenario {
	configFn(s.builder)
	return s
}

func (s *ReplicaSetListControllerScenario) the_replicaset_list_controller_is_instantiated() *ReplicaSetListControllerScenario {
	theme := theme.NewDefaultTheme()
	s.controller = NewReplicaSetListController(s.builder.GetClientset(), theme, "test-cluster", nil)
	return s
}

func (s *ReplicaSetListControllerScenario) select_next_replicaset() *ReplicaSetListControllerScenario {
	s.controller.replicasetView.SelectNext()
	return s
}

func (s *ReplicaSetListControllerScenario) select_prev_replicaset() *ReplicaSetListControllerScenario {
	s.controller.replicasetView.SelectPrev()
	return s
}

func (s *ReplicaSetListControllerScenario) the_replicaset_list_should_be(assertFn func([]models.ReplicaSet)) *ReplicaSetListControllerScenario {
	rs := s.controller.GetReplicaSets()
	assertFn(rs)
	return s
}

func (s *ReplicaSetListControllerScenario) the_selected_replicaset_should_be(assertFn func(*models.ReplicaSet)) *ReplicaSetListControllerScenario {
	assertFn(s.controller.replicasetView.GetSelected())
	return s
}

func (s *ReplicaSetListControllerScenario) refresh_replicasets() *ReplicaSetListControllerScenario {
	cmd := s.controller.refreshReplicaSets()
	if cmd != nil {
		cmd()
	}
	return s
}

func (s *ReplicaSetListControllerScenario) a_new_replicaset_is_added_to_cluster(name, namespace string) *ReplicaSetListControllerScenario {
	s.builder.WithReplicaSet(name, namespace)
	maxAttempts := 10
	for i := 0; i < maxAttempts; i++ {
		rs := s.controller.GetReplicaSets()
		for _, r := range rs {
			if r.Name == name && r.Namespace == namespace {
				return s
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
	s.t.Errorf("ReplicaSet %s in namespace %s was not detected by watch after %d attempts", name, namespace, maxAttempts)
	return s
}

func (s *ReplicaSetListControllerScenario) Cleanup() {
	if s.controller != nil {
		s.controller.Stop()
	}
	if s.builder != nil {
		s.builder.Cleanup()
	}
}
