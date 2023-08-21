package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type RepositoryAuth struct {
	Type      string `json:"type,omitempty"`
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
	PublicKey string `json:"publicKey,omitempty"`
}

type RepositorySync struct {
	Auto bool `json:"auto,omitempty"`
}

type RepositorySpec struct {
	Url     string         `json:"url,omitempty"`
	Path    string         `json:"path,omitempty"`
	Sync    RepositorySync `json:"sync,omitempty"`
	AuthRef Ref            `json:"ref,omitempty"`
}

type RepositoryStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Repository is the Schema for the circles API
type Repository struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RepositorySpec   `json:"spec,omitempty"`
	Status RepositoryStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RepositoryList contains a list of Repository
type RepositoryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Repository `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Repository{}, &RepositoryList{})
}