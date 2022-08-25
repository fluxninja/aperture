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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"

	"github.com/fluxninja/aperture/operator/api/v1alpha1"
)

var _ = Describe("Sidecar container for Agent", func() {
	var (
		probe = v1alpha1.Probe{
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
					Path: test,
				},
			},
		}
	)

	Context("Pod without agent container and default Aperture instance", func() {
		It("returns correct Sidecar container", func() {
			instance := &v1alpha1.Aperture{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.ApertureSpec{
					Sidecar: v1alpha1.SidecarSpec{
						Enabled: true,
					},
					Agent: v1alpha1.AgentSpec{
						CommonSpec: v1alpha1.CommonSpec{
							ServerPort: 80,
						},
						Image: v1alpha1.Image{
							Registry:   "gcr.io/devel-309501/cf-fn",
							Repository: "aperture-agent",
							Tag:        "latest",
							PullPolicy: "IfNotPresent",
						},
					},
				},
			}

			container := corev1.Container{}

			expected := corev1.Container{
				Name:            agentServiceName,
				Image:           "gcr.io/devel-309501/cf-fn/aperture-agent:latest",
				ImagePullPolicy: corev1.PullIfNotPresent,
				SecurityContext: &corev1.SecurityContext{},
				Command:         nil,
				Args:            nil,
				Resources:       corev1.ResourceRequirements{},
				Ports: []corev1.ContainerPort{
					{
						Name:          "grpc",
						ContainerPort: 80,
					},
					{
						Name:          "grpc-otel",
						ContainerPort: 4317,
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
								APIVersion: v1Version,
								FieldPath:  "spec.nodeName",
							},
						},
					},
					{
						Name: "APERTURE_AGENT_SERVICE_DISCOVERY_KUBERNETES_POD_NAME",
						ValueFrom: &corev1.EnvVarSource{
							FieldRef: &corev1.ObjectFieldSelector{
								APIVersion: v1Version,
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

			agentContainer(instance.DeepCopy(), &container, "")
			Expect(container).To(Equal(expected))
		})
	})

	Context("Pod with agent container and default Aperture instance", func() {
		It("returns correct Sidecar container", func() {
			instance := &v1alpha1.Aperture{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.ApertureSpec{
					Sidecar: v1alpha1.SidecarSpec{
						Enabled: true,
					},
					Agent: v1alpha1.AgentSpec{
						CommonSpec: v1alpha1.CommonSpec{
							ServerPort: 80,
						},
						Image: v1alpha1.Image{
							Registry:   "gcr.io/devel-309501/cf-fn",
							Repository: "aperture-agent",
							Tag:        "latest",
							PullPolicy: "IfNotPresent",
						},
					},
				},
			}

			container := corev1.Container{
				Name:            agentServiceName,
				Image:           "auto",
				ImagePullPolicy: corev1.PullAlways,
				SecurityContext: &corev1.SecurityContext{
					RunAsUser: pointer.Int64Ptr(1001),
				},
				Command: testArray,
				Args:    testArray,
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
				WorkingDir: test,
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      test,
						MountPath: test,
					},
				},
			}

			expected := corev1.Container{
				Name:            agentServiceName,
				Image:           "gcr.io/devel-309501/cf-fn/aperture-agent:latest",
				ImagePullPolicy: corev1.PullAlways,
				SecurityContext: &corev1.SecurityContext{
					RunAsUser: pointer.Int64Ptr(1001),
				},
				Command: testArray,
				Args:    testArray,
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceCPU: resource.MustParse("1"),
					},
				},
				Ports: []corev1.ContainerPort{
					{
						Name:          "grpc",
						ContainerPort: 80,
					},
					{
						Name:          "grpc-otel",
						ContainerPort: 4317,
					},
				},
				LivenessProbe: &corev1.Probe{
					InitialDelaySeconds: 10,
				},
				ReadinessProbe: &corev1.Probe{
					InitialDelaySeconds: 10,
				},
				Lifecycle:  nil,
				WorkingDir: test,
				Env: []corev1.EnvVar{
					{
						Name: "NODE_NAME",
						ValueFrom: &corev1.EnvVarSource{
							FieldRef: &corev1.ObjectFieldSelector{
								APIVersion: v1Version,
								FieldPath:  "spec.nodeName",
							},
						},
					},
					{
						Name: "APERTURE_AGENT_SERVICE_DISCOVERY_KUBERNETES_POD_NAME",
						ValueFrom: &corev1.EnvVarSource{
							FieldRef: &corev1.ObjectFieldSelector{
								APIVersion: v1Version,
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
						Name:      test,
						MountPath: test,
					},
					{
						Name:      "aperture-agent-config",
						MountPath: "/etc/aperture/aperture-agent/config",
					},
				},
			}

			agentContainer(instance.DeepCopy(), &container, "")
			Expect(container).To(Equal(expected))
		})
	})

	Context("Pod without agent container and all Aperture instance variables", func() {
		It("returns correct Sidecar container", func() {
			instance := &v1alpha1.Aperture{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.ApertureSpec{
					Sidecar: v1alpha1.SidecarSpec{
						Enabled: true,
					},
					Agent: v1alpha1.AgentSpec{
						AgentGroup: test,
						CommonSpec: v1alpha1.CommonSpec{
							ServerPort:     8000,
							LivenessProbe:  probe,
							ReadinessProbe: probe,
							Resources:      resourceRequirement,
							ContainerSecurityContext: v1alpha1.ContainerSecurityContext{
								Enabled:                true,
								RunAsUser:              pointer.Int64Ptr(0),
								RunAsNonRootUser:       pointer.BoolPtr(false),
								ReadOnlyRootFilesystem: pointer.BoolPtr(false),
							},
							Command:        testArray,
							Args:           testArray,
							LifecycleHooks: lifecycle,
							ExtraEnvVars: []corev1.EnvVar{
								{
									Name:  test,
									Value: test,
								},
							},
							ExtraEnvVarsCM:     test,
							ExtraEnvVarsSecret: test,
							ExtraVolumeMounts: []corev1.VolumeMount{
								{
									Name:      test,
									MountPath: test,
								},
							},
						},
						Image: v1alpha1.Image{
							Registry:    "gcr.io/devel-309501/cf-fn",
							Repository:  "aperture-agent",
							Tag:         "latest",
							PullPolicy:  "IfNotPresent",
							PullSecrets: testArray,
						},
					},
				},
			}

			container := corev1.Container{
				EnvFrom: []corev1.EnvFromSource{
					{
						ConfigMapRef: &corev1.ConfigMapEnvSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: testTwo,
							},
						},
					},
					{
						SecretRef: &corev1.SecretEnvSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: testTwo,
							},
						},
					},
				},
			}

			expected := corev1.Container{
				Name:            agentServiceName,
				Image:           "gcr.io/devel-309501/cf-fn/aperture-agent:latest",
				ImagePullPolicy: corev1.PullIfNotPresent,
				SecurityContext: &corev1.SecurityContext{
					RunAsUser:              pointer.Int64Ptr(0),
					RunAsNonRoot:           pointer.BoolPtr(false),
					ReadOnlyRootFilesystem: pointer.BoolPtr(false),
				},
				Command:   testArray,
				Args:      testArray,
				Resources: resourceRequirement,
				Ports: []corev1.ContainerPort{
					{
						Name:          "grpc",
						ContainerPort: 8000,
					},
					{
						Name:          "grpc-otel",
						ContainerPort: 4317,
					},
				},
				LivenessProbe: &corev1.Probe{
					ProbeHandler: corev1.ProbeHandler{
						HTTPGet: &corev1.HTTPGetAction{
							Path:   "/v1/status/liveness",
							Port:   intstr.FromString("grpc"),
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
							Path:   "/v1/status/readiness",
							Port:   intstr.FromString("grpc"),
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
						Name:  test,
						Value: test,
					},
					{
						Name: "NODE_NAME",
						ValueFrom: &corev1.EnvVarSource{
							FieldRef: &corev1.ObjectFieldSelector{
								APIVersion: v1Version,
								FieldPath:  "spec.nodeName",
							},
						},
					},
					{
						Name:  "APERTURE_AGENT_AGENT_INFO_AGENT_GROUP",
						Value: test,
					},
					{
						Name: "APERTURE_AGENT_SERVICE_DISCOVERY_KUBERNETES_POD_NAME",
						ValueFrom: &corev1.EnvVarSource{
							FieldRef: &corev1.ObjectFieldSelector{
								APIVersion: v1Version,
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
								Name: testTwo,
							},
						},
					},
					{
						SecretRef: &corev1.SecretEnvSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: testTwo,
							},
						},
					},
					{
						ConfigMapRef: &corev1.ConfigMapEnvSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: test,
							},
						},
					},
					{
						SecretRef: &corev1.SecretEnvSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: test,
							},
						},
					},
				},
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      test,
						MountPath: test,
					},
					{
						Name:      "aperture-agent-config",
						MountPath: "/etc/aperture/aperture-agent/config",
					},
				},
			}

			agentContainer(instance.DeepCopy(), &container, "")
			Expect(container).To(Equal(expected))
		})
	})
})

