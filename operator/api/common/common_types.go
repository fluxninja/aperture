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
package common

import (
	"github.com/fluxninja/aperture/v2/extensions/fluxninja/extconfig"
	sentry "github.com/fluxninja/aperture/v2/extensions/sentry/config"
	alertmgrconfig "github.com/fluxninja/aperture/v2/pkg/alert-manager/config"
	"github.com/fluxninja/aperture/v2/pkg/config"
	kubernetes "github.com/fluxninja/aperture/v2/pkg/discovery/kubernetes/config"
	googletoken "github.com/fluxninja/aperture/v2/pkg/google/config"
	jobs "github.com/fluxninja/aperture/v2/pkg/jobs/config"
	"github.com/fluxninja/aperture/v2/pkg/metrics"
	"github.com/fluxninja/aperture/v2/pkg/net/grpc"
	"github.com/fluxninja/aperture/v2/pkg/net/grpcgateway"
	"github.com/fluxninja/aperture/v2/pkg/net/http"
	"github.com/fluxninja/aperture/v2/pkg/net/listener"
	"github.com/fluxninja/aperture/v2/pkg/net/tlsconfig"
	"github.com/fluxninja/aperture/v2/pkg/profilers"
	watchdogconfig "github.com/fluxninja/aperture/v2/pkg/watchdog/config"

	corev1 "k8s.io/api/core/v1"
)

// Image defines the Registry, Repository, Tag, PullPolicy, PullSecrets and Debug.
type Image struct {
	// The registry of the image
	//+kubebuilder:validation:Optional
	Registry string `json:"registry" default:"docker.io/fluxninja"`

	// The tag (version) of the image
	//+kubebuilder:validation:Optional
	Tag string `json:"tag"`

	// The digest (sha) of the image
	//+kubebuilder:validation:Optional
	Digest string `json:"digest"`

	// The ImagePullPolicy of the image
	//+kubebuilder:validation:Optional
	PullPolicy string `json:"pullPolicy" default:"IfNotPresent" validate:"oneof=Never Always IfNotPresent"`

	// The PullSecrets for the image
	//+kubebuilder:validation:Optional
	PullSecrets []string `json:"pullSecrets,omitempty"`
}

// AgentImage defines Image spec for Aperture Agent.
type AgentImage struct {
	// Image specs for Agent
	Image `json:",inline"`

	// The repository of the image
	//+kubebuilder:validation:Optional
	Repository string `json:"repository" default:"aperture-agent"`
}

// ControllerImage defines Image spec for Aperture Controller.
type ControllerImage struct {
	// Image specs for Controller
	Image `json:",inline"`

	// The repository of the image
	//+kubebuilder:validation:Optional
	Repository string `json:"repository" default:"aperture-controller"`
}

// Probe defines Enabled, InitialDelaySeconds, PeriodSeconds, TimeoutSeconds, FailureThreshold and SuccessThreshold for probes like livenessProbe.
type Probe struct {
	// Enable probe on agent containers
	Enabled bool `json:"enabled" default:"true"`

	// Initial delay seconds for probe
	//+kubebuilder:validation:Optional
	InitialDelaySeconds int32 `json:"initialDelaySeconds" default:"15" validate:"gte=0"`

	// Period delay seconds for probe
	//+kubebuilder:validation:Optional
	PeriodSeconds int32 `json:"periodSeconds" default:"15" validate:"gte=1"`

	// Timeout delay seconds for probe
	//+kubebuilder:validation:Optional
	TimeoutSeconds int32 `json:"timeoutSeconds" default:"5" validate:"gte=1"`

	// Failure threshold for probe
	//+kubebuilder:validation:Optional
	FailureThreshold int32 `json:"failureThreshold" default:"6" validate:"gte=1"`

	// Success threshold for probe
	//+kubebuilder:validation:Optional
	SuccessThreshold int32 `json:"successThreshold" default:"1" validate:"gte=1"`
}

// PodSecurityContext defines Enabled and FsGroup for the Pods' security context.
type PodSecurityContext struct {
	// Enable PodSecurityContext on Pod
	Enabled bool `json:"enabled" default:"false"`

	// fsGroup to define the Group ID for the Pod
	//+kubebuilder:validation:Optional
	FsGroup int64 `json:"fsGroup" default:"1000" validate:"gte=0"`
}

