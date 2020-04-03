package util

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	androidv1alpha1 "github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1"
	"github.com/tinyzimmer/android-farm-operator/pkg/util/errors"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ReconcileDevice reconciles an AndroidDevice CR with the cluster. The checkUpdate
// function provided will be called if an update is required. If the function returns
// false, the request is requeued.
func ReconcileDevice(reqLogger logr.Logger, c client.Client, device *androidv1alpha1.AndroidDevice, checkUpdate func(string) bool) error {
	if err := SetCreationSpecAnnotation(&device.ObjectMeta, device); err != nil {
		return err
	}
	found := &androidv1alpha1.AndroidDevice{}
	if err := c.Get(context.TODO(), types.NamespacedName{Name: device.Name, Namespace: device.Namespace}, found); err != nil {
		// Return API error
		if client.IgnoreNotFound(err) != nil {
			return err
		}
		// Create the device
		reqLogger.Info("Creating new device", "Device.Name", device.Name, "Device.Namespace", device.Namespace)
		if err := c.Create(context.TODO(), device); err != nil {
			return err
		}
		return nil
	}

	// Check the found device spec
	if !CreationSpecsEqual(device.ObjectMeta, found.ObjectMeta) {
		// Check if we are allowed to update
		if cont := checkUpdate(device.Annotations[androidv1alpha1.DeviceConfigSHAAnnotation]); !cont {
			return errors.NewRequeueError("Device is not ready to be updated", 3)
		}
		// We need to update the device
		reqLogger.Info("Device annotation spec has changed, updating", "Device.Name", device.Name, "Device.Namespace", device.Namespace)
		// will requeue the farm that made us
		if err := c.Delete(context.TODO(), found); err != nil {
			return err
		}
	}

	return nil
}

// ReconcileService will reconcile a provided service spec with the cluster.
func ReconcileService(reqLogger logr.Logger, c client.Client, svc *corev1.Service) error {
	if err := SetCreationSpecAnnotation(&svc.ObjectMeta, svc); err != nil {
		return err
	}
	found := &corev1.Service{}
	if err := c.Get(context.TODO(), types.NamespacedName{Name: svc.Name, Namespace: svc.Namespace}, found); err != nil {
		// Return API error
		if client.IgnoreNotFound(err) != nil {
			return err
		}
		// Create the service
		reqLogger.Info("Creating new service", "Service.Name", svc.Name, "Service.Namespace", svc.Namespace)
		if err := c.Create(context.TODO(), svc); err != nil {
			return err
		}
		return nil
	}

	// Check the found service spec
	if !CreationSpecsEqual(svc.ObjectMeta, found.ObjectMeta) {
		// We need to update the service
		reqLogger.Info("Service annotation spec has changed, updating", "Service.Name", svc.Name, "Service.Namespace", svc.Namespace)
		svc.Spec.ClusterIP = found.Spec.ClusterIP
		found.Spec = svc.Spec
		if err := c.Update(context.TODO(), found); err != nil {
			return err
		}
		return nil
	}

	return nil
}

// ReconcilePod will reconcile a given pod definition with the cluster. If a pod
// with the same name exists but has a different configuration, the pod will be
// deleted and requeued. If the found pod has a deletion timestamp (e.g. it is still terminating)
// then the request will also be requued.
func ReconcilePod(reqLogger logr.Logger, c client.Client, pod *corev1.Pod) (bool, error) {
	if err := SetCreationSpecAnnotation(&pod.ObjectMeta, pod); err != nil {
		return false, err
	}
	found := &corev1.Pod{}
	if err := c.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found); err != nil {
		// Return API error
		if client.IgnoreNotFound(err) != nil {
			return false, err
		}
		// Create the Pod
		reqLogger.Info("Creating new Pod", "Pod.Name", pod.Name, "Pod.Namespace", pod.Namespace)
		if err := c.Create(context.TODO(), pod); err != nil {
			return false, err
		}
		return true, nil
	}

	// Check if the found pod is in the middle of terminating
	if found.GetDeletionTimestamp() != nil {
		return false, errors.NewRequeueError("Existing pod is still being terminated, requeuing", 3)
	}

	// Check the found pod spec
	if !CreationSpecsEqual(pod.ObjectMeta, found.ObjectMeta) {
		// We need to delete the pod and return a requeue
		if err := c.Delete(context.TODO(), found); err != nil {
			return false, err
		}
		return false, errors.NewRequeueError("Pod spec has changed, recreating", 3)
	}

	return false, nil
}

