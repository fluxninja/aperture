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

package controller

import (
	"fmt"

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

	controller "github.com/fluxninja/aperture/v2/cmd/aperture-controller/config"
	. "github.com/fluxninja/aperture/v2/operator/controllers"

	"github.com/fluxninja/aperture/v2/operator/api/common"
	controllerv1alpha1 "github.com/fluxninja/aperture/v2/operator/api/controller/v1alpha1"
	"github.com/fluxninja/aperture/v2/pkg/net/listener"
	otelconfig "github.com/fluxninja/aperture/v2/pkg/otelcollector/config"
)

var _ = Describe("Controller Deployment", func() {
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
		hostAliases = []corev1.HostAlias{
			{
				Hostnames: TestArray,
			},
		}
	)

	Context("Instance with default parameters", func() {
		It("returns correct deployment for Controller", func() {
			selectorLabels := map[string]string{
				"app.kubernetes.io/name":       AppName,
				"app.kubernetes.io/instance":   ControllerName,
				"app.kubernetes.io/managed-by": OperatorName,
				"app.kubernetes.io/component":  ControllerServiceName,
			}

			instance := &controllerv1alpha1.Controller{
				ObjectMeta: metav1.ObjectMeta{
					Name:      ControllerName,
					Namespace: AppName,
				},
				Spec: controllerv1alpha1.ControllerSpec{
					ConfigSpec: controllerv1alpha1.ControllerConfigSpec{
						CommonConfigSpec: common.CommonConfigSpec{
							Server: common.ServerConfigSpec{
								Listener: listener.ListenerConfig{
									Addr: ":80",
								},
							},
						},
						OTel: controller.ControllerOTelConfig{
							CommonOTelConfig: otelconfig.CommonOTelConfig{
								Ports: otelconfig.PortsConfig{
									DebugPort:       8888,
									HealthCheckPort: 13133,
									PprofPort:       1777,
									ZpagesPort:      55679,
								},
							},
						},
					},
					Image: common.ControllerImage{
						Image: common.Image{
							Registry:   "docker.io/fluxninja",
							Tag:        "latest",
							PullPolicy: "IfNotPresent",
						},
						Repository: "aperture-controller",
					},
				},
			}
			expected := &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      ControllerServiceName,
					Namespace: AppName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       AppName,
						"app.kubernetes.io/instance":   instance.GetName(),
						"app.kubernetes.io/managed-by": OperatorName,
						"app.kubernetes.io/component":  ControllerServiceName,
					},
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion:         "fluxninja.com/v1alpha1",
							Name:               instance.GetName(),
							Kind:               "Controller",
							Controller:         pointer.Bool(true),
							BlockOwnerDeletion: pointer.Bool(true),
						},
					},
					Annotations: nil,
				},
				Spec: appsv1.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: selectorLabels,
					},
					Strategy: appsv1.DeploymentStrategy{
						Type: "Recreate",
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Annotations: map[string]string{
								SidecarAnnotationKey: "false",
							},
							Labels: map[string]string{
								"app.kubernetes.io/name":       AppName,
								"app.kubernetes.io/instance":   instance.GetName(),
								"app.kubernetes.io/managed-by": OperatorName,
								"app.kubernetes.io/component":  ControllerServiceName,
							},
						},
						Spec: corev1.PodSpec{
							ServiceAccountName:            ControllerServiceName,
							HostAliases:                   nil,
							ImagePullSecrets:              []corev1.LocalObjectReference{},
							NodeSelector:                  nil,
							Tolerations:                   nil,
							SecurityContext:               &corev1.PodSecurityContext{},
							TerminationGracePeriodSeconds: nil,
							InitContainers:                nil,
							Containers: []corev1.Container{
								{
									Name:            ControllerServiceName,
									Image:           "docker.io/fluxninja/aperture-controller:latest",
									ImagePullPolicy: corev1.PullIfNotPresent,
									SecurityContext: &corev1.SecurityContext{},
									Command:         nil,
									Args:            nil,
									Env: []corev1.EnvVar{
										{
											Name: "APERTURE_CONTROLLER_SERVICE_DISCOVERY_KUBERNETES_NODE_NAME",
											ValueFrom: &corev1.EnvVarSource{
												FieldRef: &corev1.ObjectFieldSelector{
													APIVersion: V1Version,
													FieldPath:  "spec.nodeName",
												},
											},
										},
										{
											Name:  "APERTURE_CONTROLLER_NAMESPACE",
											Value: AppName,
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
											Name:      "aperture-controller-config",
											MountPath: "/etc/aperture/aperture-controller/config",
										},
										{
											Name:      "server-cert",
											MountPath: "/etc/aperture/aperture-controller/certs",
											ReadOnly:  true,
										},
									},
								},
							},
							Volumes: []corev1.Volume{
								{
									Name: "aperture-controller-config",
									VolumeSource: corev1.VolumeSource{
										ConfigMap: &corev1.ConfigMapVolumeSource{
											DefaultMode: pointer.Int32(420),
											LocalObjectReference: corev1.LocalObjectReference{
												Name: ControllerServiceName,
											},
										},
									},
								},
								{
									Name: "server-cert",
									VolumeSource: corev1.VolumeSource{
										Secret: &corev1.SecretVolumeSource{
											DefaultMode: pointer.Int32(420),
											SecretName:  fmt.Sprintf("%s-controller-cert", instance.GetName()),
										},
									},
								},
							},
						},
					},
				},
			}

			result, err := deploymentForController(instance.DeepCopy(), logr.Logger{}, scheme.Scheme)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Spec.Template.Spec.Containers).To(Equal(expected.Spec.Template.Spec.Containers))
		})
	})

	Context("Instance with all values", func() {
		It("returns correct deployment for Controller", func() {
			selectorLabels := map[string]string{
				"app.kubernetes.io/name":       AppName,
				"app.kubernetes.io/instance":   ControllerName,
				"app.kubernetes.io/managed-by": OperatorName,
				"app.kubernetes.io/component":  ControllerServiceName,
			}

			instance := &controllerv1alpha1.Controller{
				ObjectMeta: metav1.ObjectMeta{
					Name:      ControllerName,
					Namespace: AppName,
				},
				Spec: controllerv1alpha1.ControllerSpec{
					ConfigSpec: controllerv1alpha1.ControllerConfigSpec{
						CommonConfigSpec: common.CommonConfigSpec{
							Server: common.ServerConfigSpec{
								Listener: listener.ListenerConfig{
									Addr: ":80",
								},
							},
						},
						OTel: controller.ControllerOTelConfig{
							CommonOTelConfig: otelconfig.CommonOTelConfig{
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
							RunAsGroup:             0,
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
					Image: common.ControllerImage{
						Image: common.Image{
							Registry:    "docker.io/fluxninja",
							Tag:         "latest",
							PullPolicy:  "IfNotPresent",
							PullSecrets: TestArray,
						},
						Repository: "aperture-controller",
					},
					HostAliases: hostAliases,
				},
			}
			expected := &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      ControllerServiceName,
					Namespace: AppName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       AppName,
						"app.kubernetes.io/instance":   instance.GetName(),
						"app.kubernetes.io/managed-by": OperatorName,
						"app.kubernetes.io/component":  ControllerServiceName,
						Test:                           Test,
					},
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion:         "fluxninja.com/v1alpha1",
							Name:               instance.GetName(),
							Kind:               "Controller",
							Controller:         pointer.Bool(true),
							BlockOwnerDeletion: pointer.Bool(true),
						},
					},
					Annotations: TestMap,
				},
				Spec: appsv1.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: selectorLabels,
					},
					Strategy: appsv1.DeploymentStrategy{
						Type: "Recreate",
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Annotations: map[string]string{
								Test:                 Test,
								SidecarAnnotationKey: "false",
							},
							Labels: map[string]string{
								"app.kubernetes.io/name":       AppName,
								"app.kubernetes.io/instance":   instance.GetName(),
								"app.kubernetes.io/managed-by": OperatorName,
								"app.kubernetes.io/component":  ControllerServiceName,
								Test:                           Test,
							},
						},
						Spec: corev1.PodSpec{
							ServiceAccountName: ControllerServiceName,
							HostAliases:        hostAliases,
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
									Name:            ControllerServiceName,
									Image:           "docker.io/fluxninja/aperture-controller:latest",
									ImagePullPolicy: corev1.PullIfNotPresent,
									SecurityContext: &corev1.SecurityContext{
										RunAsUser:              pointer.Int64(0),
										RunAsGroup:             pointer.Int64(0),
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
											Name: "APERTURE_CONTROLLER_SERVICE_DISCOVERY_KUBERNETES_NODE_NAME",
											ValueFrom: &corev1.EnvVarSource{
												FieldRef: &corev1.ObjectFieldSelector{
													APIVersion: V1Version,
													FieldPath:  "spec.nodeName",
												},
											},
										},
										{
											Name:  "APERTURE_CONTROLLER_NAMESPACE",
											Value: AppName,
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
												Path:   "/v1/status/system/liveness",
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
												Path:   "/v1/status/system/readiness",
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
											Name:      "aperture-controller-config",
											MountPath: "/etc/aperture/aperture-controller/config",
										},
										{
											Name:      "server-cert",
											MountPath: "/etc/aperture/aperture-controller/certs",
											ReadOnly:  true,
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
									Name: "aperture-controller-config",
									VolumeSource: corev1.VolumeSource{
										ConfigMap: &corev1.ConfigMapVolumeSource{
											DefaultMode: pointer.Int32(420),
											LocalObjectReference: corev1.LocalObjectReference{
												Name: ControllerServiceName,
											},
										},
									},
								},
								{
									Name: "server-cert",
									VolumeSource: corev1.VolumeSource{
										Secret: &corev1.SecretVolumeSource{
											DefaultMode: pointer.Int32(420),
											SecretName:  fmt.Sprintf("%s-controller-cert", instance.GetName()),
										},
									},
								},
							},
						},
					},
				},
			}

			result, err := deploymentForController(instance.DeepCopy(), logr.Logger{}, scheme.Scheme)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Spec.Template.Spec.ImagePullSecrets).To(Equal(expected.Spec.Template.Spec.ImagePullSecrets))
		})
	})
})

var _ = Describe("Test Deployment Mutate", func() {
	It("Mutate should update required fields only", func() {
		expected := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{},
			Spec: appsv1.DeploymentSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels:      map[string]string{},
						Annotations: map[string]string{},
					},
					Spec: corev1.PodSpec{
						ServiceAccountName:        Test,
						ImagePullSecrets:          []corev1.LocalObjectReference{},
						HostAliases:               []corev1.HostAlias{},
						Affinity:                  &corev1.Affinity{},
						NodeSelector:              map[string]string{},
						Tolerations:               []corev1.Toleration{},
						PriorityClassName:         Test,
						TopologySpreadConstraints: []corev1.TopologySpreadConstraint{},
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

		dep := &appsv1.Deployment{}
		err := deploymentMutate(dep, expected.Spec)()

		Expect(err).NotTo(HaveOccurred())
		Expect(dep).To(Equal(expected))
	})
})