// ContainerSecurityContext defines Enabled, RunAsUser, RunAsNonRootUser and ReadOnlyRootFilesystem for the containers' security context.
type ContainerSecurityContext struct {
	// Enable ContainerSecurityContext on containers
	Enabled bool `json:"enabled" default:"false"`

	// Set containers' Security Context runAsUser
	//+kubebuilder:validation:Optional
	RunAsUser int64 `json:"runAsUser" default:"1000" validate:"gte=0"`

	// Set containers' Security Context runAsGroup
	//+kubebuilder:validation:Optional
	RunAsGroup int64 `json:"runAsGroup" default:"1000" validate:"gte=0"`

	// Set containers' Security Context runAsNonRoot
	//+kubebuilder:validation:Optional
	RunAsNonRootUser bool `json:"runAsNonRoot" default:"true"`

	// Set agent containers' Security Context runAsNonRoot
	//+kubebuilder:validation:Optional
	ReadOnlyRootFilesystem bool `json:"readOnlyRootFilesystem" default:"false"`
}

// CommonSpec defines the desired the common state of Agent and Controller.
type CommonSpec struct {
	// Labels to add to all deployed objects
	//+mapType=atomic
	//+kubebuilder:validation:Optional
	Labels map[string]string `json:"labels,omitempty"`

	// Annotations to add to all deployed objects
	//+kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// Configuration for Agent or Controller service
	//+kubebuilder:validation:Optional
	Service Service `json:"service"`

	// ServiceAccountSpec defines the the configuration pf Service account for Agent or Controller
	//+kubebuilder:validation:Optional
	ServiceAccountSpec ServiceAccountSpec `json:"serviceAccount"`

	// livenessProbe related configuration
	//+kubebuilder:validation:Optional
	LivenessProbe Probe `json:"livenessProbe"`

	// readinessProbe related configuration
	//+kubebuilder:validation:Optional
	ReadinessProbe Probe `json:"readinessProbe"`

	// Custom livenessProbe that overrides the default one
	//+kubebuilder:validation:Optional
	CustomLivenessProbe *corev1.Probe `json:"customLivenessProbe,omitempty"`

	// Custom readinessProbe that overrides the default one
	//+kubebuilder:validation:Optional
	CustomReadinessProbe *corev1.Probe `json:"customReadinessProbe,omitempty"`

	// MinReadySeconds to be applied to the deployment and daemonset.
	//+kubebuilder:validation:Optional
	MinReadySeconds int32 `json:"minReadySeconds,omitempty" default:"30" validate:"gte=0"`

	// Resource requests and limits
	//+kubebuilder:validation:Optional
	Resources corev1.ResourceRequirements `json:"resources"`

	// Configure Pods' Security Context
	//+kubebuilder:validation:Optional
	PodSecurityContext PodSecurityContext `json:"podSecurityContext"`

	// Configure Containers' Security Context
	//+kubebuilder:validation:Optional
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
	//+kubebuilder:validation:Optional
	//+mapType=atomic
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// Tolerations for pods assignment
	//+kubebuilder:validation:Optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`

	// Seconds pod needs to terminate gracefully
	//+kubebuilder:validation:Optional
	TerminationGracePeriodSeconds int64 `json:"terminationGracePeriodSeconds" validate:"gte=0"`

	// For the container(s) to automate configuration before or after startup
	//+kubebuilder:validation:Optional
	LifecycleHooks *corev1.Lifecycle `json:"lifecycleHooks,omitempty" validate:"omitempty"`

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

	// Secrets
	//+kubebuilder:validation:Optional
	Secrets Secrets `json:"secrets"`
}

// Secrets for Agent or Controller.
type Secrets struct {
	// FluxNinja extension.
	//+kubebuilder:validation:Optional
	FluxNinjaExtension APIKeySecret `json:"fluxNinjaExtension"`
}

// APIKeySecret defines fields required for creation/usage of secret for the ApiKey of Agent and Controller.
type APIKeySecret struct {
	// Create new secret or not
	Create bool `json:"create" default:"false"`

	// Secret details
	//+kubebuilder:validation:Optional
	SecretKeyRef SecretKeyRef `json:"secretKeyRef"`

	// Value for the ApiKey
	Value string `json:"value"`
}

