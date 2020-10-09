package cronjobs

import (
	"context"
	"testing"

	batchv1beta1 "k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/cache"
)

// func TestCreate(t *testing.T) {
// 	config := kube.GetLocalConfig()

// 	// create the clientset
// 	client, err := kube.NewClient(config)

// 	gocronJob, err := New("my-cronjob", "default", "hello-world", "*/1 * * * *")
// 	p := &batchv1beta1.CronJob{ObjectMeta: gocronJob.ObjectMeta, Spec: gocronJob.CronJobSpec}
// 	_, err = client.BatchV1beta1().CronJobs("default").Create(context.TODO(), p, metav1.CreateOptions{})
// 	if err != nil {
// 		t.Errorf("Error %v", err)
// 	}
// 	cronJobList, err := client.BatchV1beta1().CronJobs("default").List(context.TODO(), metav1.ListOptions{})
// 	if err != nil {
// 		t.Errorf("Error %v", err)
// 	}
// 	for _, cronJob := range cronJobList.Items {
// 		t.Logf("%v", cronJob)
// 	}
// }

// TestGoJob demonstrates how to use a fake client with SharedInformerFactory in tests.
func TestFakeJob(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create the fake client.
	client := fake.NewSimpleClientset()

	// We will create an informer that writes added jobs to a channel.
	cronJobs := make(chan *batchv1beta1.CronJob, 1)
	informers := informers.NewSharedInformerFactory(client, 0)
	cronJobInformer := informers.Batch().V1().Jobs().Informer()
	cronJobInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			cronJob := obj.(*batchv1beta1.CronJob)
			t.Logf("cronJob added: %s/%s", cronJob.Namespace, cronJob.Name)
			cronJobs <- cronJob
		},
	})

	// Make sure informers are running.
	informers.Start(ctx.Done())

	// This is not required in tests, but it serves as a proof-of-concept by
	// ensuring that the informer goroutine have warmed up and called List before
	// we send any events to it.
	cache.WaitForCacheSync(ctx.Done(), cronJobInformer.HasSynced)

	// Inject an event into the fake client.
	gocronJob, err := New("my-cronjob", "test-ns", "hello-world", "*/1 * * * *")
	if err != nil {
		t.Fatalf("error instancing cronJob: %v", err)
	}
	p := &batchv1beta1.CronJob{ObjectMeta: gocronJob.ObjectMeta, Spec: gocronJob.CronJobSpec}
	_, err = client.BatchV1beta1().CronJobs("test-ns").Create(context.TODO(), p, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("error injecting cronJob add: %v", err)
	}
	cronJobList, err := client.BatchV1beta1().CronJobs("test-ns").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		t.Errorf("Error %v", err)
	}
	for _, cronJob := range cronJobList.Items {
		t.Logf("%v", cronJob)
	}

	// select {
	// case cronJob := <-cronJobs:
	// 	t.Logf("Got cronJob from channel: %s/%s", cronJob.Namespace, cronJob.Name)
	// case <-time.After(wait.ForeverTestTimeout):
	// 	t.Error("Informer did not get the added cronJob")
	// }
}
