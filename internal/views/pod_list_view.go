package views

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/kevholditch/vigilant/internal/models"
	"github.com/kevholditch/vigilant/internal/theme"
)

// PodListView represents the pod list view
type PodListView struct {
	pods        []models.Pod
	selected    int
	width       int
	height      int
	theme       *theme.Theme
	clusterName string
}

// NewPodListView creates a new pod list view
func NewPodListView(pods []models.Pod, theme *theme.Theme, clusterName string) *PodListView {
	return &PodListView{
		pods:        pods,
		selected:    0,
		theme:       theme,
		clusterName: clusterName,
	}
}

// SetSize sets the view dimensions
func (plv *PodListView) SetSize(width, height int) {
	plv.width = width
	plv.height = height
}

// SelectNext moves selection to next pod
func (plv *PodListView) SelectNext() {
	if plv.selected < len(plv.pods)-1 {
		plv.selected++
	}
}

// SelectPrev moves selection to previous pod
func (plv *PodListView) SelectPrev() {
	if plv.selected > 0 {
		plv.selected--
	}
}

// GetSelected returns the currently selected pod
func (plv *PodListView) GetSelected() *models.Pod {
	if len(plv.pods) == 0 {
		return nil
	}
	return &plv.pods[plv.selected]
}

// UpdatePods updates the pods data
func (plv *PodListView) UpdatePods(pods []models.Pod) {
	plv.pods = pods
	// Reset selection if current selection is out of bounds
	if plv.selected >= len(plv.pods) {
		plv.selected = 0
	}
}

// Render renders the complete pod list view
func (plv *PodListView) Render() string {
	if plv.width == 0 || plv.height == 0 {
		return ""
	}

	// Pod table
	table := plv.renderTable()

	// Status bar
	statusBar := plv.renderStatusBar()

	// Combine all components
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		table,
		statusBar,
	)

	return content
}

// renderTable renders the pod table
func (plv *PodListView) renderTable() string {
	if len(plv.pods) == 0 {
		return lipgloss.NewStyle().Foreground(plv.theme.TextMuted).Render("No pods found")
	}

	// Create table headers
	headers := []string{"NAME", "NAMESPACE", "STATUS", "READY", "RESTARTS", "AGE", "IP", "NODE"}

	// Create table rows
	var rows [][]string
	for _, pod := range plv.pods {
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
		BorderStyle(lipgloss.NewStyle().Foreground(plv.theme.Primary)).
		StyleFunc(func(row, col int) lipgloss.Style {

			isSelected := row == plv.selected

			var style lipgloss.Style
			if isSelected {
				style = plv.theme.TableSelectedStyle
			} else if (row-1)%2 == 1 {
				// Alternate row style
				style = plv.theme.TableRowAltStyle
			} else {
				style = plv.theme.TableRowStyle
			}

			// Apply status styling for status column (col 2)
			podIndex := row - 1
			if col == 2 && !isSelected && podIndex >= 0 && podIndex < len(plv.pods) {
				pod := plv.pods[podIndex]
				style = style.Inherit(plv.theme.GetStatusStyle(pod.Status))
			}

			return style
		})

	// available height for table is parent height - status bar height - table border - table header
	tableHeight := plv.height - 1 - 3 // 1 for status bar, 3 for table overhead(border+header)
	if tableHeight < 0 {
		tableHeight = 0
	}
	t.Height(tableHeight)

	return t.Render()
}

// renderStatusBar renders the status bar at the bottom
func (plv *PodListView) renderStatusBar() string {
	statusText := fmt.Sprintf("Total: %d pods | Press 'd' to describe | Press 'l' to view logs", len(plv.pods))
	return plv.theme.StatusBarStyle.Width(plv.width).Render(statusText)
}

// Pods returns the list of pods (for testing)
func (plv *PodListView) Pods() []models.Pod {
	return plv.pods
}
