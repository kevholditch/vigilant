package theme

import "github.com/charmbracelet/lipgloss"

// Colors - Cyberpunk theme
var (
	// Primary colors
	Primary   = lipgloss.Color("#00ff88") // Bright cyan-green
	Secondary = lipgloss.Color("#ff0088") // Hot pink
	Accent    = lipgloss.Color("#ffaa00") // Orange
	Warning   = lipgloss.Color("#ff4400") // Red-orange
	Error     = lipgloss.Color("#ff0044") // Red
	Success   = lipgloss.Color("#00ff44") // Green

	// Background colors
	BgPrimary   = lipgloss.Color("#0a0a0a") // Dark background
	BgSecondary = lipgloss.Color("#1a1a1a") // Slightly lighter
	BgTertiary  = lipgloss.Color("#2a2a2a") // Even lighter

	// Text colors
	TextPrimary   = lipgloss.Color("#ffffff") // White
	TextSecondary = lipgloss.Color("#cccccc") // Light gray
	TextMuted     = lipgloss.Color("#888888") // Gray
	TextInverse   = lipgloss.Color("#000000") // Black
)

// Styles
var (
	// Header styles
	HeaderStyle = lipgloss.NewStyle().
			Background(Primary).
			Foreground(TextInverse).
			Bold(true).
			Padding(0, 1)

	ClusterNameStyle = lipgloss.NewStyle().
				Foreground(Primary).
				Bold(true).
				MarginRight(2)

	ClusterVersionStyle = lipgloss.NewStyle().
				Foreground(TextSecondary).
				Italic(true)

	// Table styles
	TableHeaderStyle = lipgloss.NewStyle().
				Background(BgSecondary).
				Foreground(Primary).
				Bold(true).
				Padding(0, 1)

	TableRowStyle = lipgloss.NewStyle().
			Foreground(TextPrimary).
			Padding(0, 1)

	TableRowAltStyle = lipgloss.NewStyle().
				Background(BgSecondary).
				Foreground(TextPrimary).
				Padding(0, 1)

	TableSelectedStyle = lipgloss.NewStyle().
				Background(Primary).
				Foreground(TextInverse).
				Bold(true).
				Padding(0, 1)

	// Status styles
	StatusRunningStyle = lipgloss.NewStyle().
				Foreground(Success).
				Bold(true)

	StatusPendingStyle = lipgloss.NewStyle().
				Foreground(Warning).
				Bold(true)

	StatusFailedStyle = lipgloss.NewStyle().
				Foreground(Error).
				Bold(true)

	StatusSucceededStyle = lipgloss.NewStyle().
				Foreground(Accent).
				Bold(true)

	// Border styles
	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary).
			Padding(0, 1)

	// Status bar
	StatusBarStyle = lipgloss.NewStyle().
			Background(BgSecondary).
			Foreground(TextSecondary).
			Padding(0, 1)
)

// GetStatusStyle returns the appropriate style for a pod status
func GetStatusStyle(status string) lipgloss.Style {
	switch status {
	case "Running":
		return StatusRunningStyle
	case "Pending":
		return StatusPendingStyle
	case "Failed":
		return StatusFailedStyle
	case "Succeeded":
		return StatusSucceededStyle
	default:
		return lipgloss.NewStyle().Foreground(TextMuted)
	}
}
