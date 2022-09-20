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

	"github.com/fluxninja/aperture/operator/api/v1alpha1"
	"github.com/fluxninja/aperture/pkg/net/listener"
	"github.com/fluxninja/aperture/pkg/otel"
)

var _ = Describe("Controller Deployment", func() {
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
		hostAliases = []corev1.HostAlias{
			{
				Hostnames: testArray,
			},
		}
	)

	Context("Instance with default parameters", func() {
		It("returns correct deployment for Controller", func() {
			selectorLabels := map[string]string{
				"app.kubernetes.io/name":       appName,
				"app.kubernetes.io/instance":   appName,
				"app.kubernetes.io/managed-by": operatorName,
				"app.kubernetes.io/component":  controllerServiceName,
			}

			instance := &v1alpha1.Controller{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.ControllerSpec{
					ConfigSpec: v1alpha1.ControllerConfigSpec{
						CommonConfigSpec: v1alpha1.CommonConfigSpec{
							Server: v1alpha1.ServerConfigSpec{
								ListenerConfig: listener.ListenerConfig{
									Addr: ":80",
								},
							},
							Otel: otel.OtelConfig{
								GRPCAddr: ":4317",
								HTTPAddr: ":4318",
							},
						},
					},
					OperatorImage: v1alpha1.OperatorImage{
						Image: v1alpha1.Image{
							Registry:   "docker.io/fluxninja",
							Tag:        "latest",
							PullPolicy: "IfNotPresent",
						},
						Repository: "aperture-operator",
					},
					Image: v1alpha1.ControllerImage{
						Image: v1alpha1.Image{
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
					Name:      controllerServiceName,
					Namespace: appName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       appName,
						"app.kubernetes.io/instance":   instance.GetName(),
						"app.kubernetes.io/managed-by": operatorName,
						"app.kubernetes.io/component":  controllerServiceName,
					},
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion:         "fluxninja.com/v1alpha1",
							Name:               instance.GetName(),
							Kind:               "Controller",
							Controller:         pointer.BoolPtr(true),
							BlockOwnerDeletion: pointer.BoolPtr(true),
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
								sidecarAnnotationKey: "false",
							},
							Labels: map[string]string{
								"app.kubernetes.io/name":       appName,
								"app.kubernetes.io/instance":   instance.GetName(),
								"app.kubernetes.io/managed-by": operatorName,
								"app.kubernetes.io/component":  controllerServiceName,
							},
						},
						Spec: corev1.PodSpec{
							ServiceAccountName:            controllerServiceName,
							HostAliases:                   nil,
							ImagePullSecrets:              []corev1.LocalObjectReference{},
							NodeSelector:                  nil,
							Tolerations:                   nil,
							SecurityContext:               &corev1.PodSecurityContext{},
							TerminationGracePeriodSeconds: nil,
							InitContainers:                nil,
							Containers: []corev1.Container{
								{
									Name:            "policy-watcher",
									Image:           "docker.io/fluxninja/aperture-operator:latest",
									ImagePullPolicy: corev1.PullIfNotPresent,
									SecurityContext: &corev1.SecurityContext{},
									Command: []string{
										"/aperture-operator",
										"--policy",
										"--health-probe-bind-address=:9091",
										"--metrics-bind-address=127.0.0.1:9090",
									},
									Args: []string{
										"--leader-elect=True",
									},
									Env: []corev1.EnvVar{
										{
											Name:  "APERTURE_CONTROLLER_NAMESPACE",
											Value: appName,
										},
									},
									TerminationMessagePath:   "/dev/termination-log",
									TerminationMessagePolicy: corev1.TerminationMessageReadFile,
									LivenessProbe: &corev1.Probe{
										ProbeHandler: corev1.ProbeHandler{
											HTTPGet: &corev1.HTTPGetAction{
												Path: "/healthz",
												Port: intstr.FromInt(9091),
											},
										},
										FailureThreshold:    3,
										InitialDelaySeconds: 10,
										PeriodSeconds:       10,
										SuccessThreshold:    1,
										TimeoutSeconds:      1,
									},
									ReadinessProbe: &corev1.Probe{
										ProbeHandler: corev1.ProbeHandler{
											HTTPGet: &corev1.HTTPGetAction{
												Path: "/readyz",
												Port: intstr.FromInt(9091),
											},
										},
										FailureThreshold:    3,
										InitialDelaySeconds: 10,
										PeriodSeconds:       10,
										SuccessThreshold:    1,
										TimeoutSeconds:      1,
									},
									VolumeMounts: []corev1.VolumeMount{
										{
											Name:      "etc-aperture-policies",
											MountPath: policyFilePath,
											ReadOnly:  false,
										},
									},
								},
								{
									Name:            controllerServiceName,
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
													APIVersion: v1Version,
													FieldPath:  "spec.nodeName",
												},
											},
										},
									},
									EnvFrom:   []corev1.EnvFromSource{},
									Resources: corev1.ResourceRequirements{},
									Ports: []corev1.ContainerPort{
										{
											Name:          server,
											ContainerPort: 80,
											Protocol:      corev1.ProtocolTCP,
										},
										{
											Name:          grpcOtel,
											ContainerPort: 4317,
											Protocol:      corev1.ProtocolTCP,
										},
										{
											Name:          httpOtel,
											ContainerPort: 4318,
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
											Name:      "etc-aperture-policies",
											MountPath: policyFilePath,
											ReadOnly:  true,
										},
										{
											Name:      "etc-aperture-classification",
											MountPath: "/etc/aperture/aperture-controller/classifiers",
											ReadOnly:  true,
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
											DefaultMode: pointer.Int32Ptr(420),
											LocalObjectReference: corev1.LocalObjectReference{
												Name: controllerServiceName,
											},
										},
									},
								},
								{
									Name: "etc-aperture-policies",
									VolumeSource: corev1.VolumeSource{
										EmptyDir: &corev1.EmptyDirVolumeSource{},
									},
								},
								{
									Name: "etc-aperture-classification",
									VolumeSource: corev1.VolumeSource{
										ConfigMap: &corev1.ConfigMapVolumeSource{
											DefaultMode: pointer.Int32Ptr(420),
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "classification",
											},
											Optional: pointer.BoolPtr(true),
										},
									},
								},
								{
									Name: "server-cert",
									VolumeSource: corev1.VolumeSource{
										Secret: &corev1.SecretVolumeSource{
											DefaultMode: pointer.Int32Ptr(420),
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
				"app.kubernetes.io/name":       appName,
				"app.kubernetes.io/instance":   appName,
				"app.kubernetes.io/managed-by": operatorName,
				"app.kubernetes.io/component":  controllerServiceName,
			}

			instance := &v1alpha1.Controller{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.ControllerSpec{
					ConfigSpec: v1alpha1.ControllerConfigSpec{
						CommonConfigSpec: v1alpha1.CommonConfigSpec{
							Server: v1alpha1.ServerConfigSpec{
								ListenerConfig: listener.ListenerConfig{
									Addr: ":80",
								},
							},
							Otel: otel.OtelConfig{
								GRPCAddr: ":4317",
								HTTPAddr: ":4318",
							},
						},
					},
					CommonSpec: v1alpha1.CommonSpec{
						Labels:         testMap,
						Annotations:    testMap,
						LivenessProbe:  probe,
						ReadinessProbe: probe,
						Resources:      resourceRequirement,
						PodSecurityContext: v1alpha1.PodSecurityContext{
							Enabled: true,
							FsGroup: 1001,
						},
						ContainerSecurityContext: v1alpha1.ContainerSecurityContext{
							Enabled:                true,
							RunAsUser:              0,
							RunAsNonRootUser:       false,
							ReadOnlyRootFilesystem: false,
						},
						Command:                       testArray,
						Args:                          testArray,
						PodLabels:                     testMap,
						PodAnnotations:                testMap,
						NodeSelector:                  testMap,
						Tolerations:                   tolerations,
						TerminationGracePeriodSeconds: 10,
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
					Image: v1alpha1.ControllerImage{
						Image: v1alpha1.Image{
							Registry:    "docker.io/fluxninja",
							Tag:         "latest",
							PullPolicy:  "IfNotPresent",
							PullSecrets: testArray,
						},
						Repository: "aperture-controller",
					},
					OperatorImage: v1alpha1.OperatorImage{
						Image: v1alpha1.Image{
							Registry:    "docker.io/fluxninja",
							Tag:         "latest",
							PullPolicy:  "IfNotPresent",
							PullSecrets: testArrayTwo,
						},
						Repository: "aperture-operator",
					},
					HostAliases: hostAliases,
				},
			}
			expected := &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      controllerServiceName,
					Namespace: appName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       appName,
						"app.kubernetes.io/instance":   instance.GetName(),
						"app.kubernetes.io/managed-by": operatorName,
						"app.kubernetes.io/component":  controllerServiceName,
						test:                           test,
					},
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion:         "fluxninja.com/v1alpha1",
							Name:               instance.GetName(),
							Kind:               "Controller",
							Controller:         pointer.BoolPtr(true),
							BlockOwnerDeletion: pointer.BoolPtr(true),
						},
					},
					Annotations: testMap,
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
								test:                 test,
								sidecarAnnotationKey: "false",
							},
							Labels: map[string]string{
								"app.kubernetes.io/name":       appName,
								"app.kubernetes.io/instance":   instance.GetName(),
								"app.kubernetes.io/managed-by": operatorName,
								"app.kubernetes.io/component":  controllerServiceName,
								test:                           test,
							},
						},
						Spec: corev1.PodSpec{
							ServiceAccountName: controllerServiceName,
							HostAliases:        hostAliases,
							ImagePullSecrets: []corev1.LocalObjectReference{
								{
									Name: testTwo,
								},
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
									Name:            "policy-watcher",
									Image:           "docker.io/fluxninja/aperture-operator:latest",
									ImagePullPolicy: corev1.PullIfNotPresent,
									SecurityContext: &corev1.SecurityContext{
										RunAsUser:              pointer.Int64Ptr(0),
										RunAsNonRoot:           pointer.BoolPtr(false),
										ReadOnlyRootFilesystem: pointer.BoolPtr(false),
									},
									Command: []string{
										"/aperture-operator",
										"--policy",
										"--health-probe-bind-address=:9091",
										"--metrics-bind-address=127.0.0.1:9090",
									},
									Args: []string{
										"--leader-elect=True",
									},
									Env: []corev1.EnvVar{
										{
											Name:  "APERTURE_CONTROLLER_NAMESPACE",
											Value: appName,
										},
									},
									TerminationMessagePath:   "/dev/termination-log",
									TerminationMessagePolicy: corev1.TerminationMessageReadFile,
									LivenessProbe: &corev1.Probe{
										ProbeHandler: corev1.ProbeHandler{
											HTTPGet: &corev1.HTTPGetAction{
												Path: "/healthz",
												Port: intstr.FromInt(9091),
											},
										},
										FailureThreshold:    3,
										InitialDelaySeconds: 10,
										PeriodSeconds:       10,
										SuccessThreshold:    1,
										TimeoutSeconds:      1,
									},
									ReadinessProbe: &corev1.Probe{
										ProbeHandler: corev1.ProbeHandler{
											HTTPGet: &corev1.HTTPGetAction{
												Path: "/readyz",
												Port: intstr.FromInt(9091),
											},
										},
										FailureThreshold:    3,
										InitialDelaySeconds: 10,
										PeriodSeconds:       10,
										SuccessThreshold:    1,
										TimeoutSeconds:      1,
									},
									VolumeMounts: []corev1.VolumeMount{
										{
											Name:      "etc-aperture-policies",
											MountPath: policyFilePath,
											ReadOnly:  false,
										},
									},
								},
								{
									Name:            controllerServiceName,
									Image:           "docker.io/fluxninja/aperture-controller:latest",
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
											Name: "APERTURE_CONTROLLER_SERVICE_DISCOVERY_KUBERNETES_NODE_NAME",
											ValueFrom: &corev1.EnvVarSource{
												FieldRef: &corev1.ObjectFieldSelector{
													APIVersion: v1Version,
													FieldPath:  "spec.nodeName",
												},
											},
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
											Name:          server,
											ContainerPort: 80,
											Protocol:      corev1.ProtocolTCP,
										},
										{
											Name:          grpcOtel,
											ContainerPort: 4317,
											Protocol:      corev1.ProtocolTCP,
										},
										{
											Name:          httpOtel,
											ContainerPort: 4318,
											Protocol:      corev1.ProtocolTCP,
										},
									},
									TerminationMessagePath:   "/dev/termination-log",
									TerminationMessagePolicy: corev1.TerminationMessageReadFile,
									LivenessProbe: &corev1.Probe{
										ProbeHandler: corev1.ProbeHandler{
											HTTPGet: &corev1.HTTPGetAction{
												Path:   "/v1/status/liveness",
												Port:   intstr.FromString(server),
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
												Port:   intstr.FromString(server),
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
											Name:      "aperture-controller-config",
											MountPath: "/etc/aperture/aperture-controller/config",
										},
										{
											Name:      "etc-aperture-policies",
											MountPath: policyFilePath,
											ReadOnly:  true,
										},
										{
											Name:      "etc-aperture-classification",
											MountPath: "/etc/aperture/aperture-controller/classifiers",
											ReadOnly:  true,
										},
										{
											Name:      "server-cert",
											MountPath: "/etc/aperture/aperture-controller/certs",
											ReadOnly:  true,
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
									Name: "aperture-controller-config",
									VolumeSource: corev1.VolumeSource{
										ConfigMap: &corev1.ConfigMapVolumeSource{
											DefaultMode: pointer.Int32Ptr(420),
											LocalObjectReference: corev1.LocalObjectReference{
												Name: controllerServiceName,
											},
										},
									},
								},
								{
									Name: "etc-aperture-policies",
									VolumeSource: corev1.VolumeSource{
										EmptyDir: &corev1.EmptyDirVolumeSource{},
									},
								},
								{
									Name: "etc-aperture-classification",
									VolumeSource: corev1.VolumeSource{
										ConfigMap: &corev1.ConfigMapVolumeSource{
											DefaultMode: pointer.Int32Ptr(420),
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "classification",
											},
											Optional: pointer.BoolPtr(true),
										},
									},
								},
								{
									Name: "server-cert",
									VolumeSource: corev1.VolumeSource{
										Secret: &corev1.SecretVolumeSource{
											DefaultMode: pointer.Int32Ptr(420),
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
						ServiceAccountName:        test,
						ImagePullSecrets:          []corev1.LocalObjectReference{},
						HostAliases:               []corev1.HostAlias{},
						Affinity:                  &corev1.Affinity{},
						NodeSelector:              map[string]string{},
						Tolerations:               []corev1.Toleration{},
						PriorityClassName:         test,
						TopologySpreadConstraints: []corev1.TopologySpreadConstraint{},
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

		dep := &appsv1.Deployment{}
		err := deploymentMutate(dep, expected.Spec)()

		Expect(err).NotTo(HaveOccurred())
		Expect(dep).To(Equal(expected))
	})
})
