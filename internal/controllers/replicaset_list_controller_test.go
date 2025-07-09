package controllers

import (
	"testing"

	"github.com/kevholditch/vigilant/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestReplicaSetListController(t *testing.T) {
	t.Run("should_show_no_replicasets_when_cluster_is_empty", func(t *testing.T) {
		s := NewReplicaSetListControllerScenario(t)
		defer s.Cleanup()
		s.Given().
			ConfigureCluster(func(builder *ClusterBuilder) {
				// No ReplicaSets
			}).
			When().
			the_replicaset_list_controller_is_instantiated().
			Then().
			the_replicaset_list_should_be(func(rs []models.ReplicaSet) {
				assert.Empty(t, rs)
			})
	})

	t.Run("should_list_multiple_replicasets_in_different_namespaces", func(t *testing.T) {
		s := NewReplicaSetListControllerScenario(t)
		defer s.Cleanup()
		s.Given().
			ConfigureCluster(func(builder *ClusterBuilder) {
				builder.WithReplicaSet("rs-a", "ns1").WithReplicaSet("rs-b", "ns2")
			}).
			When().
			the_replicaset_list_controller_is_instantiated().
			Then().
			the_replicaset_list_should_be(func(rs []models.ReplicaSet) {
				assert.Len(t, rs, 2)
				assert.ElementsMatch(t, []string{"rs-a", "rs-b"}, []string{rs[0].Name, rs[1].Name})
			})
	})

	t.Run("should_select_next_and_prev_replicaset", func(t *testing.T) {
		s := NewReplicaSetListControllerScenario(t)
		defer s.Cleanup()
		s.Given().
			ConfigureCluster(func(builder *ClusterBuilder) {
				builder.WithReplicaSet("rs-a", "ns1").WithReplicaSet("rs-b", "ns1")
			}).
			When().
			the_replicaset_list_controller_is_instantiated().
			select_next_replicaset().
			Then().
			the_selected_replicaset_should_be(func(rs *models.ReplicaSet) {
				assert.Equal(t, "rs-b", rs.Name)
			}).
			and().
			select_prev_replicaset().
			the_selected_replicaset_should_be(func(rs *models.ReplicaSet) {
				assert.Equal(t, "rs-a", rs.Name)
			})
	})

	t.Run("should_refresh_replicasets_when_new_replicaset_is_added", func(t *testing.T) {
		s := NewReplicaSetListControllerScenario(t)
		defer s.Cleanup()
		s.Given().
			ConfigureCluster(func(builder *ClusterBuilder) {
				builder.WithReplicaSet("rs-a", "ns1")
			}).
			When().
			the_replicaset_list_controller_is_instantiated().
			Then().
			the_replicaset_list_should_be(func(rs []models.ReplicaSet) {
				assert.Len(t, rs, 1)
				assert.Equal(t, "rs-a", rs[0].Name)
			}).
			and().
			ConfigureCluster(func(builder *ClusterBuilder) {
				builder.WithReplicaSet("rs-b", "ns1")
			}).
			refresh_replicasets().
			Then().
			the_replicaset_list_should_be(func(rs []models.ReplicaSet) {
				assert.Len(t, rs, 2)
				assert.ElementsMatch(t, []string{"rs-a", "rs-b"}, []string{rs[0].Name, rs[1].Name})
			})
	})

	t.Run("should_automatically_detect_new_replicasets_via_watch", func(t *testing.T) {
		s := NewReplicaSetListControllerScenario(t)
		defer s.Cleanup()
		s.Given().
			ConfigureCluster(func(builder *ClusterBuilder) {
				builder.WithReplicaSet("rs-a", "ns1")
			}).
			the_replicaset_list_controller_is_instantiated().
			Then().
			the_replicaset_list_should_be(func(rs []models.ReplicaSet) {
				assert.Len(t, rs, 1)
				assert.Equal(t, "rs-a", rs[0].Name)
			}).
			When().
			a_new_replicaset_is_added_to_cluster("rs-b", "ns1").
			Then().
			the_replicaset_list_should_be(func(rs []models.ReplicaSet) {
				assert.Len(t, rs, 2)
				assert.ElementsMatch(t, []string{"rs-a", "rs-b"}, []string{rs[0].Name, rs[1].Name})
			})
	})
}
