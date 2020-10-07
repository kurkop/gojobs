package jobs

import (
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GoJob is a Kubernetes Job
type GoJob struct {
	metav1.ObjectMeta
	Spec batchv1.JobSpec
}

const (
	defaultNamespace = "default"
)

// New creates a basic Job
func New(name, namespace, image string) (*GoJob, error) {
	generateName := name + "-"
	if namespace == "" {
		namespace = defaultNamespace
	}
	w := GoJob{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: generateName,
			// Name:      name,
			Namespace: namespace,
		},
		Spec: batchv1.JobSpec{
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: generateName,
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  name,
							Image: image,
						},
					},
					RestartPolicy: apiv1.RestartPolicyNever,
				},
			},
		},
	}
	return &w, nil
}

// Repository interface to handle GobJob methods
type Repository interface {
	Get(name, namespace string) (*GoJob, error)
	Create(name, namespace, image string) (*GoJob, error)
	Delete(name, namespace string) error
}
