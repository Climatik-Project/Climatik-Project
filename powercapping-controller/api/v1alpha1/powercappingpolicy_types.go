/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// +groupName=climatik.io
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PowerCappingPolicySpec defines the desired state of PowerCappingPolicy
type PowerCappingPolicySpec struct {
	PowerCapLimit    int                  `json:"powerCapLimit"`
	Selector         metav1.LabelSelector `json:"selector"`
	CappingThreshold int                  `json:"cappingThreshold"`
	CustomAlgorithms []CustomAlgorithm    `json:"customAlgorithms,omitempty"`
}

// CustomAlgorithm defines a custom recommender algorithm
type CustomAlgorithm struct {
	Name       string            `json:"name"`
	Parameters map[string]string `json:"parameters,omitempty"`
}

// PowerCappingPolicyStatus defines the observed state of PowerCappingPolicy
type PowerCappingPolicyStatus struct {
	CappingActionRequired bool        `json:"cappingActionRequired"`
	LastUpdated           metav1.Time `json:"lastUpdated"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// PowerCappingPolicy is the Schema for the powercappingpolicies API
type PowerCappingPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PowerCappingPolicySpec   `json:"spec,omitempty"`
	Status PowerCappingPolicyStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PowerCappingPolicyList contains a list of PowerCappingPolicy
type PowerCappingPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PowerCappingPolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PowerCappingPolicy{}, &PowerCappingPolicyList{})
}