var _ = Describe("Pod modification for Agent", func() {
	Context("Pod without agent container and default Aperture Instance", func() {
		It("returns correct Pod", func() {
			instance := &v1alpha1.Aperture{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.ApertureSpec{
					Sidecar: v1alpha1.SidecarSpec{
						Enabled: true,
					},
					Agent: v1alpha1.AgentSpec{
						CommonSpec: v1alpha1.CommonSpec{
							ServerPort: 80,
						},
						Image: v1alpha1.Image{
							Registry:   "gcr.io/devel-309501/cf-fn",
							Repository: "aperture-agent",
							Tag:        "latest",
							PullPolicy: "IfNotPresent",
						},
					},
				},
			}

			pod := &corev1.Pod{}

			expected := &corev1.Pod{
				Spec: corev1.PodSpec{
					ImagePullSecrets: []corev1.LocalObjectReference{},
					Containers: []corev1.Container{
						{
							Name:            agentServiceName,
							Image:           "gcr.io/devel-309501/cf-fn/aperture-agent:latest",
							ImagePullPolicy: corev1.PullIfNotPresent,
							SecurityContext: &corev1.SecurityContext{},
							Command:         nil,
							Args:            nil,
							Resources:       corev1.ResourceRequirements{},
							Ports: []corev1.ContainerPort{
								{
									Name:          "grpc",
									ContainerPort: 80,
								},
								{
									Name:          "grpc-otel",
									ContainerPort: 4317,
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
											APIVersion: v1Version,
											FieldPath:  "spec.nodeName",
										},
									},
								},
								{
									Name: "APERTURE_AGENT_SERVICE_DISCOVERY_KUBERNETES_POD_NAME",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											APIVersion: v1Version,
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
									DefaultMode: pointer.Int32Ptr(420),
									LocalObjectReference: corev1.LocalObjectReference{
										Name: agentServiceName,
									},
								},
							},
						},
					},
				},
			}

			agentPod(instance, pod)
			Expect(&pod).To(Equal(&expected))
		})
	})

	Context("Pod with agent container and all Aperture Instance", func() {
		It("returns correct Pod", func() {
			instance := &v1alpha1.Aperture{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.ApertureSpec{
					ImagePullSecrets: testArray,
					Sidecar: v1alpha1.SidecarSpec{
						Enabled: true,
					},
					Agent: v1alpha1.AgentSpec{
						CommonSpec: v1alpha1.CommonSpec{
							ServerPort: 80,
							InitContainers: []corev1.Container{
								{
									Name: test,
								},
							},
							ExtraVolumes: []corev1.Volume{
								{
									Name: test,
									VolumeSource: corev1.VolumeSource{
										EmptyDir: &corev1.EmptyDirVolumeSource{},
									},
								},
							},
						},
						Image: v1alpha1.Image{
							Registry:    "gcr.io/devel-309501/cf-fn",
							Repository:  "aperture-agent",
							Tag:         "latest",
							PullPolicy:  "IfNotPresent",
							PullSecrets: testArrayTwo,
						},
					},
				},
			}

			pod := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						agentGroupKey: test,
					},
				},
				Spec: corev1.PodSpec{
					InitContainers: []corev1.Container{
						{
							Name: testTwo,
						},
					},
					Containers: []corev1.Container{
						{
							Name:            agentServiceName,
							Image:           "auto",
							ImagePullPolicy: corev1.PullNever,
						},
					},
				},
			}

			expected := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						agentGroupKey: test,
					},
				},
				Spec: corev1.PodSpec{
					ImagePullSecrets: []corev1.LocalObjectReference{
						{
							Name: testTwo,
						},
						{
							Name: test,
						},
					},
					InitContainers: []corev1.Container{
						{
							Name: testTwo,
						},
						{
							Name: test,
						},
					},
					Containers: []corev1.Container{
						{
							Name:            agentServiceName,
							Image:           "gcr.io/devel-309501/cf-fn/aperture-agent:latest",
							ImagePullPolicy: corev1.PullNever,
							SecurityContext: &corev1.SecurityContext{},
							Command:         nil,
							Args:            nil,
							Resources:       corev1.ResourceRequirements{},
							Ports: []corev1.ContainerPort{
								{
									Name:          "grpc",
									ContainerPort: 80,
								},
								{
									Name:          "grpc-otel",
									ContainerPort: 4317,
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
											APIVersion: v1Version,
											FieldPath:  "spec.nodeName",
										},
									},
								},
								{
									Name:  "APERTURE_AGENT_AGENT_INFO_AGENT_GROUP",
									Value: test,
								},
								{
									Name: "APERTURE_AGENT_SERVICE_DISCOVERY_KUBERNETES_POD_NAME",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											APIVersion: v1Version,
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
							Name: test,
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
						{
							Name: "aperture-agent-config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									DefaultMode: pointer.Int32Ptr(420),
									LocalObjectReference: corev1.LocalObjectReference{
										Name: agentServiceName,
									},
								},
							},
						},
					},
				},
			}

			agentPod(instance, pod)
			Expect(&pod).To(Equal(&expected))
		})
	})
})
