package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AndroidDeviceConfigSpec defines the desired state of AndroidDeviceConfig
type AndroidDeviceConfigSpec struct {
	// The docker image to use for emulator devices
	DockerImage string `json:"dockerImage,omitempty"`
	// The pull policy to use for emulator pods
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`
	// Pull secrets required for the docker image.
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`
	// The ADB port that the emulator listens on. Defaults to 5555.
	// A sidecar will be spawned within emulator pods that redirects external
	// traffic to this port.
	ADBPort int32 `json:"adbPort,omitempty"`
	// An optional command to run when starting an emulator image
	Command []string `json:"command,omitempty"`
	// Any arguments to pass to the above command.
	Args []string `json:"args,omitempty"`
	// Extra port mappings to apply to the emulator pods.
	ExtraPorts []corev1.ContainerPort `json:"extraPorts,omitempty"`
	// Extra environment variables to supply to the emulator pods.
	ExtraEnvVars []corev1.EnvVar `json:"extraEnvVars,omitempty"`
	// Whether to mount the kvm device to the pods, will require that the operator
	// can launch privileged pods.
	KVMEnabled bool `json:"kvmEnabled,omitempty"`
	// A list of volume configurations to apply to the emulator pods.
	Volumes []Volume `json:"volumes,omitempty"`
	// Resource restraints to place on the emulators.
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	// A list of AndroidJobTemplates to execute against new instances.
	// TODO: Very very very beta
	StartupJobTemplates []string `json:"startupJobTemplates,omitempty"`
	// Configuration for the tcp redirection side car
	TCPRedir *TCPRedirConfig `json:"tcpRedir,omitempty"`
}

// Volume represents a volume configuration for the emulator.
type Volume struct {
	// A prefix to apply to PVCs created for devices using this configuration.
	VolumePrefix string `json:"volumePrefix"`
	// Where to mount the volume in emulator pods.
	MountPoint string `json:"mountPoint"`
	// A PVC spec to use for creating the emulator volumes.
	PVCSpec corev1.PersistentVolumeClaimSpec `json:"pvcSpec"`
}

type TCPRedirConfig struct {
	// Whether to run a sidecar with emulator pods that redirects TCP traffic on the adb port
	// to the emulator adb server listening on the loopback interface. This is required
	// for the image used in this repository, but if you are using an image that
	// exposes ADB on all interfaces itself, this is not required.
	Enabled bool `json:"enabled,omitempty"`
	// Image is the repository to download the image from.
	// Defaults to quay.io/tinyzimmer/goredir whose source is in this repository.
	Image string `json:"image,omitempty"`
	// The pull policy to attach to deployments using this image.
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`
	// Any pull secrets required for downloading the image.
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`
}

// AndroidDeviceConfigStatus defines the observed state of AndroidDeviceConfig
type AndroidDeviceConfigStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AndroidDeviceConfig is the Schema for the androiddeviceconfigs API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=androiddeviceconfigs,scope=Cluster
type AndroidDeviceConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AndroidDeviceConfigSpec   `json:"spec,omitempty"`
	Status AndroidDeviceConfigStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AndroidDeviceConfigList contains a list of AndroidDeviceConfig
type AndroidDeviceConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AndroidDeviceConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AndroidDeviceConfig{}, &AndroidDeviceConfigList{})
}
