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

package agent

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

	"github.com/fluxninja/aperture/cmd/aperture-agent/agent"
	agentv1alpha1 "github.com/fluxninja/aperture/operator/api/agent/v1alpha1"
	"github.com/fluxninja/aperture/operator/api/common"
	. "github.com/fluxninja/aperture/operator/controllers"
	"github.com/fluxninja/aperture/pkg/distcache"
	"github.com/fluxninja/aperture/pkg/net/listener"
	otelconfig "github.com/fluxninja/aperture/pkg/otelcollector/config"
)

var _ = Describe("Agent DaemonSet", func() {
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
		tolerations = []corev1.Toleration{
			{
				Key:      Test,
				Value:    Test,
				Operator: corev1.TolerationOpEqual,
			},
		}
		lifecycle = &corev1.Lifecycle{
			PreStop: &corev1.LifecycleHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: Test,
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
								MatchLabels: TestMap,
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
				"app.kubernetes.io/name":       AppName,
				"app.kubernetes.io/instance":   AppName,
				"app.kubernetes.io/managed-by": OperatorName,
				"app.kubernetes.io/component":  AgentServiceName,
			}

			instance := &agentv1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: agentv1alpha1.AgentSpec{
					ConfigSpec: agentv1alpha1.AgentConfigSpec{
						CommonConfigSpec: common.CommonConfigSpec{
							Server: common.ServerConfigSpec{
								ListenerConfig: listener.ListenerConfig{
									Addr: ":80",
								},
							},
						},
						DistCache: distcache.DistCacheConfig{
							BindAddr:           ":3320",
							MemberlistBindAddr: ":3322",
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
			expected := &appsv1.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AgentServiceName,
					Namespace: AppName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       AppName,
						"app.kubernetes.io/instance":   instance.GetName(),
						"app.kubernetes.io/managed-by": OperatorName,
						"app.kubernetes.io/component":  AgentServiceName,
					},
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion:         "fluxninja.com/v1alpha1",
							Name:               instance.GetName(),
							Kind:               "Agent",
							Controller:         pointer.Bool(true),
							BlockOwnerDeletion: pointer.Bool(true),
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
								"app.kubernetes.io/name":       AppName,
								"app.kubernetes.io/instance":   instance.GetName(),
								"app.kubernetes.io/managed-by": OperatorName,
								"app.kubernetes.io/component":  AgentServiceName,
							},
						},
						Spec: corev1.PodSpec{
							ServiceAccountName:            AgentServiceName,
							ImagePullSecrets:              []corev1.LocalObjectReference{},
							NodeSelector:                  nil,
							Tolerations:                   nil,
							SecurityContext:               &corev1.PodSecurityContext{},
							TerminationGracePeriodSeconds: pointer.Int64(0),
							InitContainers:                nil,
							Containers: []corev1.Container{
								{
									Name:            AgentServiceName,
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
													APIVersion: V1Version,
													FieldPath:  "spec.nodeName",
												},
											},
										},
										{
											Name: "APERTURE_AGENT_SERVICE_DISCOVERY_KUBERNETES_NODE_NAME",
											ValueFrom: &corev1.EnvVarSource{
												FieldRef: &corev1.ObjectFieldSelector{
													APIVersion: V1Version,
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
											DefaultMode: pointer.Int32(420),
											LocalObjectReference: corev1.LocalObjectReference{
												Name: AgentServiceName,
											},
										},
									},
								},
							},
						},
					},
				},
			}

			result, err := daemonsetForAgent(instance.DeepCopy(), logr.Logger{}, scheme.Scheme)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(expected))
		})
	})

	Context("Instance with all values", func() {
		It("returns correct daemonset for Agent", func() {
			selectorLabels := map[string]string{
				"app.kubernetes.io/name":       AppName,
				"app.kubernetes.io/instance":   AppName,
				"app.kubernetes.io/managed-by": OperatorName,
				"app.kubernetes.io/component":  AgentServiceName,
			}

			instance := &agentv1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: agentv1alpha1.AgentSpec{
					ConfigSpec: agentv1alpha1.AgentConfigSpec{
						CommonConfigSpec: common.CommonConfigSpec{
							Server: common.ServerConfigSpec{
								ListenerConfig: listener.ListenerConfig{
									Addr: ":80",
								},
							},
						},
						DistCache: distcache.DistCacheConfig{
							BindAddr:           ":3320",
							MemberlistBindAddr: ":3322",
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
						Labels:         TestMap,
						Annotations:    TestMap,
						LivenessProbe:  probe,
						ReadinessProbe: probe,
						Resources:      resourceRequirement,
						PodSecurityContext: common.PodSecurityContext{
							Enabled: true,
							FsGroup: 1001,
						},
						ContainerSecurityContext: common.ContainerSecurityContext{
							Enabled:                true,
							RunAsUser:              0,
							RunAsNonRootUser:       false,
							ReadOnlyRootFilesystem: false,
						},
						Command:                       TestArray,
						Args:                          TestArray,
						PodLabels:                     TestMap,
						PodAnnotations:                TestMap,
						NodeSelector:                  TestMap,
						Tolerations:                   tolerations,
						TerminationGracePeriodSeconds: 10,
						LifecycleHooks:                lifecycle,
						ExtraEnvVars: []corev1.EnvVar{
							{
								Name:  Test,
								Value: Test,
							},
						},
						ExtraEnvVarsCM:     Test,
						ExtraEnvVarsSecret: Test,
						ExtraVolumes: []corev1.Volume{
							{
								Name: Test,
								VolumeSource: corev1.VolumeSource{
									EmptyDir: &corev1.EmptyDirVolumeSource{},
								},
							},
						},
						ExtraVolumeMounts: []corev1.VolumeMount{
							{
								Name:      Test,
								MountPath: Test,
							},
						},
						Sidecars: []corev1.Container{
							{
								Name: Test,
							},
						},
						InitContainers: []corev1.Container{
							{
								Name: Test,
							},
						},
						ServiceAccountSpec: common.ServiceAccountSpec{
							Create: true,
						},
						Affinity: affinity,
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
			expected := &appsv1.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AgentServiceName,
					Namespace: AppName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       AppName,
						"app.kubernetes.io/instance":   instance.GetName(),
						"app.kubernetes.io/managed-by": OperatorName,
						"app.kubernetes.io/component":  AgentServiceName,
						Test:                           Test,
					},
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion:         "fluxninja.com/v1alpha1",
							Name:               instance.GetName(),
							Kind:               "Agent",
							Controller:         pointer.Bool(true),
							BlockOwnerDeletion: pointer.Bool(true),
						},
					},
					Annotations: TestMap,
				},
				Spec: appsv1.DaemonSetSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: selectorLabels,
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Annotations: map[string]string{
								Test: Test,
							},
							Labels: map[string]string{
								"app.kubernetes.io/name":       AppName,
								"app.kubernetes.io/instance":   instance.GetName(),
								"app.kubernetes.io/managed-by": OperatorName,
								"app.kubernetes.io/component":  AgentServiceName,
								Test:                           Test,
							},
						},
						Spec: corev1.PodSpec{
							ServiceAccountName: AgentServiceName,
							ImagePullSecrets: []corev1.LocalObjectReference{
								{
									Name: Test,
								},
							},
							NodeSelector: TestMap,
							Affinity:     affinity,
							Tolerations:  tolerations,
							SecurityContext: &corev1.PodSecurityContext{
								FSGroup: pointer.Int64(1001),
							},
							TerminationGracePeriodSeconds: pointer.Int64(10),
							InitContainers: []corev1.Container{
								{
									Name: Test,
								},
							},
							Containers: []corev1.Container{
								{
									Name:            AgentServiceName,
									Image:           "docker.io/fluxninja/aperture-agent:latest",
									ImagePullPolicy: corev1.PullIfNotPresent,
									SecurityContext: &corev1.SecurityContext{
										RunAsUser:              pointer.Int64(0),
										RunAsNonRoot:           pointer.Bool(false),
										ReadOnlyRootFilesystem: pointer.Bool(false),
									},
									Command: TestArray,
									Args:    TestArray,
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
											Name: "APERTURE_AGENT_SERVICE_DISCOVERY_KUBERNETES_NODE_NAME",
											ValueFrom: &corev1.EnvVarSource{
												FieldRef: &corev1.ObjectFieldSelector{
													APIVersion: V1Version,
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
									Resources: resourceRequirement,
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
									TerminationMessagePath:   "/dev/termination-log",
									TerminationMessagePolicy: corev1.TerminationMessageReadFile,
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
								},
								{
									Name: Test,
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
					},
				},
			}

			result, err := daemonsetForAgent(instance.DeepCopy(), logr.Logger{}, scheme.Scheme)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(expected))
		})
	})
})

var _ = Describe("Test DaemonSet Mutate", func() {
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
						ServiceAccountName: Test,
						ImagePullSecrets:   []corev1.LocalObjectReference{},
						NodeSelector:       map[string]string{},
						Tolerations:        []corev1.Toleration{},
						SecurityContext: &corev1.PodSecurityContext{
							FSGroup: pointer.Int64(1001),
						},
						InitContainers: []corev1.Container{},
						Containers: []corev1.Container{
							{
								Name: Test,
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
