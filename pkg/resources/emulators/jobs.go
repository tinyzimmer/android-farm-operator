package emulators

//
// import (
// 	"context"
// 	"fmt"
//
// 	androidv1alpha1 "github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1"
// 	"github.com/tinyzimmer/android-farm-operator/pkg/util"
// 	corev1 "k8s.io/api/core/v1"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/apimachinery/pkg/types"
// 	"sigs.k8s.io/controller-runtime/pkg/client"
// )
//
// // reconcileStatupJobs is super WIP. It reconciles startup jobs for devices
// // after they've been created.
// // TODO: This method should verify boot readiness.
// func reconcileStartupJobs(c client.Client, config *androidv1alpha1.AndroidDeviceConfig, pod *corev1.Pod) error {
// 	var labels map[string]string
// 	if pod.Labels == nil {
// 		labels = make(map[string]string)
// 	} else {
// 		labels = pod.Labels
// 	}
//
// 	// we need to fetch the actual pod so we can use its UUID
// 	runningDevice := &corev1.Pod{}
// 	if err := c.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, runningDevice); err != nil {
// 		return err
// 	}
//
// 	// set labels for the jobs
// 	labels["targetDevice"] = pod.Name
// 	labels["targetUUID"] = string(runningDevice.GetUID())
//
// 	for _, startupJob := range config.Spec.StartupJobTemplates {
// 		job := startupJobForDevice(runningDevice, startupJob, labels)
// 		found := &androidv1alpha1.AndroidJob{}
// 		if err := c.Get(context.TODO(), types.NamespacedName{Name: job.Name, Namespace: job.Namespace}, found); err != nil {
// 			// If an API failure return
// 			if client.IgnoreNotFound(err) != nil {
// 				return err
// 			}
// 			// If doesn't exist, go ahead and create it
// 			if err := c.Create(context.TODO(), job); err != nil {
// 				return err
// 			}
// 			// Move on to the next job
// 			continue
// 		}
// 		// The job exists, let's make sure it was for the right device
// 		if targetUUID, ok := found.Labels["targetUUID"]; ok {
// 			if targetUUID != string(runningDevice.GetUID()) {
// 				// This is garbage, should have been collected with a pod, delete it
// 				if err := c.Delete(context.TODO(), found); err != nil {
// 					return err
// 				}
// 				// Make a new one
// 				if err := c.Create(context.TODO(), job); err != nil {
// 					return err
// 				}
// 			}
// 		} else {
// 			// No TargetUUID was present, again we delete and make a new one
// 			if err := c.Delete(context.TODO(), found); err != nil {
// 				return err
// 			}
// 			// Make a new one
// 			if err := c.Create(context.TODO(), job); err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	return nil
// }
//
// // startupJobForDevice returns a startup job for a device that the job controller
// // can execute.
// func startupJobForDevice(pod *corev1.Pod, startupJob string, labels map[string]string) *androidv1alpha1.AndroidJob {
// 	return &androidv1alpha1.AndroidJob{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      fmt.Sprintf("%s-%s", pod.Name, startupJob),
// 			Namespace: pod.Namespace,
// 			Labels:    labels,
// 			OwnerReferences: []metav1.OwnerReference{
// 				{
// 					APIVersion:         pod.APIVersion,
// 					Kind:               pod.Kind,
// 					Name:               pod.Name,
// 					UID:                pod.UID,
// 					Controller:         util.BoolPointer(true),
// 					BlockOwnerDeletion: util.BoolPointer(true),
// 				},
// 			},
// 		},
// 		Spec: androidv1alpha1.AndroidJobSpec{
// 			DeviceName:  pod.Name,
// 			JobTemplate: startupJob,
// 		},
// 	}
// }
