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

package controllers

import (
	corev1 "k8s.io/api/core/v1"

	"github.com/fluxninja/aperture/operator/api/v1alpha1"
)

// agentContainer prepares Sidecar container for the Agent based on the received parameters.
func agentContainer(instance *v1alpha1.Agent, container *corev1.Container, agentGroup string) {
	spec := instance.Spec
	livenessProbe, readinessProbe := containerProbes(spec.CommonSpec)
	container.Name = agentServiceName

	if container.Image == "" || container.Image == "auto" {
		container.Image = imageString(spec.Image)
	}

	if container.ImagePullPolicy == "" {
		container.ImagePullPolicy = corev1.PullPolicy(spec.Image.PullPolicy)
	}

	if container.SecurityContext == nil {
		container.SecurityContext = containerSecurityContext(spec.ContainerSecurityContext)
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

	// Not allowing Port override until it is supported by Agent
	container.Ports = []corev1.ContainerPort{
		{
			Name:          "grpc",
			ContainerPort: int32(spec.ServerPort),
		},
		{
			Name:          "grpc-otel",
			ContainerPort: 4317,
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

	container.Env = mergeEnvVars(agentEnv(instance, agentGroup), container.Env)
	container.EnvFrom = mergeEnvFromSources(containerEnvFrom(spec.CommonSpec), container.EnvFrom)
	container.VolumeMounts = mergeVolumeMounts(agentVolumeMounts(spec), container.VolumeMounts)
}

// agentPod updates the received Pod spec to add Sidecar for the Agent.
func agentPod(instance *v1alpha1.Agent, pod *corev1.Pod) {
	apec := instance.Spec
	agentGroup := ""
	if pod.Annotations != nil {
		agentGroup = pod.Annotations[agentGroupKey]
	}

	container := corev1.Container{}
	var containerIndex int
	appendContainer := true
	for index, cont := range pod.Spec.Containers {
		if cont.Name == agentServiceName {
			container = cont
			containerIndex = index
			appendContainer = false
		}
	}

	agentContainer(instance, &container, agentGroup)
	if appendContainer {
		pod.Spec.Containers = append(pod.Spec.Containers, container)
	} else {
		pod.Spec.Containers[containerIndex] = container
	}

	pod.Spec.ImagePullSecrets = mergeImagePullSecrets(imagePullSecrets(apec.Image), pod.Spec.ImagePullSecrets)
	pod.Spec.InitContainers = mergeContainers(apec.InitContainers, pod.Spec.InitContainers)
	pod.Spec.Volumes = mergeVolumes(agentVolumes(apec), pod.Spec.Volumes)
}
