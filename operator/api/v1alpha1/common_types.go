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
	corev1 "k8s.io/api/core/v1"
)

// Image defines the Registry, Repository, Tag, PullPolicy, PullSecrets and Debug.
type Image struct {
	// The registry of the image
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:="docker.io/fluxninja"
	Registry string `json:"registry"`

	// The repository of the image
	Repository string `json:"repository"`

	// The tag (version) of the image
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:=latest
	Tag string `json:"tag"`

	// The ImagePullPolicy of the image
	//+kubebuilder:validation:Optional
	//+kubebuilder:validation:Enum=Never;Always;IfNotPresent
	//+kubebuilder:default:="IfNotPresent"
	PullPolicy string `json:"pullPolicy"`

	// The PullSecrets for the image
	//+kubebuilder:validation:Optional
	PullSecrets []string `json:"pullSecrets,omitempty"`
}

// Log defines logger configuration.
type Log struct {
	// Log level
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:="info"
	Level string `json:"level"`

	// Log filename
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:="stderr"
	File string `json:"file"`

	// Flag for non-blocking
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:=true
	NonBlocking bool `json:"nonBlocking"`

	// Flag for pretty console
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:=false
	PrettyConsole bool `json:"prettyConsole"`
}

// Probe defines Enabled, InitialDelaySeconds, PeriodSeconds, TimeoutSeconds, FailureThreshold and SuccessThreshold for probes like livenessProbe.
type Probe struct {
	// Enable probe on agent containers
	Enabled bool `json:"enabled"`

	// Initial delay seconds for probe
	//+kubebuilder:validation:Optional
	//+kubebuilder:validation:Minimum:=0
	//+kubebuilder:default:=15
	InitialDelaySeconds int32 `json:"initialDelaySeconds"`

	// Period delay seconds for probe
	//+kubebuilder:validation:Optional
	//+kubebuilder:validation:Minimum:=1
	//+kubebuilder:default:=15
	PeriodSeconds int32 `json:"periodSeconds"`

	// Timeout delay seconds for probe
	//+kubebuilder:validation:Optional
	//+kubebuilder:validation:Minimum:=1
	//+kubebuilder:default:=5
	TimeoutSeconds int32 `json:"timeoutSeconds"`

	// Failure threshold for probe
	//+kubebuilder:validation:Optional
	//+kubebuilder:validation:Minimum:=1
	//+kubebuilder:default:=6
	FailureThreshold int32 `json:"failureThreshold"`

	// Success threshold for probe
	//+kubebuilder:validation:Optional
	//+kubebuilder:validation:Minimum:=1
	//+kubebuilder:default:=1
	SuccessThreshold int32 `json:"successThreshold"`
}

// AgentEtcdSpec defines Endpoints and LeaseTtl of etcd used by Aperture Agent.
type AgentEtcdSpec struct {
	// Etcd endpoints
	Endpoints []string `json:"endpoints,omitempty"`

	// Etcd leaseTtl
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:="60s"
	LeaseTTL string `json:"leaseTtl"`
}

// ControllerEtcdSpec defines Endpoints and LeaseTtl of etcd used by Aperture Controller.
type ControllerEtcdSpec struct {
	// Etcd endpoints
	//+kubebuilder:validation:Optional
	Endpoints []string `json:"endpoints,omitempty"`

	// Etcd leaseTtl
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:="60s"
	LeaseTTL string `json:"leaseTtl"`
}

// PrometheusSpec defines parameters required for Prometheus connection.
type PrometheusSpec struct {
	// Address for Prometheus
	Address string `json:"address"`
}

// PodSecurityContext defines Enabled and FsGroup for the Pods' security context.
type PodSecurityContext struct {
	// Enable PodSecurityContext on Pod
	Enabled bool `json:"enabled"`

	// fsGroup to define the Group ID for the Pod
	//+kubebuilder:validation:minimum=0
	//+kubebuilder:validation:Optional
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:=1001
	FsGroup *int64 `json:"fsGroup"`
}

// ContainerSecurityContext defines Enabled, RunAsUser, RunAsNonRootUser and ReadOnlyRootFilesystem for the containers' security context.
type ContainerSecurityContext struct {
	// Enable ContainerSecurityContext on containers
	//+kubebuilder:validation:Optional
	Enabled bool `json:"enabled"`

	// Set containers' Security Context runAsUser
	//+kubebuilder:validation:Optional
	//+kubebuilder:validation:minimum:=0
	//+kubebuilder:default:=1001
	RunAsUser *int64 `json:"runAsUser"`

	// Set containers' Security Context runAsNonRoot
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:=false
	RunAsNonRootUser *bool `json:"runAsNonRoot"`

	// Set agent containers' Security Context runAsNonRoot
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:=false
	ReadOnlyRootFilesystem *bool `json:"readOnlyRootFilesystem"`
}

