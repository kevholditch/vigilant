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
	separator := lipgloss.NewStyle().
		Foreground(h.theme.TextMuted).
		Background(h.theme.BgSecondary).
		SetString(" | ").
		String()

	// --- Content ---
	clusterInfo := lipgloss.NewStyle().
		Foreground(h.theme.TextMuted).
		Background(h.theme.BgSecondary).
		SetString(fmt.Sprintf("☸️ %s", model.ClusterName)).String()

	k8sInfo := lipgloss.NewStyle().
		Foreground(h.theme.TextPrimary).
		Background(h.theme.BgSecondary).
		SetString(fmt.Sprintf("K8s: %s", model.KubernetesVersion)).String()

	controlPlaneInfo := lipgloss.NewStyle().
		Foreground(h.theme.TextPrimary).
		Background(h.theme.BgSecondary).
		SetString(fmt.Sprintf("🕹️ CP %d", model.ControlPlaneNodes)).String()

	workerInfo := lipgloss.NewStyle().
		Foreground(h.theme.TextPrimary).
		Background(h.theme.BgSecondary).
		SetString(fmt.Sprintf("👷 W %d", model.WorkerNodes)).String()

	viewTextStyled := lipgloss.NewStyle().
		Foreground(h.theme.TextPrimary).
		Background(h.theme.BgSecondary).
		SetString(viewText).
		String()

	content := lipgloss.JoinHorizontal(
		lipgloss.Bottom,
		clusterInfo,
		separator,
		viewTextStyled,
		separator,
		k8sInfo,
		separator,
		controlPlaneInfo,
		separator,
		workerInfo,
	)

	// --- Layout ---
	bar := lipgloss.NewStyle().
		Background(h.theme.BgSecondary).
		Foreground(h.theme.TextPrimary).
		Width(h.width).
		Align(lipgloss.Center).
		Padding(0, 1).
		Render(content)

	return bar
}
