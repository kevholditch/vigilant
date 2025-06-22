package app

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kevholditch/vigilant/internal/controllers"
	"github.com/kevholditch/vigilant/internal/theme"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// App represents the main application
type App struct {
	clientset         *kubernetes.Clientset
	width             int
	height            int
	theme             *theme.Theme
	currentController controllers.Controller
	headerController  *controllers.HeaderController
}

// NewApp creates a new application instance
func NewApp() *App {
	clientset, err := newClientSet()
	if err != nil {
		log.Fatal("error creating Kubernetes client: %v", err)
	}

	// Create theme
	theme := theme.NewDefaultTheme()

	app := &App{
		clientset: clientset,
		theme:     theme,
	}

	// Initialize the controller
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
	// Set up the pod controller that manages both list and describe views
	a.currentController = controllers.NewPodController(a.clientset, a.theme, "") // Pass empty string for clusterName if needed

	// Set up the header controller
	a.headerController = controllers.NewHeaderController(a.theme, a.clientset)
}

// Run starts the application
func (a *App) Run() error {
	fmt.Println("Starting Vigilant...")
	fmt.Println("Press 'q' to quit, arrow keys to navigate, 'd' to describe pod")

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
	return nil
}

// Update handles messages and updates the application state
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return a, tea.Quit
		default:
			// Delegate to the current controller
			if a.currentController != nil {
				return a, a.currentController.HandleKey(msg)
			}
		}
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
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
	viewDisplayHeight := a.height - headerHeight

	var viewContent string
	if a.currentController != nil {
		viewContent = a.currentController.Render(a.width, viewDisplayHeight)
	} else {
		viewContent = "No controller available"
	}

	return lipgloss.JoinVertical(lipgloss.Left, header, viewContent)
}
