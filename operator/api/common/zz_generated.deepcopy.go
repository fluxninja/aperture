//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package common

import (
	"k8s.io/api/core/v1"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *APIKeySecret) DeepCopyInto(out *APIKeySecret) {
	*out = *in
	out.SecretKeyRef = in.SecretKeyRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new APIKeySecret.
func (in *APIKeySecret) DeepCopy() *APIKeySecret {
	if in == nil {
		return nil
	}
	out := new(APIKeySecret)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *APIKeySecretSpec) DeepCopyInto(out *APIKeySecretSpec) {
	*out = *in
	out.Agent = in.Agent
	out.Controller = in.Controller
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new APIKeySecretSpec.
func (in *APIKeySecretSpec) DeepCopy() *APIKeySecretSpec {
	if in == nil {
		return nil
	}
	out := new(APIKeySecretSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AgentImage) DeepCopyInto(out *AgentImage) {
	*out = *in
	in.Image.DeepCopyInto(&out.Image)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AgentImage.
func (in *AgentImage) DeepCopy() *AgentImage {
	if in == nil {
		return nil
	}
	out := new(AgentImage)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BundledExtensionsSpec) DeepCopyInto(out *BundledExtensionsSpec) {
	*out = *in
	in.FluxNinja.DeepCopyInto(&out.FluxNinja)
	out.Sentry = in.Sentry
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BundledExtensionsSpec.
func (in *BundledExtensionsSpec) DeepCopy() *BundledExtensionsSpec {
	if in == nil {
		return nil
	}
	out := new(BundledExtensionsSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClientConfigSpec) DeepCopyInto(out *ClientConfigSpec) {
	*out = *in
	in.Proxy.DeepCopyInto(&out.Proxy)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClientConfigSpec.
func (in *ClientConfigSpec) DeepCopy() *ClientConfigSpec {
	if in == nil {
		return nil
	}
	out := new(ClientConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CommonConfigSpec) DeepCopyInto(out *CommonConfigSpec) {
	*out = *in
	in.Client.DeepCopyInto(&out.Client)
	in.Liveness.DeepCopyInto(&out.Liveness)
	in.Readiness.DeepCopyInto(&out.Readiness)
	in.Log.DeepCopyInto(&out.Log)
	out.Metrics = in.Metrics
	out.Profilers = in.Profilers
	in.TokenSource.DeepCopyInto(&out.TokenSource)
	in.Server.DeepCopyInto(&out.Server)
	in.Watchdog.DeepCopyInto(&out.Watchdog)
	in.Alertmanagers.DeepCopyInto(&out.Alertmanagers)
	in.BundledExtensionsSpec.DeepCopyInto(&out.BundledExtensionsSpec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CommonConfigSpec.
func (in *CommonConfigSpec) DeepCopy() *CommonConfigSpec {
	if in == nil {
		return nil
	}
	out := new(CommonConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CommonSpec) DeepCopyInto(out *CommonSpec) {
	*out = *in
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	in.Service.DeepCopyInto(&out.Service)
	in.ServiceAccountSpec.DeepCopyInto(&out.ServiceAccountSpec)
	out.LivenessProbe = in.LivenessProbe
	out.ReadinessProbe = in.ReadinessProbe
	if in.CustomLivenessProbe != nil {
		in, out := &in.CustomLivenessProbe, &out.CustomLivenessProbe
		*out = new(v1.Probe)
		(*in).DeepCopyInto(*out)
	}
	if in.CustomReadinessProbe != nil {
		in, out := &in.CustomReadinessProbe, &out.CustomReadinessProbe
		*out = new(v1.Probe)
		(*in).DeepCopyInto(*out)
	}
	in.Resources.DeepCopyInto(&out.Resources)
	out.PodSecurityContext = in.PodSecurityContext
	out.ContainerSecurityContext = in.ContainerSecurityContext
	if in.Command != nil {
		in, out := &in.Command, &out.Command
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Args != nil {
		in, out := &in.Args, &out.Args
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.PodLabels != nil {
		in, out := &in.PodLabels, &out.PodLabels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.PodAnnotations != nil {
		in, out := &in.PodAnnotations, &out.PodAnnotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Affinity != nil {
		in, out := &in.Affinity, &out.Affinity
		*out = new(v1.Affinity)
		(*in).DeepCopyInto(*out)
	}
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Tolerations != nil {
		in, out := &in.Tolerations, &out.Tolerations
		*out = make([]v1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.LifecycleHooks != nil {
		in, out := &in.LifecycleHooks, &out.LifecycleHooks
		*out = new(v1.Lifecycle)
		(*in).DeepCopyInto(*out)
	}
	if in.ExtraEnvVars != nil {
		in, out := &in.ExtraEnvVars, &out.ExtraEnvVars
		*out = make([]v1.EnvVar, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.ExtraVolumes != nil {
		in, out := &in.ExtraVolumes, &out.ExtraVolumes
		*out = make([]v1.Volume, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.ExtraVolumeMounts != nil {
		in, out := &in.ExtraVolumeMounts, &out.ExtraVolumeMounts
		*out = make([]v1.VolumeMount, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Sidecars != nil {
		in, out := &in.Sidecars, &out.Sidecars
		*out = make([]v1.Container, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.InitContainers != nil {
		in, out := &in.InitContainers, &out.InitContainers
		*out = make([]v1.Container, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	out.Secrets = in.Secrets
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CommonSpec.
func (in *CommonSpec) DeepCopy() *CommonSpec {
	if in == nil {
		return nil
	}
	out := new(CommonSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ContainerSecurityContext) DeepCopyInto(out *ContainerSecurityContext) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ContainerSecurityContext.
func (in *ContainerSecurityContext) DeepCopy() *ContainerSecurityContext {
	if in == nil {
		return nil
	}
	out := new(ContainerSecurityContext)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ControllerClientCertConfig) DeepCopyInto(out *ControllerClientCertConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ControllerClientCertConfig.
func (in *ControllerClientCertConfig) DeepCopy() *ControllerClientCertConfig {
	if in == nil {
		return nil
	}
	out := new(ControllerClientCertConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ControllerImage) DeepCopyInto(out *ControllerImage) {
	*out = *in
	in.Image.DeepCopyInto(&out.Image)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ControllerImage.
func (in *ControllerImage) DeepCopy() *ControllerImage {
	if in == nil {
		return nil
	}
	out := new(ControllerImage)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Image) DeepCopyInto(out *Image) {
	*out = *in
	if in.PullSecrets != nil {
		in, out := &in.PullSecrets, &out.PullSecrets
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Image.
func (in *Image) DeepCopy() *Image {
	if in == nil {
		return nil
	}
	out := new(Image)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PodSecurityContext) DeepCopyInto(out *PodSecurityContext) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodSecurityContext.
func (in *PodSecurityContext) DeepCopy() *PodSecurityContext {
	if in == nil {
		return nil
	}
	out := new(PodSecurityContext)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Probe) DeepCopyInto(out *Probe) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Probe.
func (in *Probe) DeepCopy() *Probe {
	if in == nil {
		return nil
	}
	out := new(Probe)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ProbeConfigSpec) DeepCopyInto(out *ProbeConfigSpec) {
	*out = *in
	out.Scheduler = in.Scheduler
	in.Service.DeepCopyInto(&out.Service)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ProbeConfigSpec.
func (in *ProbeConfigSpec) DeepCopy() *ProbeConfigSpec {
	if in == nil {
		return nil
	}
	out := new(ProbeConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretKeyRef) DeepCopyInto(out *SecretKeyRef) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretKeyRef.
func (in *SecretKeyRef) DeepCopy() *SecretKeyRef {
	if in == nil {
		return nil
	}
	out := new(SecretKeyRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Secrets) DeepCopyInto(out *Secrets) {
	*out = *in
	out.FluxNinjaExtension = in.FluxNinjaExtension
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Secrets.
func (in *Secrets) DeepCopy() *Secrets {
	if in == nil {
		return nil
	}
	out := new(Secrets)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServerConfigSpec) DeepCopyInto(out *ServerConfigSpec) {
	*out = *in
	in.Listener.DeepCopyInto(&out.Listener)
	in.Grpc.DeepCopyInto(&out.Grpc)
	out.GrpcGateway = in.GrpcGateway
	in.HTTP.DeepCopyInto(&out.HTTP)
	out.TLS = in.TLS
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServerConfigSpec.
func (in *ServerConfigSpec) DeepCopy() *ServerConfigSpec {
	if in == nil {
		return nil
	}
	out := new(ServerConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Service) DeepCopyInto(out *Service) {
	*out = *in
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Service.
func (in *Service) DeepCopy() *Service {
	if in == nil {
		return nil
	}
	out := new(Service)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceAccountSpec) DeepCopyInto(out *ServiceAccountSpec) {
	*out = *in
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceAccountSpec.
func (in *ServiceAccountSpec) DeepCopy() *ServiceAccountSpec {
	if in == nil {
		return nil
	}
	out := new(ServiceAccountSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceDiscoverySpec) DeepCopyInto(out *ServiceDiscoverySpec) {
	*out = *in
	out.KubernetesDiscoveryConfig = in.KubernetesDiscoveryConfig
	in.StaticDiscoveryConfig.DeepCopyInto(&out.StaticDiscoveryConfig)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceDiscoverySpec.
func (in *ServiceDiscoverySpec) DeepCopy() *ServiceDiscoverySpec {
	if in == nil {
		return nil
	}
	out := new(ServiceDiscoverySpec)
	in.DeepCopyInto(out)
	return out
}