package v1alpha1

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

// ProviderNodeSelector returns the node selector to use when scheduling ADB
// and provider deployments.
func (f *DeviceGroup) ProviderNodeSelector() map[string]string {
	if f.HostUSB != nil {
		if f.HostUSB.NodeName != "" {
			return map[string]string{
				"kubernetes.io/hostname": f.HostUSB.NodeName,
			}
		}
	}
	return nil
}

// MaxUSBDevices returns the maximum number of usb devices we expect to run
// on a given node.
func (f *DeviceGroup) MaxUSBDevices() int32 {
	if f.HostUSB != nil {
		if f.HostUSB.MaxDevices != 0 {
			return f.HostUSB.MaxDevices
		}
		return 1
	}
	return 0
}

// GetProviderName returns the name of the provider for this device group.
func (f *DeviceGroup) GetProviderName() string {
	return fmt.Sprintf("provider-%s", f.Name)
}

// ProviderNoCleanup returns true if providers in this group should persist
// device state.
func (f *DeviceGroup) ProviderNoCleanup() bool {
	if f.Provider != nil {
		return f.Provider.PersistDeviceState
	}
	return false
}

// GetProviderStartPort returns the starting port to use for the device group's
// provider instance.
func (f *DeviceGroup) GetProviderStartPort() int32 {
	if f.Provider == nil || f.Provider.StartPort == 0 {
		return int32(15000)
	}
	return f.Provider.StartPort
}

// ADBPodSecurityContext returns the pod security context to use for adb deployments
// in this farm.
func (s *STFConfig) ADBPodSecurityContext(group *DeviceGroup) *corev1.PodSecurityContext {
	if group.HostUSB != nil {
		return &corev1.PodSecurityContext{
			RunAsNonRoot: &falseVal,
		}
	}
	return &corev1.PodSecurityContext{
		RunAsUser: &defaultRunUser,
	}
}

// ADBContainerSecurityContext returns the container security context to use for
// adb deployments in this farm.
func (s *STFConfig) ADBContainerSecurityContext(group *DeviceGroup) *corev1.SecurityContext {
	if group.HostUSB != nil {
		return &corev1.SecurityContext{
			Privileged: &trueVal,
		}
	}
	return nil
}

// ADBImage returns the image to use for adb servers in this farm.
func (s *STFConfig) ADBImage(group *DeviceGroup) string {
	if group.ADB != nil && group.ADB.Image != "" {
		return group.ADB.Image
	}
	if s.ADB != nil && s.ADB.Image != "" {
		return s.ADB.Image
	}
	return "quay.io/tinyzimmer/adbmon:latest"
}

// ADBImagePullPolicy returns the pull policy for adb images in this farm
func (s *STFConfig) ADBImagePullPolicy(group *DeviceGroup) corev1.PullPolicy {
	if group.ADB != nil && group.ADB.ImagePullPolicy != "" {
		return group.ADB.ImagePullPolicy
	}
	if s.ADB != nil && s.ADB.ImagePullPolicy != "" {
		return s.ADB.ImagePullPolicy
	}
	return corev1.PullIfNotPresent
}

// ADBSidecarContainer returns the container definition for the provider adb sidecars.
func (s *STFConfig) ADBSidecarContainer(providerName string, group *DeviceGroup) corev1.Container {
	container := corev1.Container{
		Name:            "adb",
		ImagePullPolicy: s.ADBImagePullPolicy(group),
		Image:           s.ADBImage(group),
		Ports:           []corev1.ContainerPort{{Name: "adb-server", ContainerPort: 5037}},
		SecurityContext: s.ADBContainerSecurityContext(group),
		Resources:       s.ADBResourceRequirements(group),
		Env: []corev1.EnvVar{
			{
				Name:  "HOME",
				Value: "/tmp",
			},
		},
	}
	if container.Image == "quay.io/tinyzimmer/adbmon:latest" {
		if group.IsEmulatedGroup() {
			container.Args = []string{"--no-usb", "--provider", providerName}
		} else if group.IsUSBGroup() {
			container.Args = []string{"--provider", providerName}
		}
	}
	if group.IsUSBGroup() {
		container.VolumeMounts = []corev1.VolumeMount{
			{
				Name:      "usb",
				MountPath: "/dev/bus/usb",
			},
		}
	}
	return container
}
