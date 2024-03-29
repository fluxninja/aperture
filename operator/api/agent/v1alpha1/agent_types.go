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

// +kubebuilder:object:generate=true
// +groupName=fluxninja.com
package v1alpha1

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	agent "github.com/fluxninja/aperture/v2/cmd/aperture-agent/config"
	"github.com/fluxninja/aperture/v2/operator/api"
	"github.com/fluxninja/aperture/v2/operator/api/common"
	afconfig "github.com/fluxninja/aperture/v2/pkg/agent-functions/config"
	agentinfo "github.com/fluxninja/aperture/v2/pkg/agent-info"
	distcache "github.com/fluxninja/aperture/v2/pkg/dist-cache/config"
	"github.com/fluxninja/aperture/v2/pkg/etcd"
	"github.com/fluxninja/aperture/v2/pkg/net/http"
	objectstorageconfig "github.com/fluxninja/aperture/v2/pkg/objectstorage/config"
	peers "github.com/fluxninja/aperture/v2/pkg/peers/config"
	autoscalek8sconfig "github.com/fluxninja/aperture/v2/pkg/policies/autoscale/kubernetes/config"
	preview "github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/service/preview/config"
	prometheus "github.com/fluxninja/aperture/v2/pkg/prometheus/config"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AgentSpec defines the desired state for the Agent.
type AgentSpec struct {
	// CommonSpec defines the common state between Agent and Controller
	common.CommonSpec `json:",inline"`

	// Image configuration
	//+kubebuilder:validation:Optional
	Image common.AgentImage `json:"image"`

	// Sidecar defines the desired state of Sidecar setup for Agent
	//+kubebuilder:validation:Optional
	Sidecar SidecarSpec `json:"sidecar"`

	// Agent Configuration
	//+kubebuilder:validation:Optional
	ConfigSpec AgentConfigSpec `json:"config"`

	// ControllerClientCertConfig configuration.
	//+kubebuilder:validation:Optional
	ControllerClientCertConfig common.ControllerClientCertConfig `json:"controller_client_cert"`

	// DeploymentConfigSpec defines the deployment configuration of the agent.
	// This is an experimental feature. Only DaemonSet is supported.
	//+kubebuilder:validation:Optional
	DeploymentConfigSpec DeploymentConfigSpec `json:"deployment_config"`

	// NameOverride overrides the name of the resources created for agent.
	// This is an experimental feature.
	//+kubebuilder:validation:Optional
	NameOverride string `json:"name_override"`
}

// DeploymentConfigSpec defines the deployment configuration of the agent.
type DeploymentConfigSpec struct {
	// Type of the deployment.
	//+kubebuilder:validation:Optional
	Type string `json:"type" default:"DaemonSet" validate:"oneof=deployment Deployment DaemonSet daemonset"`

	// Number of replicas when type is set to Deployment.
	//+kubebuilder:validation:Optional
	Replicas int32 `json:"replicas,omitempty" default:"1" validate:"gt=0"`

	// TopologySpreadConstraints to be applied to the deployment.
	//+kubebuilder:validation:Optional
	TopologySpreadConstraints []corev1.TopologySpreadConstraint `json:"topology_spread_constraints,omitempty"`

	// Strategy to be applied to the deployment upgrades.
	//+kubebuilder:validation:Optional
	Strategy appsv1.DeploymentStrategy `json:"strategy,omitempty"`
}

// AgentConfigSpec holds agent configuration.
type AgentConfigSpec struct {
	// CommonConfigSpec
	//+kubebuilder:validation:Optional
	common.CommonConfigSpec `json:",inline"`

	// Etcd configuration.
	//+kubebuilder:validation:Optional
	Etcd etcd.EtcdConfig `json:"etcd"`

	// Prometheus configuration.
	//+kubebuilder:validation:Optional
	Prometheus prometheus.PrometheusConfig `json:"prometheus"`

	// AgentInfo configuration.
	//+kubebuilder:validation:Optional
	AgentInfo agentinfo.AgentInfoConfig `json:"agent_info"`

	// DistCache configuration.
	//+kubebuilder:validation:Optional
	DistCache distcache.DistCacheConfig `json:"dist_cache"`

	// Kubernetes client configuration.
	//+kubebuilder:validation:Optional
	KubernetesClient http.HTTPClientConfig `json:"kubernetes_client"`

	// Peer discovery configuration.
	//+kubebuilder:validation:Optional
	PeerDiscovery peers.PeerDiscoveryConfig `json:"peer_discovery"`

	// FlowControl configuration.
	//+kubebuilder:validation:Optional
	FlowControl FlowControlConfigSpec `json:"flow_control"`

	// AutoScale configuration.
	//+kubebuilder:validation:Optional
	AutoScale AutoScaleConfigSpec `json:"auto_scale"`

	// Service Discovery configuration.
	//+kubebuilder:validation:Optional
	ServiceDiscoverySpec common.ServiceDiscoverySpec `json:"service_discovery"`

	// OTel configuration.
	//+kubebuilder:validation:Optional
	OTel agent.AgentOTelConfig `json:"otel"`

	// Agent functions configuration.
	//+kubebuilder:validation:Optional
	AgentFunctions afconfig.AgentFunctionsConfig `json:"agent_functions"`

	// ObjectStorage configuration.
	//+kubebuilder:validation:Optional
	ObjectStorage objectstorageconfig.ObjectStorageConfig `json:"object_storage"`
}

// FlowControlConfigSpec holds flow control configuration.
type FlowControlConfigSpec struct {
	// FlowPreviewConfig holds flow preview configuration.
	//+kubebuilder:validation:Optional
	FlowPreviewConfig preview.FlowPreviewConfig `json:"preview_service"`
}

// AutoScaleConfigSpec holds auto-scale configuration.
type AutoScaleConfigSpec struct {
	// AutoScaleKubernetesConfig holds auto-scale kubernetes configuration.
	//+kubebuilder:validation:Optional
	AutoScaleKubernetesConfig autoscalek8sconfig.AutoScaleKubernetesConfig `json:"kubernetes"`
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
	api.SchemeBuilder.Register(&Agent{}, &AgentList{})
}
