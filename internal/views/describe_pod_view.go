package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/kevholditch/vigilant/internal/models"
	"github.com/kevholditch/vigilant/internal/theme"
)

// DescribePodView represents the pod description view
type DescribePodView struct {
	pod     *models.Pod
	theme   *theme.Theme
	width   int
	height  int
	scrollY int
}

// NewDescribePodView creates a new describe pod view
func NewDescribePodView(pod *models.Pod, theme *theme.Theme) *DescribePodView {
	return &DescribePodView{
		pod:     pod,
		theme:   theme,
		scrollY: 0,
	}
}

// SetSize sets the view dimensions
func (dpv *DescribePodView) SetSize(width, height int) {
	dpv.width = width
	dpv.height = height
}

// UpdatePod updates the pod data
func (dpv *DescribePodView) UpdatePod(pod *models.Pod) {
	dpv.pod = pod
}

// ScrollUp scrolls the view up
func (dpv *DescribePodView) ScrollUp() {
	if dpv.scrollY > 0 {
		dpv.scrollY--
	}
}

// ScrollDown scrolls the view down
func (dpv *DescribePodView) ScrollDown() {
	dpv.scrollY++
}

// ScrollPageUp scrolls the view up by a page
func (dpv *DescribePodView) ScrollPageUp() {
	dpv.scrollY -= dpv.height / 2
	if dpv.scrollY < 0 {
		dpv.scrollY = 0
	}
}

// ScrollPageDown scrolls the view down by a page
func (dpv *DescribePodView) ScrollPageDown() {
	dpv.scrollY += dpv.height / 2
}

// ScrollToTop scrolls to the top of the view
func (dpv *DescribePodView) ScrollToTop() {
	dpv.scrollY = 0
}

// ScrollToBottom scrolls to the bottom of the view
func (dpv *DescribePodView) ScrollToBottom() {
	// This will be calculated in the render method
}

// Render renders the describe pod view
func (dpv *DescribePodView) Render() string {
	if dpv.width == 0 || dpv.height == 0 {
		return ""
	}

	content := dpv.renderContent()
	lines := strings.Split(content, "\n")

	// Calculate max scroll
	maxScroll := len(lines) - dpv.height
	if maxScroll < 0 {
		maxScroll = 0
	}

	// Clamp scroll position
	if dpv.scrollY > maxScroll {
		dpv.scrollY = maxScroll
	}

	// Get visible lines
	start := dpv.scrollY
	end := start + dpv.height
	if end > len(lines) {
		end = len(lines)
	}

	if start >= len(lines) {
		return lipgloss.NewStyle().Foreground(dpv.theme.TextMuted).Render("No content to display")
	}

	visibleLines := lines[start:end]
	return strings.Join(visibleLines, "\n")
}

// renderContent renders the full pod description content
func (dpv *DescribePodView) renderContent() string {
	if dpv.pod == nil {
		return lipgloss.NewStyle().Foreground(dpv.theme.Error).Render("No pod data available")
	}

	p := dpv.pod

	var sections []string

	// Basic information
	basicInfo := fmt.Sprintf(`Name:         %s
Namespace:    %s
Status:       %s
Age:          %s
IP:           %s
Node:         %s`, p.Name, p.Namespace, p.Status, p.FormatAge(), p.IP, p.Node)
	sections = append(sections, lipgloss.NewStyle().Foreground(dpv.theme.Primary).Bold(true).Render("Basic Information"), basicInfo)

	// Container information
	containerInfo := fmt.Sprintf(`Ready:        %s
Restarts:     %d`, p.Ready, p.Restarts)
	sections = append(sections, lipgloss.NewStyle().Foreground(dpv.theme.Primary).Bold(true).Render("Container Information"), containerInfo)

	// Status details
	statusDetails := dpv.renderStatusDetails(p)
	if statusDetails != "" {
		sections = append(sections, lipgloss.NewStyle().Foreground(dpv.theme.Primary).Bold(true).Render("Status Details"), statusDetails)
	}

	// Network information
	networkInfo := fmt.Sprintf(`Pod IP:       %s
Node:         %s`, p.IP, p.Node)
	sections = append(sections, lipgloss.NewStyle().Foreground(dpv.theme.Primary).Bold(true).Render("Network Information"), networkInfo)

	return strings.Join(sections, "\n\n")
}

// renderStatusDetails renders detailed status information
func (dpv *DescribePodView) renderStatusDetails(p *models.Pod) string {
	var details []string

	switch p.Status {
	case "Running":
		details = append(details, lipgloss.NewStyle().Foreground(dpv.theme.Success).Render("✓ Pod is running and healthy"))
	case "Pending":
		details = append(details, lipgloss.NewStyle().Foreground(dpv.theme.Warning).Render("⏳ Pod is pending - waiting for resources or scheduling"))
	case "Succeeded":
		details = append(details, lipgloss.NewStyle().Foreground(dpv.theme.Accent).Render("✓ Pod completed successfully"))
	case "Failed":
		details = append(details, lipgloss.NewStyle().Foreground(dpv.theme.Error).Render("✗ Pod failed to run"))
	case "Unknown":
		details = append(details, lipgloss.NewStyle().Foreground(dpv.theme.TextMuted).Render("? Pod status is unknown"))
	default:
		details = append(details, lipgloss.NewStyle().Foreground(dpv.theme.TextMuted).Render("? Unknown status"))
	}

	// Add restart information if applicable
	if p.Restarts > 0 {
		restartInfo := fmt.Sprintf("🔄 Pod has restarted %d times", p.Restarts)
		if p.Restarts > 5 {
			details = append(details, lipgloss.NewStyle().Foreground(dpv.theme.Warning).Render(restartInfo))
		} else {
			details = append(details, lipgloss.NewStyle().Foreground(dpv.theme.TextSecondary).Render(restartInfo))
		}
	}

	return strings.Join(details, "\n")
}
