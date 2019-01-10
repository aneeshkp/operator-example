package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AputturkindSpec defines the desired state of Aputturkind
type AputturkindSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	Count int32  `json:"count"`
	Group string `json:"group"`
	Image string `json:"image"`
	Port  int32  `json:"port"`
}

// AputturkindStatus defines the observed state of Aputturkind
type AputturkindStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	PodNames []string `json:"podnames"`
	AppGroup string   `json:"appgroup"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Aputturkind is the Schema for the aputturkinds API
// +k8s:openapi-gen=true
type Aputturkind struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AputturkindSpec   `json:"spec,omitempty"`
	Status AputturkindStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AputturkindList contains a list of Aputturkind
type AputturkindList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Aputturkind `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Aputturkind{}, &AputturkindList{})
}
