package jobs

import (
	"context"
	"testing"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/cache"
)

// func TestCreate(t *testing.T) {
// 	config := kube.GetLocalConfig()

// 	// create the clientset
// 	client, err := kube.NewClient(config)

// 	gojob, err := New("my-job", "", "hello-world")
// 	if err != nil {
// 		t.Fatalf("error instancing job: %v", err)
// 	}
// 	p := &batchv1.Job{ObjectMeta: gojob.ObjectMeta, Spec: gojob.Spec}
// 	newJob, err := client.BatchV1().Jobs("default").Create(context.TODO(), p, metav1.CreateOptions{})
// 	t.Logf("Job name: %v", newJob.ObjectMeta.GetName())
// 	//client.BatchV1().Jobs("default").Get(context.TODO())
// 	if err != nil {
// 		t.Fatalf("error injecting job add: %v", err)
// 	}
// }

// TestGoJob demonstrates how to use a fake client with SharedInformerFactory in tests.
func TestFakeJob(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create the fake client.
	client := fake.NewSimpleClientset()

	// We will create an informer that writes added jobs to a channel.
	jobs := make(chan *batchv1.Job, 1)
	informers := informers.NewSharedInformerFactory(client, 0)
	jobInformer := informers.Batch().V1().Jobs().Informer()
	jobInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			job := obj.(*batchv1.Job)
			t.Logf("job added: %s/%s", job.Namespace, job.Name)
			jobs <- job
		},
	})

	// Make sure informers are running.
	informers.Start(ctx.Done())

	// This is not required in tests, but it serves as a proof-of-concept by
	// ensuring that the informer goroutine have warmed up and called List before
	// we send any events to it.
	cache.WaitForCacheSync(ctx.Done(), jobInformer.HasSynced)

	// Inject an event into the fake client.
	gojob, err := New("", "my-job", "test-ns", "hello-world")
	if err != nil {
		t.Fatalf("error instancing job: %v", err)
	}
	p := &batchv1.Job{ObjectMeta: gojob.ObjectMeta, Spec: gojob.Spec}
	_, err = client.BatchV1().Jobs("test-ns").Create(context.TODO(), p, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("error injecting job add: %v", err)
	}
	jobList, err := client.BatchV1().Jobs("test-ns").List(context.TODO(), metav1.ListOptions{})

	for _, job := range jobList.Items {
		t.Logf("%v", job)
	}

	select {
	case job := <-jobs:
		t.Logf("Got job from channel: %s/%s", job.Namespace, job.Name)
	case <-time.After(wait.ForeverTestTimeout):
		t.Error("Informer did not get the added job")
	}
}
