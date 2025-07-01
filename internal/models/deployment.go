package models

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Deployment represents a Kubernetes deployment
type Deployment struct {
	Name      string
	Namespace string
	Status    string
	Ready     string
	UpToDate  int
	Available int
	Age       time.Duration
	Strategy  string
	Image     string
}

// GetDeployment fetches a single deployment by name and namespace
func GetDeployment(clientset *kubernetes.Clientset, namespace, name string) (*Deployment, error) {
	k8sDeployment, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not get deployment %s in namespace %s: %w", name, namespace, err)
	}

	deployment := ToDeploymentModel(*k8sDeployment)
	return &deployment, nil
}

// ToDeploymentModel converts a Kubernetes API deployment object to our internal Deployment model
func ToDeploymentModel(d appsv1.Deployment) Deployment {
	ready := d.Status.ReadyReplicas
	available := d.Status.AvailableReplicas
	upToDate := d.Status.UpdatedReplicas

	// Determine status
	status := "Unknown"
	if d.Status.Replicas == 0 {
		status = "Scaled to 0"
	} else if ready == d.Status.Replicas {
		status = "Ready"
	} else if available > 0 {
		status = "Available"
	} else {
		status = "Not Ready"
	}

	// Get strategy type
	strategy := "RollingUpdate"
	if d.Spec.Strategy.Type == appsv1.RecreateDeploymentStrategyType {
		strategy = "Recreate"
	}

	// Get image from first container
	image := "N/A"
	if len(d.Spec.Template.Spec.Containers) > 0 {
		image = d.Spec.Template.Spec.Containers[0].Image
	}

	return Deployment{
		Name:      d.Name,
		Namespace: d.Namespace,
		Status:    status,
		Ready:     fmt.Sprintf("%d/%d", ready, d.Status.Replicas),
		UpToDate:  int(upToDate),
		Available: int(available),
		Age:       time.Since(d.CreationTimestamp.Time),
		Strategy:  strategy,
		Image:     image,
	}
}

// FormatAge formats the age duration to a human-readable string
func (d Deployment) FormatAge() string {
	if d.Age < time.Minute {
		return "<1m"
	}
	if d.Age < time.Hour {
		return d.Age.Round(time.Minute).String()
	}
	if d.Age < 24*time.Hour {
		return d.Age.Round(time.Hour).String()
	}
	return d.Age.Round(24 * time.Hour).String()
}
