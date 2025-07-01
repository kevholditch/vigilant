package views

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/kevholditch/vigilant/internal/models"
	"github.com/kevholditch/vigilant/internal/theme"
)

// DeploymentListView represents the deployment list view
type DeploymentListView struct {
	deployments []models.Deployment
	selected    int
	width       int
	height      int
	theme       *theme.Theme
	clusterName string
}

// NewDeploymentListView creates a new deployment list view
func NewDeploymentListView(deployments []models.Deployment, theme *theme.Theme, clusterName string) *DeploymentListView {
	return &DeploymentListView{
		deployments: deployments,
		selected:    0,
		theme:       theme,
		clusterName: clusterName,
	}
}

// SetSize sets the view dimensions
func (dlv *DeploymentListView) SetSize(width, height int) {
	dlv.width = width
	dlv.height = height
}

// SelectNext moves selection to next deployment
func (dlv *DeploymentListView) SelectNext() {
	if dlv.selected < len(dlv.deployments)-1 {
		dlv.selected++
	}
}

// SelectPrev moves selection to previous deployment
func (dlv *DeploymentListView) SelectPrev() {
	if dlv.selected > 0 {
		dlv.selected--
	}
}

// GetSelected returns the currently selected deployment
func (dlv *DeploymentListView) GetSelected() *models.Deployment {
	if len(dlv.deployments) == 0 {
		return nil
	}
	return &dlv.deployments[dlv.selected]
}

// UpdateDeployments updates the deployments data
func (dlv *DeploymentListView) UpdateDeployments(deployments []models.Deployment) {
	dlv.deployments = deployments
	// Reset selection if current selection is out of bounds
	if dlv.selected >= len(dlv.deployments) {
		dlv.selected = 0
	}
}

// Render renders the complete deployment list view
func (dlv *DeploymentListView) Render() string {
	if dlv.width == 0 || dlv.height == 0 {
		return ""
	}

	// Deployment table
	table := dlv.renderTable()

	// Status bar
	statusBar := dlv.renderStatusBar()

	// Combine all components
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		table,
		statusBar,
	)

	return content
}

// renderTable renders the deployment table
func (dlv *DeploymentListView) renderTable() string {
	if len(dlv.deployments) == 0 {
		return lipgloss.NewStyle().Foreground(dlv.theme.TextMuted).Render("No deployments found")
	}

	// Create table headers
	headers := []string{"NAME", "NAMESPACE", "STATUS", "READY", "UP-TO-DATE", "AVAILABLE", "AGE", "STRATEGY", "IMAGE"}

	// Create table rows
	var rows [][]string
	for _, deployment := range dlv.deployments {
		row := []string{
			deployment.Name,
			deployment.Namespace,
			deployment.Status,
			deployment.Ready,
			fmt.Sprintf("%d", deployment.UpToDate),
			fmt.Sprintf("%d", deployment.Available),
			deployment.FormatAge(),
			deployment.Strategy,
			deployment.Image,
		}
		rows = append(rows, row)
	}

	// Create the table
	t := table.New().
		Headers(headers...).
		Rows(rows...).
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(dlv.theme.Primary)).
		StyleFunc(func(row, col int) lipgloss.Style {

			isSelected := row == dlv.selected

			var style lipgloss.Style
			if isSelected {
				style = dlv.theme.TableSelectedStyle
			} else if (row-1)%2 == 1 {
				// Alternate row style
				style = dlv.theme.TableRowAltStyle
			} else {
				style = dlv.theme.TableRowStyle
			}

			// Apply status styling for status column (col 2)
			deploymentIndex := row - 1
			if col == 2 && !isSelected && deploymentIndex >= 0 && deploymentIndex < len(dlv.deployments) {
				deployment := dlv.deployments[deploymentIndex]
				style = style.Inherit(dlv.theme.GetStatusStyle(deployment.Status))
			}

			return style
		})

	// available height for table is parent height - status bar height - table border - table header
	tableHeight := dlv.height - 1 - 3 // 1 for status bar, 3 for table overhead(border+header)
	if tableHeight < 0 {
		tableHeight = 0
	}
	t.Height(tableHeight)

	return t.Render()
}

// renderStatusBar renders the status bar at the bottom
func (dlv *DeploymentListView) renderStatusBar() string {
	statusText := fmt.Sprintf("Total: %d deployments | Press 'd' to describe", len(dlv.deployments))
	return dlv.theme.StatusBarStyle.Width(dlv.width).Render(statusText)
}

// Deployments returns the list of deployments (for testing)
func (dlv *DeploymentListView) Deployments() []models.Deployment {
	return dlv.deployments
}
