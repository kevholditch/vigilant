package models

import (
	"context"
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Pod represents a Kubernetes pod
type Pod struct {
	Name      string
	Namespace string
	Status    string
	Ready     string
	Restarts  int
	Age       time.Duration
	IP        string
	Node      string
}

// GetPods fetches a list of pods from the Kubernetes cluster
// (Removed: now handled by the controller)

// GetPod fetches a single pod by name and namespace
func GetPod(clientset *kubernetes.Clientset, namespace, name string) (*Pod, error) {
	k8sPod, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not get pod %s in namespace %s: %w", name, namespace, err)
	}

	pod := ToPodModel(*k8sPod)
	return &pod, nil
}

// ToPodModel converts a Kubernetes API pod object to our internal Pod model
func ToPodModel(p v1.Pod) Pod {
	restarts := 0
	readyContainers := 0
	for _, cs := range p.Status.ContainerStatuses {
		restarts += int(cs.RestartCount)
		if cs.Ready {
			readyContainers++
		}
	}

	return Pod{
		Name:      p.Name,
		Namespace: p.Namespace,
		Status:    string(p.Status.Phase),
		Ready:     fmt.Sprintf("%d/%d", readyContainers, len(p.Spec.Containers)),
		Restarts:  restarts,
		Age:       time.Since(p.CreationTimestamp.Time),
		IP:        p.Status.PodIP,
		Node:      p.Spec.NodeName,
	}
}

// FormatAge formats the age duration to a human-readable string
func (p Pod) FormatAge() string {
	if p.Age < time.Minute {
		return "<1m"
	}
	if p.Age < time.Hour {
		return p.Age.Round(time.Minute).String()
	}
	if p.Age < 24*time.Hour {
		return p.Age.Round(time.Hour).String()
	}
	return p.Age.Round(24 * time.Hour).String()
}
