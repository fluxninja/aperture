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
}