// FluxNinjaPluginSpec defines the parameters for FluxNinja Plugin.
type FluxNinjaPluginSpec struct {
	// Enabled the FluxNinja plugin with Aperture
	//+kubebuilder:validation:Optional
	Enabled bool `json:"enabled"`

	// FluxNinja cloud instance address
	//+kubebuilder:validation:Optional
	Endpoint string `json:"endpoint"`

	// Specifies how often to send heartbeats to the FluxNinja cloud
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:="30s"
	HeartbeatsInterval string `json:"heartbeatsInterval"`

	// tls configuration to communicate with FluxNinja cloud
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:={insecure:false,insecureSkipVerify:false}
	TLS TLSSpec `json:"tls"`

	// API Key secret references for Agent or Controller
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:={create:true}
	APIKeySecret APIKeySecret `json:"apiKeySecret"`
}

// CommonSpec defines the desired the common state of Agent and Controller.
type CommonSpec struct {
	// Labels to add to all deployed objects
	//+mapType=atomic
	//+kubebuilder:validation:Optional
	//+operator-sdk:csv:customresourcedefinitions:type=spec
	Labels map[string]string `json:"labels,omitempty"`

	// Annotations to add to all deployed objects
	//+kubebuilder:validation:Optional
	//+operator-sdk:csv:customresourcedefinitions:type=spec
	Annotations map[string]string `json:"annotations,omitempty"`

	// FluxNinjaPlugin defines the parameters for FluxNinja plugin with Agent or Controller
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:={enabled:false,heartbeatsInterval:"30s"}
	//+operator-sdk:csv:customresourcedefinitions:type=spec
	FluxNinjaPlugin FluxNinjaPluginSpec `json:"fluxninjaPlugin"`

	// OtelConfig is the configuration for the OTEL collector
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:={grpcAddr:":4317"}
	//+operator-sdk:csv:customresourcedefinitions:type=spec
	OtelConfig OtelConfig `json:"otelConfig"`

	// Configuration for Agent or Controller service
	//+kubebuilder:validation:Optional
	Service Service `json:"service"`

	// Server port for the Agent
	//+kubebuilder:default:=80
	//+kubebuilder:validation:Optional
	//+kubebuilder:validation:Maximum:=65535
	//+kubebuilder:validation:Minimum:=1
	ServerPort int32 `json:"serverPort"`

	// Log related configurations
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:={prettyConsole:false,nonBlocking:true,level:"info",file:"stderr"}
	Log Log `json:"log"`

	// ServiceAccountSpec defines the the configuration pf Service account for Agent or Controller
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:={create:true,automountServiceAccountToken:true}
	ServiceAccountSpec ServiceAccountSpec `json:"serviceAccount"`

	// livenessProbe related configuration
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:={enabled:true,initialDelaySeconds:15,periodSeconds:15,timeoutSeconds:5,failureThreshold:6,successThreshold:1}
	LivenessProbe Probe `json:"livenessProbe"`

	// readinessProbe related configuration
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:={enabled:true,initialDelaySeconds:15,periodSeconds:15,timeoutSeconds:5,failureThreshold:6,successThreshold:1}
	ReadinessProbe Probe `json:"readinessProbe"`

	// Custom livenessProbe that overrides the default one
	//+kubebuilder:validation:Optional
	CustomLivenessProbe *corev1.Probe `json:"customLivenessProbe,omitempty"`

	// Custom readinessProbe that overrides the default one
	//+kubebuilder:validation:Optional
	CustomReadinessProbe *corev1.Probe `json:"customReadinessProbe,omitempty"`

	// Resource requests and limits
	//+kubebuilder:validation:Optional
	Resources corev1.ResourceRequirements `json:"resources"`

	// Configure Pods' Security Context
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:={enabled:false,fsGroup:1001}
	PodSecurityContext PodSecurityContext `json:"podSecurityContext"`

	// Configure Containers' Security Context
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:={enabled:false,runAsUser:1001,runAsNonRoot:false,readOnlyRootFilesystem:false}
	ContainerSecurityContext ContainerSecurityContext `json:"containerSecurityContext"`

	// Override default container command
	//+kubebuilder:validation:Optional
	Command []string `json:"command,omitempty"`

	// Override default container args
	//+kubebuilder:validation:Optional
	Args []string `json:"args,omitempty"`

	// Extra labels for pods
	//+mapType=atomic
	//+kubebuilder:validation:Optional
	PodLabels map[string]string `json:"podLabels,omitempty"`

	// Extra Annotations for pods
	//+kubebuilder:validation:Optional
	PodAnnotations map[string]string `json:"podAnnotations,omitempty"`

	// Affinity for pods assignment.
	//+kubebuilder:validation:Optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// Node labels for pods assignment
	//+mapType=atomic
	//+kubebuilder:validation:Optional
	//+mapType=atomic
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// Tolerations for pods assignment
	//+kubebuilder:validation:Optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`

	// Seconds Redmine pod needs to terminate gracefully
	//+kubebuilder:validation:Minimum:=0
	//+kubebuilder:validation:Optional
	TerminationGracePeriodSeconds *int64 `json:"terminationGracePeriodSeconds"`

	// For the container(s) to automate configuration before or after startup
	//+kubebuilder:validation:Optional
	LifecycleHooks *corev1.Lifecycle `json:"lifecycleHooks"`

	// Array with extra environment variables to add
	//+kubebuilder:validation:Optional
	ExtraEnvVars []corev1.EnvVar `json:"extraEnvVars,omitempty"`

	// Name of existing ConfigMap containing extra env vars
	//+kubebuilder:validation:Optional
	ExtraEnvVarsCM string `json:"extraEnvVarsCM"`

	// Name of existing Secret containing extra env vars
	//+kubebuilder:validation:Optional
	ExtraEnvVarsSecret string `json:"extraEnvVarsSecret"`

	// Optionally specify extra list of additional volumes for the pod(s)
	//+kubebuilder:validation:Optional
	ExtraVolumes []corev1.Volume `json:"extraVolumes,omitempty"`

	// Optionally specify extra list of additional volumeMounts
	//+kubebuilder:validation:Optional
	ExtraVolumeMounts []corev1.VolumeMount `json:"extraVolumeMounts,omitempty"`

	// Add additional sidecar containers
	//+kubebuilder:validation:Optional
	Sidecars []corev1.Container `json:"sidecars,omitempty"`

	// Add additional init containers
	//+kubebuilder:validation:Optional
	InitContainers []corev1.Container `json:"initContainers,omitempty"`
}

