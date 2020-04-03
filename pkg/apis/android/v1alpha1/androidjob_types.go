package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// MAJOR WIP

type JobStatus string

const (
	StatusPending  JobStatus = "Pending"
	StatusComplete JobStatus = "Complete"
	StatusFailed   JobStatus = "Failed"
)

// AndroidJobSpec defines the desired state of AndroidJob
type AndroidJobSpec struct {
	DeviceName              string            `json:"deviceName,omitempty"`
	DeviceSelector          map[string]string `json:"deviceSelector,omitempty"`
	JobTemplate             string            `json:"jobTemplate"`
	TTLSecondsAfterCreation *int              `json:"ttlSecondsAfterCreation,omitempty"`
}

// AndroidJobStatus defines the observed state of AndroidJob
type AndroidJobStatus struct {
	// JobStatus is a map of device name to device status
	JobStatus map[string]DeviceJobStatus `json:"jobStatus,omitempty"`
}

// DeviceJobStatus defines the state of the job for a single device
type DeviceJobStatus struct {
	// Status is the current status of the job
	Status JobStatus `json:"jobStatus,omitempty"`
	// Message may contain extra information about the status of the job
	Message string `json:"message,omitempty"`
}

func (a *AndroidJob) DeviceNamespacedName() types.NamespacedName {
	return types.NamespacedName{Name: a.Spec.DeviceName, Namespace: a.Namespace}
}

func (a *AndroidJob) TemplateNamespacedName() types.NamespacedName {
	return types.NamespacedName{Name: a.Spec.JobTemplate, Namespace: metav1.NamespaceAll}
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AndroidJob is the Schema for the androidjobs API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=androidjobs,scope=Namespaced
type AndroidJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AndroidJobSpec   `json:"spec,omitempty"`
	Status AndroidJobStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AndroidJobList contains a list of AndroidJob
type AndroidJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AndroidJob `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AndroidJob{}, &AndroidJobList{})
}
