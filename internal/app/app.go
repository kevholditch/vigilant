package app

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kevholditch/vigilant/internal/controllers"
	"github.com/kevholditch/vigilant/internal/theme"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// App represents the main application
type App struct {
	clientset            *kubernetes.Clientset
	width                int
	height               int
	theme                *theme.Theme
	currentController    controllers.Controller
	headerController     *controllers.HeaderController
	commandBarController *controllers.CommandBarController
	controllerRegistry   *controllers.ControllerRegistry
}

// NewApp creates a new application instance
func NewApp() *App {
	clientset, err := newClientSet()
	if err != nil {
		log.Fatal(fmt.Sprintf("error creating Kubernetes client: %v", err))
	}

	// Create theme
	theme := theme.NewDefaultTheme()

	app := &App{
		clientset: clientset,
		theme:     theme,
	}

	// Initialize the controllers
	app.initializeControllers()

	return app
}

func newClientSet() (*kubernetes.Clientset, error) {

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error getting user home dir: %v", err)
	}
	kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, fmt.Errorf("error getting Kubernetes config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error creating Kubernetes client: %v", err)
	}

	return clientset, nil
}

// initializeControllers sets up the controllers
func (a *App) initializeControllers() {
	a.buildRegistry()
	if controller, exists := a.controllerRegistry.GetController("pods"); exists {
		a.currentController = controller
	}
	a.headerController = controllers.NewHeaderController(a.theme, a.clientset)
	availableResources := a.controllerRegistry.GetAvailableResources()
	a.commandBarController = controllers.NewCommandBarController(a.clientset, a.theme, "", availableResources, a.handleViewSwitch)
}

func (a *App) buildRegistry() {
	a.controllerRegistry = controllers.NewControllerRegistry(a.clientset, a.theme)
	a.controllerRegistry.Register("pods", func(clientset *kubernetes.Clientset, theme *controllers.Theme) controllers.Controller {
		return controllers.NewPodController(clientset, theme, "")
	})
	a.controllerRegistry.Register("deployments", func(clientset *kubernetes.Clientset, theme *controllers.Theme) controllers.Controller {
		return controllers.NewDeploymentController(clientset, theme, "")
	})
}

// handleViewSwitch handles switching between different views
func (a *App) handleViewSwitch(resource string) tea.Cmd {
	return func() tea.Msg {
		if controller, exists := a.controllerRegistry.GetController(resource); exists {
			a.currentController = controller
		}
		return nil
	}
}

// tickMsg is sent periodically to check for updates
type tickMsg time.Time

// tick returns a command that sends a tick message
func tick() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(time.Second)
		return tickMsg(time.Now())
	}
}

// Run starts the application
func (a *App) Run() error {
	fmt.Println("Starting Vigilant...")
	fmt.Println("Press 'q' to quit, ':' to open command bar, arrow keys to navigate")

	// Create the bubble tea program
	p := tea.NewProgram(
		a,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Turn on mouse support so we can track the mouse wheel
	)

	// Run the program
	_, err := p.Run()
	return err
}

// Init initializes the application
func (a *App) Init() tea.Cmd {
	return tick()
}

// Update handles messages and updates the application state
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return a, tea.Quit
		case ":":
			// Activate command bar
			a.commandBarController.Activate()
			return a, nil
		default:
			// Check if command bar is active first
			if a.commandBarController.IsActive() {
				return a, a.commandBarController.HandleKey(msg)
			}

			// Delegate to the current controller
			if a.currentController != nil {
				return a, a.currentController.HandleKey(msg)
			}
		}
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
	case controllers.UpdateMsg:
		// Update received, trigger a re-render
		return a, tick()
	case tickMsg:
		// Check if we need to update the view
		if a.currentController != nil {
			// Try to get update channel from the current controller
			if updateableController, ok := a.currentController.(controllers.UpdateableController); ok {
				updateChan := updateableController.GetUpdateChannel()
				if updateChan != nil {
					// Force a re-render by calling Render
					headerHeight := a.headerController.GetHeight()
					commandBarHeight := a.getCommandBarHeight()
					a.currentController.Render(a.width, a.height-headerHeight-commandBarHeight)
				}
			}
		}
		return a, tick()
	}
	return a, nil
}

// View renders the application
func (a *App) View() string {
	if a.width == 0 || a.height == 0 {
		return "Loading..."
	}

	var viewText string
	if a.currentController != nil {
		viewText = a.currentController.ActionText()
	} else {
		viewText = "No controller available"
	}

	header := a.headerController.Render(a.width, viewText)
	headerHeight := a.headerController.GetHeight()

	// Render command bar
	commandBar := a.commandBarController.Render(a.width, 0)
	commandBarHeight := a.getCommandBarHeight()

	// Calculate available height for main content
	viewDisplayHeight := a.height - headerHeight - commandBarHeight

	var viewContent string
	if a.currentController != nil {
		viewContent = a.currentController.Render(a.width, viewDisplayHeight)
	} else {
		viewContent = "No controller available"
	}

	// Combine all components
	var components []string
	components = append(components, header)
	if commandBar != "" {
		components = append(components, commandBar)
	}
	components = append(components, viewContent)

	return lipgloss.JoinVertical(lipgloss.Left, components...)
}

// getCommandBarHeight returns the height of the command bar
func (a *App) getCommandBarHeight() int {
	if a.commandBarController.IsActive() {
		// Command bar takes up space when active
		return 3 // Approximate height for command bar + suggestions
	}
	return 0
}
