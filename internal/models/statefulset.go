package models

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// StatefulSet represents a Kubernetes StatefulSet
type StatefulSet struct {
	Name      string
	Namespace string
	Status    string
	Ready     string
	Replicas  int
	Available int
	Age       time.Duration
	Image     string
}

// GetStatefulSet fetches a single StatefulSet by name and namespace
func GetStatefulSet(clientset *kubernetes.Clientset, namespace, name string) (*StatefulSet, error) {
	k8sStatefulSet, err := clientset.AppsV1().StatefulSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not get statefulset %s in namespace %s: %w", name, namespace, err)
	}

	statefulset := ToStatefulSetModel(*k8sStatefulSet)
	return &statefulset, nil
}

// ToStatefulSetModel converts a Kubernetes API StatefulSet object to our internal StatefulSet model
func ToStatefulSetModel(ss appsv1.StatefulSet) StatefulSet {
	ready := ss.Status.ReadyReplicas
	available := ss.Status.CurrentReplicas
	replicas := ss.Status.Replicas
	status := "Unknown"
	if replicas == 0 {
		status = "Scaled to 0"
	} else if ready == replicas {
		status = "Ready"
	} else if available > 0 {
		status = "Available"
	} else {
		status = "Not Ready"
	}
	image := "N/A"
	if len(ss.Spec.Template.Spec.Containers) > 0 {
		image = ss.Spec.Template.Spec.Containers[0].Image
	}
	return StatefulSet{
		Name:      ss.Name,
		Namespace: ss.Namespace,
		Status:    status,
		Ready:     fmt.Sprintf("%d/%d", ready, replicas),
		Replicas:  int(replicas),
		Available: int(available),
		Age:       time.Since(ss.CreationTimestamp.Time),
		Image:     image,
	}
}

// FormatAge formats the age duration to a human-readable string
func (s StatefulSet) FormatAge() string {
	if s.Age < time.Minute {
		return "<1m"
	}
	if s.Age < time.Hour {
		return s.Age.Round(time.Minute).String()
	}
	if s.Age < 24*time.Hour {
		return s.Age.Round(time.Hour).String()
	}
	return s.Age.Round(24 * time.Hour).String()
}
