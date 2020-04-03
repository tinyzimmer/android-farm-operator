package emulators

import (
	"context"
	"fmt"

	androidv1alpha1 "github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// appendKVMVolume attaches the kvm device to an emulator pod
// TODO: This could be done through getters in the API
func appendKVMVolume(namePrefix string, volumes []corev1.Volume, volumeMounts []corev1.VolumeMount) ([]corev1.Volume, []corev1.VolumeMount) {
	volName := fmt.Sprintf("%s-kvm", namePrefix)
	volumes = append(volumes, corev1.Volume{
		Name: volName,
		VolumeSource: corev1.VolumeSource{
			HostPath: &corev1.HostPathVolumeSource{
				Path: "/dev/kvm",
			},
		},
	})
	volumeMounts = append(volumeMounts, corev1.VolumeMount{
		Name:      volName,
		MountPath: "/dev/kvm",
	})
	return volumes, volumeMounts
}

// appendPVCToPod will append a persistent volume claim to an emulator pod.
// TODO: This could also be done through getters in the API
func appendPVCToPod(pod *corev1.Pod, pvc *corev1.PersistentVolumeClaim, vol androidv1alpha1.Volume) *corev1.Pod {
	pod.Spec.Volumes = append(pod.Spec.Volumes, corev1.Volume{
		Name: pvc.Name,
		VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: pvc.Name,
			},
		},
	})
	pod.Spec.Containers[0].VolumeMounts = append(pod.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
		Name:      pvc.Name,
		MountPath: vol.MountPoint,
	})
	return pod
}

// reconcilePVCForPod will ensure a PVC for an emulator pod based off the user
// provided spec.
// TODO: Make generic and move to util package
func reconcilePVCForPod(reqLogger logr.Logger, c client.Client, pod *corev1.Pod, volume androidv1alpha1.Volume, labels map[string]string) (*corev1.PersistentVolumeClaim, error) {
	volName := fmt.Sprintf("%s%s", volume.VolumePrefix, pod.Name)
	pvc := &corev1.PersistentVolumeClaim{}
	if err := c.Get(context.TODO(), types.NamespacedName{Name: volName, Namespace: pod.Namespace}, pvc); err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("Creating new PVC for emulator pod", "Pod.Name", pod.Name, "Pod.Namespace", pod.Namespace, "PVC.Name", volName)
			pvc = &corev1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name:      volName,
					Namespace: pod.Namespace,
					Labels:    labels,
				},
				Spec: volume.PVCSpec,
			}
			if err := c.Create(context.TODO(), pvc); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		reqLogger.Info("Using existing PVC for emulator pod", "Pod.Name", pod.Name, "Pod.Namespace", pod.Namespace, "PVC.Name", pvc.Name)
	}
	return pvc, nil
}
