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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AgentSpec defines the desired state for the Agent.
type AgentSpec struct {
	// CommonSpec defines the common state between Agent and Controller
	CommonSpec `json:",inline"`

	// Port for the Agent's distributed cache service
	//+kubebuilder:default:=3320
	//+kubebuilder:validation:Optional
	//+kubebuilder:validation:Maximum:=65535
	//+kubebuilder:validation:Minimum:=1
	DistributedCachePort int32 `json:"distributedCachePort"`

	// Port for the Agent's member list service
	//+kubebuilder:default:=3322
	//+kubebuilder:validation:Optional
	//+kubebuilder:validation:Maximum:=65535
	//+kubebuilder:validation:Minimum:=1
	MemberListPort int32 `json:"memberListPort"`

	// Image configuration
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:={tag:"latest",pullPolicy:"IfNotPresent",registry:"docker.io/fluxninja",repository:"aperture-agent"}
	Image Image `json:"image"`

	// AgentGroup name for the Agent
	//+kubebuilder:validation:Optional
	AgentGroup string `json:"agentGroup"`

	// Sidecar defines the desired state of Sidecar setup for Agent
	//+kubebuilder:validation:Optional
	//+operator-sdk:csv:customresourcedefinitions:type=spec
	Sidecar SidecarSpec `json:"sidecar"`

	// Batch prerollup processor configuration.
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:={timeout:1000000000,sendBatchSize:10000}
	BatchPrerollup Batch `json:"batchPrerollup"`

	// Batch postrollup processor configuration.
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:={timeout:1000000000,sendBatchSize:10000}
	BatchPostrollup Batch `json:"batchPostrollup"`

	// Batch metrics/fast processor configuration.
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:={timeout:1000000000,sendBatchSize:1000}
	BatchMetricsFast Batch `json:"batchMetricsFast"`
}

// AgentStatus defines the observed state of Agent.
type AgentStatus struct {
	Resources string `json:"resources,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Resources",type=string,JSONPath=`.status.resources`
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// Agent is the Schema for the agents API.
type Agent struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	//+kubebuilder:default:={serverPort:80}
	Spec   AgentSpec   `json:"spec,omitempty"`
	Status AgentStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AgentList contains a list of Agent.
type AgentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Agent `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Agent{}, &AgentList{})
}
