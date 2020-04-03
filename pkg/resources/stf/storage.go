package stf

import (
	"context"

	androidv1alpha1 "github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1"
	"github.com/tinyzimmer/android-farm-operator/pkg/util/builders"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *STFReconciler) reconcileStorage(reqLogger logr.Logger, instance *androidv1alpha1.AndroidFarm) error {
	// Reconcile storage volume if needed
	if volume := instance.STFStorageVolumeClaim(); volume != nil {
		reqLogger.Info("Ensuring PVC for STF storage")
		foundPVC := &corev1.PersistentVolumeClaim{}
		if err := r.client.Get(context.TODO(), types.NamespacedName{Name: volume.Name, Namespace: volume.Namespace}, foundPVC); err != nil {
			if client.IgnoreNotFound(err) != nil {
				return err
			}
			if err := r.client.Create(context.TODO(), volume); err != nil {
				return err
			}
		}
	}

	return builders.NewDeploymentBuilder(reqLogger, instance, "storage").
		WithReplicas(instance.STFConfig().StorageReplicas()).
		WithVolumes(instance.STFStorageVolumes(), instance.STFStorageVolumeMounts()).
		WithPort("stf-storage", 3000).
		WithResourceRequirements(instance.STFConfig().StorageResourceRequirements()).
		WithPodSecurityContext(instance.STFConfig().PodSecurityContext()).
		WithContainerSecurityContext(instance.STFConfig().ContainerSecurityContext()).
		WithArgs([]string{
			"storage-temp",
			"--port", "3000",
			"--save-dir", "/data",
		}).
		WithService("ClusterIP").
		WithWait().
		Reconcile(r.client)
}

func (r *STFReconciler) reconcileAPKStorage(reqLogger logr.Logger, instance *androidv1alpha1.AndroidFarm) error {
	return builders.NewDeploymentBuilder(reqLogger, instance, "apk-storage").
		WithPort("stf-apk-storage", 3000).
		WithResourceRequirements(instance.STFConfig().APKStorageResourceRequirements()).
		WithArgs([]string{
			"storage-plugin-apk",
			"--port", "3000",
			// "--storage-url", fmt.Sprintf("%s://%s", instance.STFConfig().GetHTTPScheme(), instance.STFConfig().GetAppHostname()),
			"--storage-url", instance.InternalStorageURL(),
		}).
		WithService("ClusterIP").
		WithWait().
		Reconcile(r.client)
}

func (r *STFReconciler) reconcileImageStorage(reqLogger logr.Logger, instance *androidv1alpha1.AndroidFarm) error {
	return builders.NewDeploymentBuilder(reqLogger, instance, "img-storage").
		WithPort("img-storage", 3000).
		WithResourceRequirements(instance.STFConfig().ImgStorageResourceRequirements()).
		WithArgs([]string{
			"storage-plugin-image",
			"--port", "3000",
			// "--storage-url", fmt.Sprintf("%s://%s", instance.STFConfig().GetHTTPScheme(), instance.STFConfig().GetAppHostname()),
			"--storage-url", instance.InternalStorageURL(),
		}).
		WithService("ClusterIP").
		WithWait().
		Reconcile(r.client)
}
