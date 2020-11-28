package config

import (
	"github.com/kurkop/gojobs/shared/kube"
	"k8s.io/client-go/kubernetes"
)

var (
	// KubeClient points to a kubernetes.Clientset
	KubeClient *kubernetes.Clientset
)

// KubeConnect opens a new kubernetes Clientset
func KubeConnect() {
	config := kube.GetLocalConfig()

	// create the clientset
	KubeClient, _ = kube.NewClient(config)
}
