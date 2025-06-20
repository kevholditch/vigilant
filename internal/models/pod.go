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
func GetPods(clientset *kubernetes.Clientset) ([]Pod, error) {
	// Fetch pods from all namespaces
	podList, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not list pods: %w", err)
	}

	var pods []Pod
	for _, k8sPod := range podList.Items {
		pods = append(pods, toPodModel(k8sPod))
	}
	return pods, nil
}

// toPodModel converts a Kubernetes API pod object to our internal Pod model
func toPodModel(p v1.Pod) Pod {
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
