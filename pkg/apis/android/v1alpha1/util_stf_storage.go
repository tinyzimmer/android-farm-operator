package v1alpha1

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// STFStorageVolumeClaim returns the PVC definition that just be provisioned
// with the STF storage service, or nil if no PVCSpec is provided.
func (a *AndroidFarm) STFStorageVolumeClaim() *corev1.PersistentVolumeClaim {
	if a.STFConfig() != nil {
		if a.STFConfig().Storage != nil {
			if a.STFConfig().Storage.PVCSpec != nil {
				return &corev1.PersistentVolumeClaim{
					ObjectMeta: metav1.ObjectMeta{
						Name:            fmt.Sprintf("%s-stf-storage", a.GetName()),
						Namespace:       a.STFConfig().GetNamespace(),
						Labels:          a.STFComponentLabels("storage"),
						OwnerReferences: a.OwnerReferences(),
					},
					Spec: *a.STFConfig().Storage.PVCSpec,
				}
			}
		}
	}
	return nil
}

// STFStorageVolumes returns the deployment volume definitions to be used for
// the STF storage deployment, or an empty volume slice if no PVCspec is provided.
func (a *AndroidFarm) STFStorageVolumes() []corev1.Volume {
	if a.STFConfig() != nil {
		if a.STFConfig().Storage != nil {
			if a.STFConfig().Storage.PVCSpec != nil {
				return []corev1.Volume{
					{
						Name: "stf-storage",
						VolumeSource: corev1.VolumeSource{
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
								ClaimName: fmt.Sprintf("%s-stf-storage", a.GetName()),
							},
						},
					},
				}
			}
		}
	}
	return []corev1.Volume{
		{
			Name: "stf-storage",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		},
	}
}

// STFStorageVolumes returns the deployment volume mounts to be used for
// the STF storage deployment, or an empty slice if no PVCspec is provided.
func (a *AndroidFarm) STFStorageVolumeMounts() []corev1.VolumeMount {
	return []corev1.VolumeMount{{
		Name:      "stf-storage",
		MountPath: "/data",
	}}
}

func (a *AndroidFarm) InternalProxyHost() string {
	return fmt.Sprintf("%s.%s.svc", a.STFComponentName("traefik"), a.STFConfig().GetNamespace())
}

// InternalStorageURL returns the cluster local storage URL
func (a *AndroidFarm) InternalStorageURL() string {
	return fmt.Sprintf("http://%s:8880", a.InternalProxyHost())
}

// StorageReplicas returns the number of replicas to run in the storage deployment.
func (s *STFConfig) StorageReplicas() int32 {
	if s.Storage != nil {
		if s.Storage.Replicas != 0 {
			return s.Storage.Replicas
		}
	}
	return 1
}
