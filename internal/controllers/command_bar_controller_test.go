package controllers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandBarController(t *testing.T) {
	t.Run("should_activate_and_deactivate_command_bar", func(t *testing.T) {
		s := NewCommandBarControllerScenario(t).
			WithAvailableResources("pods", "deployments")
		defer s.Cleanup()
		s.Given().
			the_command_bar_controller_is_instantiated().
			Then().
			the_command_bar_should_be_active(func(active bool) {
				assert.False(t, active)
			}).
			When().
			the_command_bar_is_activated().
			Then().
			the_command_bar_should_be_active(func(active bool) {
				assert.True(t, active)
			}).
			When().
			the_user_presses_escape().
			Then().
			the_command_bar_should_be_active(func(active bool) {
				assert.False(t, active)
			})
	})

	t.Run("should_handle_character_input", func(t *testing.T) {
		s := NewCommandBarControllerScenario(t).
			WithAvailableResources("pods", "deployments")
		defer s.Cleanup()
		s.Given().
			the_command_bar_controller_is_instantiated().
			the_command_bar_is_activated().
			When().
			the_user_types('p').
			the_user_types('o').
			the_user_types('d').
			Then().
			the_input_should_be(func(input string) {
				assert.Equal(t, "pod", input)
			})
	})

	t.Run("should_handle_backspace", func(t *testing.T) {
		s := NewCommandBarControllerScenario(t).
			WithAvailableResources("pods", "deployments")
		defer s.Cleanup()
		s.Given().
			the_command_bar_controller_is_instantiated().
			the_command_bar_is_activated().
			the_user_types('p').
			the_user_types('o').
			When().
			the_user_presses_backspace().
			Then().
			the_input_should_be(func(input string) {
				assert.Equal(t, "p", input)
			})
	})

	t.Run("should_handle_tab_completion", func(t *testing.T) {
		s := NewCommandBarControllerScenario(t).
			WithAvailableResources("pods", "deployments")
		defer s.Cleanup()
		s.Given().
			the_command_bar_controller_is_instantiated().
			the_command_bar_is_activated().
			the_user_types('p').
			When().
			the_user_presses_tab().
			Then().
			the_input_should_be(func(input string) {
				assert.Equal(t, "pods", input)
			})
	})

	t.Run("should_handle_arrow_keys_for_suggestions", func(t *testing.T) {
		s := NewCommandBarControllerScenario(t).
			WithAvailableResources("pods", "deployments")
		defer s.Cleanup()
		s.Given().
			the_command_bar_controller_is_instantiated().
			the_command_bar_is_activated().
			Then().
			the_suggestions_should_be(func(suggestions []string) {
				assert.Len(t, suggestions, 2)
				assert.Equal(t, "pods", suggestions[0])
				assert.Equal(t, "deployments", suggestions[1])
			}).
			When().
			the_user_presses_down_arrow().
			Then().
			the_selected_suggestion_should_be(func(selected string) {
				assert.Equal(t, "deployments", selected)
			})
	})

	t.Run("should_show_filtered_suggestions_when_typing", func(t *testing.T) {
		s := NewCommandBarControllerScenario(t).
			WithAvailableResources("pods", "deployments")
		defer s.Cleanup()
		s.Given().
			the_command_bar_controller_is_instantiated().
			the_command_bar_is_activated().
			When().
			the_user_types('d').
			Then().
			the_suggestions_should_be(func(suggestions []string) {
				assert.Len(t, suggestions, 1)
				assert.Equal(t, "deployments", suggestions[0])
			})
	})

	t.Run("should_handle_cursor_movement", func(t *testing.T) {
		s := NewCommandBarControllerScenario(t).
			WithAvailableResources("pods", "deployments")
		defer s.Cleanup()
		s.Given().
			the_command_bar_controller_is_instantiated().
			the_command_bar_is_activated().
			the_user_types('p').
			the_user_types('o').
			the_user_types('d').
			When().
			the_user_presses_left_arrow().
			the_user_presses_left_arrow().
			the_user_types('x').
			Then().
			the_input_should_be(func(input string) {
				assert.Equal(t, "pxod", input)
			})
	})

	t.Run("should_handle_home_and_end_keys", func(t *testing.T) {
		s := NewCommandBarControllerScenario(t).
			WithAvailableResources("pods", "deployments")
		defer s.Cleanup()
		s.Given().
			the_command_bar_controller_is_instantiated().
			the_command_bar_is_activated().
			the_user_types('p').
			the_user_types('o').
			the_user_types('d').
			When().
			the_user_presses_home().
			the_user_types('x').
			Then().
			the_input_should_be(func(input string) {
				assert.Equal(t, "xpod", input)
			}).
			When().
			the_user_presses_end().
			the_user_types('s').
			Then().
			the_input_should_be(func(input string) {
				assert.Equal(t, "xpods", input)
			})
	})

	t.Run("should_deactivate_on_enter", func(t *testing.T) {
		s := NewCommandBarControllerScenario(t).
			WithAvailableResources("pods", "deployments")
		defer s.Cleanup()
		s.Given().
			the_command_bar_controller_is_instantiated().
			the_command_bar_is_activated().
			the_user_types('p').
			the_user_types('o').
			the_user_types('d').
			When().
			the_user_presses_enter().
			Then().
			the_command_bar_should_be_active(func(active bool) {
				assert.False(t, active)
			}).
			and().
			the_input_should_be(func(input string) {
				assert.Equal(t, "", input)
			})
	})

	t.Run("should_filter_similar_resource_names", func(t *testing.T) {
		s := NewCommandBarControllerScenario(t).
			WithAvailableResources("pod", "pods", "podman", "podinfo", "deployments")
		defer s.Cleanup()
		s.Given().
			the_command_bar_controller_is_instantiated().
			the_command_bar_is_activated().
			Then().
			the_suggestions_should_be(func(suggestions []string) {
				assert.ElementsMatch(t, []string{"pod", "pods", "podman", "podinfo", "deployments"}, suggestions)
			})
		// Type 'p' - all pod* options
		s.When().the_user_types('p').Then().the_suggestions_should_be(func(suggestions []string) {
			assert.ElementsMatch(t, []string{"pod", "pods", "podman", "podinfo"}, suggestions)
		})
		// Type 'o' - still all pod* options
		s.When().the_user_types('o').Then().the_suggestions_should_be(func(suggestions []string) {
			assert.ElementsMatch(t, []string{"pod", "pods", "podman", "podinfo"}, suggestions)
		})
		// Type 'd' - only pod, pods, podman, podinfo
		s.When().the_user_types('d').Then().the_suggestions_should_be(func(suggestions []string) {
			assert.ElementsMatch(t, []string{"pod", "pods", "podman", "podinfo"}, suggestions)
		})
		// Type 'm' - only podman
		s.When().the_user_types('m').Then().the_suggestions_should_be(func(suggestions []string) {
			assert.Equal(t, []string{"podman"}, suggestions)
		})
		// Backspace to 'pod' - pod, pods, podman, podinfo
		s.When().the_user_presses_backspace().Then().the_suggestions_should_be(func(suggestions []string) {
			assert.ElementsMatch(t, []string{"pod", "pods", "podman", "podinfo"}, suggestions)
		})
		// Type 'i' - only podinfo
		s.When().the_user_types('i').Then().the_suggestions_should_be(func(suggestions []string) {
			assert.Equal(t, []string{"podinfo"}, suggestions)
		})
		// Backspace to 'pod' - pod, pods, podman, podinfo
		s.When().the_user_presses_backspace().Then().the_suggestions_should_be(func(suggestions []string) {
			assert.ElementsMatch(t, []string{"pod", "pods", "podman", "podinfo"}, suggestions)
		})
		// Type 's' - only pods
		s.When().the_user_types('s').Then().the_suggestions_should_be(func(suggestions []string) {
			assert.Equal(t, []string{"pods"}, suggestions)
		})
	})
}
