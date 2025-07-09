package controllers

import (
	"testing"

	"github.com/kevholditch/vigilant/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestStatefulSetListController(t *testing.T) {
	t.Run("should_show_no_statefulsets_when_cluster_is_empty", func(t *testing.T) {
		s := NewStatefulSetListControllerScenario(t)
		defer s.Cleanup()
		s.Given().
			ConfigureCluster(func(builder *ClusterBuilder) {
				// No StatefulSets
			}).
			When().
			the_statefulset_list_controller_is_instantiated().
			Then().
			the_statefulset_list_should_be(func(ss []models.StatefulSet) {
				assert.Empty(t, ss)
			})
	})

	t.Run("should_list_multiple_statefulsets_in_different_namespaces", func(t *testing.T) {
		s := NewStatefulSetListControllerScenario(t)
		defer s.Cleanup()
		s.Given().
			ConfigureCluster(func(builder *ClusterBuilder) {
				builder.WithStatefulSet("ss-a", "ns1").WithStatefulSet("ss-b", "ns2")
			}).
			When().
			the_statefulset_list_controller_is_instantiated().
			Then().
			the_statefulset_list_should_be(func(ss []models.StatefulSet) {
				assert.Len(t, ss, 2)
				assert.ElementsMatch(t, []string{"ss-a", "ss-b"}, []string{ss[0].Name, ss[1].Name})
			})
	})

	t.Run("should_select_next_and_prev_statefulset", func(t *testing.T) {
		s := NewStatefulSetListControllerScenario(t)
		defer s.Cleanup()
		s.Given().
			ConfigureCluster(func(builder *ClusterBuilder) {
				builder.WithStatefulSet("ss-a", "ns1").WithStatefulSet("ss-b", "ns1")
			}).
			When().
			the_statefulset_list_controller_is_instantiated().
			select_next_statefulset().
			Then().
			the_selected_statefulset_should_be(func(ss *models.StatefulSet) {
				assert.Equal(t, "ss-b", ss.Name)
			}).
			and().
			select_prev_statefulset().
			the_selected_statefulset_should_be(func(ss *models.StatefulSet) {
				assert.Equal(t, "ss-a", ss.Name)
			})
	})

	t.Run("should_refresh_statefulsets_when_new_statefulset_is_added", func(t *testing.T) {
		s := NewStatefulSetListControllerScenario(t)
		defer s.Cleanup()
		s.Given().
			ConfigureCluster(func(builder *ClusterBuilder) {
				builder.WithStatefulSet("ss-a", "ns1")
			}).
			When().
			the_statefulset_list_controller_is_instantiated().
			Then().
			the_statefulset_list_should_be(func(ss []models.StatefulSet) {
				assert.Len(t, ss, 1)
				assert.Equal(t, "ss-a", ss[0].Name)
			}).
			and().
			ConfigureCluster(func(builder *ClusterBuilder) {
				builder.WithStatefulSet("ss-b", "ns1")
			}).
			refresh_statefulsets().
			Then().
			the_statefulset_list_should_be(func(ss []models.StatefulSet) {
				assert.Len(t, ss, 2)
				assert.ElementsMatch(t, []string{"ss-a", "ss-b"}, []string{ss[0].Name, ss[1].Name})
			})
	})

	t.Run("should_automatically_detect_new_statefulsets_via_watch", func(t *testing.T) {
		s := NewStatefulSetListControllerScenario(t)
		defer s.Cleanup()
		s.Given().
			ConfigureCluster(func(builder *ClusterBuilder) {
				builder.WithStatefulSet("ss-a", "ns1")
			}).
			the_statefulset_list_controller_is_instantiated().
			Then().
			the_statefulset_list_should_be(func(ss []models.StatefulSet) {
				assert.Len(t, ss, 1)
				assert.Equal(t, "ss-a", ss[0].Name)
			}).
			When().
			a_new_statefulset_is_added_to_cluster("ss-b", "ns1").
			Then().
			the_statefulset_list_should_be(func(ss []models.StatefulSet) {
				assert.Len(t, ss, 2)
				assert.ElementsMatch(t, []string{"ss-a", "ss-b"}, []string{ss[0].Name, ss[1].Name})
			})
	})
}
