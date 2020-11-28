package inkube

import (
	"context"
	"log"
	"sync"

	"github.com/kurkop/gojobs/internal/cronjobs"
	"k8s.io/client-go/kubernetes"

	batchv1beta1 "k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type goCronJobsRepository struct {
	goCronJobs map[string]cronjobs.GoCronJob
	client     *kubernetes.Clientset
}

var (
	goCronJobsOnce     sync.Once
	goCronJobsInstance *goCronJobsRepository
)

// NewGoCronJobsRepository Instance new repository
func NewGoCronJobsRepository(client *kubernetes.Clientset) cronjobs.Repository {
	goCronJobsOnce.Do(func() {
		goCronJobsInstance = &goCronJobsRepository{
			goCronJobs: make(map[string]cronjobs.GoCronJob),
			client:     client,
		}
	})
	return goCronJobsInstance
}

func (m *goCronJobsRepository) Create(name, image, schedule string) (*cronjobs.GoCronJob, error) {
	newJob, err := cronjobs.New(name, image, schedule)
	if err != nil {
		log.Fatalf("error instancing job: %v", err)
		return nil, err
	}
	p := &batchv1beta1.CronJob{ObjectMeta: newJob.ObjectMeta, Spec: newJob.CronJobSpec}
	jobCreated, err := m.client.BatchV1beta1().CronJobs(cronjobs.DefaultNamespace).Create(context.TODO(), p, metav1.CreateOptions{})
	if err != nil {
		log.Fatalf("error creating job: %v", err)
		return nil, err
	}
	jobCreatedName := jobCreated.ObjectMeta.GetName()
	jobGot := cronjobs.GoCronJob{ObjectMeta: jobCreated.ObjectMeta, CronJobSpec: jobCreated.Spec}
	m.goCronJobs[jobCreatedName] = jobGot
	return &jobGot, nil
}

func (m *goCronJobsRepository) Get(name string) (*cronjobs.GoCronJob, error) {
	getJob, err := m.client.BatchV1beta1().CronJobs(cronjobs.DefaultNamespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	jobGot := cronjobs.GoCronJob{ObjectMeta: getJob.ObjectMeta, CronJobSpec: getJob.Spec}
	m.goCronJobs[name] = jobGot
	return &jobGot, nil
}

func (m *goCronJobsRepository) GetAll() (cronJobsGot *cronjobs.GoCronJobList, err error) {
	getJobs, err := m.client.BatchV1beta1().CronJobs(cronjobs.DefaultNamespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	cronJobsGot = &cronjobs.GoCronJobList{}
	for _, item := range getJobs.Items {
		cronJobsGot.Items = append(cronJobsGot.Items, cronjobs.GoCronJob{ObjectMeta: item.ObjectMeta, CronJobSpec: item.Spec})
	}
	return
}

func (m *goCronJobsRepository) Update(name string, jobSpec batchv1beta1.CronJobSpec) error {
	currentJob, err := m.Get(name)
	if err != nil {
		return err
	}
	p := &batchv1beta1.CronJob{ObjectMeta: currentJob.ObjectMeta, Spec: jobSpec}
	jobUpdated, err := m.client.BatchV1beta1().CronJobs(cronjobs.DefaultNamespace).Update(context.TODO(), p, metav1.UpdateOptions{})
	if err != nil {
		log.Fatalf("error updating job: %v", err)
		return err
	}
	jobUpdatedName := jobUpdated.ObjectMeta.GetName()
	jobGot := cronjobs.GoCronJob{ObjectMeta: jobUpdated.ObjectMeta, CronJobSpec: jobUpdated.Spec}
	m.goCronJobs[jobUpdatedName] = jobGot
	return nil
}

func (m *goCronJobsRepository) Delete(name string) error {
	err := m.client.BatchV1beta1().CronJobs(cronjobs.DefaultNamespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	if _, ok := m.goCronJobs[name]; ok {
		delete(m.goCronJobs, name)
	}
	return nil
}
