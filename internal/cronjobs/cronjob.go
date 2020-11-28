package cronjobs

import (
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GoCronJob is a Kubernetes Job
type GoCronJob struct {
	metav1.ObjectMeta
	CronJobSpec batchv1beta1.CronJobSpec
}

// GoCronJobList is a slice of GoCronJob
type GoCronJobList struct {
	Items []GoCronJob `json:"items"`
}

const (
	DefaultNamespace = "gojobs"
)

// New creates a basic Job
func New(name, image, schedule string) (*GoCronJob, error) {
	generateName := name + "-"
	// var sucessHistoryLimit, failedHistoryLimit int32
	var sucessHistoryLimit int32 = 20
	var failedHistoryLimit int32 = 10
	w := GoCronJob{
		ObjectMeta: metav1.ObjectMeta{
			// GenerateName: generateName,
			Name:      name,
			Namespace: DefaultNamespace,
		},
		CronJobSpec: batchv1beta1.CronJobSpec{
			Schedule:                   schedule,
			SuccessfulJobsHistoryLimit: &sucessHistoryLimit,
			FailedJobsHistoryLimit:     &failedHistoryLimit,
			JobTemplate: batchv1beta1.JobTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: generateName,
				},
				Spec: batchv1.JobSpec{
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{
							Containers: []apiv1.Container{
								{
									Name:  name,
									Image: image,
								},
							},
							RestartPolicy: v1.RestartPolicyNever,
						},
					},
				},
			},
		},
	}
	return &w, nil
}

// Repository interface to handle GobJob methods
type Repository interface {
	Get(name string) (*GoCronJob, error)
	GetAll() (cronJobsGot *GoCronJobList, err error)
	Create(name, image, schedule string) (*GoCronJob, error)
	Update(name string, jobSpec batchv1beta1.CronJobSpec) error
	Delete(name string) error
}
