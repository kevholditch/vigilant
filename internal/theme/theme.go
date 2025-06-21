package theme

import "github.com/charmbracelet/lipgloss"

// Theme represents a complete UI theme with all colors and styles
type Theme struct {
	// Colors
	Primary   lipgloss.Color
	Secondary lipgloss.Color
	Accent    lipgloss.Color
	Warning   lipgloss.Color
	Error     lipgloss.Color
	Success   lipgloss.Color
	Purple    lipgloss.Color

	// Background colors
	BgPrimary   lipgloss.Color
	BgSecondary lipgloss.Color
	BgTertiary  lipgloss.Color

	// Text colors
	TextPrimary   lipgloss.Color
	TextSecondary lipgloss.Color
	TextMuted     lipgloss.Color
	TextInverse   lipgloss.Color

	// Styles
	HeaderStyle          lipgloss.Style
	ClusterNameStyle     lipgloss.Style
	ClusterVersionStyle  lipgloss.Style
	TableHeaderStyle     lipgloss.Style
	TableRowStyle        lipgloss.Style
	TableRowAltStyle     lipgloss.Style
	TableSelectedStyle   lipgloss.Style
	StatusRunningStyle   lipgloss.Style
	StatusPendingStyle   lipgloss.Style
	StatusFailedStyle    lipgloss.Style
	StatusSucceededStyle lipgloss.Style
	BorderStyle          lipgloss.Style
	StatusBarStyle       lipgloss.Style
}

// NewDefaultTheme creates a new theme with the default cyberpunk colors
func NewDefaultTheme() *Theme {
	theme := &Theme{
		// Colors - Cyberpunk theme
		Primary:   lipgloss.Color("#00ff88"), // Bright cyan-green
		Secondary: lipgloss.Color("#ff0088"), // Hot pink
		Accent:    lipgloss.Color("#ffaa00"), // Orange
		Warning:   lipgloss.Color("#ff4400"), // Red-orange
		Error:     lipgloss.Color("#ff0044"), // Red
		Success:   lipgloss.Color("#00ff44"), // Green
		Purple:    lipgloss.Color("#bd93f9"), // A nice purple from the Dracula theme

		// Background colors
		BgPrimary:   lipgloss.Color("#0a0a0a"), // Dark background
		BgSecondary: lipgloss.Color("#1a1a1a"), // Slightly lighter
		BgTertiary:  lipgloss.Color("#2a2a2a"), // Even lighter

		// Text colors
		TextPrimary:   lipgloss.Color("#ffffff"), // White
		TextSecondary: lipgloss.Color("#cccccc"), // Light gray
		TextMuted:     lipgloss.Color("#888888"), // Gray
		TextInverse:   lipgloss.Color("#000000"), // Black
	}

	// Initialize styles
	theme.HeaderStyle = lipgloss.NewStyle().
		Background(theme.Primary).
		Foreground(theme.TextInverse).
		Bold(true).
		Padding(0, 1)

	theme.ClusterNameStyle = lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true).
		MarginRight(2)

	theme.ClusterVersionStyle = lipgloss.NewStyle().
		Foreground(theme.TextSecondary).
		Italic(true)

	theme.TableHeaderStyle = lipgloss.NewStyle().
		Background(theme.BgSecondary).
		Foreground(theme.Primary).
		Bold(true).
		Padding(0, 1)

	theme.TableRowStyle = lipgloss.NewStyle().
		Foreground(theme.TextPrimary).
		Padding(0, 1)

	theme.TableRowAltStyle = lipgloss.NewStyle().
		Background(theme.BgSecondary).
		Foreground(theme.TextPrimary).
		Padding(0, 1)

	theme.TableSelectedStyle = lipgloss.NewStyle().
		Background(theme.Primary).
		Foreground(theme.TextInverse).
		Bold(true).
		Padding(0, 1)

	theme.StatusRunningStyle = lipgloss.NewStyle().
		Foreground(theme.Success).
		Bold(true)

	theme.StatusPendingStyle = lipgloss.NewStyle().
		Foreground(theme.Warning).
		Bold(true)

	theme.StatusFailedStyle = lipgloss.NewStyle().
		Foreground(theme.Error).
		Bold(true)

	theme.StatusSucceededStyle = lipgloss.NewStyle().
		Foreground(theme.Accent).
		Bold(true)

	theme.BorderStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Primary).
		Padding(0, 1)

	theme.StatusBarStyle = lipgloss.NewStyle().
		Background(theme.BgSecondary).
		Foreground(theme.TextSecondary).
		Padding(0, 1)

	return theme
}

// GetStatusStyle returns the appropriate style for a pod status
func (t *Theme) GetStatusStyle(status string) lipgloss.Style {
	switch status {
	case "Running":
		return t.StatusRunningStyle
	case "Pending":
		return t.StatusPendingStyle
	case "Failed":
		return t.StatusFailedStyle
	case "Succeeded":
		return t.StatusSucceededStyle
	default:
		return lipgloss.NewStyle().Foreground(t.TextMuted)
	}
}

// Legacy compatibility - keeping the old global variables for now
// These will be removed in a future update

// Colors - Cyberpunk theme
var (
	// Primary colors
	Primary   = lipgloss.Color("#00ff88") // Bright cyan-green
	Secondary = lipgloss.Color("#ff0088") // Hot pink
	Accent    = lipgloss.Color("#ffaa00") // Orange
	Warning   = lipgloss.Color("#ff4400") // Red-orange
	Error     = lipgloss.Color("#ff0044") // Red
	Success   = lipgloss.Color("#00ff44") // Green
	Purple    = lipgloss.Color("#bd93f9") // A nice purple from the Dracula theme

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

// GetStatusStyle returns the appropriate style for a pod status (legacy function)
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
