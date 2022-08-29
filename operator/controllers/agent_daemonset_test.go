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
	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/utils/pointer"

	"github.com/fluxninja/aperture/operator/api/v1alpha1"
)

var _ = Describe("Agent Daemonset", func() {
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
		tolerations = []corev1.Toleration{
			{
				Key:      test,
				Value:    test,
				Operator: corev1.TolerationOpEqual,
			},
		}
		lifecycle = &corev1.Lifecycle{
			PreStop: &corev1.LifecycleHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: test,
				},
			},
		}
		affinity = &corev1.Affinity{
			PodAffinity: &corev1.PodAffinity{
				PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
					{
						Weight: 1,
						PodAffinityTerm: corev1.PodAffinityTerm{
							LabelSelector: &metav1.LabelSelector{
								MatchLabels: testMap,
							},
						},
					},
				},
			},
		}
	)

	Context("Instance with default parameters", func() {
		It("returns correct daemonset for Agent", func() {
			selectorLabels := map[string]string{
				"app.kubernetes.io/name":       appName,
				"app.kubernetes.io/instance":   appName,
				"app.kubernetes.io/managed-by": operatorName,
				"app.kubernetes.io/component":  agentServiceName,
			}

			instance := &v1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.AgentSpec{
					CommonSpec: v1alpha1.CommonSpec{
						ServerPort: 80,
					},
					DistributedCachePort: 3320,
					MemberListPort:       3322,
					Image: v1alpha1.Image{
						Registry:   "docker.io/fluxninja",
						Repository: "aperture-agent",
						Tag:        "latest",
						PullPolicy: "IfNotPresent",
					},
				},
			}
			expected := &appsv1.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      agentServiceName,
					Namespace: appName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       appName,
						"app.kubernetes.io/instance":   instance.GetName(),
						"app.kubernetes.io/managed-by": operatorName,
						"app.kubernetes.io/component":  agentServiceName,
					},
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion:         "fluxninja.com/v1alpha1",
							Name:               instance.GetName(),
							Kind:               "Agent",
							Controller:         pointer.BoolPtr(true),
							BlockOwnerDeletion: pointer.BoolPtr(true),
						},
					},
					Annotations: nil,
				},
				Spec: appsv1.DaemonSetSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: selectorLabels,
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Annotations: nil,
							Labels: map[string]string{
								"app.kubernetes.io/name":       appName,
								"app.kubernetes.io/instance":   instance.GetName(),
								"app.kubernetes.io/managed-by": operatorName,
								"app.kubernetes.io/component":  agentServiceName,
							},
						},
						Spec: corev1.PodSpec{
							ServiceAccountName:            agentServiceName,
							ImagePullSecrets:              []corev1.LocalObjectReference{},
							NodeSelector:                  nil,
							Tolerations:                   nil,
							SecurityContext:               &corev1.PodSecurityContext{},
							TerminationGracePeriodSeconds: nil,
							InitContainers:                nil,
							Containers: []corev1.Container{
								{
									Name:            agentServiceName,
									Image:           "docker.io/fluxninja/aperture-agent:latest",
									ImagePullPolicy: corev1.PullIfNotPresent,
									SecurityContext: &corev1.SecurityContext{},
									Command:         nil,
									Args:            nil,
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
											Name: "APERTURE_AGENT_SERVICE_DISCOVERY_KUBERNETES_NODE_NAME",
											ValueFrom: &corev1.EnvVarSource{
												FieldRef: &corev1.ObjectFieldSelector{
													APIVersion: v1Version,
													FieldPath:  "spec.nodeName",
												},
											},
										},
										{
											Name:  "APERTURE_AGENT_SERVICE_DISCOVERY_KUBERNETES_DISCOVERY_ENABLED",
											Value: "true",
										},
									},
									EnvFrom:   []corev1.EnvFromSource{},
									Resources: corev1.ResourceRequirements{},
									Ports: []corev1.ContainerPort{
										{
											Name:          "grpc",
											ContainerPort: 80,
											Protocol:      corev1.ProtocolTCP,
										},
										{
											Name:          "grpc-otel",
											ContainerPort: 4317,
											Protocol:      corev1.ProtocolTCP,
										},
										{
											Name:          "dist-cache",
											ContainerPort: 3320,
											Protocol:      corev1.ProtocolTCP,
										},
										{
											Name:          "memberlist",
											ContainerPort: 3322,
											Protocol:      corev1.ProtocolTCP,
										},
									},
									TerminationMessagePath:   "/dev/termination-log",
									TerminationMessagePolicy: corev1.TerminationMessageReadFile,
									LivenessProbe:            nil,
									ReadinessProbe:           nil,
									Lifecycle:                nil,
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
					},
				},
			}

			result, _ := daemonsetForAgent(instance.DeepCopy(), logr.Logger{}, scheme.Scheme)
			Expect(result).To(Equal(expected))
		})
	})

	Context("Instance with all values", func() {
		It("returns correct daemonset for Agent", func() {
			selectorLabels := map[string]string{
				"app.kubernetes.io/name":       appName,
				"app.kubernetes.io/instance":   appName,
				"app.kubernetes.io/managed-by": operatorName,
				"app.kubernetes.io/component":  agentServiceName,
			}

			instance := &v1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.AgentSpec{
					DistributedCachePort: 3320,
					MemberListPort:       3322,
					CommonSpec: v1alpha1.CommonSpec{
						Labels:         testMap,
						Annotations:    testMap,
						ServerPort:     80,
						LivenessProbe:  probe,
						ReadinessProbe: probe,
						Resources:      resourceRequirement,
						PodSecurityContext: v1alpha1.PodSecurityContext{
							Enabled: true,
							FsGroup: pointer.Int64Ptr(1001),
						},
						ContainerSecurityContext: v1alpha1.ContainerSecurityContext{
							Enabled:                true,
							RunAsUser:              pointer.Int64Ptr(0),
							RunAsNonRootUser:       pointer.BoolPtr(false),
							ReadOnlyRootFilesystem: pointer.BoolPtr(false),
						},
						Command:                       testArray,
						Args:                          testArray,
						PodLabels:                     testMap,
						PodAnnotations:                testMap,
						NodeSelector:                  testMap,
						Tolerations:                   tolerations,
						TerminationGracePeriodSeconds: pointer.Int64Ptr(10),
						LifecycleHooks:                lifecycle,
						ExtraEnvVars: []corev1.EnvVar{
							{
								Name:  test,
								Value: test,
							},
						},
						ExtraEnvVarsCM:     test,
						ExtraEnvVarsSecret: test,
						ExtraVolumes: []corev1.Volume{
							{
								Name: test,
								VolumeSource: corev1.VolumeSource{
									EmptyDir: &corev1.EmptyDirVolumeSource{},
								},
							},
						},
						ExtraVolumeMounts: []corev1.VolumeMount{
							{
								Name:      test,
								MountPath: test,
							},
						},
						Sidecars: []corev1.Container{
							{
								Name: test,
							},
						},
						InitContainers: []corev1.Container{
							{
								Name: test,
							},
						},
						ServiceAccountSpec: v1alpha1.ServiceAccountSpec{
							Create: true,
						},
						Affinity: affinity,
					},
					Image: v1alpha1.Image{
						Registry:    "docker.io/fluxninja",
						Repository:  "aperture-agent",
						Tag:         "latest",
						PullPolicy:  "IfNotPresent",
						PullSecrets: testArray,
					},
				},
			}
			expected := &appsv1.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      agentServiceName,
					Namespace: appName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       appName,
						"app.kubernetes.io/instance":   instance.GetName(),
						"app.kubernetes.io/managed-by": operatorName,
						"app.kubernetes.io/component":  agentServiceName,
						test:                           test,
					},
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion:         "fluxninja.com/v1alpha1",
							Name:               instance.GetName(),
							Kind:               "Agent",
							Controller:         pointer.BoolPtr(true),
							BlockOwnerDeletion: pointer.BoolPtr(true),
						},
					},
					Annotations: testMap,
				},
				Spec: appsv1.DaemonSetSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: selectorLabels,
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Annotations: map[string]string{
								test: test,
							},
							Labels: map[string]string{
								"app.kubernetes.io/name":       appName,
								"app.kubernetes.io/instance":   instance.GetName(),
								"app.kubernetes.io/managed-by": operatorName,
								"app.kubernetes.io/component":  agentServiceName,
								test:                           test,
							},
						},
						Spec: corev1.PodSpec{
							ServiceAccountName: agentServiceName,
							ImagePullSecrets: []corev1.LocalObjectReference{
								{
									Name: test,
								},
							},
							NodeSelector: testMap,
							Affinity:     affinity,
							Tolerations:  tolerations,
							SecurityContext: &corev1.PodSecurityContext{
								FSGroup: pointer.Int64Ptr(1001),
							},
							TerminationGracePeriodSeconds: pointer.Int64Ptr(10),
							InitContainers: []corev1.Container{
								{
									Name: test,
								},
							},
							Containers: []corev1.Container{
								{
									Name:            agentServiceName,
									Image:           "docker.io/fluxninja/aperture-agent:latest",
									ImagePullPolicy: corev1.PullIfNotPresent,
									SecurityContext: &corev1.SecurityContext{
										RunAsUser:              pointer.Int64Ptr(0),
										RunAsNonRoot:           pointer.BoolPtr(false),
										ReadOnlyRootFilesystem: pointer.BoolPtr(false),
									},
									Command: testArray,
									Args:    testArray,
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
											Name: "APERTURE_AGENT_SERVICE_DISCOVERY_KUBERNETES_NODE_NAME",
											ValueFrom: &corev1.EnvVarSource{
												FieldRef: &corev1.ObjectFieldSelector{
													APIVersion: v1Version,
													FieldPath:  "spec.nodeName",
												},
											},
										},
										{
											Name:  "APERTURE_AGENT_SERVICE_DISCOVERY_KUBERNETES_DISCOVERY_ENABLED",
											Value: "true",
										},
									},
									EnvFrom: []corev1.EnvFromSource{
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
									Resources: resourceRequirement,
									Ports: []corev1.ContainerPort{
										{
											Name:          "grpc",
											ContainerPort: 80,
											Protocol:      corev1.ProtocolTCP,
										},
										{
											Name:          "grpc-otel",
											ContainerPort: 4317,
											Protocol:      corev1.ProtocolTCP,
										},
										{
											Name:          "dist-cache",
											ContainerPort: 3320,
											Protocol:      corev1.ProtocolTCP,
										},
										{
											Name:          "memberlist",
											ContainerPort: 3322,
											Protocol:      corev1.ProtocolTCP,
										},
									},
									TerminationMessagePath:   "/dev/termination-log",
									TerminationMessagePolicy: corev1.TerminationMessageReadFile,
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
								},
								{
									Name: test,
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
					},
				},
			}

			result, _ := daemonsetForAgent(instance.DeepCopy(), logr.Logger{}, scheme.Scheme)
			Expect(result).To(Equal(expected))
		})
	})
})

var _ = Describe("Test Daemonset Mutate", func() {
	It("Mutate should update required fields only", func() {
		expected := &appsv1.DaemonSet{
			ObjectMeta: metav1.ObjectMeta{},
			Spec: appsv1.DaemonSetSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels:      map[string]string{},
						Annotations: map[string]string{},
					},
					Spec: corev1.PodSpec{
						ServiceAccountName: test,
						ImagePullSecrets:   []corev1.LocalObjectReference{},
						NodeSelector:       map[string]string{},
						Tolerations:        []corev1.Toleration{},
						SecurityContext: &corev1.PodSecurityContext{
							FSGroup: pointer.Int64Ptr(1001),
						},
						InitContainers: []corev1.Container{},
						Containers: []corev1.Container{
							{
								Name: test,
							},
						},
						Volumes: []corev1.Volume{},
					},
				},
			},
		}

		dms := &appsv1.DaemonSet{}
		err := daemonsetMutate(dms, expected.Spec)()
		Expect(err).NotTo(HaveOccurred())
		Expect(dms).To(Equal(expected))
	})
})
