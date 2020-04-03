package androidjob

import (
	"context"
	"time"

	androidv1alpha1 "github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var ttlLog = logf.Log.WithName("controller_androidjob_ttl")

func runTTLController(c client.Client) {
	ticker := time.NewTicker(time.Duration(10) * time.Second)
	for range ticker.C {
		if err := checkTTLs(c); err != nil {
			ttlLog.Error(err, "Failed to run TTL check for AndroidJobs")
		}
	}
}

func checkTTLs(c client.Client) error {
	jobList := &androidv1alpha1.AndroidJobList{}
	if err := c.List(context.TODO(), jobList, client.InNamespace(metav1.NamespaceAll)); err != nil {
		return err
	}
	for _, job := range jobList.Items {
		if job.Spec.TTLSecondsAfterCreation != nil {
			ttlSeconds := time.Duration(*job.Spec.TTLSecondsAfterCreation) * time.Second
			if time.Since(job.GetCreationTimestamp().Time) >= ttlSeconds {
				ttlLog.Info("Job is past TTL since creation", "Job.Name", job.Name, "Job.Namespace", job.Namespace)
				if err := c.Delete(context.TODO(), &job); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