// SecretKeyRef defines fields for details of the ApiKey secret.
type SecretKeyRef struct {
	// Name of the secret
	//+kubebuilder:validation:Optional
	Name string `json:"name"`

	// Key of the secret in Data
	//+kubebuilder:validation:Optional
	Key string `json:"key" default:"apiKey"`
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

// CommonConfigSpec defines common configuration for agent and controller.
type CommonConfigSpec struct {
	// Client configuration such as proxy settings.
	//+kubebuilder:validation:Optional
	Client ClientConfigSpec `json:"client"`

	// Liveness probe configuration.
	//+kubebuilder:validation:Optional
	Liveness ProbeConfigSpec `json:"liveness"`

	// Readiness probe configuration.
	//+kubebuilder:validation:Optional
	Readiness ProbeConfigSpec `json:"readiness"`

	// Log configuration.
	//+kubebuilder:validation:Optional
	Log config.LogConfig `json:"log"`

	// Metrics configuration.
	//+kubebuilder:validation:Optional
	Metrics metrics.MetricsConfig `json:"metrics"`

	// Profilers configuration.
	//+kubebuilder:validation:Optional
	Profilers profilers.ProfilersConfig `json:"profilers"`

	// Google Token Source configuration
	//+kubebuilder:validation:Optional
	TokenSource googletoken.Config `json:"token_source"`

	// Server configuration.
	//+kubebuilder:validation:Optional
	Server ServerConfigSpec `json:"server"`

	// Watchdog configuration.
	//+kubebuilder:validation:Optional
	Watchdog watchdogconfig.WatchdogConfig `json:"watchdog"`

	// Alert Managers configuration.
	//+kubebuilder:validation:Optional
	Alertmanagers alertmgrconfig.AlertManagerConfig `json:"alertmanagers,omitempty"`

	// BundledExtensionsSpec defines configuration for bundled extensions.
	//+kubebuilder:validation:Optional
	BundledExtensionsSpec `json:",inline"`
}

// ServerConfigSpec configures main server.
type ServerConfigSpec struct {
	// Listener configuration.
	//+kubebuilder:validation:Optional
	Listener listener.ListenerConfig `json:"listener"`

	// GRPC server configuration.
	//+kubebuilder:validation:Optional
	Grpc grpc.GRPCServerConfig `json:"grpc"`

	// GRPC Gateway configuration.
	//+kubebuilder:validation:Optional
	GrpcGateway grpcgateway.GRPCGatewayConfig `json:"grpc_gateway"`

	// HTTP server configuration.
	//+kubebuilder:validation:Optional
	HTTP http.HTTPServerConfig `json:"http"`

	// TLS configuration.
	//+kubebuilder:validation:Optional
	TLS tlsconfig.ServerTLSConfig `json:"tls"`
}

// ProbeConfigSpec defines liveness and readiness probe configuration.
type ProbeConfigSpec struct {
	// Scheduler settings.
	//+kubebuilder:validation:Optional
	Scheduler jobs.JobGroupConfig `json:"scheduler"`

	// Service settings.
	//+kubebuilder:validation:Optional
	Service jobs.JobConfig `json:"service"`
}

// ClientConfigSpec defines client configuration.
type ClientConfigSpec struct {
	// Proxy settings for the client.
	//+kubebuilder:validation:Optional
	Proxy http.ProxyConfig `json:"proxy"`
}

// BundledExtensionsSpec defines configuration for bundled extensions.
type BundledExtensionsSpec struct {
	// FluxNinja extension configuration.
	//+kubebuilder:validation:Optional
	FluxNinja extconfig.FluxNinjaExtensionConfig `json:"fluxninja"`

	// Sentry extension configuration.
	//+kubebuilder:validation:Optional
	Sentry sentry.SentryConfig `json:"sentry"`
}

// ServiceDiscoverySpec defines configuration for Service discovery.
type ServiceDiscoverySpec struct {
	// KubernetesDiscoveryConfig for Kubernetes service discovery.
	KubernetesDiscoveryConfig kubernetes.KubernetesDiscoveryConfig `json:"kubernetes"`
}

// ControllerClientCertConfig defines configuration for client certificate for Controller.
type ControllerClientCertConfig struct {
	// ConfigMapName is the name of the ConfigMap containing the client certificate which will be mounted at '/etc/aperture/aperture-agent/certs' path with given key name.
	//+kubebuilder:validation:Optional
	ConfigMapName string `json:"configMapName"`

	// ClientCertKeyName is the key name of the client certificate in the ConfigMap.
	//+kubebuilder:validation:Optional
	ClientCertKeyName string `json:"clientCertKeyName" default:"controller-ca.pem"`
}
