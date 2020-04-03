package emulators

import (
	androidv1alpha1 "github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1"
	"github.com/tinyzimmer/android-farm-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// serviceForDeviceGroup returns a headless service definition for the given farm
// group and configuration.
func serviceForDeviceGroup(farm *androidv1alpha1.AndroidFarm, group *androidv1alpha1.DeviceGroup, conf *androidv1alpha1.AndroidDeviceConfig) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      group.GetSubdomain(),
			Namespace: group.GetNamespace(),
			Labels:    util.DeviceFarmLabels(farm, group),
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         farm.APIVersion,
					Kind:               farm.Kind,
					Name:               farm.Name,
					UID:                farm.UID,
					Controller:         util.BoolPointer(true),
					BlockOwnerDeletion: util.BoolPointer(true),
				},
			},
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: "None",
			Ports:     conf.GetServicePorts(),
			Selector:  util.DeviceFarmLabels(farm, group),
		},
	}
}
