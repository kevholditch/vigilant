package models

import "time"

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

// GetSamplePods returns hardcoded sample pod data
func GetSamplePods() []Pod {
	return []Pod{
		{
			Name:      "nginx-deployment-7d4c8b5c9",
			Namespace: "default",
			Status:    "Running",
			Ready:     "1/1",
			Restarts:  0,
			Age:       2 * time.Hour,
			IP:        "10.244.0.15",
			Node:      "worker-1",
		},
		{
			Name:      "redis-master-0",
			Namespace: "default",
			Status:    "Running",
			Ready:     "1/1",
			Restarts:  2,
			Age:       1 * time.Hour,
			IP:        "10.244.0.16",
			Node:      "worker-2",
		},
		{
			Name:      "postgres-0",
			Namespace: "database",
			Status:    "Running",
			Ready:     "1/1",
			Restarts:  0,
			Age:       5 * time.Hour,
			IP:        "10.244.0.17",
			Node:      "worker-1",
		},
		{
			Name:      "api-server-7d4c8b5c9",
			Namespace: "default",
			Status:    "Pending",
			Ready:     "0/1",
			Restarts:  0,
			Age:       30 * time.Minute,
			IP:        "",
			Node:      "",
		},
		{
			Name:      "cron-job-123456",
			Namespace: "default",
			Status:    "Succeeded",
			Ready:     "0/1",
			Restarts:  0,
			Age:       10 * time.Minute,
			IP:        "10.244.0.18",
			Node:      "worker-3",
		},
		{
			Name:      "webhook-handler",
			Namespace: "default",
			Status:    "Failed",
			Ready:     "0/1",
			Restarts:  3,
			Age:       45 * time.Minute,
			IP:        "10.244.0.19",
			Node:      "worker-1",
		},
		{
			Name:      "monitoring-grafana",
			Namespace: "monitoring",
			Status:    "Running",
			Ready:     "1/1",
			Restarts:  1,
			Age:       3 * time.Hour,
			IP:        "10.244.0.20",
			Node:      "worker-2",
		},
		{
			Name:      "monitoring-prometheus",
			Namespace: "monitoring",
			Status:    "Running",
			Ready:     "1/1",
			Restarts:  0,
			Age:       3 * time.Hour,
			IP:        "10.244.0.21",
			Node:      "worker-3",
		},
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
