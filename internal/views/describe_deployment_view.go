package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/kevholditch/vigilant/internal/models"
	"github.com/kevholditch/vigilant/internal/theme"
)

// DescribeDeploymentView represents the describe deployment view
type DescribeDeploymentView struct {
	deployment *models.Deployment
	theme      *theme.Theme
	width      int
	height     int
	scrollY    int
}

// NewDescribeDeploymentView creates a new describe deployment view
func NewDescribeDeploymentView(deployment *models.Deployment, theme *theme.Theme) *DescribeDeploymentView {
	return &DescribeDeploymentView{
		deployment: deployment,
		theme:      theme,
		scrollY:    0,
	}
}

// SetSize sets the view dimensions
func (ddv *DescribeDeploymentView) SetSize(width, height int) {
	ddv.width = width
	ddv.height = height
}

// ScrollUp scrolls the view up
func (ddv *DescribeDeploymentView) ScrollUp() {
	if ddv.scrollY > 0 {
		ddv.scrollY--
	}
}

// ScrollDown scrolls the view down
func (ddv *DescribeDeploymentView) ScrollDown() {
	ddv.scrollY++
}

// ScrollPageUp scrolls the view up by a page
func (ddv *DescribeDeploymentView) ScrollPageUp() {
	ddv.scrollY -= ddv.height / 2
	if ddv.scrollY < 0 {
		ddv.scrollY = 0
	}
}

// ScrollPageDown scrolls the view down by a page
func (ddv *DescribeDeploymentView) ScrollPageDown() {
	ddv.scrollY += ddv.height / 2
}

// ScrollToTop scrolls to the top of the view
func (ddv *DescribeDeploymentView) ScrollToTop() {
	ddv.scrollY = 0
}

// ScrollToBottom scrolls to the bottom of the view
func (ddv *DescribeDeploymentView) ScrollToBottom() {
	// This will be calculated in the render method
}

// UpdateDeployment updates the deployment data
func (ddv *DescribeDeploymentView) UpdateDeployment(deployment *models.Deployment) {
	ddv.deployment = deployment
}

// Render renders the describe deployment view
func (ddv *DescribeDeploymentView) Render() string {
	if ddv.width == 0 || ddv.height == 0 {
		return ""
	}

	content := ddv.renderContent()
	lines := strings.Split(content, "\n")

	// Calculate max scroll
	maxScroll := len(lines) - ddv.height
	if maxScroll < 0 {
		maxScroll = 0
	}

	// Clamp scroll position
	if ddv.scrollY > maxScroll {
		ddv.scrollY = maxScroll
	}

	// Get visible lines
	start := ddv.scrollY
	end := start + ddv.height
	if end > len(lines) {
		end = len(lines)
	}

	if start >= len(lines) {
		return lipgloss.NewStyle().Foreground(ddv.theme.TextMuted).Render("No content to display")
	}

	visibleLines := lines[start:end]
	return strings.Join(visibleLines, "\n")
}

// renderContent renders the full deployment description content
func (ddv *DescribeDeploymentView) renderContent() string {
	if ddv.deployment == nil {
		return lipgloss.NewStyle().Foreground(ddv.theme.Error).Render("No deployment data available")
	}

	d := ddv.deployment

	var sections []string

	// Basic information
	basicInfo := fmt.Sprintf(`Name:         %s
Namespace:    %s
Status:       %s
Age:          %s
Strategy:     %s`, d.Name, d.Namespace, d.Status, d.FormatAge(), d.Strategy)
	sections = append(sections, lipgloss.NewStyle().Foreground(ddv.theme.Primary).Bold(true).Render("Basic Information"), basicInfo)

	// Replica information
	replicaInfo := fmt.Sprintf(`Ready:        %s
Up-to-date:   %d
Available:    %d`, d.Ready, d.UpToDate, d.Available)
	sections = append(sections, lipgloss.NewStyle().Foreground(ddv.theme.Primary).Bold(true).Render("Replica Information"), replicaInfo)

	// Container information
	containerInfo := fmt.Sprintf(`Image:        %s`, d.Image)
	sections = append(sections, lipgloss.NewStyle().Foreground(ddv.theme.Primary).Bold(true).Render("Container Information"), containerInfo)

	// Status details
	statusDetails := ddv.renderStatusDetails(d)
	if statusDetails != "" {
		sections = append(sections, lipgloss.NewStyle().Foreground(ddv.theme.Primary).Bold(true).Render("Status Details"), statusDetails)
	}

	return strings.Join(sections, "\n\n")
}

// renderStatusDetails renders detailed status information
func (ddv *DescribeDeploymentView) renderStatusDetails(d *models.Deployment) string {
	var details []string

	switch d.Status {
	case "Ready":
		details = append(details, lipgloss.NewStyle().Foreground(ddv.theme.Success).Render("✓ All replicas are ready and available"))
	case "Available":
		details = append(details, lipgloss.NewStyle().Foreground(ddv.theme.Warning).Render("⚠ Some replicas are available but not all are ready"))
	case "Not Ready":
		details = append(details, lipgloss.NewStyle().Foreground(ddv.theme.Error).Render("✗ No replicas are ready or available"))
	case "Scaled to 0":
		details = append(details, lipgloss.NewStyle().Foreground(ddv.theme.TextMuted).Render("○ Deployment is scaled to 0 replicas"))
	default:
		details = append(details, lipgloss.NewStyle().Foreground(ddv.theme.TextMuted).Render("? Unknown status"))
	}

	return strings.Join(details, "\n")
}
