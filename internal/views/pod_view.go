package views

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
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
func NewPodView(pods []models.Pod) *PodView {
	return &PodView{
		pods:     pods,
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

	// Create table headers
	headers := []string{"NAME", "NAMESPACE", "STATUS", "READY", "RESTARTS", "AGE", "IP", "NODE"}

	// Create table rows
	var rows [][]string
	for _, pod := range pv.pods {
		row := []string{
			pod.Name,
			pod.Namespace,
			pod.Status,
			pod.Ready,
			fmt.Sprintf("%d", pod.Restarts),
			pod.FormatAge(),
			pod.IP,
			pod.Node,
		}
		rows = append(rows, row)
	}

	// Create the table
	t := table.New().
		Headers(headers...).
		Rows(rows...).
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(theme.Primary)).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == 0 {
				// Header row
				return theme.TableHeaderStyle
			}

			// Data rows
			isSelected := row-1 == pv.selected
			isAlt := (row-1)%2 == 1

			var style lipgloss.Style
			if isSelected {
				style = theme.TableSelectedStyle
			} else if isAlt {
				style = theme.TableRowAltStyle
			} else {
				style = theme.TableRowStyle
			}

			// Apply status styling for status column (col 2)
			if col == 2 && !isSelected && row-1 >= 0 && row-1 < len(pv.pods) {
				pod := pv.pods[row-1]
				style = style.Inherit(theme.GetStatusStyle(pod.Status))
			}

			return style
		})

	return t.Render()
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
