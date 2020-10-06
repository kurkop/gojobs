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
	p := &batchv1.Job{ObjectMeta: metav1.ObjectMeta{Name: "my-job"}}
	_, err := client.BatchV1().Jobs("test-ns").Create(context.TODO(), p, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("error injecting job add: %v", err)
	}

	select {
	case job := <-jobs:
		t.Logf("Got job from channel: %s/%s", job.Namespace, job.Name)
	case <-time.After(wait.ForeverTestTimeout):
		t.Error("Informer did not get the added job")
	}
}
