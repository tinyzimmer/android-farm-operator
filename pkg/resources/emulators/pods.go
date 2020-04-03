package emulators

import (
	androidv1alpha1 "github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1"
	"github.com/tinyzimmer/android-farm-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// newPodForCR returns a pod for a given emulator device and its configuration.
// A simple tcp redirection container is run as a sidecard to forward ADB traffic
// to the listener on the local pod's loopback interface.
func newPodForDevice(device *androidv1alpha1.AndroidDevice, conf *androidv1alpha1.AndroidDeviceConfig) *corev1.Pod {
	// TODO : Could write some more getters to handle most of this logic.
	volumes := []corev1.Volume{}
	volumeMounts := []corev1.VolumeMount{}
	securityContext := &corev1.SecurityContext{}
	podSecurityContext := &corev1.PodSecurityContext{}
	if conf.IsKVMEnabled() {
		volumes, volumeMounts = appendKVMVolume(device.Name, volumes, volumeMounts)
		securityContext.Privileged = util.BoolPointer(true)
		podSecurityContext.RunAsNonRoot = util.BoolPointer(false)
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      device.Name,
			Namespace: device.Namespace,
			Labels:    util.DeviceLabels(device),
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         device.APIVersion,
					Kind:               device.Kind,
					Name:               device.GetName(),
					UID:                device.GetUID(),
					Controller:         util.BoolPointer(true),
					BlockOwnerDeletion: util.BoolPointer(true),
				},
			},
		},
		Spec: corev1.PodSpec{
			Hostname:         device.Spec.Hostname,
			Subdomain:        device.Spec.Subdomain,
			ImagePullSecrets: conf.GetImagePullSecrets(),
			Volumes:          volumes,
			SecurityContext:  podSecurityContext,
			Containers: append([]corev1.Container{
				{
					Name:            device.GetName(),
					Command:         conf.Spec.Command,
					Args:            conf.Spec.Args,
					Image:           conf.Spec.DockerImage,
					ImagePullPolicy: conf.Spec.ImagePullPolicy,
					VolumeMounts:    volumeMounts,
					Ports:           conf.GetPorts(),
					SecurityContext: securityContext,
					Env:             conf.GetEnvVars(),
					Lifecycle: &corev1.Lifecycle{
						PreStop: &corev1.Handler{
							Exec: &corev1.ExecAction{
								Command: []string{"bash", "-c", "adb shell reboot -p"},
							},
						},
					},
					Resources: conf.Spec.Resources,
				},
			}, conf.GetSidecars()...),
		},
	}
}
