package controllers

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kevholditch/vigilant/internal/theme"
)

type CommandBarControllerScenario struct {
	t                  *testing.T
	controller         *CommandBarController
	theme              *theme.Theme
	availableResources []string
}

func NewCommandBarControllerScenario(t *testing.T) *CommandBarControllerScenario {
	theme := theme.NewDefaultTheme()
	return &CommandBarControllerScenario{
		t:     t,
		theme: theme,
	}
}

func (s *CommandBarControllerScenario) WithAvailableResources(resources ...string) *CommandBarControllerScenario {
	s.availableResources = resources
	return s
}

func (s *CommandBarControllerScenario) Given() *CommandBarControllerScenario { return s }
func (s *CommandBarControllerScenario) When() *CommandBarControllerScenario  { return s }
func (s *CommandBarControllerScenario) Then() *CommandBarControllerScenario  { return s }
func (s *CommandBarControllerScenario) and() *CommandBarControllerScenario   { return s }

func (s *CommandBarControllerScenario) the_command_bar_controller_is_instantiated() *CommandBarControllerScenario {
	s.controller = NewCommandBarController(nil, s.theme, "", s.availableResources, nil)
	return s
}

func (s *CommandBarControllerScenario) the_command_bar_is_activated() *CommandBarControllerScenario {
	s.controller.Activate()
	return s
}

func (s *CommandBarControllerScenario) the_user_types(char rune) *CommandBarControllerScenario {
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{char}}
	s.controller.HandleKey(msg)
	return s
}

func (s *CommandBarControllerScenario) the_user_presses_backspace() *CommandBarControllerScenario {
	msg := tea.KeyMsg{Type: tea.KeyBackspace}
	s.controller.HandleKey(msg)
	return s
}

func (s *CommandBarControllerScenario) the_user_presses_tab() *CommandBarControllerScenario {
	msg := tea.KeyMsg{Type: tea.KeyTab}
	s.controller.HandleKey(msg)
	return s
}

func (s *CommandBarControllerScenario) the_user_presses_escape() *CommandBarControllerScenario {
	msg := tea.KeyMsg{Type: tea.KeyEscape}
	s.controller.HandleKey(msg)
	return s
}

func (s *CommandBarControllerScenario) the_user_presses_enter() *CommandBarControllerScenario {
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	s.controller.HandleKey(msg)
	return s
}

func (s *CommandBarControllerScenario) the_user_presses_down_arrow() *CommandBarControllerScenario {
	msg := tea.KeyMsg{Type: tea.KeyDown}
	s.controller.HandleKey(msg)
	return s
}

func (s *CommandBarControllerScenario) the_user_presses_up_arrow() *CommandBarControllerScenario {
	msg := tea.KeyMsg{Type: tea.KeyUp}
	s.controller.HandleKey(msg)
	return s
}

func (s *CommandBarControllerScenario) the_user_presses_left_arrow() *CommandBarControllerScenario {
	msg := tea.KeyMsg{Type: tea.KeyLeft}
	s.controller.HandleKey(msg)
	return s
}

func (s *CommandBarControllerScenario) the_user_presses_right_arrow() *CommandBarControllerScenario {
	msg := tea.KeyMsg{Type: tea.KeyRight}
	s.controller.HandleKey(msg)
	return s
}

func (s *CommandBarControllerScenario) the_user_presses_home() *CommandBarControllerScenario {
	msg := tea.KeyMsg{Type: tea.KeyHome}
	s.controller.HandleKey(msg)
	return s
}

func (s *CommandBarControllerScenario) the_user_presses_end() *CommandBarControllerScenario {
	msg := tea.KeyMsg{Type: tea.KeyEnd}
	s.controller.HandleKey(msg)
	return s
}

func (s *CommandBarControllerScenario) the_command_bar_should_be_active(assertFn func(bool)) *CommandBarControllerScenario {
	assertFn(s.controller.IsActive())
	return s
}

func (s *CommandBarControllerScenario) the_input_should_be(assertFn func(string)) *CommandBarControllerScenario {
	assertFn(s.controller.commandBarView.GetInput())
	return s
}

func (s *CommandBarControllerScenario) the_suggestions_should_be(assertFn func([]string)) *CommandBarControllerScenario {
	suggestions := s.controller.commandBarView.GetSuggestions()
	assertFn(suggestions)
	return s
}

func (s *CommandBarControllerScenario) the_selected_suggestion_should_be(assertFn func(string)) *CommandBarControllerScenario {
	selected := s.controller.commandBarView.GetSelectedSuggestion()
	assertFn(selected)
	return s
}

func (s *CommandBarControllerScenario) Cleanup() {
	// No cleanup needed for command bar controller as it doesn't have external resources
}
