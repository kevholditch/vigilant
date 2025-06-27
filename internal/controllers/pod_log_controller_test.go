package controllers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPodLogController(t *testing.T) {
	t.Run("should_show_correct_action_text", func(t *testing.T) {
		s := NewPodLogControllerScenario(t)
		defer s.Cleanup()
		s.Given().
			with_pod_name("test-pod").
			with_namespace("default").
			When().
			the_pod_log_controller_is_instantiated().
			Then().
			the_action_text_should_be(func(actionText string) {
				assert.Equal(t, "Viewing logs for pod test-pod", actionText)
			})
	})

	t.Run("should_scroll_up_and_down_with_predictable_content", func(t *testing.T) {
		s := NewPodLogControllerScenario(t)
		defer s.Cleanup()
		s.Given().
			with_content("line1\nline2\nline3\nline4\nline5\nline6\nline7\nline8\nline9\nline10").
			with_view_size(80, 5).
			When().
			the_pod_log_controller_is_instantiated().
			scroll_down().
			Then().
			the_scroll_position_should_be(func(scrollY int) {
				assert.Equal(t, 1, scrollY)
			}).
			and().
			scroll_down().
			Then().
			the_scroll_position_should_be(func(scrollY int) {
				assert.Equal(t, 2, scrollY)
			}).
			and().
			scroll_up().
			Then().
			the_scroll_position_should_be(func(scrollY int) {
				assert.Equal(t, 1, scrollY)
			})
	})

	t.Run("should_page_up_and_down_with_predictable_content", func(t *testing.T) {
		s := NewPodLogControllerScenario(t)
		defer s.Cleanup()
		s.Given().
			with_content("line1\nline2\nline3\nline4\nline5\nline6\nline7\nline8\nline9\nline10").
			with_view_size(80, 5).
			When().
			the_pod_log_controller_is_instantiated().
			page_down().
			Then().
			the_scroll_position_should_be(func(scrollY int) {
				// With height 5, available height is 5-4=1, so page down should move by 1
				assert.Equal(t, 1, scrollY)
			}).
			and().
			page_down().
			Then().
			the_scroll_position_should_be(func(scrollY int) {
				assert.Equal(t, 2, scrollY)
			}).
			and().
			page_up().
			Then().
			the_scroll_position_should_be(func(scrollY int) {
				assert.Equal(t, 1, scrollY)
			})
	})

	t.Run("should_go_to_start_and_end_with_predictable_content", func(t *testing.T) {
		s := NewPodLogControllerScenario(t)
		defer s.Cleanup()
		s.Given().
			with_content("line1\nline2\nline3\nline4\nline5\nline6\nline7\nline8\nline9\nline10").
			with_view_size(80, 5).
			When().
			the_pod_log_controller_is_instantiated().
			go_to_end().
			Then().
			the_scroll_position_should_be(func(scrollY int) {
				// 10 lines - 1 available height = 9 max scroll
				assert.Equal(t, 9, scrollY)
			}).
			and().
			go_to_start().
			Then().
			the_scroll_position_should_be(func(scrollY int) {
				assert.Equal(t, 0, scrollY)
			})
	})

	t.Run("should_not_scroll_below_zero", func(t *testing.T) {
		s := NewPodLogControllerScenario(t)
		defer s.Cleanup()
		s.Given().
			with_content("line1\nline2\nline3").
			with_view_size(80, 5).
			When().
			the_pod_log_controller_is_instantiated().
			scroll_up().
			Then().
			the_scroll_position_should_be(func(scrollY int) {
				assert.Equal(t, 0, scrollY)
			})
	})

	t.Run("should_not_scroll_above_max", func(t *testing.T) {
		s := NewPodLogControllerScenario(t)
		defer s.Cleanup()
		s.Given().
			with_content("line1\nline2\nline3\nline4\nline5").
			with_view_size(80, 5).
			When().
			the_pod_log_controller_is_instantiated().
			go_to_end().
			Then().
			the_scroll_position_should_be(func(scrollY int) {
				// 5 lines - 1 available height = 4 max scroll
				assert.Equal(t, 4, scrollY)
			}).
			and().
			scroll_down().
			Then().
			the_scroll_position_should_be(func(scrollY int) {
				// Should remain at max position
				assert.Equal(t, 4, scrollY)
			})
	})

	t.Run("should_reset_scroll_position_on_refresh", func(t *testing.T) {
		s := NewPodLogControllerScenario(t)
		defer s.Cleanup()
		s.Given().
			with_content("line1\nline2\nline3\nline4\nline5").
			with_view_size(80, 5).
			When().
			the_pod_log_controller_is_instantiated().
			scroll_down().
			Then().
			the_scroll_position_should_be(func(scrollY int) {
				assert.Equal(t, 1, scrollY)
			}).
			and().
			refresh_logs().
			Then().
			the_scroll_position_should_be(func(scrollY int) {
				assert.Equal(t, 0, scrollY)
			})
	})
}
