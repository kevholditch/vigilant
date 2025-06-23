package controllers

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

// ClusterBuilder provides a fluent interface for building Kubernetes test clusters
// It is now independent and can be reused in any test.
type ClusterBuilder struct {
	t         *testing.T
	env       *envtest.Environment
	clientset *kubernetes.Clientset
}

// NewClusterBuilder creates a new cluster builder and starts the envtest environment
func NewClusterBuilder(t *testing.T) *ClusterBuilder {
	env := &envtest.Environment{
		CRDDirectoryPaths:     []string{},
		ErrorIfCRDPathMissing: false,
	}
	cfg, err := env.Start()
	require.NoError(t, err)
	clientset, err := kubernetes.NewForConfig(cfg)
	require.NoError(t, err)
	return &ClusterBuilder{
		t:         t,
		env:       env,
		clientset: clientset,
	}
}

// WithWorkerNodes creates the specified number of worker nodes
func (cb *ClusterBuilder) WithWorkerNodes(count int) *ClusterBuilder {
	for i := 1; i <= count; i++ {
		workerNode := &corev1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name:   fmt.Sprintf("worker-node-%d", i),
				Labels: map[string]string{},
			},
		}
		_, err := cb.clientset.CoreV1().Nodes().Create(context.TODO(), workerNode, metav1.CreateOptions{})
		require.NoError(cb.t, err)
	}
	return cb
}

// WithControlPlaneNodes creates the specified number of control plane nodes
func (cb *ClusterBuilder) WithControlPlaneNodes(count int) *ClusterBuilder {
	for i := 1; i <= count; i++ {
		controlPlaneNode := &corev1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("control-plane-node-%d", i),
				Labels: map[string]string{
					"node-role.kubernetes.io/control-plane": "",
				},
			},
		}
		_, err := cb.clientset.CoreV1().Nodes().Create(context.TODO(), controlPlaneNode, metav1.CreateOptions{})
		require.NoError(cb.t, err)
	}
	return cb
}

// WithMasterNodes creates the specified number of legacy master nodes
func (cb *ClusterBuilder) WithMasterNodes(count int) *ClusterBuilder {
	for i := 1; i <= count; i++ {
		masterNode := &corev1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("master-node-%d", i),
				Labels: map[string]string{
					"node-role.kubernetes.io/master": "",
				},
			},
		}
		_, err := cb.clientset.CoreV1().Nodes().Create(context.TODO(), masterNode, metav1.CreateOptions{})
		require.NoError(cb.t, err)
	}
	return cb
}

// Cleanup stops the envtest environment
func (cb *ClusterBuilder) Cleanup() {
	if cb.env != nil {
		cb.env.Stop()
	}
}

// GetClientset returns the clientset
func (cb *ClusterBuilder) GetClientset() *kubernetes.Clientset {
	return cb.clientset
}

// GetEnv returns the envtest environment
func (cb *ClusterBuilder) GetEnv() *envtest.Environment {
	return cb.env
}
