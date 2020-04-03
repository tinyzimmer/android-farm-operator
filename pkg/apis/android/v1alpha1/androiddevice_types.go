package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AndroidDeviceSpec defines the desired state of AndroidDevice
type AndroidDeviceSpec struct {
	// A reference to an AndroidDeviceConfig to use for the emulators in this group.
	ConfigRef *corev1.LocalObjectReference `json:"configRef,omitempty"`
	// Any overrides to the config represented by the ConfigRef. Any values supplied here
	// will be merged into the found AndroidDeviceConfig, with fields in this object
	// taking precedence over existing ones in the AndroidDeviceConfig.
	DeviceConfig *AndroidDeviceConfigSpec `json:"deviceConfig,omitempty"`
	// A hostname to apply to the device (used by AndroidFarm controller)
	Hostname string `json:"hostname,omitempty"`
	// A subdomain to apply to the device (used by AndroidFarm controller)
	Subdomain string `json:"subdomain,omitempty"`
}

// AndroidDeviceStatus defines the observed state of AndroidDevice
type AndroidDeviceStatus struct {
	State string `json:"state"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AndroidDevice is the Schema for the androiddevices API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=androiddevices,scope=Namespaced
type AndroidDevice struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AndroidDeviceSpec   `json:"spec,omitempty"`
	Status AndroidDeviceStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AndroidDeviceList contains a list of AndroidDevice
type AndroidDeviceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AndroidDevice `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AndroidDevice{}, &AndroidDeviceList{})
}
