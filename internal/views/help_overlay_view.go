package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/kevholditch/vigilant/internal/theme"
)

// KeyBinding represents a key binding with description
type KeyBinding struct {
	Key         string
	Description string
}

// HelpOverlayView represents the help overlay view
type HelpOverlayView struct {
	theme  *theme.Theme
	width  int
	height int
}

// NewHelpOverlayView creates a new help overlay view
func NewHelpOverlayView(theme *theme.Theme) *HelpOverlayView {
	return &HelpOverlayView{
		theme: theme,
	}
}

// SetSize sets the view dimensions
func (hov *HelpOverlayView) SetSize(width, height int) {
	hov.width = width
	hov.height = height
}

// Render renders the help overlay view
func (hov *HelpOverlayView) Render(width, height int, keyBindings []KeyBinding) string {
	hov.SetSize(width, height)

	// Create a semi-transparent overlay effect
	overlayStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("0")).
		Foreground(lipgloss.Color("7")).
		Width(width).
		Height(height)

	// Create the help modal content
	modalContent := hov.renderHelpContent(keyBindings)

	// Center the modal on screen
	modalWidth := 60
	modalHeight := lipgloss.Height(modalContent)

	// Calculate centering
	leftMargin := (width - modalWidth) / 2
	topMargin := (height - modalHeight) / 2

	// Create modal container with border
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(hov.theme.Primary).
		Background(hov.theme.BgPrimary).
		Foreground(hov.theme.TextPrimary).
		Padding(1, 2).
		Width(modalWidth).
		Margin(topMargin, leftMargin)

	// Render the modal
	modal := modalStyle.Render(modalContent)

	// Create the overlay with the modal
	overlay := overlayStyle.Render(modal)

	return overlay
}

// renderHelpContent renders the help content
func (hov *HelpOverlayView) renderHelpContent(keyBindings []KeyBinding) string {
	var sections []string

	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground(hov.theme.Primary).
		Bold(true).
		Align(lipgloss.Center)

	title := titleStyle.Render("Vigilant Help")
	sections = append(sections, title)
	sections = append(sections, "")

	// Navigation (always shown)
	navigationSection := hov.renderSection("Navigation", []string{
		"↑/↓ or j/k    Navigate through items",
		"q              Quit the application",
		"h              Show this help",
		"esc            Close help / Go back",
		":              Open command bar",
	})
	sections = append(sections, navigationSection)

	// Context-specific commands
	if len(keyBindings) > 0 {
		contextSection := hov.renderKeyBindingsSection("Current View Commands", keyBindings)
		sections = append(sections, contextSection)
	}

	// Command Bar (always shown)
	commandBarSection := hov.renderSection("Command Bar", []string{
		"Enter          Execute command",
		"Tab            Cycle through suggestions",
		"↑/↓           Navigate suggestions",
		"Ctrl+C         Cancel command",
	})
	sections = append(sections, commandBarSection)

	// Footer
	footerStyle := lipgloss.NewStyle().
		Foreground(hov.theme.TextMuted).
		Align(lipgloss.Center).
		Italic(true)

	footer := footerStyle.Render("Press 'esc' to close")
	sections = append(sections, "")
	sections = append(sections, footer)

	return strings.Join(sections, "\n")
}

// renderSection renders a help section
func (hov *HelpOverlayView) renderSection(title string, items []string) string {
	var lines []string

	// Section title
	titleStyle := lipgloss.NewStyle().
		Foreground(hov.theme.Primary).
		Bold(true)

	lines = append(lines, titleStyle.Render(title))

	// Section items
	for _, item := range items {
		itemStyle := lipgloss.NewStyle().
			Foreground(hov.theme.TextPrimary).
			PaddingLeft(2)

		lines = append(lines, itemStyle.Render(item))
	}

	return strings.Join(lines, "\n")
}

// renderKeyBindingsSection renders a key bindings section
func (hov *HelpOverlayView) renderKeyBindingsSection(title string, keyBindings []KeyBinding) string {
	var lines []string

	// Section title
	titleStyle := lipgloss.NewStyle().
		Foreground(hov.theme.Primary).
		Bold(true)

	lines = append(lines, titleStyle.Render(title))

	// Section items
	for _, binding := range keyBindings {
		itemStyle := lipgloss.NewStyle().
			Foreground(hov.theme.TextPrimary).
			PaddingLeft(2)

		// Format: "key    description"
		formatted := fmt.Sprintf("%-15s %s", binding.Key, binding.Description)
		lines = append(lines, itemStyle.Render(formatted))
	}

	return strings.Join(lines, "\n")
}
