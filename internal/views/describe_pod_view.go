package views

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/kevholditch/vigilant/internal/models"
	"github.com/kevholditch/vigilant/internal/theme"
)

// DescribePodView represents the pod description view
type DescribePodView struct {
	pod     *models.Pod
	content string
	scrollY int
	width   int
	height  int
	theme   *theme.Theme
}

// NewDescribePodView creates a new describe pod view
func NewDescribePodView(pod *models.Pod, theme *theme.Theme) *DescribePodView {
	view := &DescribePodView{
		pod:   pod,
		theme: theme,
	}
	view.loadPodDescription()
	return view
}

// SetSize sets the view dimensions
func (dpv *DescribePodView) SetSize(width, height int) {
	dpv.width = width
	dpv.height = height
}

// UpdatePod updates the pod data and reloads the description
func (dpv *DescribePodView) UpdatePod(pod *models.Pod) {
	dpv.pod = pod
	dpv.scrollY = 0 // Reset scroll position
	dpv.loadPodDescription()
}

// ScrollUp moves the view up
func (dpv *DescribePodView) ScrollUp() {
	if dpv.scrollY > 0 {
		dpv.scrollY--
	}
}

// ScrollDown moves the view down
func (dpv *DescribePodView) ScrollDown() {
	// Calculate max scroll based on content height
	contentLines := strings.Split(dpv.content, "\n")
	availableHeight := dpv.height - 4 // Header + status bar + borders

	if availableHeight <= 0 {
		return
	}

	maxScroll := len(contentLines) - availableHeight
	if maxScroll < 0 {
		maxScroll = 0
	}

	if dpv.scrollY < maxScroll {
		dpv.scrollY++
	}
}

// ScrollPageUp moves the view up by a page
func (dpv *DescribePodView) ScrollPageUp() {
	availableHeight := dpv.height - 4 // Header + status bar + borders
	if availableHeight <= 0 {
		return
	}

	// Move up by available height (one page)
	dpv.scrollY -= availableHeight
	if dpv.scrollY < 0 {
		dpv.scrollY = 0
	}
}

// ScrollPageDown moves the view down by a page
func (dpv *DescribePodView) ScrollPageDown() {
	// Calculate max scroll based on content height
	contentLines := strings.Split(dpv.content, "\n")
	availableHeight := dpv.height - 4 // Header + status bar + borders

	if availableHeight <= 0 {
		return
	}

	maxScroll := len(contentLines) - availableHeight
	if maxScroll < 0 {
		maxScroll = 0
	}

	// Move down by available height (one page)
	dpv.scrollY += availableHeight
	if dpv.scrollY > maxScroll {
		dpv.scrollY = maxScroll
	}
}

// ScrollToTop moves the view to the top
func (dpv *DescribePodView) ScrollToTop() {
	dpv.scrollY = 0
}

// ScrollToBottom moves the view to the bottom
func (dpv *DescribePodView) ScrollToBottom() {
	// Calculate max scroll based on content height
	contentLines := strings.Split(dpv.content, "\n")
	availableHeight := dpv.height - 4 // Header + status bar + borders

	if availableHeight <= 0 {
		return
	}

	maxScroll := len(contentLines) - availableHeight
	if maxScroll < 0 {
		maxScroll = 0
	}

	dpv.scrollY = maxScroll
}

// loadPodDescription loads the pod description using kubectl
func (dpv *DescribePodView) loadPodDescription() {
	cmd := exec.Command("kubectl", "describe", "pod", dpv.pod.Name, "-n", dpv.pod.Namespace)
	output, err := cmd.Output()
	if err != nil {
		dpv.content = fmt.Sprintf("Error getting pod description: %v", err)
		return
	}
	dpv.content = string(output)
}

// Render renders the complete describe pod view
func (dpv *DescribePodView) Render() string {
	if dpv.width == 0 || dpv.height == 0 {
		return ""
	}

	// Content area
	content := dpv.renderContent()

	// Status bar
	statusBar := dpv.renderStatusBar()

	// Combine all components
	viewContent := lipgloss.JoinVertical(
		lipgloss.Left,
		content,
		statusBar,
	)

	return viewContent
}

// renderContent renders the scrollable content area
func (dpv *DescribePodView) renderContent() string {
	// Split content into lines
	lines := strings.Split(dpv.content, "\n")

	// Calculate available height for content (subtract status bar)
	availableHeight := dpv.height - 1 // 1 for status bar
	if availableHeight < 0 {
		availableHeight = 0
	}

	if availableHeight <= 0 {
		return lipgloss.NewStyle().Foreground(dpv.theme.TextMuted).Render("Window too small")
	}

	// Ensure scroll position is valid
	if dpv.scrollY < 0 {
		dpv.scrollY = 0
	}

	maxScroll := len(lines) - availableHeight
	if maxScroll < 0 {
		maxScroll = 0
	}

	if dpv.scrollY > maxScroll {
		dpv.scrollY = maxScroll
	}

	// Apply scrolling
	startLine := dpv.scrollY
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
		Foreground(dpv.theme.TextPrimary).
		Background(dpv.theme.BgPrimary).
		Width(dpv.width).
		Height(availableHeight)

	return contentStyle.Render(content)
}

// renderStatusBar renders the status bar at the bottom
func (dpv *DescribePodView) renderStatusBar() string {
	statusText := fmt.Sprintf("Pod: %s | Esc: return | ↑/↓: scroll | PgUp/PgDn/Ctrl+u/Ctrl+d: page | g/G: top/bottom", dpv.pod.Name)
	return dpv.theme.StatusBarStyle.Width(dpv.width).Render(statusText)
}
