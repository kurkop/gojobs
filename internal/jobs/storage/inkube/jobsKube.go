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

func (m *goJobsRepository) Create(name, namespace, image string) (*jobs.GoJob, error) {
	newJob, err := jobs.New(name, namespace, image)
	if err != nil {
		log.Fatalf("error instancing job: %v", err)
		return nil, err
	}
	p := &batchv1.Job{ObjectMeta: newJob.ObjectMeta, Spec: newJob.Spec}
	jobCreated, err := m.client.BatchV1().Jobs(namespace).Create(context.TODO(), p, metav1.CreateOptions{})
	if err != nil {
		log.Fatalf("error creating job: %v", err)
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

// func (m *goJobsRepository) GetOne(filter *bson.M) (*jobs.GoJob, error) {
// 	var goJob jobs.GoJob
// 	opts := options.FindOne()
// 	err := m.client.FindOne(context.TODO(), filter, opts).Decode(&goJob)
// 	if err != nil {
// 		log.Fatal(err)
// 		return nil, err
// 	}
// 	return &goJob, nil
// }

// func (m *goJobsRepository) Find(filter *bson.M) (*mongo.Cursor, error) {
// 	findOptions := options.Find()
// 	cur, err := m.client.Find(context.TODO(), filter, findOptions)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return cur, nil
// }

// func (m *goJobsRepository) Delete(ID string) error {
// 	goJobID, _ := primitive.ObjectIDFromHex(ID)
// 	filter := bson.M{
// 		"_id": goJobID,
// 	}
// 	opts := options.Delete()
// 	_, err := m.client.DeleteOne(context.TODO(), filter, opts)
// 	if err != nil {
// 		log.Fatal(err)
// 		return err
// 	}
// 	return nil
// }

// func (m *goJobsRepository) Save(goJob jobs.GoJob) error {
// 	opts := options.Update().SetUpsert(true)
// 	filter := bson.D{{"_id", goJob.ID}}
// 	// data, err := json.Marshal(goJob)
// 	// fmt.Println(data)
// 	primitive.NewObjectID()

// 	update := bson.D{{"$set", goJob}}
// 	result, err := m.client.UpdateOne(context.TODO(), filter, update, opts)
// 	if err != nil {
// 		log.Fatal(err)
// 		return err
// 	}

// 	if result.MatchedCount != 0 {
// 		log.Println("matched and replaced an existing document")
// 	}
// 	if result.UpsertedCount != 0 {
// 		log.Printf("inserted a new document with ID %v\n", result.UpsertedID)
// 	}

// 	return nil
// }
