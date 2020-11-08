package inkube

import (
	"context"
	"log"
	"sync"

	"github.com/kurkop/gojob/internal/jobs"
	"k8s.io/client-go/kubernetes"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type goJobsRepository struct {
	goJobs map[string]jobs.GoJob
	client *kubernetes.Clientset
}

var (
	goJobsOnce     sync.Once
	goJobsInstance *goJobsRepository
)

// NewGoJobsRepository Instance new repository
func NewGoJobsRepository(client *kubernetes.Clientset) jobs.Repository {
	goJobsOnce.Do(func() {
		goJobsInstance = &goJobsRepository{
			goJobs: make(map[string]jobs.GoJob),
			client: client,
		}
	})
	return goJobsInstance
}

func (m *goJobsRepository) Create(name, generateName, namespace, image string) (*jobs.GoJob, error) {
	newJob, err := jobs.New(name, generateName, namespace, image)
	if err != nil {
		log.Printf("error instancing job: %v", err)
		return nil, err
	}
	p := &batchv1.Job{ObjectMeta: newJob.ObjectMeta, Spec: newJob.Spec}
	jobCreated, err := m.client.BatchV1().Jobs(namespace).Create(context.TODO(), p, metav1.CreateOptions{})
	if err != nil {
		log.Printf("error creating job: %v", err)
		return nil, err
	}
	jobCreatedName := jobCreated.ObjectMeta.GetName()
	jobGot := jobs.GoJob{ObjectMeta: jobCreated.ObjectMeta, Spec: jobCreated.Spec}
	m.goJobs[jobCreatedName] = jobGot
	return &jobGot, nil
}

func (m *goJobsRepository) Get(name, namespace string) (*jobs.GoJob, error) {
	getJob, err := m.client.BatchV1().Jobs(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	// log.Printf("GetJob Spec", getJob.Spec)
	jobGot := jobs.GoJob{ObjectMeta: getJob.ObjectMeta, Spec: getJob.Spec}
	m.goJobs[name] = jobGot
	return &jobGot, nil
}

func (m *goJobsRepository) Update(name, namespace string, jobSpec batchv1.JobSpec) error {
	currentJob, err := m.Get(name, namespace)
	if err != nil {
		return err
	}
	p := &batchv1.Job{ObjectMeta: currentJob.ObjectMeta, Spec: jobSpec}
	jobUpdated, err := m.client.BatchV1().Jobs(namespace).Update(context.TODO(), p, metav1.UpdateOptions{})
	if err != nil {
		log.Fatalf("error updating job: %v", err)
		return err
	}
	jobUpdatedName := jobUpdated.ObjectMeta.GetName()
	jobGot := jobs.GoJob{ObjectMeta: jobUpdated.ObjectMeta, Spec: jobUpdated.Spec}
	m.goJobs[jobUpdatedName] = jobGot
	return nil
}

func (m *goJobsRepository) Delete(name, namespace string) error {
	err := m.client.BatchV1().Jobs(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	if _, ok := m.goJobs[name]; ok {
		delete(m.goJobs, name)
	}
	return nil
}
