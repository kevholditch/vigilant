package app

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kevholditch/vigilant/internal/models"
	"github.com/kevholditch/vigilant/internal/views"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// App represents the main application
type App struct {
	clientset *kubernetes.Clientset
	podView   *views.PodView
	width     int
	height    int
}

// NewApp creates a new application instance
func NewApp() *App {
	// --- Kubernetes Client Setup ---
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("error getting user home dir: %v", err)
	}
	kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		log.Fatalf("error getting Kubernetes config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("error creating Kubernetes client: %v", err)
	}
	// --- End Kubernetes Client Setup ---

	pods, err := models.GetPods(clientset)
	if err != nil {
		log.Fatalf("error getting pods: %v", err)
	}

	podView := views.NewPodView(pods)

	return &App{
		clientset: clientset,
		podView:   podView,
	}
}

// Run starts the application
func (a *App) Run() error {
	fmt.Println("Starting Vigilant...")
	fmt.Println("Press 'q' to quit, arrow keys to navigate")

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
		case "up", "k":
			a.podView.SelectPrev()
		case "down", "j":
			a.podView.SelectNext()
		}
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.podView.SetSize(msg.Width, msg.Height)
	}
	return a, nil
}

// View renders the application
func (a *App) View() string {
	if a.width == 0 || a.height == 0 {
		return "Loading..."
	}
	return a.podView.Render()
}
