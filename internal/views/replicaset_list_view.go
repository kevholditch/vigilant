package views

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/kevholditch/vigilant/internal/models"
	"github.com/kevholditch/vigilant/internal/theme"
)

// ReplicaSetListView represents the ReplicaSet list view
type ReplicaSetListView struct {
	replicasets []models.ReplicaSet
	selected    int
	width       int
	height      int
	theme       *theme.Theme
	clusterName string
}

// NewReplicaSetListView creates a new ReplicaSet list view
func NewReplicaSetListView(replicasets []models.ReplicaSet, theme *theme.Theme, clusterName string) *ReplicaSetListView {
	return &ReplicaSetListView{
		replicasets: replicasets,
		selected:    0,
		theme:       theme,
		clusterName: clusterName,
	}
}

// SetSize sets the view dimensions
func (rlv *ReplicaSetListView) SetSize(width, height int) {
	rlv.width = width
	rlv.height = height
}

// SelectNext moves selection to next ReplicaSet
func (rlv *ReplicaSetListView) SelectNext() {
	if rlv.selected < len(rlv.replicasets)-1 {
		rlv.selected++
	}
}

// SelectPrev moves selection to previous ReplicaSet
func (rlv *ReplicaSetListView) SelectPrev() {
	if rlv.selected > 0 {
		rlv.selected--
	}
}

// GetSelected returns the currently selected ReplicaSet
func (rlv *ReplicaSetListView) GetSelected() *models.ReplicaSet {
	if len(rlv.replicasets) == 0 {
		return nil
	}
	return &rlv.replicasets[rlv.selected]
}

// UpdateReplicaSets updates the ReplicaSets data
func (rlv *ReplicaSetListView) UpdateReplicaSets(replicasets []models.ReplicaSet) {
	rlv.replicasets = replicasets
	if rlv.selected >= len(rlv.replicasets) {
		rlv.selected = 0
	}
}

// Render renders the complete ReplicaSet list view
func (rlv *ReplicaSetListView) Render() string {
	if rlv.width == 0 || rlv.height == 0 {
		return ""
	}
	table := rlv.renderTable()
	statusBar := rlv.renderStatusBar()
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		table,
		statusBar,
	)
	return content
}

// renderTable renders the ReplicaSet table
func (rlv *ReplicaSetListView) renderTable() string {
	if len(rlv.replicasets) == 0 {
		return lipgloss.NewStyle().Foreground(rlv.theme.TextMuted).Render("No ReplicaSets found")
	}
	headers := []string{"NAME", "NAMESPACE", "STATUS", "READY", "REPLICAS", "AVAILABLE", "AGE", "IMAGE"}
	var rows [][]string
	for _, rs := range rlv.replicasets {
		row := []string{
			rs.Name,
			rs.Namespace,
			rs.Status,
			rs.Ready,
			fmt.Sprintf("%d", rs.Replicas),
			fmt.Sprintf("%d", rs.Available),
			rs.FormatAge(),
			rs.Image,
		}
		rows = append(rows, row)
	}
	t := table.New().
		Headers(headers...).
		Rows(rows...).
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(rlv.theme.Primary)).
		StyleFunc(func(row, col int) lipgloss.Style {
			isSelected := row == rlv.selected
			var style lipgloss.Style
			if isSelected {
				style = rlv.theme.TableSelectedStyle
			} else if (row-1)%2 == 1 {
				style = rlv.theme.TableRowAltStyle
			} else {
				style = rlv.theme.TableRowStyle
			}
			if col == 2 && !isSelected && row-1 >= 0 && row-1 < len(rlv.replicasets) {
				rs := rlv.replicasets[row-1]
				style = style.Inherit(rlv.theme.GetStatusStyle(rs.Status))
			}
			return style
		})
	tableHeight := rlv.height - 1 - 3
	if tableHeight < 0 {
		tableHeight = 0
	}
	t.Height(tableHeight)
	return t.Render()
}

// renderStatusBar renders the status bar at the bottom
func (rlv *ReplicaSetListView) renderStatusBar() string {
	statusText := fmt.Sprintf("Total: %d ReplicaSets | Press 'd' to describe", len(rlv.replicasets))
	return rlv.theme.StatusBarStyle.Width(rlv.width).Render(statusText)
}

// ReplicaSets returns the list of ReplicaSets (for testing)
func (rlv *ReplicaSetListView) ReplicaSets() []models.ReplicaSet {
	return rlv.replicasets
}
