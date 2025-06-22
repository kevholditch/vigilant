package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kevholditch/vigilant/internal/controllers"
	"github.com/kevholditch/vigilant/internal/models"
	"github.com/kevholditch/vigilant/internal/theme"
	"github.com/kevholditch/vigilant/internal/views"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// ViewType represents the current view type
type ViewType int

const (
	PodListView ViewType = iota
	DescribePodView
)

// ViewTransitionMsg represents a message to transition between views
type ViewTransitionMsg struct {
	ToView ViewType
}

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
	controlPlaneNodes int
	workerNodes       int
	currentController controllers.Controller
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

	// --- Get Node Counts ---
	// Get control plane nodes using label selector
	controlPlaneNodesList, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{
		LabelSelector: "node-role.kubernetes.io/control-plane",
	})
	if err != nil {
		log.Printf("could not list control plane nodes: %v", err)
	}
	controlPlaneNodes := len(controlPlaneNodesList.Items)

	// If no control plane nodes found with the new label, try the old master label
	if controlPlaneNodes == 0 {
		masterNodesList, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{
			LabelSelector: "node-role.kubernetes.io/master",
		})
		if err != nil {
			log.Printf("could not list master nodes: %v", err)
		}
		controlPlaneNodes = len(masterNodesList.Items)
	}

	// Get worker nodes using label selector (nodes without control-plane or master labels)
	workerNodesList, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{
		LabelSelector: "!node-role.kubernetes.io/control-plane,!node-role.kubernetes.io/master",
	})
	if err != nil {
		log.Printf("could not list worker nodes: %v", err)
	}
	workerNodes := len(workerNodesList.Items)
	// --- End Node Counts ---

	pods, err := models.GetPods(clientset)
	if err != nil {
		log.Fatalf("error getting pods: %v", err)
	}

	// Create theme
	theme := theme.NewDefaultTheme()

	podView := views.NewPodView(pods, theme, clusterName)

	app := &App{
		clientset:         clientset,
		podView:           podView,
		currentView:       PodListView,
		theme:             theme,
		clusterName:       clusterName,
		controlPlaneNodes: controlPlaneNodes,
		workerNodes:       workerNodes,
	}

	if k8sVersion != nil {
		app.kubernetesVersion = k8sVersion.String()
	}

	// Initialize the default controller
	app.initializeControllers()

	return app
}

// initializeControllers sets up the controllers for different views
func (a *App) initializeControllers() {
	// Set up the default pod list controller
	a.currentController = controllers.NewPodListController(a.podView, a.handleDescribePod)
}

// handleDescribePod handles the transition to describe pod view
func (a *App) handleDescribePod(podView *views.PodView) tea.Cmd {
	return func() tea.Msg {
		selectedPod := podView.GetSelected()
		if selectedPod != nil {
			return ViewTransitionMsg{ToView: DescribePodView}
		}
		return nil
	}
}

// handleBackToList handles the transition back to pod list view
func (a *App) handleBackToList() tea.Cmd {
	return func() tea.Msg {
		return ViewTransitionMsg{ToView: PodListView}
	}
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
	case ViewTransitionMsg:
		switch msg.ToView {
		case DescribePodView:
			selectedPod := a.podView.GetSelected()
			if selectedPod != nil {
				a.describePodView = views.NewDescribePodView(selectedPod, a.theme)
				a.describePodView.SetSize(a.width, a.height)
				a.currentView = DescribePodView
				a.currentController = controllers.NewDescribePodController(a.describePodView, a.handleBackToList)
			}
		case PodListView:
			a.currentView = PodListView
			a.describePodView = nil
			a.currentController = controllers.NewPodListController(a.podView, a.handleDescribePod)
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
	// --- Styles ---
	barStyle := lipgloss.NewStyle().
		Background(a.theme.BgSecondary).
		Foreground(a.theme.TextPrimary).
		Padding(0, 1)

	separator := lipgloss.NewStyle().
		Foreground(a.theme.TextMuted).
		SetString(" | ").
		String()

	// --- Content ---
	clusterInfo := fmt.Sprintf("☸️ %s", a.clusterName)

	var viewText string
	switch a.currentView {
	case PodListView:
		viewText = "🔍 Viewing pods"
	case DescribePodView:
		selectedPod := a.podView.GetSelected()
		if selectedPod != nil {
			viewText = fmt.Sprintf("🔍 Describing pod %s", selectedPod.Name)
		} else {
			viewText = "🔍 Describing pod"
		}
	}

	controlPlaneInfo := fmt.Sprintf("🕹️ CP %d", a.controlPlaneNodes)
	workerInfo := fmt.Sprintf("👷 W %d", a.workerNodes)
	k8sInfo := fmt.Sprintf("K8s: %s", a.kubernetesVersion)

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
		Width(a.width).
		Align(lipgloss.Center).
		Render(barStyle.Render(content))

	return bar
}
