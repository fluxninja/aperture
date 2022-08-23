/*
Copyright 2022.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ApertureSpec defines the desired state of Aperture Agent and Controller.
type ApertureSpec struct {
	// The registry of all the images
	//+kubebuilder:validation:Optional
	//+operator-sdk:csv:customresourcedefinitions:type=spec
	ImageRegistry string `json:"imageRegistry"`

	// The PullSecrets for all the images
	//+kubebuilder:validation:Optional
	//+operator-sdk:csv:customresourcedefinitions:type=spec
	ImagePullSecrets []string `json:"imagePullSecrets,omitempty"`

	// Labels to add to all deployed objects
	//+mapType=atomic
	//+kubebuilder:validation:Optional
	//+operator-sdk:csv:customresourcedefinitions:type=spec
	Labels map[string]string `json:"labels,omitempty"`

	// Annotations to add to all deployed objects
	//+kubebuilder:validation:Optional
	//+operator-sdk:csv:customresourcedefinitions:type=spec
	Annotations map[string]string `json:"annotations,omitempty"`

	// FluxNinjaPlugin defines the parameters for FluxNinja plugin with Aperture
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:={enabled:false,heartbeatsInterval:"30s"}
	//+operator-sdk:csv:customresourcedefinitions:type=spec
	FluxNinjaPlugin FluxNinjaPluginSpec `json:"fluxninjaPlugin"`

	// Agent defines the desired state of Agent
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:={serverPort:80}
	//+operator-sdk:csv:customresourcedefinitions:type=spec
	Agent AgentSpec `json:"agent"`

	// Controller defines the desired state of Controller
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:={serverPort:80}
	//+operator-sdk:csv:customresourcedefinitions:type=spec
	Controller ControllerSpec `json:"controller"`

	// Service defines the desired state of Services for Agent and Controller
	//+kubebuilder:validation:Optional
	//+operator-sdk:csv:customresourcedefinitions:type=spec
	Service ServiceSpec `json:"service"`

	// Sidecar defines the desired state of Sidecar setup for Agent
	//+kubebuilder:validation:Optional
	//+operator-sdk:csv:customresourcedefinitions:type=spec
	Sidecar SidecarSpec `json:"sidecar"`

	// Etcd parameters for Agent and Controller
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:={leaseTtl:"60s"}
	//+operator-sdk:csv:customresourcedefinitions:type=spec
	Etcd EtcdSpec `json:"etcd"`

	// Prometheus parameters for Agent and Controller
	//+kubebuilder:validation:Optional
	//+operator-sdk:csv:customresourcedefinitions:type=spec
	Prometheus PrometheusSpec `json:"prometheus"`
}

// ApertureStatus defines the observed state of Aperture.
type ApertureStatus struct {
	Resources string `json:"resources,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Resources",type=string,JSONPath=`.status.resources`
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// Aperture is the Schema for the Aperture API.
type Aperture struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	//+kubebuilder:default:={agent:{serverPort:80},controller:{serverPort:80}}
	Spec   ApertureSpec   `json:"spec,omitempty"`
	Status ApertureStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ApertureList contains a list of Aperture.
type ApertureList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Aperture `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Aperture{}, &ApertureList{})
}
