package controllers

import (
	"testing"

	"github.com/kevholditch/vigilant/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestPodListController(t *testing.T) {
	t.Run("should_show_no_pods_when_cluster_is_empty", func(t *testing.T) {
		s := NewPodListControllerScenario(t)
		defer s.Cleanup()
		s.Given().
			ConfigureCluster(func(builder *ClusterBuilder) {
				// No pods
			}).
			When().
			the_pod_list_controller_is_instantiated().
			Then().
			the_pod_list_should_be(func(pods []models.Pod) {
				assert.Empty(t, pods)
			})
	})

	t.Run("should_list_multiple_pods_in_different_namespaces", func(t *testing.T) {
		s := NewPodListControllerScenario(t)
		defer s.Cleanup()
		s.Given().
			ConfigureCluster(func(builder *ClusterBuilder) {
				builder.WithPod("pod-a", "ns1").WithPod("pod-b", "ns2")
			}).
			When().
			the_pod_list_controller_is_instantiated().
			Then().
			the_pod_list_should_be(func(pods []models.Pod) {
				assert.Len(t, pods, 2)
				assert.ElementsMatch(t, []string{"pod-a", "pod-b"}, []string{pods[0].Name, pods[1].Name})
			})
	})

	t.Run("should_select_next_and_prev_pod", func(t *testing.T) {
		s := NewPodListControllerScenario(t)
		defer s.Cleanup()
		s.Given().
			ConfigureCluster(func(builder *ClusterBuilder) {
				builder.WithPod("pod-a", "ns1").WithPod("pod-b", "ns1")
			}).
			When().
			the_pod_list_controller_is_instantiated().
			select_next_pod().
			Then().
			the_selected_pod_should_be(func(pod *models.Pod) {
				assert.Equal(t, "pod-b", pod.Name)
			}).
			and().
			select_prev_pod().
			the_selected_pod_should_be(func(pod *models.Pod) {
				assert.Equal(t, "pod-a", pod.Name)
			})
	})

	t.Run("should_refresh_pods_when_new_pod_is_added", func(t *testing.T) {
		s := NewPodListControllerScenario(t)
		defer s.Cleanup()
		s.Given().
			ConfigureCluster(func(builder *ClusterBuilder) {
				builder.WithPod("pod-a", "ns1")
			}).
			When().
			the_pod_list_controller_is_instantiated().
			Then().
			the_pod_list_should_be(func(pods []models.Pod) {
				assert.Len(t, pods, 1)
				assert.Equal(t, "pod-a", pods[0].Name)
			}).
			and().
			ConfigureCluster(func(builder *ClusterBuilder) {
				builder.WithPod("pod-b", "ns1")
			}).
			refresh_pods().
			Then().
			the_pod_list_should_be(func(pods []models.Pod) {
				assert.Len(t, pods, 2)
				assert.ElementsMatch(t, []string{"pod-a", "pod-b"}, []string{pods[0].Name, pods[1].Name})
			})
	})

	t.Run("should_automatically_detect_new_pods_via_watch", func(t *testing.T) {
		s := NewPodListControllerScenario(t)
		defer s.Cleanup()
		s.Given().
			ConfigureCluster(func(builder *ClusterBuilder) {
				builder.WithPod("pod-a", "ns1")
			}).
			the_pod_list_controller_is_instantiated().
			Then().
			the_pod_list_should_be(func(pods []models.Pod) {
				assert.Len(t, pods, 1)
				assert.Equal(t, "pod-a", pods[0].Name)
			}).
			When().
			a_new_pod_is_added_to_cluster("pod-b", "ns1").
			Then().
			the_pod_list_should_be(func(pods []models.Pod) {
				assert.Len(t, pods, 2)
				assert.ElementsMatch(t, []string{"pod-a", "pod-b"}, []string{pods[0].Name, pods[1].Name})
			})
	})
}
