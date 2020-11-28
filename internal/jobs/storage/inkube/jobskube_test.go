package inkube

import (
	"testing"

	"github.com/kurkop/gojobs/shared/kube"
)

func TestCreate(t *testing.T) {
	config := kube.GetLocalConfig()

	// create the clientset
	client, err := kube.NewClient(config)

	goJobRepo := NewGoJobsRepository(client)
	goJobCreated, err := goJobRepo.Create("", "job-test", "hello-world")
	if err != nil {
		t.Fatalf("error creating job: %v", err)
	}
	t.Logf("goJob created %v", goJobCreated)
	goJobGot, err := goJobRepo.Get(goJobCreated.GetName())
	if err != nil {
		t.Fatalf("error getting job: %v", err)
	}
	goJobsGot, err := goJobRepo.GetAll()
	t.Logf("goJobs got %v", goJobsGot)
	goJobRepo.Delete(goJobGot.GetName())
}

// // TestFakeJob demonstrates how to use a fake client with SharedInformerFactory in tests.
// func TestFakeCreateGetDelete(t *testing.T) {
// 	// Create the fake client.
// 	client := fake.NewSimpleClientset()

// 	goJobRepo := NewGoJobsRepository(client)
// 	goJobCreated, err := goJobRepo.Create("", "job-test", "hello-world")
// 	if err != nil {
// 		t.Fatalf("error creating job: %v", err)
// 	}
// 	t.Logf("goJob created %v", goJobCreated)
// 	goJobGot, err := goJobRepo.Get(goJobCreated.GetName())
// 	if err != nil {
// 		t.Fatalf("error getting job: %v", err)
// 	}
// 	t.Logf("goJob got %v", goJobGot)
// 	goJobRepo.Delete(goJobGot.GetName(), goJobGot.GetNamespace())
// }
