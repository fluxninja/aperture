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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"

	agent "github.com/fluxninja/aperture/cmd/aperture-agent/config"
	agentv1alpha1 "github.com/fluxninja/aperture/operator/api/agent/v1alpha1"
	"github.com/fluxninja/aperture/operator/api/common"
	. "github.com/fluxninja/aperture/operator/controllers"
	"github.com/fluxninja/aperture/pkg/agentinfo"
	"github.com/fluxninja/aperture/pkg/distcache"
	"github.com/fluxninja/aperture/pkg/net/listener"
	otelconfig "github.com/fluxninja/aperture/pkg/otelcollector/config"
)

var _ = Describe("Sidecar container for Agent", func() {
	var (
		probe = common.Probe{
			Enabled:             true,
			InitialDelaySeconds: 15,
			PeriodSeconds:       15,
			TimeoutSeconds:      5,
			FailureThreshold:    6,
			SuccessThreshold:    1,
		}
		resourceRequirement = corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				corev1.ResourceCPU: resource.MustParse("0.25"),
			},
			Requests: corev1.ResourceList{
				corev1.ResourceCPU: resource.MustParse("1"),
			},
		}
		lifecycle = &corev1.Lifecycle{
			PreStop: &corev1.LifecycleHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: Test,
				},
			},
		}
	)

	Context("Pod without agent container and default Aperture instance", func() {
		It("returns correct Sidecar container", func() {
			instance := &agentv1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: agentv1alpha1.AgentSpec{
					Sidecar: agentv1alpha1.SidecarSpec{
						Enabled: true,
					},
					ConfigSpec: agentv1alpha1.AgentConfigSpec{
						CommonConfigSpec: common.CommonConfigSpec{
							Server: common.ServerConfigSpec{
								ListenerConfig: listener.ListenerConfig{
									Addr: ":80",
								},
							},
						},
						DistCache: distcache.DistCacheConfig{
							MemberlistBindAddr: ":3322",
							BindAddr:           ":3320",
						},
						OTEL: agent.AgentOTELConfig{
							CommonOTELConfig: otelconfig.CommonOTELConfig{
								Ports: otelconfig.PortsConfig{
									DebugPort:       8888,
									HealthCheckPort: 13133,
									PprofPort:       1777,
									ZpagesPort:      55679,
								},
							},
						},
					},
					Image: common.AgentImage{
						Image: common.Image{
							Registry:   "docker.io/fluxninja",
							Tag:        "latest",
							PullPolicy: "IfNotPresent",
						},
						Repository: "aperture-agent",
					},
				},
			}

			container := corev1.Container{}

			expected := corev1.Container{
				Name:            AgentServiceName,
				Image:           "docker.io/fluxninja/aperture-agent:latest",
				ImagePullPolicy: corev1.PullIfNotPresent,
				SecurityContext: &corev1.SecurityContext{},
				Command:         nil,
				Args:            nil,
				Resources:       corev1.ResourceRequirements{},
				Ports: []corev1.ContainerPort{
					{
						Name:          Server,
						ContainerPort: 80,
						Protocol:      corev1.ProtocolTCP,
					},
					{
						Name:          DistCache,
						ContainerPort: 3320,
						Protocol:      corev1.ProtocolTCP,
					},
					{
						Name:          MemberList,
						ContainerPort: 3322,
						Protocol:      corev1.ProtocolTCP,
					},
					{
						Name:          OtelDebugPort,
						ContainerPort: 8888,
						Protocol:      corev1.ProtocolTCP,
					},
					{
						Name:          OtelHealthcheckPort,
						ContainerPort: 13133,
						Protocol:      corev1.ProtocolTCP,
					},
					{
						Name:          OtelPprofPort,
						ContainerPort: 1777,
						Protocol:      corev1.ProtocolTCP,
					},
					{
						Name:          OtelZpagesPort,
						ContainerPort: 55679,
						Protocol:      corev1.ProtocolTCP,
					},
				},
				LivenessProbe:  nil,
				ReadinessProbe: nil,
				Lifecycle:      nil,
				Env: []corev1.EnvVar{
					{
						Name: "NODE_NAME",
						ValueFrom: &corev1.EnvVarSource{
							FieldRef: &corev1.ObjectFieldSelector{
								APIVersion: V1Version,
								FieldPath:  "spec.nodeName",
							},
						},
					},
					{
						Name: "APERTURE_AGENT_SERVICE_DISCOVERY_KUBERNETES_POD_NAME",
						ValueFrom: &corev1.EnvVarSource{
							FieldRef: &corev1.ObjectFieldSelector{
								APIVersion: V1Version,
								FieldPath:  "metadata.name",
							},
						},
					},
					{
						Name:  "APERTURE_AGENT_SERVICE_DISCOVERY_KUBERNETES_DISCOVERY_ENABLED",
						Value: "false",
					},
				},
				EnvFrom: []corev1.EnvFromSource{},
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      "aperture-agent-config",
						MountPath: "/etc/aperture/aperture-agent/config",
					},
				},
			}

			err := agentContainer(instance.DeepCopy(), &container, "")

			Expect(err).NotTo(HaveOccurred())
			Expect(container).To(Equal(expected))
		})
	})

	Context("Pod with agent container and default Aperture instance", func() {
		It("returns correct Sidecar container", func() {
			instance := &agentv1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: agentv1alpha1.AgentSpec{
					Sidecar: agentv1alpha1.SidecarSpec{
						Enabled: true,
					},
					ConfigSpec: agentv1alpha1.AgentConfigSpec{
						CommonConfigSpec: common.CommonConfigSpec{
							Server: common.ServerConfigSpec{
								ListenerConfig: listener.ListenerConfig{
									Addr: ":80",
								},
							},
						},
						DistCache: distcache.DistCacheConfig{
							MemberlistBindAddr: ":3322",
							BindAddr:           ":3320",
						},
						OTEL: agent.AgentOTELConfig{
							CommonOTELConfig: otelconfig.CommonOTELConfig{
								Ports: otelconfig.PortsConfig{
									DebugPort:       8888,
									HealthCheckPort: 13133,
									PprofPort:       1777,
									ZpagesPort:      55679,
								},
							},
						},
					},
					Image: common.AgentImage{
						Image: common.Image{
							Registry:   "docker.io/fluxninja",
							Tag:        "latest",
							PullPolicy: "IfNotPresent",
						},
						Repository: "aperture-agent",
					},
				},
			}

			container := corev1.Container{
				Name:            AgentServiceName,
				Image:           "auto",
				ImagePullPolicy: corev1.PullAlways,
				SecurityContext: &corev1.SecurityContext{
					RunAsUser: pointer.Int64(1001),
				},
				Command: TestArray,
				Args:    TestArray,
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceCPU: resource.MustParse("1"),
					},
				},
				LivenessProbe: &corev1.Probe{
					InitialDelaySeconds: 10,
				},
				ReadinessProbe: &corev1.Probe{
					InitialDelaySeconds: 10,
				},
				WorkingDir: Test,
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      Test,
						MountPath: Test,
					},
				},
			}

			expected := corev1.Container{
				Name:            AgentServiceName,
				Image:           "docker.io/fluxninja/aperture-agent:latest",
				ImagePullPolicy: corev1.PullAlways,
				SecurityContext: &corev1.SecurityContext{
					RunAsUser: pointer.Int64(1001),
				},
				Command: TestArray,
				Args:    TestArray,
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceCPU: resource.MustParse("1"),
					},
				},
				Ports: []corev1.ContainerPort{
					{
						Name:          Server,
						ContainerPort: 80,
						Protocol:      corev1.ProtocolTCP,
					},
					{
						Name:          DistCache,
						ContainerPort: 3320,
						Protocol:      corev1.ProtocolTCP,
					},
					{
						Name:          MemberList,
						ContainerPort: 3322,
						Protocol:      corev1.ProtocolTCP,
					},
					{
						Name:          OtelDebugPort,
						ContainerPort: 8888,
						Protocol:      corev1.ProtocolTCP,
					},
					{
						Name:          OtelHealthcheckPort,
						ContainerPort: 13133,
						Protocol:      corev1.ProtocolTCP,
					},
					{
						Name:          OtelPprofPort,
						ContainerPort: 1777,
						Protocol:      corev1.ProtocolTCP,
					},
					{
						Name:          OtelZpagesPort,
						ContainerPort: 55679,
						Protocol:      corev1.ProtocolTCP,
					},
				},
				LivenessProbe: &corev1.Probe{
					InitialDelaySeconds: 10,
				},
				ReadinessProbe: &corev1.Probe{
					InitialDelaySeconds: 10,
				},
				Lifecycle:  nil,
				WorkingDir: Test,
				Env: []corev1.EnvVar{
					{
						Name: "NODE_NAME",
						ValueFrom: &corev1.EnvVarSource{
							FieldRef: &corev1.ObjectFieldSelector{
								APIVersion: V1Version,
								FieldPath:  "spec.nodeName",
							},
						},
					},
					{
						Name: "APERTURE_AGENT_SERVICE_DISCOVERY_KUBERNETES_POD_NAME",
						ValueFrom: &corev1.EnvVarSource{
							FieldRef: &corev1.ObjectFieldSelector{
								APIVersion: V1Version,
								FieldPath:  "metadata.name",
							},
						},
					},
					{
						Name:  "APERTURE_AGENT_SERVICE_DISCOVERY_KUBERNETES_DISCOVERY_ENABLED",
						Value: "false",
					},
				},
				EnvFrom: []corev1.EnvFromSource{},
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      Test,
						MountPath: Test,
					},
					{
						Name:      "aperture-agent-config",
						MountPath: "/etc/aperture/aperture-agent/config",
					},
				},
			}

			err := agentContainer(instance.DeepCopy(), &container, "")

			Expect(err).NotTo(HaveOccurred())
			Expect(container).To(Equal(expected))
		})
	})

	Context("Pod without agent container and all Aperture instance variables", func() {
		It("returns correct Sidecar container", func() {
			instance := &agentv1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: agentv1alpha1.AgentSpec{
					Sidecar: agentv1alpha1.SidecarSpec{
						Enabled: true,
					},
					ConfigSpec: agentv1alpha1.AgentConfigSpec{
						AgentInfo: agentinfo.AgentInfoConfig{
							AgentGroup: "test",
						},
						CommonConfigSpec: common.CommonConfigSpec{
							Server: common.ServerConfigSpec{
								ListenerConfig: listener.ListenerConfig{
									Addr: ":8000",
								},
							},
						},
						DistCache: distcache.DistCacheConfig{
							MemberlistBindAddr: ":3322",
							BindAddr:           ":3320",
						},
						OTEL: agent.AgentOTELConfig{
							CommonOTELConfig: otelconfig.CommonOTELConfig{
								Ports: otelconfig.PortsConfig{
									DebugPort:       8888,
									HealthCheckPort: 13133,
									PprofPort:       1777,
									ZpagesPort:      55679,
								},
							},
						},
					},
					CommonSpec: common.CommonSpec{
						LivenessProbe:  probe,
						ReadinessProbe: probe,
						Resources:      resourceRequirement,
						ContainerSecurityContext: common.ContainerSecurityContext{
							Enabled:                true,
							RunAsUser:              0,
							RunAsNonRootUser:       false,
							ReadOnlyRootFilesystem: false,
						},
						Command:        TestArray,
						Args:           TestArray,
						LifecycleHooks: lifecycle,
						ExtraEnvVars: []corev1.EnvVar{
							{
								Name:  Test,
								Value: Test,
							},
						},
						ExtraEnvVarsCM:     Test,
						ExtraEnvVarsSecret: Test,
						ExtraVolumeMounts: []corev1.VolumeMount{
							{
								Name:      Test,
								MountPath: Test,
							},
						},
					},
					Image: common.AgentImage{
						Image: common.Image{
							Registry:    "docker.io/fluxninja",
							Tag:         "latest",
							PullPolicy:  "IfNotPresent",
							PullSecrets: TestArray,
						},
						Repository: "aperture-agent",
					},
				},
			}

			container := corev1.Container{
				EnvFrom: []corev1.EnvFromSource{
					{
						ConfigMapRef: &corev1.ConfigMapEnvSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: TestTwo,
							},
						},
					},
					{
						SecretRef: &corev1.SecretEnvSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: TestTwo,
							},
						},
					},
				},
			}

			expected := corev1.Container{
				Name:            AgentServiceName,
				Image:           "docker.io/fluxninja/aperture-agent:latest",
				ImagePullPolicy: corev1.PullIfNotPresent,
				SecurityContext: &corev1.SecurityContext{
					RunAsUser:              pointer.Int64(0),
					RunAsNonRoot:           pointer.Bool(false),
					ReadOnlyRootFilesystem: pointer.Bool(false),
				},
				Command:   TestArray,
				Args:      TestArray,
				Resources: resourceRequirement,
				Ports: []corev1.ContainerPort{
					{
						Name:          Server,
						ContainerPort: 8000,
						Protocol:      corev1.ProtocolTCP,
					},
					{
						Name:          DistCache,
						ContainerPort: 3320,
						Protocol:      corev1.ProtocolTCP,
					},
					{
						Name:          MemberList,
						ContainerPort: 3322,
						Protocol:      corev1.ProtocolTCP,
					},
					{
						Name:          OtelDebugPort,
						ContainerPort: 8888,
						Protocol:      corev1.ProtocolTCP,
					},
					{
						Name:          OtelHealthcheckPort,
						ContainerPort: 13133,
						Protocol:      corev1.ProtocolTCP,
					},
					{
						Name:          OtelPprofPort,
						ContainerPort: 1777,
						Protocol:      corev1.ProtocolTCP,
					},
					{
						Name:          OtelZpagesPort,
						ContainerPort: 55679,
						Protocol:      corev1.ProtocolTCP,
					},
				},
				LivenessProbe: &corev1.Probe{
					ProbeHandler: corev1.ProbeHandler{
						HTTPGet: &corev1.HTTPGetAction{
							Path:   "/v1/status/subsystem/liveness",
							Port:   intstr.FromString(Server),
							Scheme: corev1.URISchemeHTTP,
						},
					},
					InitialDelaySeconds: probe.InitialDelaySeconds,
					TimeoutSeconds:      probe.TimeoutSeconds,
					PeriodSeconds:       probe.PeriodSeconds,
					FailureThreshold:    probe.FailureThreshold,
					SuccessThreshold:    probe.SuccessThreshold,
				},
				ReadinessProbe: &corev1.Probe{
					ProbeHandler: corev1.ProbeHandler{
						HTTPGet: &corev1.HTTPGetAction{
							Path:   "/v1/status/subsystem/readiness",
							Port:   intstr.FromString(Server),
							Scheme: corev1.URISchemeHTTP,
						},
					},
					InitialDelaySeconds: probe.InitialDelaySeconds,
					TimeoutSeconds:      probe.TimeoutSeconds,
					PeriodSeconds:       probe.PeriodSeconds,
					FailureThreshold:    probe.FailureThreshold,
					SuccessThreshold:    probe.SuccessThreshold,
				},
				Lifecycle: lifecycle,
				Env: []corev1.EnvVar{
					{
						Name:  Test,
						Value: Test,
					},
					{
						Name: "NODE_NAME",
						ValueFrom: &corev1.EnvVarSource{
							FieldRef: &corev1.ObjectFieldSelector{
								APIVersion: V1Version,
								FieldPath:  "spec.nodeName",
							},
						},
					},
					{
						Name:  "APERTURE_AGENT_AGENT_INFO_AGENT_GROUP",
						Value: Test,
					},
					{
						Name: "APERTURE_AGENT_SERVICE_DISCOVERY_KUBERNETES_POD_NAME",
						ValueFrom: &corev1.EnvVarSource{
							FieldRef: &corev1.ObjectFieldSelector{
								APIVersion: V1Version,
								FieldPath:  "metadata.name",
							},
						},
					},
					{
						Name:  "APERTURE_AGENT_SERVICE_DISCOVERY_KUBERNETES_DISCOVERY_ENABLED",
						Value: "false",
					},
				},
				EnvFrom: []corev1.EnvFromSource{
					{
						ConfigMapRef: &corev1.ConfigMapEnvSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: TestTwo,
							},
						},
					},
					{
						SecretRef: &corev1.SecretEnvSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: TestTwo,
							},
						},
					},
					{
						ConfigMapRef: &corev1.ConfigMapEnvSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: Test,
							},
						},
					},
					{
						SecretRef: &corev1.SecretEnvSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: Test,
							},
						},
					},
				},
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      Test,
						MountPath: Test,
					},
					{
						Name:      "aperture-agent-config",
						MountPath: "/etc/aperture/aperture-agent/config",
					},
				},
			}

			err := agentContainer(instance.DeepCopy(), &container, "")

			Expect(err).NotTo(HaveOccurred())
			Expect(container).To(Equal(expected))
		})
	})
})

