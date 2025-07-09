package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/kevholditch/vigilant/internal/models"
	"github.com/kevholditch/vigilant/internal/theme"
)

// DescribeReplicaSetView represents the ReplicaSet description view
type DescribeReplicaSetView struct {
	replicaset *models.ReplicaSet
	theme      *theme.Theme
	width      int
	height     int
	scrollY    int
}

// NewDescribeReplicaSetView creates a new describe ReplicaSet view
func NewDescribeReplicaSetView(replicaset *models.ReplicaSet, theme *theme.Theme) *DescribeReplicaSetView {
	return &DescribeReplicaSetView{
		replicaset: replicaset,
		theme:      theme,
		scrollY:    0,
	}
}

// SetSize sets the view dimensions
func (drv *DescribeReplicaSetView) SetSize(width, height int) {
	drv.width = width
	drv.height = height
}

// UpdateReplicaSet updates the ReplicaSet data
func (drv *DescribeReplicaSetView) UpdateReplicaSet(rs *models.ReplicaSet) {
	drv.replicaset = rs
}

// ScrollUp scrolls the view up
func (drv *DescribeReplicaSetView) ScrollUp() {
	if drv.scrollY > 0 {
		drv.scrollY--
	}
}

// ScrollDown scrolls the view down
func (drv *DescribeReplicaSetView) ScrollDown() {
	drv.scrollY++
}

// ScrollPageUp scrolls the view up by a page
func (drv *DescribeReplicaSetView) ScrollPageUp() {
	drv.scrollY -= drv.height / 2
	if drv.scrollY < 0 {
		drv.scrollY = 0
	}
}

// ScrollPageDown scrolls the view down by a page
func (drv *DescribeReplicaSetView) ScrollPageDown() {
	drv.scrollY += drv.height / 2
}

// ScrollToTop scrolls to the top of the view
func (drv *DescribeReplicaSetView) ScrollToTop() {
	drv.scrollY = 0
}

// ScrollToBottom scrolls to the bottom of the view
func (drv *DescribeReplicaSetView) ScrollToBottom() {
	// This will be calculated in the render method
}

// Render renders the describe ReplicaSet view
func (drv *DescribeReplicaSetView) Render() string {
	if drv.width == 0 || drv.height == 0 {
		return ""
	}
	content := drv.renderContent()
	lines := strings.Split(content, "\n")
	maxScroll := len(lines) - drv.height
	if maxScroll < 0 {
		maxScroll = 0
	}
	if drv.scrollY > maxScroll {
		drv.scrollY = maxScroll
	}
	start := drv.scrollY
	end := start + drv.height
	if end > len(lines) {
		end = len(lines)
	}
	if start >= len(lines) {
		return lipgloss.NewStyle().Foreground(drv.theme.TextMuted).Render("No content to display")
	}
	visibleLines := lines[start:end]
	return strings.Join(visibleLines, "\n")
}

// renderContent renders the full ReplicaSet description content
func (drv *DescribeReplicaSetView) renderContent() string {
	if drv.replicaset == nil {
		return lipgloss.NewStyle().Foreground(drv.theme.Error).Render("No ReplicaSet data available")
	}
	rs := drv.replicaset
	var sections []string
	basicInfo := fmt.Sprintf(`Name:         %s
Namespace:    %s
Status:       %s
Age:          %s
Image:        %s`, rs.Name, rs.Namespace, rs.Status, rs.FormatAge(), rs.Image)
	sections = append(sections, lipgloss.NewStyle().Foreground(drv.theme.Primary).Bold(true).Render("Basic Information"), basicInfo)
	replicaInfo := fmt.Sprintf(`Ready:        %s
Replicas:     %d
Available:    %d`, rs.Ready, rs.Replicas, rs.Available)
	sections = append(sections, lipgloss.NewStyle().Foreground(drv.theme.Primary).Bold(true).Render("Replica Information"), replicaInfo)
	return strings.Join(sections, "\n\n")
}
