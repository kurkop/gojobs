package inkube

import (
	"context"
	"log"
	"sync"

	"github.com/kurkop/gojobs/internal/jobs"
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

func (m *goJobsRepository) Create(name, generateName, image string) (*jobs.GoJob, error) {
	newJob, err := jobs.New(name, generateName, image)
	if err != nil {
		log.Printf("error instancing job: %v", err)
		return nil, err
	}
	p := &batchv1.Job{ObjectMeta: newJob.ObjectMeta, Spec: newJob.Spec}
	jobCreated, err := m.client.BatchV1().Jobs(jobs.DefaultNamespace).Create(context.TODO(), p, metav1.CreateOptions{})
	if err != nil {
		log.Printf("error creating job: %v", err)
		return nil, err
	}
	jobCreatedName := jobCreated.ObjectMeta.GetName()
	jobGot := jobs.GoJob{ObjectMeta: jobCreated.ObjectMeta, Spec: jobCreated.Spec}
	m.goJobs[jobCreatedName] = jobGot
	return &jobGot, nil
}

func (m *goJobsRepository) Get(name string) (*jobs.GoJob, error) {
	getJob, err := m.client.BatchV1().Jobs(jobs.DefaultNamespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	// log.Printf("GetJob Spec", getJob.Spec)
	jobGot := jobs.GoJob{ObjectMeta: getJob.ObjectMeta, Spec: getJob.Spec}
	m.goJobs[name] = jobGot
	return &jobGot, nil
}

func (m *goJobsRepository) GetAll() (jobsGot *jobs.GoJobList, err error) {
	getJobs, err := m.client.BatchV1().Jobs(jobs.DefaultNamespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	jobsGot = &jobs.GoJobList{}
	for _, item := range getJobs.Items {
		log.Printf("Item ObjectMeta: %v", item.ObjectMeta)
		log.Printf("Item Spec: %v", item.Spec)
		jobsGot.Items = append(jobsGot.Items, jobs.GoJob{ObjectMeta: item.ObjectMeta, Spec: item.Spec})
	}
	return
}

func (m *goJobsRepository) Update(name string, jobSpec batchv1.JobSpec) error {
	currentJob, err := m.Get(name)
	if err != nil {
		return err
	}
	p := &batchv1.Job{ObjectMeta: currentJob.ObjectMeta, Spec: jobSpec}
	jobUpdated, err := m.client.BatchV1().Jobs(jobs.DefaultNamespace).Update(context.TODO(), p, metav1.UpdateOptions{})
	if err != nil {
		log.Fatalf("error updating job: %v", err)
		return err
	}
	jobUpdatedName := jobUpdated.ObjectMeta.GetName()
	jobGot := jobs.GoJob{ObjectMeta: jobUpdated.ObjectMeta, Spec: jobUpdated.Spec}
	m.goJobs[jobUpdatedName] = jobGot
	return nil
}

func (m *goJobsRepository) Delete(name string) error {
	err := m.client.BatchV1().Jobs(jobs.DefaultNamespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	if _, ok := m.goJobs[name]; ok {
		delete(m.goJobs, name)
	}
	return nil
}
