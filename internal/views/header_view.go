package views

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/kevholditch/vigilant/internal/models"
	"github.com/kevholditch/vigilant/internal/theme"
)

// HeaderView represents the header view component
type HeaderView struct {
	theme *theme.Theme
	width int
}

// NewHeaderView creates a new header view
func NewHeaderView(theme *theme.Theme) *HeaderView {
	return &HeaderView{
		theme: theme,
	}
}

// SetSize sets the width of the header view
func (h *HeaderView) SetSize(width int) {
	h.width = width
}

// Render renders the header view using the HeaderModel and viewText
func (h *HeaderView) Render(model *models.HeaderModel, viewText string) string {
	// --- Styles ---
	barStyle := lipgloss.NewStyle().
		Background(h.theme.BgSecondary).
		Foreground(h.theme.TextPrimary).
		Padding(0, 1)

	separator := lipgloss.NewStyle().
		Foreground(h.theme.TextMuted).
		SetString(" | ").
		String()

	// --- Content ---
	clusterInfo := fmt.Sprintf("☸️ %s", model.ClusterName)
	controlPlaneInfo := fmt.Sprintf("🕹️ CP %d", model.ControlPlaneNodes)
	workerInfo := fmt.Sprintf("👷 W %d", model.WorkerNodes)
	k8sInfo := fmt.Sprintf("K8s: %s", model.KubernetesVersion)

	// --- Assembly ---
	content := lipgloss.JoinHorizontal(
		lipgloss.Bottom,
		clusterInfo,
		separator,
		viewText,
		separator,
		k8sInfo,
		separator,
		controlPlaneInfo,
		separator,
		workerInfo,
	)

	// --- Layout ---
	bar := lipgloss.NewStyle().
		Width(h.width).
		Align(lipgloss.Center).
		Render(barStyle.Render(content))

	return bar
}
