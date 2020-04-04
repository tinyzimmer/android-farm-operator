package v1alpha1

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

// Checksum returns the SHA256 hash of the current android device.
// Since this instance will often get merged with another, this is helpful
// for storing the state of the merged configuration.
func (c *AndroidDeviceConfig) Checksum() (string, error) {
	out, err := json.Marshal(c.Spec)
	if err != nil {
		return "", err
	}
	h := sha256.New()
	if _, err := h.Write(out); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// GetEnvVars returns the environment variables that should be used in pods
// using this device configuration.
func (c *AndroidDeviceConfig) GetEnvVars() []corev1.EnvVar {
	envVars := []corev1.EnvVar{
		{
			Name:  "ANDROID_ARCH",
			Value: "x86",
		},
		{
			Name: "POD_IP",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "status.podIP",
				},
			},
		},
	}
	if c.Spec.ExtraEnvVars != nil {
		envVars = append(envVars, c.Spec.ExtraEnvVars...)
	}
	return envVars
}

// IsKVMEnabled returns true if the device requires KVM virtualization
// TODO: This is probably not even necessary, they all do.
func (c *AndroidDeviceConfig) IsKVMEnabled() bool {
	return c.Spec.KVMEnabled
}

// GetPorts returns the container ports to be used for pods using this
// configuration.
func (c *AndroidDeviceConfig) GetPorts() []corev1.ContainerPort {
	ports := make([]corev1.ContainerPort, 0)
	if c.Spec.TCPRedir == nil || !c.Spec.TCPRedir.Enabled {
		ports = append(ports, c.GetADBContainerPort()...)
	}
	if c.Spec.ExtraPorts != nil {
		ports = append(ports, c.Spec.ExtraPorts...)
	}
	return ports
}

// GetADBContainerPort returns the adb port to be used by the redir sidecar.
func (c *AndroidDeviceConfig) GetADBContainerPort() []corev1.ContainerPort {
	return []corev1.ContainerPort{
		{
			Name:          "adb",
			ContainerPort: c.GetADBPort(),
		},
	}
}

// GetServicePorts returns the service ports that should be used for the headless
// service in front of a device group using this configuration.
func (c *AndroidDeviceConfig) GetServicePorts() []corev1.ServicePort {
	ports := []corev1.ServicePort{
		{Name: "adb", Port: c.GetADBPort()},
	}
	for _, port := range c.Spec.ExtraPorts {
		ports = append(ports, corev1.ServicePort{
			Name: port.Name,
			Port: port.ContainerPort,
		})
	}
	return ports
}

// GetADBPort returns the ADB port to be used by devices using this configuration.
func (c *AndroidDeviceConfig) GetADBPort() int32 {
	if c.Spec.ADBPort == 0 {
		return 5555
	}
	return c.Spec.ADBPort
}

// GetImagePullSecrets will return all image pull secrets for emulator pods using
// this configuration
func (c *AndroidDeviceConfig) GetImagePullSecrets() []corev1.LocalObjectReference {
	secrets := make([]corev1.LocalObjectReference, 0)
	if c.Spec.ImagePullSecrets != nil {
		secrets = append(secrets, c.Spec.ImagePullSecrets...)
	}
	if c.Spec.TCPRedir != nil && c.Spec.TCPRedir.ImagePullSecrets != nil {
		secrets = append(secrets, c.Spec.TCPRedir.ImagePullSecrets...)
	}
	return secrets
}

// GetRedirImage will return the configured redir image for these emulators.
func (c *AndroidDeviceConfig) GetRedirImage() string {
	if c.Spec.TCPRedir != nil && c.Spec.TCPRedir.Image != "" {
		return c.Spec.TCPRedir.Image
	}
	return "quay.io/tinyzimmer/goredir:latest"
}

// GetSidecars returns sidecars, if any, to run in the emulator pods.
func (c *AndroidDeviceConfig) GetSidecars() []corev1.Container {
	if c.Spec.TCPRedir != nil && c.Spec.TCPRedir.Enabled {
		return []corev1.Container{
			{
				Name:            "redir",
				Image:           c.GetRedirImage(),
				Ports:           c.GetADBContainerPort(),
				ImagePullPolicy: "IfNotPresent",
				Env:             c.GetEnvVars(),
				Args: []string{
					fmt.Sprintf("-target=127.0.0.1:%d", c.GetADBPort()),
					"-host=$(POD_IP)",
				},
			},
		}
	}
	return []corev1.Container{}
}

// MergeInto merges the contents of this spec into a provided configuration,
// granting precedence to the target spec. A pointer to the final product is
// returned.
func (c *AndroidDeviceConfigSpec) MergeInto(target AndroidDeviceConfigSpec) (*AndroidDeviceConfigSpec, error) {
	jb, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jb, &target)
	if err != nil {
		return nil, err
	}
	return &target, nil
}
