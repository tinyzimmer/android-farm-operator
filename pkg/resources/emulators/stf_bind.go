package emulators

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	androidv1alpha1 "github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1"
	"github.com/tinyzimmer/android-farm-operator/pkg/util"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// reconcileSTFBinding will ensure a job is run that binds a freshly booted
// emulator to its stf provider.
func reconcileSTFBinding(reqLogger logr.Logger, c client.Client, device *androidv1alpha1.AndroidDevice, pod *corev1.Pod) error {
	// If no adb server annotation, screw it
	if device.GetAnnotations() == nil {
		return nil
	}
	adbServer, ok := device.Annotations[androidv1alpha1.STFProviderAnnotation]
	if !ok {
		return nil
	}

	// check if connected already
	if connected, ok := device.Annotations[androidv1alpha1.ADBConnectedAnnotation]; ok {
		if connected == "true" {
			return nil
		}
	}

	// fetch the farm so we can use details about it
	farm, err := device.GetFarm(c)
	if err != nil {
		return err
	}

	adbPort, err := util.GetPodADBPort(*pod)
	if err != nil {
		return err
	}
	podSerial := fmt.Sprintf("%s:%d", getPodAddr(pod), adbPort)
	job := newSTFBindingJob(farm, pod, adbServer, podSerial)

	// ReconcileJob with true requeues until the job is finished
	if err := util.ReconcileJob(reqLogger, c, job, true); err != nil {
		return err
	}

	// fetch the completed job
	completedjob := &batchv1.Job{}
	if err := c.Get(context.TODO(), types.NamespacedName{Name: job.Name, Namespace: job.Namespace}, completedjob); err != nil {
		return err
	}

	// We can delete the job, and fallback to the TTL if we fail
	if err := c.Delete(context.TODO(), completedjob); err != nil {
		if client.IgnoreNotFound(err) != nil {
			reqLogger.Info("Could not clean up job, hopefully ttl will catch it")
		}
	}

	// also clean up pods
	if err := c.DeleteAllOf(context.TODO(), &corev1.Pod{}, client.InNamespace(completedjob.Namespace), client.MatchingLabels{"job-name": completedjob.Name}); err != nil {
		if client.IgnoreNotFound(err) != nil {
			reqLogger.Info("Could not clean up job pod(s), hopefully ttl will catch it")
		}
	}

	device.Annotations[androidv1alpha1.ADBConnectedAnnotation] = "true"
	device.Annotations[androidv1alpha1.ProviderSerialAnnotation] = podSerial
	return c.Update(context.TODO(), device)
}

func getPodAddr(pod *corev1.Pod) string {
	if pod.Spec.Hostname == "" || pod.Spec.Subdomain == "" {
		return pod.Status.PodIP
	}
	return fmt.Sprintf("%s.%s.%s", pod.Spec.Hostname, pod.Spec.Subdomain, pod.Namespace)
}

// We try to delete jobs after 10 minutes by default
var jobTTL int32 = 600

// newSTFBindingJob returns a new job definition for an emulator stf binding.
// Note that TTL seconds after completion is not respected properly by all
// kubernetes versions.
func newSTFBindingJob(cr *androidv1alpha1.AndroidFarm, device *corev1.Pod, adbServer, podSerial string) *batchv1.Job {
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-stf-connect", device.GetName()),
			Namespace: cr.STFConfig().GetNamespace(),
			Labels:    device.GetLabels(),
		},
		Spec: batchv1.JobSpec{
			TTLSecondsAfterFinished: &jobTTL,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: device.GetLabels(),
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: cr.STFConfig().GetServiceAccount(),
					SecurityContext:    cr.STFConfig().PodSecurityContext(),
					RestartPolicy:      "OnFailure",
					Containers: []corev1.Container{
						{
							Name:            "stf-connect",
							Image:           "quay.io/tinyzimmer/adbmon",
							ImagePullPolicy: "IfNotPresent",
							Args:            []string{"--host", adbServer, "--connect", podSerial, "--verbose"},
							SecurityContext: cr.STFConfig().ContainerSecurityContext(),
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									"cpu":    *resource.NewMilliQuantity(100, resource.DecimalSI),
									"memory": *resource.NewQuantity(128*1024*1024, resource.BinarySI),
								},
							},
						},
					},
				},
			},
		},
	}
}
