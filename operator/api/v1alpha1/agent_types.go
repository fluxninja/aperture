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

import "time"

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

// Batch defines configuration for OTEL batch processor.
type Batch struct {
	// Timeout sets the time after which a batch will be sent regardless of size.
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:=1000000000
	Timeout time.Duration `json:"timeout"`
	// SendBatchSize is the size of a batch which after hit, will trigger it to be sent.
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:=10000
	SendBatchSize uint32 `json:"sendBatchSize"`
}
