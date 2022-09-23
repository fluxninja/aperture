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

package mutatingwebhook

import (
	"fmt"

	"github.com/fluxninja/aperture/operator/controllers"

	corev1 "k8s.io/api/core/v1"

	"github.com/fluxninja/aperture/operator/api/agent/v1alpha1"
)

// agentContainer prepares Sidecar container for the Agent based on the received parameters.
func agentContainer(instance *v1alpha1.Agent, container *corev1.Container, agentGroup string) error {
	spec := instance.Spec
	probeScheme := corev1.URISchemeHTTP
	if instance.Spec.ConfigSpec.Server.TLS.Enabled {
		probeScheme = corev1.URISchemeHTTPS
	}
	livenessProbe, readinessProbe := controllers.ContainerProbes(spec.CommonSpec, probeScheme)
	container.Name = controllers.AgentServiceName

	if container.Image == "" || container.Image == "auto" {
		container.Image = controllers.ImageString(spec.Image.Image, spec.Image.Repository)
	}

	if container.ImagePullPolicy == "" {
		container.ImagePullPolicy = corev1.PullPolicy(spec.Image.PullPolicy)
	}

	if container.SecurityContext == nil {
		container.SecurityContext = controllers.ContainerSecurityContext(spec.ContainerSecurityContext)
	}

	if container.Command == nil {
		container.Command = spec.Command
	}

	if container.Args == nil {
		container.Args = spec.Args
	}

	if container.Resources.Limits == nil {
		container.Resources.Limits = spec.Resources.Limits
	}

	if container.Resources.Requests == nil {
		container.Resources.Requests = spec.Resources.Requests
	}

	serverPort, err := controllers.GetPort(spec.ConfigSpec.Server.Addr)
	if err != nil {
		return fmt.Errorf("invalid value '%v' provided for 'server.addr' config", spec.ConfigSpec.Server.Addr)
	}

	distCachePort, err := controllers.GetPort(spec.ConfigSpec.DistCache.BindAddr)
	if err != nil {
		return fmt.Errorf("invalid value '%v' provided for 'dist_cache.bind_addr' config", spec.ConfigSpec.DistCache.BindAddr)
	}

	memberListPort, err := controllers.GetPort(spec.ConfigSpec.DistCache.MemberlistBindAddr)
	if err != nil {
		return fmt.Errorf("invalid value '%v' provided for 'dist_cache.memberlist_bind_addr' config", spec.ConfigSpec.DistCache.MemberlistBindAddr)
	}

	container.Ports = []corev1.ContainerPort{
		{
			Name:          controllers.Server,
			ContainerPort: serverPort,
			Protocol:      corev1.ProtocolTCP,
		},
		{
			Name:          controllers.DistCache,
			ContainerPort: distCachePort,
			Protocol:      corev1.ProtocolTCP,
		},
		{
			Name:          controllers.MemberList,
			ContainerPort: memberListPort,
			Protocol:      corev1.ProtocolTCP,
		},
	}

	if container.LivenessProbe == nil {
		container.LivenessProbe = livenessProbe
	}

	if container.ReadinessProbe == nil {
		container.ReadinessProbe = readinessProbe
	}

	if container.Lifecycle == nil {
		container.Lifecycle = spec.LifecycleHooks
	}

	container.Env = controllers.MergeEnvVars(controllers.AgentEnv(instance, agentGroup), container.Env)
	container.EnvFrom = controllers.MergeEnvFromSources(controllers.ContainerEnvFrom(spec.CommonSpec), container.EnvFrom)
	container.VolumeMounts = controllers.MergeVolumeMounts(controllers.AgentVolumeMounts(spec), container.VolumeMounts)

	return nil
}

// agentPod updates the received Pod spec to add Sidecar for the Agent.
func agentPod(instance *v1alpha1.Agent, pod *corev1.Pod) error {
	spec := instance.Spec
	agentGroup := ""
	if pod.Annotations != nil {
		agentGroup = pod.Annotations[controllers.AgentGroupKey]
	}

	container := corev1.Container{}
	var containerIndex int
	appendContainer := true
	for index, cont := range pod.Spec.Containers {
		if cont.Name == controllers.AgentServiceName {
			container = cont
			containerIndex = index
			appendContainer = false
		}
	}

	err := agentContainer(instance, &container, agentGroup)
	if err != nil {
		return err
	}
	if appendContainer {
		pod.Spec.Containers = append(pod.Spec.Containers, container)
	} else {
		pod.Spec.Containers[containerIndex] = container
	}

	pod.Spec.ImagePullSecrets = controllers.MergeImagePullSecrets(controllers.ImagePullSecrets(spec.Image.Image), pod.Spec.ImagePullSecrets)
	pod.Spec.InitContainers = controllers.MergeContainers(spec.InitContainers, pod.Spec.InitContainers)
	pod.Spec.Volumes = controllers.MergeVolumes(controllers.AgentVolumes(spec), pod.Spec.Volumes)

	return nil
}