// Ingestion defines the fields required for configuring the FluxNinja cloud connection details.
type Ingestion struct {
	// Specifies address of FluxNinja cloud instance to connect to
	//+kubebuilder:validation:Optional
	Address string `json:"address"`

	// Specifies port of FluxNinja cloud instance to connect to
	//+kubebuilder:validation:Optional
	//+kubebuilder:validation:Maximum:=65535
	//+kubebuilder:validation:Minimum:=1
	//+kubebuilder:default:=443
	Port int `json:"port"`

	// Specifies whether to connect to the FluxNinja cloud over TLS or in plain text
	//+kubebuilder:validation:Optional
	Insecure bool `json:"insecure"`

	// Specifies whether to verify FluxNinja cloud server certificate
	//+kubebuilder:validation:Optional
	InsecureSkipVerify bool `json:"insecureSkipVerify"`
}

// APIKeySecret defines fields required for creation/usage of secret for the ApiKey of Agent and Controller.
type APIKeySecret struct {
	// Create new secret or not
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:=true
	Create bool `json:"create"`

	// Secret details
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:={key:"apiKey"}
	SecretKeyRef SecretKeyRef `json:"secretKeyRef"`

	// Value for the ApiKey
	//+kubebuilder:validation:Optional
	Value string `json:"value"`
}

// SecretKeyRef defines fields for details of the ApiKey secret.
type SecretKeyRef struct {
	// Name of the secret
	//+kubebuilder:validation:Optional
	Name string `json:"name"`

	// Key of the secret in Data
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:=apiKey
	Key string `json:"key"`
}

// TLSSpec defines fields for TLS configuration for FluxNinja Plugin.
type TLSSpec struct {
	// Specifies whether to communicate with FluxNinja cloud over TLS or in plain text
	//+kubebuilder:validation:Optional
	Insecure bool `json:"insecure"`

	// Specifies whether to verify FluxNinja cloud certificate
	//+kubebuilder:validation:Optional
	InsecureSkipVerify bool `json:"insecureSkipVerify"`

	// Alternative CA certificates bundle to use to validate FluxNinja cloud certificate
	//+kubebuilder:validation:Optional
	CAFile string `json:"caFile"`
}

// APIKeySecretSpec defines API Key secret details for Agent and Controller.
type APIKeySecretSpec struct {
	// API Key secret reference for Agent
	//+kubebuilder:validation:Optional
	Agent APIKeySecret `json:"agent"`

	// API Key secret reference for Controller
	//+kubebuilder:validation:Optional
	Controller APIKeySecret `json:"controller"`
}

// Batch defines configuration for OTEL batch processor.
type Batch struct {
	// Timeout sets the time after which a batch will be sent regardless of size.
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:="1s"
	Timeout string `json:"timeout"`

	// SendBatchSize is the size of a batch which after hit, will trigger it to be sent.
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:=10000
	SendBatchSize uint32 `json:"sendBatchSize"`
}

// OtelConfig defines the configuration for the OTEL collector.
type OtelConfig struct {
	// GRPC listener addr for OTEL Collector
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:=":4317"
	GRPCAddr string `json:"grpcAddr"`

	// HTTP listener addr for OTEL Collector
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:=":4318"
	HTTPAddr string `json:"httpAddr"`

	// Batch prerollup processor configuration.
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:={timeout:"1s",sendBatchSize:10000}
	BatchPrerollup Batch `json:"batchPrerollup"`

	// Batch postrollup processor configuration.
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:={timeout:"1s",sendBatchSize:10000}
	BatchPostrollup Batch `json:"batchPostrollup"`
}
