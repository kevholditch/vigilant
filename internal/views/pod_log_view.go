package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/kevholditch/vigilant/internal/theme"
)

// PodLogView represents the pod logs view
type PodLogView struct {
	podName   string
	namespace string
	content   string
	scrollY   int
	width     int
	height    int
	theme     *theme.Theme
}

// NewPodLogView creates a new pod log view
func NewPodLogView(podName, namespace string, theme *theme.Theme) *PodLogView {
	view := &PodLogView{
		podName:   podName,
		namespace: namespace,
		theme:     theme,
		content:   "Loading logs...",
	}
	return view
}

// SetSize sets the view dimensions
func (plv *PodLogView) SetSize(width, height int) {
	plv.width = width
	plv.height = height
}

// ScrollUp moves the view up
func (plv *PodLogView) ScrollUp() {
	if plv.scrollY > 0 {
		plv.scrollY--
	}
}

// ScrollDown moves the view down
func (plv *PodLogView) ScrollDown() {
	// Calculate max scroll based on content height
	contentLines := strings.Split(plv.content, "\n")
	availableHeight := plv.height - 4 // Header + status bar + borders

	if availableHeight <= 0 {
		return
	}

	maxScroll := len(contentLines) - availableHeight
	if maxScroll < 0 {
		maxScroll = 0
	}

	if plv.scrollY < maxScroll {
		plv.scrollY++
	}
}

// PageUp scrolls up by one page
func (plv *PodLogView) PageUp() {
	availableHeight := plv.height - 4 // Header + status bar + borders
	if availableHeight <= 0 {
		availableHeight = 1
	}
	plv.scrollY -= availableHeight
	if plv.scrollY < 0 {
		plv.scrollY = 0
	}
}

// PageDown scrolls down by one page
func (plv *PodLogView) PageDown() {
	contentLines := strings.Split(plv.content, "\n")
	availableHeight := plv.height - 4 // Header + status bar + borders
	if availableHeight <= 0 {
		availableHeight = 1
	}
	maxScroll := len(contentLines) - availableHeight
	if maxScroll < 0 {
		maxScroll = 0
	}
	plv.scrollY += availableHeight
	if plv.scrollY > maxScroll {
		plv.scrollY = maxScroll
	}
}

// GoToStart scrolls to the top
func (plv *PodLogView) GoToStart() {
	plv.scrollY = 0
}

// GoToEnd scrolls to the bottom
func (plv *PodLogView) GoToEnd() {
	contentLines := strings.Split(plv.content, "\n")
	availableHeight := plv.height - 4 // Header + status bar + borders
	if availableHeight <= 0 {
		availableHeight = 1
	}
	maxScroll := len(contentLines) - availableHeight
	if maxScroll < 0 {
		maxScroll = 0
	}
	plv.scrollY = maxScroll
}

// RefreshLogs resets scroll position (content update is handled by controller)
func (plv *PodLogView) RefreshLogs() {
	plv.scrollY = 0 // Reset scroll position
}

// UpdateContent updates the log content and resets scroll position
func (plv *PodLogView) UpdateContent(content string) {
	plv.content = content
	plv.scrollY = 0 // Reset scroll position when content changes
}

// Render renders the complete pod log view
func (plv *PodLogView) Render() string {
	if plv.width == 0 || plv.height == 0 {
		return ""
	}

	// Content area
	content := plv.renderContent()

	// Status bar
	statusBar := plv.renderStatusBar()

	// Combine all components
	viewContent := lipgloss.JoinVertical(
		lipgloss.Left,
		content,
		statusBar,
	)

	return viewContent
}

// renderContent renders the scrollable content area
func (plv *PodLogView) renderContent() string {
	// Split content into lines
	lines := strings.Split(plv.content, "\n")

	// Calculate available height for content (subtract status bar)
	availableHeight := plv.height - 1 // 1 for status bar
	if availableHeight < 0 {
		availableHeight = 0
	}

	if availableHeight <= 0 {
		return lipgloss.NewStyle().Foreground(plv.theme.TextMuted).Render("Window too small")
	}

	// Ensure scroll position is valid
	if plv.scrollY < 0 {
		plv.scrollY = 0
	}

	maxScroll := len(lines) - availableHeight
	if maxScroll < 0 {
		maxScroll = 0
	}

	if plv.scrollY > maxScroll {
		plv.scrollY = maxScroll
	}

	// Apply scrolling
	startLine := plv.scrollY
	endLine := startLine + availableHeight
	if endLine > len(lines) {
		endLine = len(lines)
	}

	// Get visible lines
	var visibleLines []string
	if startLine < len(lines) {
		visibleLines = lines[startLine:endLine]
	}

	// Join visible lines
	content := strings.Join(visibleLines, "\n")

	// Style the content
	contentStyle := lipgloss.NewStyle().
		Foreground(plv.theme.TextPrimary).
		Background(plv.theme.BgPrimary).
		Width(plv.width).
		Height(availableHeight)

	return contentStyle.Render(content)
}

// renderStatusBar renders the status bar at the bottom
func (plv *PodLogView) renderStatusBar() string {
	statusText := fmt.Sprintf("Logs: %s | Press 'Esc' to return | Use ↑/↓ to scroll | PgUp/PgDn for page scroll | g/G for start/end | Press 'r' to refresh", plv.podName)
	return plv.theme.StatusBarStyle.Width(plv.width).Render(statusText)
}

// GetScrollPosition returns the current scroll position
func (plv *PodLogView) GetScrollPosition() int {
	return plv.scrollY
}
