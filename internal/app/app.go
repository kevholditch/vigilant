package app

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kevholditch/vigilant/internal/models"
	"github.com/kevholditch/vigilant/internal/theme"
	"github.com/kevholditch/vigilant/internal/views"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// ViewType represents the current view type
type ViewType int

const (
	PodListView ViewType = iota
	DescribePodView
)

// App represents the main application
type App struct {
	clientset         *kubernetes.Clientset
	kubernetesVersion string
	clusterName       string
	podView           *views.PodView
	describePodView   *views.DescribePodView
	currentView       ViewType
	width             int
	height            int
	theme             *theme.Theme
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

	// Get the cluster name from the current context
	kubeConfig, err := clientcmd.NewDefaultClientConfigLoadingRules().Load()
	if err != nil {
		log.Fatalf("error loading kubeconfig: %v", err)
	}
	currentContext, ok := kubeConfig.Contexts[kubeConfig.CurrentContext]
	if !ok {
		log.Fatalf("current context not found in kubeconfig")
	}
	clusterName := currentContext.Cluster

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("error creating Kubernetes client: %v", err)
	}
	// --- End Kubernetes Client Setup ---

	// --- Get K8s Version ---
	k8sVersion, err := clientset.Discovery().ServerVersion()
	if err != nil {
		log.Printf("could not get Kubernetes version, leaving blank. Error: %v", err)
		k8sVersion = nil
	}
	// --- End K8s Version ---

	pods, err := models.GetPods(clientset)
	if err != nil {
		log.Fatalf("error getting pods: %v", err)
	}

	// Create theme
	theme := theme.NewDefaultTheme()

	podView := views.NewPodView(pods, theme, clusterName)

	app := &App{
		clientset:   clientset,
		podView:     podView,
		currentView: PodListView,
		theme:       theme,
		clusterName: clusterName,
	}

	if k8sVersion != nil {
		app.kubernetesVersion = k8sVersion.String()
	}
	return app
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
		case "up", "k":
			if a.currentView == PodListView {
				a.podView.SelectPrev()
			} else if a.currentView == DescribePodView && a.describePodView != nil {
				a.describePodView.ScrollUp()
			}
		case "down", "j":
			if a.currentView == PodListView {
				a.podView.SelectNext()
			} else if a.currentView == DescribePodView && a.describePodView != nil {
				a.describePodView.ScrollDown()
			}
		case "d":
			if a.currentView == PodListView {
				selectedPod := a.podView.GetSelected()
				if selectedPod != nil {
					a.describePodView = views.NewDescribePodView(selectedPod, a.theme)
					a.describePodView.SetSize(a.width, a.height)
					a.currentView = DescribePodView
				}
			}
		case "esc":
			if a.currentView == DescribePodView {
				a.currentView = PodListView
				a.describePodView = nil
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

	header := a.renderHeader()
	headerHeight := lipgloss.Height(header)
	viewDisplayHeight := a.height - headerHeight

	var viewContent string
	switch a.currentView {
	case PodListView:
		a.podView.SetSize(a.width, viewDisplayHeight)
		viewContent = a.podView.Render()
	case DescribePodView:
		if a.describePodView != nil {
			a.describePodView.SetSize(a.width, viewDisplayHeight)
			viewContent = a.describePodView.Render()
		} else {
			viewContent = "Error: Describe pod view not initialized"
		}
	default:
		viewContent = "Unknown view"
	}

	return lipgloss.JoinVertical(lipgloss.Left, header, viewContent)
}

func (a *App) renderHeader() string {
	titleStyle := lipgloss.NewStyle().
		Background(a.theme.Purple).
		Foreground(a.theme.TextPrimary).
		Padding(0, 1).
		Bold(true)

	k8sVersionStyle := lipgloss.NewStyle().
		Inherit(titleStyle).
		Background(a.theme.Primary).
		Foreground(a.theme.TextInverse)

	var viewTitle string
	switch a.currentView {
	case PodListView:
		viewTitle = "Pods on " + a.clusterName
	case DescribePodView:
		if a.describePodView != nil {
			selectedPod := a.podView.GetSelected()
			if selectedPod != nil {
				viewTitle = "Describe Pod: " + selectedPod.Name
			}
		} else {
			viewTitle = "Describe Pod"
		}
	}

	title := titleStyle.Render("👁️ Vigilant - " + viewTitle)
	version := k8sVersionStyle.Render("K8s: " + a.kubernetesVersion)

	// Spacer to push version to the right
	spacerWidth := a.width - lipgloss.Width(title) - lipgloss.Width(version)
	if spacerWidth < 0 {
		spacerWidth = 0
	}
	spacer := lipgloss.NewStyle().Width(spacerWidth).Render("")

	headerContent := lipgloss.JoinHorizontal(lipgloss.Top, title, spacer, version)
	return lipgloss.NewStyle().Width(a.width).Render(headerContent)
}