// ReconcileDeployment reconciles a given deployment configuration with the cluster.
// If wait is true, then the deployment is checked for readiness and requeued if
// not all replicas are available.
func ReconcileDeployment(reqLogger logr.Logger, c client.Client, deployment *appsv1.Deployment, wait bool) error {
	if err := SetCreationSpecAnnotation(&deployment.ObjectMeta, deployment); err != nil {
		return err
	}

	foundDeployment := &appsv1.Deployment{}
	if err := c.Get(context.TODO(), types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, foundDeployment); err != nil {
		// Return API error
		if client.IgnoreNotFound(err) != nil {
			return err
		}
		// Create the deployment
		reqLogger.Info("Creating new deployment", "Deployment.Name", deployment.Name, "Deployment.Namespace", deployment.Namespace)
		if err := c.Create(context.TODO(), deployment); err != nil {
			return err
		}
		if wait {
			return errors.NewRequeueError("Created new deployment with wait, requeing for status check", 3)
		}
		return nil
	}

	// Check the found deployment spec
	if !CreationSpecsEqual(deployment.ObjectMeta, foundDeployment.ObjectMeta) {
		// We need to update the deployment
		reqLogger.Info("Deployment annotation spec has changed, updating", "Deployment.Name", deployment.Name, "Deployment.Namespace", deployment.Namespace)
		foundDeployment.Spec = deployment.Spec
		if err := c.Update(context.TODO(), foundDeployment); err != nil {
			return err
		}
	}

	if wait {
		runningDeploy := &appsv1.Deployment{}
		if err := c.Get(context.TODO(), types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, runningDeploy); err != nil {
			return err
		}
		if runningDeploy.Status.ReadyReplicas != *deployment.Spec.Replicas {
			return errors.NewRequeueError(fmt.Sprintf("Waiting for %s to be ready", deployment.Name), 3)
		}
	}

	return nil
}

// ReconcileConfigMap will reconcile a given configmap definition with the cluster.
func ReconcileConfigMap(reqLogger logr.Logger, c client.Client, cm *corev1.ConfigMap) error {
	if err := SetCreationSpecAnnotation(&cm.ObjectMeta, cm); err != nil {
		return err
	}

	foundCM := &corev1.ConfigMap{}
	if err := c.Get(context.TODO(), types.NamespacedName{Name: cm.Name, Namespace: cm.Namespace}, foundCM); err != nil {
		// Return API error
		if client.IgnoreNotFound(err) != nil {
			return err
		}
		// Create the configmap
		reqLogger.Info("Creating new configmap", "ConfigMap.Name", cm.Name, "ConfigMap.Namespace", cm.Namespace)
		if err := c.Create(context.TODO(), cm); err != nil {
			return err
		}
		return nil
	}

	// Check the found configmap data
	if !CreationSpecsEqual(cm.ObjectMeta, foundCM.ObjectMeta) {
		// We need to update the configmap
		reqLogger.Info("Configmap data has changed, updating", "ConfigMap.Name", cm.Name, "ConfigMap.Namespace", cm.Namespace)
		foundCM.Data = cm.Data
		if err := c.Update(context.TODO(), foundCM); err != nil {
			return err
		}
	}

	return nil
}

// ReconcileJob will ensure a job runs on the cluster. If wait is true, the request
// will be requeued if the job exists and has not yet completed.
func ReconcileJob(reqLogger logr.Logger, c client.Client, job *batchv1.Job, wait bool) error {
	if err := SetCreationSpecAnnotation(&job.ObjectMeta, job); err != nil {
		return err
	}

	found := &batchv1.Job{}
	if err := c.Get(context.TODO(), types.NamespacedName{Name: job.Name, Namespace: job.Namespace}, found); err != nil {
		// Return API error
		if client.IgnoreNotFound(err) != nil {
			return err
		}
		// Create the job
		if err := c.Create(context.TODO(), job); err != nil {
			return err
		}
	}

	// Check the found job spec
	if !CreationSpecsEqual(job.ObjectMeta, found.ObjectMeta) {
		// We need to delete the job and requeue
		if err := c.Delete(context.TODO(), found); err != nil {
			return err
		}
		if wait {
			return errors.NewRequeueError("Deleting stale job and requeuing", 5)
		}
	}

	if wait {
		// Make sure the job completed
		if found.Status.Succeeded != 1 {
			return errors.NewRequeueError("Waiting for job to complete", 3)
		}
	}

	return nil
}

// ReconcileStatefulSet reconciles a statefulset configuration with the cluster.
func ReconcileStatefulSet(reqLogger logr.Logger, c client.Client, ss *appsv1.StatefulSet) error {
	if err := SetCreationSpecAnnotation(&ss.ObjectMeta, ss); err != nil {
		return err
	}

	found := &appsv1.StatefulSet{}
	if err := c.Get(context.TODO(), types.NamespacedName{Name: ss.Name, Namespace: ss.Namespace}, found); err != nil {
		// Return API error
		if client.IgnoreNotFound(err) != nil {
			return err
		}
		// Create the statefulset
		if err := c.Create(context.TODO(), ss); err != nil {
			return err
		}
		return nil
	}

	// Check the found statefulset spec
	if !CreationSpecsEqual(ss.ObjectMeta, found.ObjectMeta) {
		// We need to update the statefulset
		found.Spec = ss.Spec
		if err := c.Update(context.TODO(), found); err != nil {
			return err
		}
	}

	return nil
}
