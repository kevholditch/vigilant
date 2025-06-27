package controllers

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kevholditch/vigilant/internal/theme"
)

type PodLogControllerScenario struct {
	t          *testing.T
	builder    *ClusterBuilder
	controller *PodLogController
	podName    string
	namespace  string
	content    string
	width      int
	height     int
}

func NewPodLogControllerScenario(t *testing.T) *PodLogControllerScenario {
	builder := NewClusterBuilder(t)
	return &PodLogControllerScenario{
		t:         t,
		builder:   builder,
		podName:   "test-pod",
		namespace: "default",
		content:   "line1\nline2\nline3\nline4\nline5\nline6\nline7\nline8\nline9\nline10",
		width:     80,
		height:    20,
	}
}

func (s *PodLogControllerScenario) Given() *PodLogControllerScenario { return s }
func (s *PodLogControllerScenario) When() *PodLogControllerScenario  { return s }
func (s *PodLogControllerScenario) Then() *PodLogControllerScenario  { return s }
func (s *PodLogControllerScenario) and() *PodLogControllerScenario   { return s }

func (s *PodLogControllerScenario) ConfigureCluster(configFn func(*ClusterBuilder)) *PodLogControllerScenario {
	configFn(s.builder)
	return s
}

func (s *PodLogControllerScenario) with_pod_name(podName string) *PodLogControllerScenario {
	s.podName = podName
	return s
}

func (s *PodLogControllerScenario) with_namespace(namespace string) *PodLogControllerScenario {
	s.namespace = namespace
	return s
}

func (s *PodLogControllerScenario) with_content(content string) *PodLogControllerScenario {
	s.content = content
	return s
}

func (s *PodLogControllerScenario) with_view_size(width, height int) *PodLogControllerScenario {
	s.width = width
	s.height = height
	return s
}

func (s *PodLogControllerScenario) the_pod_log_controller_is_instantiated() *PodLogControllerScenario {
	theme := theme.NewDefaultTheme()

	// Create a test log fetcher that returns our predefined content
	testLogFetcher := func(podName, namespace string) (string, error) {
		return s.content, nil
	}

	s.controller = NewPodLogController(
		testLogFetcher,
		theme,
		s.podName,
		s.namespace,
		func() tea.Cmd { return nil },
	)

	// Set the view size after controller creation
	s.controller.podLogView.SetSize(s.width, s.height)

	return s
}

func (s *PodLogControllerScenario) scroll_up() *PodLogControllerScenario {
	s.controller.podLogView.ScrollUp()
	return s
}

func (s *PodLogControllerScenario) scroll_down() *PodLogControllerScenario {
	s.controller.podLogView.ScrollDown()
	return s
}

func (s *PodLogControllerScenario) page_up() *PodLogControllerScenario {
	s.controller.podLogView.PageUp()
	return s
}

func (s *PodLogControllerScenario) page_down() *PodLogControllerScenario {
	s.controller.podLogView.PageDown()
	return s
}

func (s *PodLogControllerScenario) go_to_start() *PodLogControllerScenario {
	s.controller.podLogView.GoToStart()
	return s
}

func (s *PodLogControllerScenario) go_to_end() *PodLogControllerScenario {
	s.controller.podLogView.GoToEnd()
	return s
}

func (s *PodLogControllerScenario) refresh_logs() *PodLogControllerScenario {
	cmd := s.controller.refreshLogs()
	if cmd != nil {
		cmd() // simulate running the command
	}
	return s
}

func (s *PodLogControllerScenario) the_scroll_position_should_be(assertFn func(int)) *PodLogControllerScenario {
	scrollY := s.controller.podLogView.GetScrollPosition()
	assertFn(scrollY)
	return s
}

func (s *PodLogControllerScenario) the_action_text_should_be(assertFn func(string)) *PodLogControllerScenario {
	actionText := s.controller.ActionText()
	assertFn(actionText)
	return s
}

func (s *PodLogControllerScenario) Cleanup() {
	if s.builder != nil {
		s.builder.Cleanup()
	}
}
