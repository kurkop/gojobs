package config

import (
	"os"

	"github.com/kurkop/gojobs/shared/kube"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	// KubeClient points to a kubernetes.Clientset
	KubeClient *kubernetes.Clientset
)

// KubeConnect opens a new kubernetes Clientset
func KubeConnect() {
	var config *rest.Config

	if os.Getenv("IN_CLUSTER") == "" {
		config = kube.GetLocalConfig()
	} else {
		config = kube.GetInClusterConfig()
	}

	// create the clientset
	KubeClient, _ = kube.NewClient(config)
}
