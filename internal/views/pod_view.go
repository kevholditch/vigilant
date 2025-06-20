package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/kevholditch/vigilant/internal/models"
	"github.com/kevholditch/vigilant/internal/theme"
)

// PodView represents the pod list view
type PodView struct {
	pods     []models.Pod
	selected int
	width    int
	height   int
}

// NewPodView creates a new pod view
func NewPodView() *PodView {
	return &PodView{
		pods:     models.GetSamplePods(),
		selected: 0,
	}
}

// SetSize sets the view dimensions
func (pv *PodView) SetSize(width, height int) {
	pv.width = width
	pv.height = height
}

// SelectNext moves selection to next pod
func (pv *PodView) SelectNext() {
	if pv.selected < len(pv.pods)-1 {
		pv.selected++
	}
}

// SelectPrev moves selection to previous pod
func (pv *PodView) SelectPrev() {
	if pv.selected > 0 {
		pv.selected--
	}
}

// GetSelected returns the currently selected pod
func (pv *PodView) GetSelected() *models.Pod {
	if len(pv.pods) == 0 {
		return nil
	}
	return &pv.pods[pv.selected]
}

// Render renders the complete pod view
func (pv *PodView) Render() string {
	if pv.width == 0 || pv.height == 0 {
		return ""
	}

	// Header with cluster info
	header := pv.renderHeader()

	// Pod table
	table := pv.renderTable()

	// Status bar
	statusBar := pv.renderStatusBar()

	// Combine all components
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		table,
		statusBar,
	)

	// Apply border and padding
	return theme.BorderStyle.Width(pv.width).Height(pv.height).Render(content)
}

// renderHeader renders the cluster information header
func (pv *PodView) renderHeader() string {
	clusterName := theme.ClusterNameStyle.Render("🚀 production-cluster")
	clusterVersion := theme.ClusterVersionStyle.Render("v1.28.0")

	headerContent := lipgloss.JoinHorizontal(
		lipgloss.Left,
		clusterName,
		clusterVersion,
	)

	return theme.HeaderStyle.Width(pv.width).Render(headerContent)
}

// renderTable renders the pod table
func (pv *PodView) renderTable() string {
	if len(pv.pods) == 0 {
		return lipgloss.NewStyle().Foreground(theme.TextMuted).Render("No pods found")
	}

	// Table headers
	headers := []string{"NAME", "NAMESPACE", "STATUS", "READY", "RESTARTS", "AGE", "IP", "NODE"}
	headerRow := pv.renderTableRow(headers, true, false)

	// Table rows
	var rows []string
	for i, pod := range pv.pods {
		isSelected := i == pv.selected
		isAlt := i%2 == 1
		row := pv.renderPodRow(pod, isSelected, isAlt)
		rows = append(rows, row)
	}

	// Combine header and rows
	tableContent := append([]string{headerRow}, rows...)
	return strings.Join(tableContent, "\n")
}

// renderTableRow renders a single table row
func (pv *PodView) renderTableRow(cells []string, isHeader, isSelected bool) string {
	var styledCells []string

	for i, cell := range cells {
		var style lipgloss.Style

		if isHeader {
			style = theme.TableHeaderStyle
		} else if isSelected {
			style = theme.TableSelectedStyle
		} else {
			style = theme.TableRowStyle
		}

		// Adjust column widths
		width := pv.getColumnWidth(i)
		styledCell := style.Copy().Height(1).Width(width).Render(cell)
		styledCells = append(styledCells, styledCell)
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, styledCells...)
}

// renderPodRow renders a pod row with proper styling
func (pv *PodView) renderPodRow(pod models.Pod, isSelected, isAlt bool) string {
	cells := []string{
		pod.Name,
		pod.Namespace,
		pod.Status,
		pod.Ready,
		fmt.Sprintf("%d", pod.Restarts),
		pod.FormatAge(),
		pod.IP,
		pod.Node,
	}

	var styledCells []string

	// Determine base style for the row
	rowStyle := theme.TableRowStyle
	if isSelected {
		rowStyle = theme.TableSelectedStyle
	} else if isAlt {
		rowStyle = theme.TableRowAltStyle
	}

	for i, cell := range cells {
		style := rowStyle.Copy()
		width := pv.getColumnWidth(i)

		// Apply special styling for the status column on unselected rows.
		if i == 2 && !isSelected {
			style = style.Inherit(theme.GetStatusStyle(pod.Status))
		}

		styledCell := style.Height(1).Width(width).Render(cell)
		styledCells = append(styledCells, styledCell)
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, styledCells...)
}

// renderStatusBar renders the status bar at the bottom
func (pv *PodView) renderStatusBar() string {
	selectedPod := pv.GetSelected()
	var statusText string

	if selectedPod != nil {
		statusText = fmt.Sprintf("Selected: %s (%s) | Total: %d pods",
			selectedPod.Name, selectedPod.Status, len(pv.pods))
	} else {
		statusText = fmt.Sprintf("Total: %d pods", len(pv.pods))
	}

	return theme.StatusBarStyle.Width(pv.width).Render(statusText)
}

// getColumnWidth returns the width for a specific column
func (pv *PodView) getColumnWidth(colIndex int) int {
	// Define column widths as percentages of total width
	widths := []int{25, 12, 10, 8, 8, 8, 15, 12} // percentages

	if colIndex >= len(widths) {
		return 10
	}

	// Calculate actual width based on percentage, accounting for border/padding
	availableWidth := pv.width - 4
	width := (availableWidth * widths[colIndex]) / 100
	if width < 3 {
		width = 3 // minimum width
	}
	return width
}
