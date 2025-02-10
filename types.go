package main

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Define the CRD Spec
type WebResourceSpec struct {
	Replicas int    `json:"replicas"`
	Image    string `json:"image"`
}

// Define the CRD Structure
type WebResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              WebResourceSpec `json:"spec"`
}

// Define a List Type for Multiple CRs
type WebResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []WebResource `json:"items"`
}
