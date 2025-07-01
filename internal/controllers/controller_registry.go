package controllers

import (
	"github.com/kevholditch/vigilant/internal/theme"
	"k8s.io/client-go/kubernetes"
)

type Theme = theme.Theme

// ControllerRegistry manages controller lifecycle with caching
type ControllerRegistry struct {
	clientset *kubernetes.Clientset
	theme     *Theme
	cache     map[string]Controller
	factories map[string]func(*kubernetes.Clientset, *Theme) Controller
}

// NewControllerRegistry creates a new registry
func NewControllerRegistry(clientset *kubernetes.Clientset, theme *Theme) *ControllerRegistry {
	return &ControllerRegistry{
		clientset: clientset,
		theme:     theme,
		cache:     make(map[string]Controller),
		factories: make(map[string]func(*kubernetes.Clientset, *Theme) Controller),
	}
}

// Register adds a controller factory to the registry
func (r *ControllerRegistry) Register(resource string, factory func(*kubernetes.Clientset, *Theme) Controller) {
	r.factories[resource] = factory
}

// GetController returns a controller instance, creating it if needed
func (r *ControllerRegistry) GetController(resource string) (Controller, bool) {
	// Check cache first
	if controller, exists := r.cache[resource]; exists {
		return controller, true
	}

	// Create new controller if factory exists
	if factory, exists := r.factories[resource]; exists {
		controller := factory(r.clientset, r.theme)
		r.cache[resource] = controller
		return controller, true
	}

	return nil, false
}

// GetAvailableResources returns all registered resource names
func (r *ControllerRegistry) GetAvailableResources() []string {
	resources := make([]string, 0, len(r.factories))
	for resource := range r.factories {
		resources = append(resources, resource)
	}
	return resources
}

// ClearCache clears the controller cache (useful for testing or memory management)
func (r *ControllerRegistry) ClearCache() {
	r.cache = make(map[string]Controller)
}
