package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/kevholditch/vigilant/internal/models"
	"github.com/kevholditch/vigilant/internal/theme"
)

// DescribeStatefulSetView represents the StatefulSet description view
type DescribeStatefulSetView struct {
	statefulset *models.StatefulSet
	theme       *theme.Theme
	width       int
	height      int
	scrollY     int
}

// NewDescribeStatefulSetView creates a new describe StatefulSet view
func NewDescribeStatefulSetView(statefulset *models.StatefulSet, theme *theme.Theme) *DescribeStatefulSetView {
	return &DescribeStatefulSetView{
		statefulset: statefulset,
		theme:       theme,
		scrollY:     0,
	}
}

// SetSize sets the view dimensions
func (dsv *DescribeStatefulSetView) SetSize(width, height int) {
	dsv.width = width
	dsv.height = height
}

// UpdateStatefulSet updates the StatefulSet data
func (dsv *DescribeStatefulSetView) UpdateStatefulSet(ss *models.StatefulSet) {
	dsv.statefulset = ss
}

// ScrollUp scrolls the view up
func (dsv *DescribeStatefulSetView) ScrollUp() {
	if dsv.scrollY > 0 {
		dsv.scrollY--
	}
}

// ScrollDown scrolls the view down
func (dsv *DescribeStatefulSetView) ScrollDown() {
	dsv.scrollY++
}

// ScrollPageUp scrolls the view up by a page
func (dsv *DescribeStatefulSetView) ScrollPageUp() {
	dsv.scrollY -= dsv.height / 2
	if dsv.scrollY < 0 {
		dsv.scrollY = 0
	}
}

// ScrollPageDown scrolls the view down by a page
func (dsv *DescribeStatefulSetView) ScrollPageDown() {
	dsv.scrollY += dsv.height / 2
}

// ScrollToTop scrolls to the top of the view
func (dsv *DescribeStatefulSetView) ScrollToTop() {
	dsv.scrollY = 0
}

// ScrollToBottom scrolls to the bottom of the view
func (dsv *DescribeStatefulSetView) ScrollToBottom() {
	// This will be calculated in the render method
}

// Render renders the describe StatefulSet view
func (dsv *DescribeStatefulSetView) Render() string {
	if dsv.width == 0 || dsv.height == 0 {
		return ""
	}
	content := dsv.renderContent()
	lines := strings.Split(content, "\n")
	maxScroll := len(lines) - dsv.height
	if maxScroll < 0 {
		maxScroll = 0
	}
	if dsv.scrollY > maxScroll {
		dsv.scrollY = maxScroll
	}
	start := dsv.scrollY
	end := start + dsv.height
	if end > len(lines) {
		end = len(lines)
	}
	if start >= len(lines) {
		return lipgloss.NewStyle().Foreground(dsv.theme.TextMuted).Render("No content to display")
	}
	visibleLines := lines[start:end]
	return strings.Join(visibleLines, "\n")
}

// renderContent renders the full StatefulSet description content
func (dsv *DescribeStatefulSetView) renderContent() string {
	if dsv.statefulset == nil {
		return lipgloss.NewStyle().Foreground(dsv.theme.Error).Render("No StatefulSet data available")
	}
	ss := dsv.statefulset
	var sections []string
	basicInfo := fmt.Sprintf(`Name:         %s
Namespace:    %s
Status:       %s
Age:          %s
Image:        %s`, ss.Name, ss.Namespace, ss.Status, ss.FormatAge(), ss.Image)
	sections = append(sections, lipgloss.NewStyle().Foreground(dsv.theme.Primary).Bold(true).Render("Basic Information"), basicInfo)
	replicaInfo := fmt.Sprintf(`Ready:        %s
Replicas:     %d
Available:    %d`, ss.Ready, ss.Replicas, ss.Available)
	sections = append(sections, lipgloss.NewStyle().Foreground(dsv.theme.Primary).Bold(true).Render("Replica Information"), replicaInfo)
	return strings.Join(sections, "\n\n")
}
