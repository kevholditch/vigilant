package models

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ReplicaSet represents a Kubernetes ReplicaSet
type ReplicaSet struct {
	Name      string
	Namespace string
	Status    string
	Ready     string
	Replicas  int
	Available int
	Age       time.Duration
	Image     string
}

// GetReplicaSet fetches a single ReplicaSet by name and namespace
func GetReplicaSet(clientset *kubernetes.Clientset, namespace, name string) (*ReplicaSet, error) {
	k8sReplicaSet, err := clientset.AppsV1().ReplicaSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not get replicaset %s in namespace %s: %w", name, namespace, err)
	}

	replicaset := ToReplicaSetModel(*k8sReplicaSet)
	return &replicaset, nil
}

// ToReplicaSetModel converts a Kubernetes API ReplicaSet object to our internal ReplicaSet model
func ToReplicaSetModel(rs appsv1.ReplicaSet) ReplicaSet {
	ready := rs.Status.ReadyReplicas
	available := rs.Status.AvailableReplicas
	replicas := rs.Status.Replicas
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
	if len(rs.Spec.Template.Spec.Containers) > 0 {
		image = rs.Spec.Template.Spec.Containers[0].Image
	}
	return ReplicaSet{
		Name:      rs.Name,
		Namespace: rs.Namespace,
		Status:    status,
		Ready:     fmt.Sprintf("%d/%d", ready, replicas),
		Replicas:  int(replicas),
		Available: int(available),
		Age:       time.Since(rs.CreationTimestamp.Time),
		Image:     image,
	}
}

// FormatAge formats the age duration to a human-readable string
func (r ReplicaSet) FormatAge() string {
	if r.Age < time.Minute {
		return "<1m"
	}
	if r.Age < time.Hour {
		return r.Age.Round(time.Minute).String()
	}
	if r.Age < 24*time.Hour {
		return r.Age.Round(time.Hour).String()
	}
	return r.Age.Round(24 * time.Hour).String()
}
