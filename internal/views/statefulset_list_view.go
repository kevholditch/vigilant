package views

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/kevholditch/vigilant/internal/models"
	"github.com/kevholditch/vigilant/internal/theme"
)

// StatefulSetListView represents the StatefulSet list view
type StatefulSetListView struct {
	statefulsets []models.StatefulSet
	selected     int
	width        int
	height       int
	theme        *theme.Theme
	clusterName  string
}

// NewStatefulSetListView creates a new StatefulSet list view
func NewStatefulSetListView(statefulsets []models.StatefulSet, theme *theme.Theme, clusterName string) *StatefulSetListView {
	return &StatefulSetListView{
		statefulsets: statefulsets,
		selected:     0,
		theme:        theme,
		clusterName:  clusterName,
	}
}

// SetSize sets the view dimensions
func (slv *StatefulSetListView) SetSize(width, height int) {
	slv.width = width
	slv.height = height
}

// SelectNext moves selection to next StatefulSet
func (slv *StatefulSetListView) SelectNext() {
	if slv.selected < len(slv.statefulsets)-1 {
		slv.selected++
	}
}

// SelectPrev moves selection to previous StatefulSet
func (slv *StatefulSetListView) SelectPrev() {
	if slv.selected > 0 {
		slv.selected--
	}
}

// GetSelected returns the currently selected StatefulSet
func (slv *StatefulSetListView) GetSelected() *models.StatefulSet {
	if len(slv.statefulsets) == 0 {
		return nil
	}
	return &slv.statefulsets[slv.selected]
}

// UpdateStatefulSets updates the StatefulSets data
func (slv *StatefulSetListView) UpdateStatefulSets(statefulsets []models.StatefulSet) {
	slv.statefulsets = statefulsets
	if slv.selected >= len(slv.statefulsets) {
		slv.selected = 0
	}
}

// Render renders the complete StatefulSet list view
func (slv *StatefulSetListView) Render() string {
	if slv.width == 0 || slv.height == 0 {
		return ""
	}
	table := slv.renderTable()
	statusBar := slv.renderStatusBar()
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		table,
		statusBar,
	)
	return content
}

// renderTable renders the StatefulSet table
func (slv *StatefulSetListView) renderTable() string {
	if len(slv.statefulsets) == 0 {
		return lipgloss.NewStyle().Foreground(slv.theme.TextMuted).Render("No StatefulSets found")
	}
	headers := []string{"NAME", "NAMESPACE", "STATUS", "READY", "REPLICAS", "AVAILABLE", "AGE", "IMAGE"}
	var rows [][]string
	for _, ss := range slv.statefulsets {
		row := []string{
			ss.Name,
			ss.Namespace,
			ss.Status,
			ss.Ready,
			fmt.Sprintf("%d", ss.Replicas),
			fmt.Sprintf("%d", ss.Available),
			ss.FormatAge(),
			ss.Image,
		}
		rows = append(rows, row)
	}
	t := table.New().
		Headers(headers...).
		Rows(rows...).
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(slv.theme.Primary)).
		StyleFunc(func(row, col int) lipgloss.Style {
			isSelected := row == slv.selected
			var style lipgloss.Style
			if isSelected {
				style = slv.theme.TableSelectedStyle
			} else if (row-1)%2 == 1 {
				style = slv.theme.TableRowAltStyle
			} else {
				style = slv.theme.TableRowStyle
			}
			if col == 2 && !isSelected && row-1 >= 0 && row-1 < len(slv.statefulsets) {
				ss := slv.statefulsets[row-1]
				style = style.Inherit(slv.theme.GetStatusStyle(ss.Status))
			}
			return style
		})
	tableHeight := slv.height - 1 - 3
	if tableHeight < 0 {
		tableHeight = 0
	}
	t.Height(tableHeight)
	return t.Render()
}

// renderStatusBar renders the status bar at the bottom
func (slv *StatefulSetListView) renderStatusBar() string {
	statusText := fmt.Sprintf("Total: %d StatefulSets | Press 'd' to describe", len(slv.statefulsets))
	return slv.theme.StatusBarStyle.Width(slv.width).Render(statusText)
}

// StatefulSets returns the list of StatefulSets (for testing)
func (slv *StatefulSetListView) StatefulSets() []models.StatefulSet {
	return slv.statefulsets
}
