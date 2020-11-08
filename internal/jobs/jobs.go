package jobs

import (
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GoJob is a Kubernetes Job
type GoJob struct {
	metav1.ObjectMeta `json:"object_meta" bson:"object_meta"`
	Spec              batchv1.JobSpec `json:"spec" bson:"spec"`
}

const (
	defaultNamespace = "default"
)

// New creates a basic Job
func New(name, generateName, namespace, image string) (*GoJob, error) {
	if namespace == "" {
		namespace = defaultNamespace
	}
	var containerName string
	if name != "" {
		containerName = name
	} else {
		containerName = generateName
	}
	w := GoJob{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: generateName,
			Name:         name,
			Namespace:    namespace,
		},
		Spec: batchv1.JobSpec{
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: generateName,
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  containerName,
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
	Create(name, generateName, namespace, image string) (*GoJob, error)
	Update(name, namespace string, jobSpec batchv1.JobSpec) error
	Delete(name, namespace string) error
}
