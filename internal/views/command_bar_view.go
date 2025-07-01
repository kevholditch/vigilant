package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/kevholditch/vigilant/internal/theme"
)

// CommandBarView represents the command bar view
type CommandBarView struct {
	theme              *theme.Theme
	width              int
	height             int
	input              string
	cursor             int
	isActive           bool
	suggestions        []string
	selectedSuggestion int
	availableResources []string
}

// NewCommandBarView creates a new command bar view
func NewCommandBarView(theme *theme.Theme, availableResources []string) *CommandBarView {
	return &CommandBarView{
		theme:              theme,
		input:              "",
		cursor:             0,
		isActive:           false,
		suggestions:        []string{},
		selectedSuggestion: 0,
		availableResources: availableResources,
	}
}

// SetSize sets the view dimensions
func (cbv *CommandBarView) SetSize(width, height int) {
	cbv.width = width
	cbv.height = height
}

// Activate activates the command bar
func (cbv *CommandBarView) Activate() {
	cbv.isActive = true
	cbv.input = ""
	cbv.cursor = 0
	cbv.updateSuggestions()
}

// Deactivate deactivates the command bar
func (cbv *CommandBarView) Deactivate() {
	cbv.isActive = false
	cbv.input = ""
	cbv.cursor = 0
	cbv.suggestions = []string{}
	cbv.selectedSuggestion = 0
}

// IsActive returns whether the command bar is active
func (cbv *CommandBarView) IsActive() bool {
	return cbv.isActive
}

// AddChar adds a character to the input
func (cbv *CommandBarView) AddChar(char rune) {
	if !cbv.isActive {
		return
	}

	// Insert character at cursor position
	if cbv.cursor == len(cbv.input) {
		cbv.input += string(char)
	} else {
		cbv.input = cbv.input[:cbv.cursor] + string(char) + cbv.input[cbv.cursor:]
	}
	cbv.cursor++
	cbv.updateSuggestions()
}

// DeleteChar deletes a character from the input
func (cbv *CommandBarView) DeleteChar() {
	if !cbv.isActive || cbv.cursor == 0 {
		return
	}

	if cbv.cursor == len(cbv.input) {
		cbv.input = cbv.input[:len(cbv.input)-1]
	} else {
		cbv.input = cbv.input[:cbv.cursor-1] + cbv.input[cbv.cursor:]
	}
	cbv.cursor--
	cbv.updateSuggestions()
}

// MoveCursorLeft moves the cursor left
func (cbv *CommandBarView) MoveCursorLeft() {
	if !cbv.isActive || cbv.cursor == 0 {
		return
	}
	cbv.cursor--
}

// MoveCursorRight moves the cursor right
func (cbv *CommandBarView) MoveCursorRight() {
	if !cbv.isActive || cbv.cursor == len(cbv.input) {
		return
	}
	cbv.cursor++
}

// MoveCursorToStart moves the cursor to the start
func (cbv *CommandBarView) MoveCursorToStart() {
	if !cbv.isActive {
		return
	}
	cbv.cursor = 0
}

// MoveCursorToEnd moves the cursor to the end
func (cbv *CommandBarView) MoveCursorToEnd() {
	if !cbv.isActive {
		return
	}
	cbv.cursor = len(cbv.input)
}

// SelectNextSuggestion selects the next suggestion
func (cbv *CommandBarView) SelectNextSuggestion() {
	if !cbv.isActive || len(cbv.suggestions) == 0 {
		return
	}
	cbv.selectedSuggestion = (cbv.selectedSuggestion + 1) % len(cbv.suggestions)
}

// SelectPrevSuggestion selects the previous suggestion
func (cbv *CommandBarView) SelectPrevSuggestion() {
	if !cbv.isActive || len(cbv.suggestions) == 0 {
		return
	}
	cbv.selectedSuggestion = (cbv.selectedSuggestion - 1 + len(cbv.suggestions)) % len(cbv.suggestions)
}

// TabComplete completes the input with the selected suggestion
func (cbv *CommandBarView) TabComplete() {
	if !cbv.isActive || len(cbv.suggestions) == 0 {
		return
	}

	selected := cbv.suggestions[cbv.selectedSuggestion]
	cbv.input = selected
	cbv.cursor = len(selected)
	cbv.updateSuggestions()
}

// GetInput returns the current input
func (cbv *CommandBarView) GetInput() string {
	return cbv.input
}

// GetSelectedSuggestion returns the currently selected suggestion
func (cbv *CommandBarView) GetSelectedSuggestion() string {
	if len(cbv.suggestions) == 0 {
		return ""
	}
	return cbv.suggestions[cbv.selectedSuggestion]
}

// updateSuggestions updates the suggestions based on current input
func (cbv *CommandBarView) updateSuggestions() {
	cbv.suggestions = []string{}
	cbv.selectedSuggestion = 0

	if cbv.input == "" {
		// Show all resources when input is empty
		cbv.suggestions = cbv.availableResources
		return
	}

	// Filter resources that match the input
	for _, resource := range cbv.availableResources {
		if strings.HasPrefix(strings.ToLower(resource), strings.ToLower(cbv.input)) {
			cbv.suggestions = append(cbv.suggestions, resource)
		}
	}
}

// Render renders the command bar view
func (cbv *CommandBarView) Render() string {
	if !cbv.isActive {
		return ""
	}

	// Command bar background
	barStyle := lipgloss.NewStyle().
		Background(cbv.theme.BgSecondary).
		Foreground(cbv.theme.TextPrimary).
		Padding(0, 1).
		Width(cbv.width)

	// Prompt
	prompt := ":"
	promptStyle := lipgloss.NewStyle().
		Foreground(cbv.theme.Primary).
		Bold(true)

	// Input area
	inputStyle := lipgloss.NewStyle().
		Foreground(cbv.theme.TextPrimary).
		Background(cbv.theme.BgPrimary).
		Padding(0, 1)

	// Create the input display with cursor
	inputDisplay := cbv.input
	if cbv.cursor < len(cbv.input) {
		// Insert cursor character
		inputDisplay = cbv.input[:cbv.cursor] + "█" + cbv.input[cbv.cursor:]
	} else {
		// Cursor at end
		inputDisplay = cbv.input + "█"
	}

	// Build the command bar content
	content := promptStyle.Render(prompt) + " " + inputStyle.Render(inputDisplay)

	// Add suggestions if available
	if len(cbv.suggestions) > 0 {
		suggestionsText := cbv.renderSuggestions()
		content += "\n" + suggestionsText
	}

	return barStyle.Render(content)
}

// renderSuggestions renders the suggestions list
func (cbv *CommandBarView) renderSuggestions() string {
	if len(cbv.suggestions) == 0 {
		return ""
	}

	var suggestionLines []string
	for i, suggestion := range cbv.suggestions {
		var style lipgloss.Style
		if i == cbv.selectedSuggestion {
			style = lipgloss.NewStyle().
				Foreground(cbv.theme.TextInverse).
				Background(cbv.theme.Primary).
				Bold(true)
		} else {
			style = lipgloss.NewStyle().
				Foreground(cbv.theme.TextSecondary)
		}

		// Add tab completion hint
		suggestionText := fmt.Sprintf("  %s", suggestion)
		if i == cbv.selectedSuggestion {
			suggestionText += " (Tab to complete)"
		}

		suggestionLines = append(suggestionLines, style.Render(suggestionText))
	}

	return strings.Join(suggestionLines, "\n")
}

// GetSuggestions returns the current suggestions (for testing)
func (cbv *CommandBarView) GetSuggestions() []string {
	return cbv.suggestions
}
