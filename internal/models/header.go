package models

type HeaderModel struct {
	ClusterName       string
	KubernetesVersion string
	ControlPlaneNodes int
	WorkerNodes       int
}
