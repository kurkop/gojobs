package inkube

import (
	"testing"

	"github.com/kurkop/gojob/shared/kube"
)

func TestCreate(t *testing.T) {
	config := kube.GetLocalConfig()

	// create the clientset
	client, err := kube.NewClient(config)

	goCronJobRepo := NewGoCronJobsRepository(client)
	goCronJobCreated, err := goCronJobRepo.Create("cronjob-test", "default", "hello-world", "*/1 * * * *")
	if err != nil {
		t.Fatalf("error creating cronJob: %v", err)
	}
	t.Logf("goCronJob created %v", goCronJobCreated)
	goCronJobGot, err := goCronJobRepo.Get(goCronJobCreated.GetName(), goCronJobCreated.GetNamespace())
	if err != nil {
		t.Fatalf("error getting cronJob: %v", err)
	}
	t.Logf("goCronJob got %v", goCronJobGot)
	goCronJobsGot, err := goCronJobRepo.GetAll(goCronJobCreated.GetNamespace())
	t.Logf("goCronJobs got %v", goCronJobsGot)
	goCronJobRepo.Delete(goCronJobGot.GetName(), goCronJobGot.GetNamespace())
}

// // TestFakeJob demonstrates how to use a fake client with SharedInformerFactory in tests.
// func TestFakeCreateGetDelete(t *testing.T) {
// 	// Create the fake client.
// 	client := fake.NewSimpleClientset()

// 	goCronJobRepo := NewGoCronJobsRepository(client)
// 	goCronJobCreated, err := goCronJobRepo.Create("cronjob-test", "default", "hello-world", "*/1 * * * *")
// 	if err != nil {
// 		t.Fatalf("error creating cronJob: %v", err)
// 	}
// 	t.Logf("goCronJob created %v", goCronJobCreated)
// 	goCronJobGot, err := goCronJobRepo.Get(goCronJobCreated.GetName(), goCronJobCreated.GetNamespace())
// 	if err != nil {
// 		t.Fatalf("error getting cronJob: %v", err)
// 	}
// 	t.Logf("goCronJob got %v", goCronJobGot)
// 	goCronJobRepo.Delete(goCronJobGot.GetName(), goCronJobGot.GetNamespace())
// }
