package util

import (
	androidv1alpha1 "github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1"
)

func IsFarmedDevice(cr *androidv1alpha1.AndroidDevice) bool {
	if cr.Labels != nil {
		if _, ok := cr.Labels[androidv1alpha1.DeviceGroupLabel]; ok {
			return true
		}
	}
	return false
}

func DeviceFarmLabels(cr *androidv1alpha1.AndroidFarm, group *androidv1alpha1.DeviceGroup) map[string]string {
	var labels map[string]string
	if cr.Labels != nil {
		labels = cr.Labels
	} else {
		labels = make(map[string]string)
	}
	labels[androidv1alpha1.DeviceFarmLabel] = cr.Name
	labels[androidv1alpha1.DeviceGroupLabel] = group.Name
	if group.Emulators != nil && group.Emulators.ConfigRef != nil && group.Emulators.ConfigRef.Name != "" {
		labels[androidv1alpha1.DeviceConfigLabel] = group.Emulators.ConfigRef.Name
	}
	return labels
}

func DeviceLabels(cr *androidv1alpha1.AndroidDevice) map[string]string {
	var labels map[string]string
	if cr.Labels != nil {
		labels = cr.Labels
	} else {
		labels = make(map[string]string)
	}
	return labels
}
