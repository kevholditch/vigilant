package controllers

import (
	"testing"
)

func TestHeaderController(t *testing.T) {
	t.Run("should correctly count worker nodes", func(t *testing.T) {
		s := NewHeaderControllerScenario(t)
		defer s.Cleanup()

		s.Given().a_kubernetes_cluster_with_3_worker_nodes()

		s.When().the_header_controller_is_instantiated().and().
			the_header_model_is_built()

		s.Then().the_header_model_worker_nodes_count_should_be(3).and().
			the_header_model_control_plane_nodes_count_should_be(0).and().
			the_header_model_kubernetes_version_should_not_be_empty().and().
			the_header_model_cluster_name_should_be_empty()
	})

	t.Run("should correctly count control plane nodes", func(t *testing.T) {
		s := NewHeaderControllerScenario(t)
		defer s.Cleanup()

		s.Given().a_kubernetes_cluster_with_2_control_plane_nodes()

		s.When().the_header_controller_is_instantiated().and().
			the_header_model_is_built()

		s.Then().the_header_model_worker_nodes_count_should_be(0).and().
			the_header_model_control_plane_nodes_count_should_be(2).and().
			the_header_model_kubernetes_version_should_not_be_empty().and().
			the_header_model_cluster_name_should_be_empty()
	})

	t.Run("should correctly count mixed nodes", func(t *testing.T) {
		s := NewHeaderControllerScenario(t)
		defer s.Cleanup()

		s.Given().a_kubernetes_cluster_with_mixed_nodes()

		s.When().the_header_controller_is_instantiated().and().
			the_header_model_is_built()

		s.Then().the_header_model_worker_nodes_count_should_be(3).and().
			the_header_model_control_plane_nodes_count_should_be(2).and().
			the_header_model_kubernetes_version_should_not_be_empty().and().
			the_header_model_cluster_name_should_be_empty()
	})

	t.Run("should handle legacy master nodes", func(t *testing.T) {
		s := NewHeaderControllerScenario(t)
		defer s.Cleanup()

		s.Given().a_kubernetes_cluster_with_legacy_master_nodes()

		s.When().the_header_controller_is_instantiated().and().
			the_header_model_is_built()

		s.Then().the_header_model_worker_nodes_count_should_be(0).and().
			the_header_model_control_plane_nodes_count_should_be(1).and().
			the_header_model_kubernetes_version_should_not_be_empty().and().
			the_header_model_cluster_name_should_be_empty()
	})

	t.Run("should handle empty cluster", func(t *testing.T) {
		s := NewHeaderControllerScenario(t)
		defer s.Cleanup()

		s.Given().an_empty_kubernetes_cluster()

		s.When().the_header_controller_is_instantiated().and().
			the_header_model_is_built()

		s.Then().the_header_model_worker_nodes_count_should_be(0).and().
			the_header_model_control_plane_nodes_count_should_be(0).and().
			the_header_model_kubernetes_version_should_not_be_empty().and().
			the_header_model_cluster_name_should_be_empty()
	})
}
