package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Major WIP

type Activity string
type ActionType string

const (
	CommandActivity  Activity = "Command"
	InstallActivity  Activity = "Install"
	WaitActivity     Activity = "Wait"
	InteractActivity Activity = "Interact"
)

const (
	ClickAction ActionType = "Click"
	TypeAction  ActionType = "TypeText"
)

// AndroidJobTemplateSpec defines the desired state of AndroidJobTemplate
type AndroidJobTemplateSpec struct {
	Actions []Action `json:"actions"`
}

type Action struct {
	Activity     Activity      `json:"activity"`
	Name         string        `json:"name,omitempty"`
	RunAsRoot    bool          `json:"runAsRoot,omitempty"`
	Commands     []string      `json:"commands,omitempty"`
	APKUrl       string        `json:"apkURL,omitempty"`
	Seconds      int           `json:"seconds,omitempty"`
	Interactions []Interaction `json:"interactions,omitempty"`
}

type Interaction struct {
	Type   ActionType `json:"type,omitempty"`
	Target string     `json:"target,omitempty"`
	Input  string     `json:"input,omitempty"`
}

// AndroidJobTemplateStatus defines the observed state of AndroidJobTemplate
type AndroidJobTemplateStatus struct{}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AndroidJobTemplate is the Schema for the androidjobtemplates API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=androidjobtemplates,scope=Cluster
type AndroidJobTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AndroidJobTemplateSpec   `json:"spec,omitempty"`
	Status AndroidJobTemplateStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AndroidJobTemplateList contains a list of AndroidJobTemplate
type AndroidJobTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AndroidJobTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AndroidJobTemplate{}, &AndroidJobTemplateList{})
}