var _ = Describe("Pod modification for Agent", func() {
	Context("Pod without agent container and default Aperture Instance", func() {
		It("returns correct Pod", func() {
			instance := &agentv1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: agentv1alpha1.AgentSpec{
					Sidecar: agentv1alpha1.SidecarSpec{
						Enabled: true,
					},
					ConfigSpec: agentv1alpha1.AgentConfigSpec{
						CommonConfigSpec: common.CommonConfigSpec{
							Server: common.ServerConfigSpec{
								ListenerConfig: listener.ListenerConfig{
									Addr: ":80",
								},
							},
						},
						DistCache: distcache.DistCacheConfig{
							MemberlistBindAddr: ":3322",
							BindAddr:           ":3320",
						},
						OTEL: agent.AgentOTELConfig{
							CommonOTELConfig: otelconfig.CommonOTELConfig{
								Ports: otelconfig.PortsConfig{
									DebugPort:       8888,
									HealthCheckPort: 13133,
									PprofPort:       1777,
									ZpagesPort:      55679,
								},
							},
						},
					},
					Image: common.AgentImage{
						Image: common.Image{
							Registry:   "docker.io/fluxninja",
							Tag:        "latest",
							PullPolicy: "IfNotPresent",
						},
						Repository: "aperture-agent",
					},
				},
			}

			pod := &corev1.Pod{}

			expected := &corev1.Pod{
				Spec: corev1.PodSpec{
					ImagePullSecrets: []corev1.LocalObjectReference{},
					Containers: []corev1.Container{
						{
							Name:            AgentServiceName,
							Image:           "docker.io/fluxninja/aperture-agent:latest",
							ImagePullPolicy: corev1.PullIfNotPresent,
							SecurityContext: &corev1.SecurityContext{},
							Command:         nil,
							Args:            nil,
							Resources:       corev1.ResourceRequirements{},
							Ports: []corev1.ContainerPort{
								{
									Name:          Server,
									ContainerPort: 80,
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          DistCache,
									ContainerPort: 3320,
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          MemberList,
									ContainerPort: 3322,
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          OtelDebugPort,
									ContainerPort: 8888,
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          OtelHealthcheckPort,
									ContainerPort: 13133,
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          OtelPprofPort,
									ContainerPort: 1777,
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          OtelZpagesPort,
									ContainerPort: 55679,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							LivenessProbe:  nil,
							ReadinessProbe: nil,
							Lifecycle:      nil,
							Env: []corev1.EnvVar{
								{
									Name: "NODE_NAME",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											APIVersion: V1Version,
											FieldPath:  "spec.nodeName",
										},
									},
								},
								{
									Name: "APERTURE_AGENT_SERVICE_DISCOVERY_KUBERNETES_POD_NAME",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											APIVersion: V1Version,
											FieldPath:  "metadata.name",
										},
									},
								},
								{
									Name:  "APERTURE_AGENT_SERVICE_DISCOVERY_KUBERNETES_DISCOVERY_ENABLED",
									Value: "false",
								},
							},
							EnvFrom: []corev1.EnvFromSource{},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "aperture-agent-config",
									MountPath: "/etc/aperture/aperture-agent/config",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "aperture-agent-config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									DefaultMode: pointer.Int32(420),
									LocalObjectReference: corev1.LocalObjectReference{
										Name: AgentServiceName,
									},
								},
							},
						},
					},
				},
			}

			err := agentPod(instance, pod)

			Expect(err).NotTo(HaveOccurred())
			Expect(&pod).To(Equal(&expected))
		})
	})

	Context("Pod with agent container and all Aperture Instance", func() {
		It("returns correct Pod", func() {
			instance := &agentv1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: agentv1alpha1.AgentSpec{
					Sidecar: agentv1alpha1.SidecarSpec{
						Enabled: true,
					},
					ConfigSpec: agentv1alpha1.AgentConfigSpec{
						CommonConfigSpec: common.CommonConfigSpec{
							Server: common.ServerConfigSpec{
								ListenerConfig: listener.ListenerConfig{
									Addr: ":80",
								},
							},
						},
						DistCache: distcache.DistCacheConfig{
							MemberlistBindAddr: ":3322",
							BindAddr:           ":3320",
						},
						OTEL: agent.AgentOTELConfig{
							CommonOTELConfig: otelconfig.CommonOTELConfig{
								Ports: otelconfig.PortsConfig{
									DebugPort:       8888,
									HealthCheckPort: 13133,
									PprofPort:       1777,
									ZpagesPort:      55679,
								},
							},
						},
					},
					CommonSpec: common.CommonSpec{
						InitContainers: []corev1.Container{
							{
								Name: Test,
							},
						},
						ExtraVolumes: []corev1.Volume{
							{
								Name: Test,
								VolumeSource: corev1.VolumeSource{
									EmptyDir: &corev1.EmptyDirVolumeSource{},
								},
							},
						},
					},
					Image: common.AgentImage{
						Image: common.Image{
							Registry:    "docker.io/fluxninja",
							Tag:         "latest",
							PullPolicy:  "IfNotPresent",
							PullSecrets: TestArrayTwo,
						},
						Repository: "aperture-agent",
					},
				},
			}

			pod := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						AgentGroupKey: Test,
					},
				},
				Spec: corev1.PodSpec{
					InitContainers: []corev1.Container{
						{
							Name: TestTwo,
						},
					},
					Containers: []corev1.Container{
						{
							Name:            AgentServiceName,
							Image:           "auto",
							ImagePullPolicy: corev1.PullNever,
						},
					},
				},
			}

			expected := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						AgentGroupKey: Test,
					},
				},
				Spec: corev1.PodSpec{
					ImagePullSecrets: []corev1.LocalObjectReference{
						{
							Name: TestTwo,
						},
						{
							Name: Test,
						},
					},
					InitContainers: []corev1.Container{
						{
							Name: TestTwo,
						},
						{
							Name: Test,
						},
					},
					Containers: []corev1.Container{
						{
							Name:            AgentServiceName,
							Image:           "docker.io/fluxninja/aperture-agent:latest",
							ImagePullPolicy: corev1.PullNever,
							SecurityContext: &corev1.SecurityContext{},
							Command:         nil,
							Args:            nil,
							Resources:       corev1.ResourceRequirements{},
							Ports: []corev1.ContainerPort{
								{
									Name:          Server,
									ContainerPort: 80,
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          DistCache,
									ContainerPort: 3320,
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          MemberList,
									ContainerPort: 3322,
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          OtelDebugPort,
									ContainerPort: 8888,
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          OtelHealthcheckPort,
									ContainerPort: 13133,
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          OtelPprofPort,
									ContainerPort: 1777,
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          OtelZpagesPort,
									ContainerPort: 55679,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							LivenessProbe:  nil,
							ReadinessProbe: nil,
							Lifecycle:      nil,
							Env: []corev1.EnvVar{
								{
									Name: "NODE_NAME",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											APIVersion: V1Version,
											FieldPath:  "spec.nodeName",
										},
									},
								},
								{
									Name:  "APERTURE_AGENT_AGENT_INFO_AGENT_GROUP",
									Value: Test,
								},
								{
									Name: "APERTURE_AGENT_SERVICE_DISCOVERY_KUBERNETES_POD_NAME",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											APIVersion: V1Version,
											FieldPath:  "metadata.name",
										},
									},
								},
								{
									Name:  "APERTURE_AGENT_SERVICE_DISCOVERY_KUBERNETES_DISCOVERY_ENABLED",
									Value: "false",
								},
							},
							EnvFrom: []corev1.EnvFromSource{},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "aperture-agent-config",
									MountPath: "/etc/aperture/aperture-agent/config",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: Test,
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
						{
							Name: "aperture-agent-config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									DefaultMode: pointer.Int32(420),
									LocalObjectReference: corev1.LocalObjectReference{
										Name: AgentServiceName,
									},
								},
							},
						},
					},
				},
			}

			err := agentPod(instance, pod)

			Expect(err).NotTo(HaveOccurred())
			Expect(&pod).To(Equal(&expected))
		})
	})
})
