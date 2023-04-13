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
	"bytes"
	"encoding/pem"
	"fmt"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"

	agentv1alpha1 "github.com/fluxninja/aperture/operator/api/agent/v1alpha1"
	"github.com/fluxninja/aperture/operator/api/common"
	controllerv1alpha1 "github.com/fluxninja/aperture/operator/api/controller/v1alpha1"
)

var _ = Describe("Tests for containerSecurityContext", func() {
	Context("When ContainerSecurityContext is not enabled", func() {
		It("returns correct SecurityContext", func() {
			instance := &agentv1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: agentv1alpha1.AgentSpec{
					CommonSpec: common.CommonSpec{
						ContainerSecurityContext: common.ContainerSecurityContext{
							Enabled: false,
						},
					},
				},
			}

			result := ContainerSecurityContext(instance.Spec.ContainerSecurityContext)
			Expect(result).To(Equal(&corev1.SecurityContext{}))
		})
	})

	Context("When ContainerSecurityContext is enabled", func() {
		It("returns correct SecurityContext", func() {
			instance := &agentv1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: agentv1alpha1.AgentSpec{
					CommonSpec: common.CommonSpec{
						ContainerSecurityContext: common.ContainerSecurityContext{
							Enabled:                true,
							RunAsUser:              0,
							RunAsNonRootUser:       false,
							ReadOnlyRootFilesystem: false,
						},
					},
				},
			}

			expected := &corev1.SecurityContext{
				RunAsUser:              pointer.Int64(0),
				RunAsNonRoot:           pointer.Bool(false),
				ReadOnlyRootFilesystem: pointer.Bool(false),
			}

			result := ContainerSecurityContext(instance.Spec.ContainerSecurityContext)
			Expect(result).To(Equal(expected))
		})
	})
})

var _ = Describe("Tests for podSecurityContext", func() {
	Context("When PodSecurityContext is not enabled", func() {
		It("returns correct SecurityContext", func() {
			instance := &agentv1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: agentv1alpha1.AgentSpec{
					CommonSpec: common.CommonSpec{
						PodSecurityContext: common.PodSecurityContext{
							Enabled: false,
						},
					},
				},
			}

			result := PodSecurityContext(instance.Spec.PodSecurityContext)
			Expect(result).To(Equal(&corev1.PodSecurityContext{}))
		})
	})

	Context("When PodSecurityContext is enabled", func() {
		It("returns correct SecurityContext", func() {
			instance := &agentv1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: agentv1alpha1.AgentSpec{
					CommonSpec: common.CommonSpec{
						PodSecurityContext: common.PodSecurityContext{
							Enabled: true,
							FsGroup: 1001,
						},
					},
				},
			}

			expected := &corev1.PodSecurityContext{
				FSGroup: pointer.Int64(1001),
			}

			result := PodSecurityContext(instance.Spec.PodSecurityContext)
			Expect(result).To(Equal(expected))
		})
	})
})

var _ = Describe("Tests for imageString", func() {
	Context("When local image registry is provided", func() {
		It("returns correct image string", func() {
			instance := &agentv1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: agentv1alpha1.AgentSpec{
					Image: common.AgentImage{
						Image: common.Image{
							Registry: Test,
							Tag:      Test,
						},
						Repository: Test,
					},
				},
			}

			result := ImageString(instance.Spec.Image.Image, instance.Spec.Image.Repository)
			Expect(result).To(Equal("test/test:test"))
		})
	})

	Context("When any image registry is not provided", func() {
		It("returns correct image string", func() {
			instance := &agentv1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: agentv1alpha1.AgentSpec{
					Image: common.AgentImage{
						Image: common.Image{
							Tag: Test,
						},
						Repository: Test,
					},
				},
			}

			result := ImageString(instance.Spec.Image.Image, instance.Spec.Image.Repository)
			Expect(result).To(Equal("test:test"))
		})
	})
})

var _ = Describe("Tests for imagePullSecrets", func() {
	Context("When only local image pullSecrets are provided", func() {
		It("returns correct imagePullSecrets", func() {
			instance := &agentv1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: agentv1alpha1.AgentSpec{
					Image: common.AgentImage{
						Image: common.Image{
							PullSecrets: TestArray,
						},
					},
				},
			}

			expected := []corev1.LocalObjectReference{
				{
					Name: Test,
				},
			}

			result := ImagePullSecrets(instance.Spec.Image.Image)
			Expect(result).To(Equal(expected))
		})
	})
})

var _ = Describe("Tests for containerEnvFrom", func() {
	Context("When only configMap is provided", func() {
		It("returns correct EnvFromSource", func() {
			instance := &agentv1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: agentv1alpha1.AgentSpec{
					CommonSpec: common.CommonSpec{
						ExtraEnvVarsCM: Test,
					},
				},
			}

			expected := []corev1.EnvFromSource{
				{
					ConfigMapRef: &corev1.ConfigMapEnvSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: Test,
						},
					},
				},
			}

			result := ContainerEnvFrom(instance.Spec.CommonSpec)
			Expect(result).To(Equal(expected))
		})
	})

	Context("When only secret is provided", func() {
		It("returns correct EnvFromSource", func() {
			instance := &agentv1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: agentv1alpha1.AgentSpec{
					CommonSpec: common.CommonSpec{
						ExtraEnvVarsSecret: Test,
					},
				},
			}

			expected := []corev1.EnvFromSource{
				{
					SecretRef: &corev1.SecretEnvSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: Test,
						},
					},
				},
			}

			result := ContainerEnvFrom(instance.Spec.CommonSpec)
			Expect(result).To(Equal(expected))
		})
	})

	Context("When both configMap and secret are provided", func() {
		It("returns correct EnvFromSource", func() {
			instance := &agentv1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: agentv1alpha1.AgentSpec{
					CommonSpec: common.CommonSpec{
						ExtraEnvVarsCM:     Test,
						ExtraEnvVarsSecret: Test,
					},
				},
			}

			expected := []corev1.EnvFromSource{
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
			}

			result := ContainerEnvFrom(instance.Spec.CommonSpec)
			Expect(result).To(Equal(expected))
		})
	})
})

var _ = Describe("Tests for containerProbes", func() {
	Context("When only livenessProbe is provided", func() {
		It("returns correct Probe", func() {
			instance := &agentv1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: agentv1alpha1.AgentSpec{
					CommonSpec: common.CommonSpec{
						LivenessProbe: common.Probe{
							Enabled:             true,
							TimeoutSeconds:      10,
							InitialDelaySeconds: 10,
							PeriodSeconds:       10,
							FailureThreshold:    5,
							SuccessThreshold:    1,
						},
					},
				},
			}

			expectedLiveness := &corev1.Probe{
				ProbeHandler: corev1.ProbeHandler{
					HTTPGet: &corev1.HTTPGetAction{
						Path:   "/v1/status/system/liveness",
						Port:   intstr.FromString(Server),
						Scheme: corev1.URISchemeHTTP,
					},
				},
				InitialDelaySeconds: 10,
				TimeoutSeconds:      10,
				PeriodSeconds:       10,
				FailureThreshold:    5,
				SuccessThreshold:    1,
			}

			var expectedReadiness *corev1.Probe

			liveness, readiness := ContainerProbes(instance.Spec.CommonSpec, corev1.URISchemeHTTP)
			Expect(liveness).To(Equal(expectedLiveness))
			Expect(readiness).To(Equal(expectedReadiness))
		})
	})

	Context("When only custom livenessProbe is provided", func() {
		It("returns correct Probe", func() {
			probe := &corev1.Probe{
				ProbeHandler: corev1.ProbeHandler{
					HTTPGet: &corev1.HTTPGetAction{
						Path: "/v1/status/system/liveness",
						Port: intstr.FromString(Server),
					},
				},
				InitialDelaySeconds: 10,
				TimeoutSeconds:      10,
				PeriodSeconds:       10,
				FailureThreshold:    5,
				SuccessThreshold:    1,
			}

			instance := &agentv1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: agentv1alpha1.AgentSpec{
					CommonSpec: common.CommonSpec{
						CustomLivenessProbe: probe,
					},
				},
			}

			expectedLiveness := probe

			var expectedReadiness *corev1.Probe

			liveness, readiness := ContainerProbes(instance.Spec.CommonSpec, corev1.URISchemeHTTP)
			Expect(liveness).To(Equal(expectedLiveness))
			Expect(readiness).To(Equal(expectedReadiness))
		})
	})

	Context("When only readinessProbe is provided", func() {
		It("returns correct Probe", func() {
			instance := &agentv1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: agentv1alpha1.AgentSpec{
					CommonSpec: common.CommonSpec{
						ReadinessProbe: common.Probe{
							Enabled:             true,
							TimeoutSeconds:      10,
							InitialDelaySeconds: 10,
							PeriodSeconds:       10,
							FailureThreshold:    5,
							SuccessThreshold:    1,
						},
					},
				},
			}

			expectedReadiness := &corev1.Probe{
				ProbeHandler: corev1.ProbeHandler{
					HTTPGet: &corev1.HTTPGetAction{
						Path:   "/v1/status/system/readiness",
						Port:   intstr.FromString(Server),
						Scheme: corev1.URISchemeHTTP,
					},
				},
				InitialDelaySeconds: 10,
				TimeoutSeconds:      10,
				PeriodSeconds:       10,
				FailureThreshold:    5,
				SuccessThreshold:    1,
			}

			var expectedLiveness *corev1.Probe

			liveness, readiness := ContainerProbes(instance.Spec.CommonSpec, corev1.URISchemeHTTP)
			Expect(liveness).To(Equal(expectedLiveness))
			Expect(readiness).To(Equal(expectedReadiness))
		})
	})

	Context("When only custom readinessProbe is provided", func() {
		It("returns correct Probe", func() {
			probe := &corev1.Probe{
				ProbeHandler: corev1.ProbeHandler{
					HTTPGet: &corev1.HTTPGetAction{
						Path: "/v1/status/system/readiness",
						Port: intstr.FromString(Server),
					},
				},
				InitialDelaySeconds: 10,
				TimeoutSeconds:      10,
				PeriodSeconds:       10,
				FailureThreshold:    5,
				SuccessThreshold:    1,
			}

			instance := &agentv1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: agentv1alpha1.AgentSpec{
					CommonSpec: common.CommonSpec{
						CustomReadinessProbe: probe,
					},
				},
			}

			expectedReadiness := probe

			var expectedLiveness *corev1.Probe

			liveness, readiness := ContainerProbes(instance.Spec.CommonSpec, corev1.URISchemeHTTP)
			Expect(liveness).To(Equal(expectedLiveness))
			Expect(readiness).To(Equal(expectedReadiness))
		})
	})

	Context("When both livenessProbe and readinessProbe are provided", func() {
		It("returns correct Probe", func() {
			instance := &agentv1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: agentv1alpha1.AgentSpec{
					CommonSpec: common.CommonSpec{
						LivenessProbe: common.Probe{
							Enabled:             true,
							InitialDelaySeconds: 10,
							TimeoutSeconds:      10,
							PeriodSeconds:       10,
							FailureThreshold:    1,
							SuccessThreshold:    1,
						},
						ReadinessProbe: common.Probe{
							Enabled:             true,
							InitialDelaySeconds: 10,
							TimeoutSeconds:      10,
							PeriodSeconds:       10,
							FailureThreshold:    1,
							SuccessThreshold:    1,
						},
					},
				},
			}

			expectedReadiness := &corev1.Probe{
				ProbeHandler: corev1.ProbeHandler{
					HTTPGet: &corev1.HTTPGetAction{
						Path:   "/v1/status/system/readiness",
						Port:   intstr.FromString(Server),
						Scheme: corev1.URISchemeHTTP,
					},
				},
				InitialDelaySeconds: 10,
				TimeoutSeconds:      10,
				PeriodSeconds:       10,
				FailureThreshold:    1,
				SuccessThreshold:    1,
			}

			expectedLiveness := &corev1.Probe{
				ProbeHandler: corev1.ProbeHandler{
					HTTPGet: &corev1.HTTPGetAction{
						Path:   "/v1/status/system/liveness",
						Port:   intstr.FromString(Server),
						Scheme: corev1.URISchemeHTTP,
					},
				},
				InitialDelaySeconds: 10,
				TimeoutSeconds:      10,
				PeriodSeconds:       10,
				FailureThreshold:    1,
				SuccessThreshold:    1,
			}

			liveness, readiness := ContainerProbes(instance.Spec.CommonSpec, corev1.URISchemeHTTP)
			Expect(liveness).To(Equal(expectedLiveness))
			Expect(readiness).To(Equal(expectedReadiness))
		})
	})
})

var _ = Describe("Tests for agentEnv", func() {
	Context("When extra Env are not provided", func() {
		It("returns correct EnvVarSource", func() {
			instance := &agentv1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: agentv1alpha1.AgentSpec{
					CommonSpec: common.CommonSpec{
						Secrets: common.Secrets{
							FluxNinjaExtension: common.APIKeySecret{
								Create: true,
								SecretKeyRef: common.SecretKeyRef{
									Name: Test,
									Key:  Test,
								},
							},
						},
					},
				},
			}

			expected := []corev1.EnvVar{
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
					Name:  "APERTURE_AGENT_SERVICE_DISCOVERY_KUBERNETES_ENABLED",
					Value: "true",
				},
				{
					Name: "APERTURE_AGENT_FLUXNINJA_API_KEY",
					ValueFrom: &corev1.EnvVarSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: Test,
							},
							Key:      Test,
							Optional: pointer.Bool(false),
						},
					},
				},
			}

			result := AgentEnv(instance, "")
			Expect(result).To(Equal(expected))
		})
	})

	Context("When extra Env are provided", func() {
		It("returns correct EnvVarSource", func() {
			instance := &agentv1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: agentv1alpha1.AgentSpec{
					CommonSpec: common.CommonSpec{
						ExtraEnvVars: []corev1.EnvVar{
							{
								Name:  Test,
								Value: Test,
							},
						},
					},
				},
			}

			expected := []corev1.EnvVar{
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
					Name:  "APERTURE_AGENT_SERVICE_DISCOVERY_KUBERNETES_ENABLED",
					Value: "true",
				},
			}

			result := AgentEnv(instance, "")
			Expect(result).To(Equal(expected))
		})
	})
})

var _ = Describe("Tests for agentVolumeMounts", func() {
	Context("When extra VolumeMounts are not provided", func() {
		It("returns correct VolumeMount", func() {
			instance := &agentv1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: agentv1alpha1.AgentSpec{},
			}

			expected := []corev1.VolumeMount{
				{
					Name:      "aperture-agent-config",
					MountPath: "/etc/aperture/aperture-agent/config",
				},
			}

			result := AgentVolumeMounts(instance.Spec)
			Expect(result).To(Equal(expected))
		})
	})

	Context("When extra VolumeMounts are provided", func() {
		It("returns correct VolumeMount", func() {
			instance := &agentv1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: agentv1alpha1.AgentSpec{
					CommonSpec: common.CommonSpec{
						ExtraVolumeMounts: []corev1.VolumeMount{
							{
								Name:      Test,
								MountPath: Test,
							},
						},
					},
				},
			}

			expected := []corev1.VolumeMount{
				{
					Name:      Test,
					MountPath: Test,
				},
				{
					Name:      "aperture-agent-config",
					MountPath: "/etc/aperture/aperture-agent/config",
				},
			}

			result := AgentVolumeMounts(instance.Spec)
			Expect(result).To(Equal(expected))
		})
	})
})

var _ = Describe("Tests for agentVolumes", func() {
	Context("When extra Volumes are not provided", func() {
		It("returns correct Volume", func() {
			instance := &agentv1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: agentv1alpha1.AgentSpec{},
			}

			expected := []corev1.Volume{
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
			}

			result := AgentVolumes(instance.Spec)
			Expect(result).To(Equal(expected))
		})
	})

	Context("When extra Volumes are provided", func() {
		It("returns correct Volume", func() {
			instance := &agentv1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: agentv1alpha1.AgentSpec{
					CommonSpec: common.CommonSpec{
						ExtraVolumes: []corev1.Volume{
							{
								Name: Test,
								VolumeSource: corev1.VolumeSource{
									EmptyDir: &corev1.EmptyDirVolumeSource{},
								},
							},
						},
					},
				},
			}

			expected := []corev1.Volume{
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
			}

			result := AgentVolumes(instance.Spec)
			Expect(result).To(Equal(expected))
		})
	})
})

var _ = Describe("Tests for controllerEnv", func() {
	Context("When extra Env are not provided", func() {
		It("returns correct EnvVarSource", func() {
			instance := &controllerv1alpha1.Controller{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: controllerv1alpha1.ControllerSpec{
					CommonSpec: common.CommonSpec{
						Secrets: common.Secrets{
							FluxNinjaExtension: common.APIKeySecret{
								Create: true,
								SecretKeyRef: common.SecretKeyRef{
									Name: Test,
									Key:  Test,
								},
							},
						},
					},
				},
			}

			expected := []corev1.EnvVar{
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
				{
					Name: "APERTURE_CONTROLLER_FLUXNINJA_API_KEY",
					ValueFrom: &corev1.EnvVarSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: Test,
							},
							Key:      Test,
							Optional: pointer.Bool(false),
						},
					},
				},
			}

			result := ControllerEnv(instance)
			Expect(result).To(Equal(expected))
		})
	})

	Context("When extra Env are provided", func() {
		It("returns correct EnvVarSource", func() {
			instance := &controllerv1alpha1.Controller{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: controllerv1alpha1.ControllerSpec{
					CommonSpec: common.CommonSpec{
						ExtraEnvVars: []corev1.EnvVar{
							{
								Name:  Test,
								Value: Test,
							},
						},
					},
				},
			}

			expected := []corev1.EnvVar{
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
			}

			result := ControllerEnv(instance)
			Expect(result).To(Equal(expected))
		})
	})
})

var _ = Describe("Tests for controllerVolumeMounts", func() {
	Context("When extra VolumeMounts are not provided", func() {
		It("returns correct VolumeMount", func() {
			instance := &controllerv1alpha1.Controller{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: controllerv1alpha1.ControllerSpec{},
			}

			expected := []corev1.VolumeMount{
				{
					Name:      "aperture-controller-config",
					MountPath: "/etc/aperture/aperture-controller/config",
				},
				{
					Name:      "server-cert",
					MountPath: "/etc/aperture/aperture-controller/certs",
					ReadOnly:  true,
				},
			}

			result := ControllerVolumeMounts(instance.Spec.CommonSpec)
			Expect(result).To(Equal(expected))
		})
	})

	Context("When extra VolumeMounts are provided", func() {
		It("returns correct VolumeMount", func() {
			instance := &controllerv1alpha1.Controller{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: controllerv1alpha1.ControllerSpec{
					CommonSpec: common.CommonSpec{
						ExtraVolumeMounts: []corev1.VolumeMount{
							{
								Name:      Test,
								MountPath: Test,
							},
						},
					},
				},
			}

			expected := []corev1.VolumeMount{
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
			}

			result := ControllerVolumeMounts(instance.Spec.CommonSpec)
			Expect(result).To(Equal(expected))
		})
	})
})

var _ = Describe("Tests for controllerVolumes", func() {
	Context("When extra Volumes are not provided", func() {
		It("returns correct Volume", func() {
			instance := &controllerv1alpha1.Controller{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: controllerv1alpha1.ControllerSpec{},
			}

			expected := []corev1.Volume{
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
			}

			result := ControllerVolumes(instance.DeepCopy())
			Expect(result).To(Equal(expected))
		})
	})

	Context("When extra Volumes are provided", func() {
		It("returns correct Volume", func() {
			instance := &controllerv1alpha1.Controller{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: controllerv1alpha1.ControllerSpec{
					CommonSpec: common.CommonSpec{
						ExtraVolumes: []corev1.Volume{
							{
								Name: Test,
								VolumeSource: corev1.VolumeSource{
									EmptyDir: &corev1.EmptyDirVolumeSource{},
								},
							},
						},
					},
				},
			}

			expected := []corev1.Volume{
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
			}

			result := ControllerVolumes(instance.DeepCopy())
			Expect(result).To(Equal(expected))
		})
	})
})

var _ = Describe("Tests for commonLabels", func() {
	Context("When global labels are not provided", func() {
		It("returns correct labels", func() {
			instance := &controllerv1alpha1.Controller{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: controllerv1alpha1.ControllerSpec{},
			}

			expected := map[string]string{
				"app.kubernetes.io/name":       AppName,
				"app.kubernetes.io/instance":   instance.GetName(),
				"app.kubernetes.io/managed-by": OperatorName,
				"app.kubernetes.io/component":  Test,
			}

			result := CommonLabels(instance.Spec.Labels, instance.GetName(), Test)
			Expect(result).To(Equal(expected))
		})
	})

	Context("When global labels are provided", func() {
		It("returns correct labels", func() {
			instance := &controllerv1alpha1.Controller{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: controllerv1alpha1.ControllerSpec{
					CommonSpec: common.CommonSpec{
						Labels: TestMap,
					},
				},
			}

			expected := map[string]string{
				"app.kubernetes.io/name":       AppName,
				"app.kubernetes.io/instance":   instance.GetName(),
				"app.kubernetes.io/managed-by": OperatorName,
				"app.kubernetes.io/component":  Test,
				Test:                           Test,
			}

			result := CommonLabels(instance.Spec.Labels, instance.GetName(), Test)
			Expect(result).To(Equal(expected))
		})
	})
})

var _ = Describe("Tests for secretName", func() {
	Context("When secret name is provided", func() {
		It("returns correct secret name", func() {
			instance := &agentv1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: agentv1alpha1.AgentSpec{
					CommonSpec: common.CommonSpec{
						Secrets: common.Secrets{
							FluxNinjaExtension: common.APIKeySecret{
								SecretKeyRef: common.SecretKeyRef{
									Name: Test,
								},
							},
						},
					},
				},
			}

			result := SecretName(AppName, "agent", &instance.Spec.Secrets.FluxNinjaExtension)
			Expect(result).To(Equal(Test))
		})
	})

	Context("When secret name is not provided for Agent", func() {
		It("returns correct secret name for Agent", func() {
			instance := &agentv1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: agentv1alpha1.AgentSpec{
					CommonSpec: common.CommonSpec{
						Secrets: common.Secrets{
							FluxNinjaExtension: common.APIKeySecret{
								SecretKeyRef: common.SecretKeyRef{},
							},
						},
					},
				},
			}

			expected := fmt.Sprintf("%s-agent-apikey", AppName)

			result := SecretName(AppName, "agent", &instance.Spec.Secrets.FluxNinjaExtension)
			Expect(result).To(Equal(expected))
		})
	})

	Context("When secret name is not provided for Controller", func() {
		It("returns correct secret name for controller", func() {
			instance := &agentv1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: agentv1alpha1.AgentSpec{
					CommonSpec: common.CommonSpec{
						Secrets: common.Secrets{
							FluxNinjaExtension: common.APIKeySecret{
								SecretKeyRef: common.SecretKeyRef{},
							},
						},
					},
				},
			}

			expected := fmt.Sprintf("%s-controller-apikey", AppName)

			result := SecretName(AppName, "controller", &instance.Spec.Secrets.FluxNinjaExtension)
			Expect(result).To(Equal(expected))
		})
	})
})

var _ = Describe("Tests for secretDataKey", func() {
	Context("When secret key is provided", func() {
		It("returns correct secret key", func() {
			instance := &agentv1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: agentv1alpha1.AgentSpec{
					CommonSpec: common.CommonSpec{
						Secrets: common.Secrets{
							FluxNinjaExtension: common.APIKeySecret{
								SecretKeyRef: common.SecretKeyRef{
									Key: Test,
								},
							},
						},
					},
				},
			}

			result := SecretDataKey(&instance.Spec.Secrets.FluxNinjaExtension.SecretKeyRef)
			Expect(result).To(Equal(Test))
		})
	})

	Context("When secret key is not provided", func() {
		It("returns correct secret key", func() {
			instance := &agentv1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: agentv1alpha1.AgentSpec{
					CommonSpec: common.CommonSpec{
						Secrets: common.Secrets{
							FluxNinjaExtension: common.APIKeySecret{
								SecretKeyRef: common.SecretKeyRef{},
							},
						},
					},
				},
			}

			result := SecretDataKey(&instance.Spec.Secrets.FluxNinjaExtension.SecretKeyRef)
			Expect(result).To(Equal(SecretKey))
		})
	})
})

var _ = Describe("Tests for checkCertificate", func() {
	Context("When certificate is provided", func() {
		It("returns true", func() {
			os.Setenv("APERTURE_OPERATOR_CERT_DIR", CertDir)
			os.Setenv("APERTURE_OPERATOR_CERT_NAME", "tls.crt")
			os.Setenv("APERTURE_OPERATOR_KEY_NAME", "tls.key")
			os.Setenv("APERTURE_OPERATOR_SERVICE_NAME", AppName)
			os.Setenv("APERTURE_OPERATOR_NAMESPACE", AppName)

			err := CheckAndGenerateCertForOperator()

			Expect(err).NotTo(HaveOccurred())
			Expect(CheckCertificate()).To(Equal(true))
		})
	})

	Context("When certificate is not provided", func() {
		It("returns false", func() {
			os.Setenv("APERTURE_OPERATOR_CERT_DIR", CertDir)
			os.Setenv("APERTURE_OPERATOR_CERT_NAME", "tls1.crt")
			os.Setenv("APERTURE_OPERATOR_SERVICE_NAME", AppName)
			os.Setenv("APERTURE_OPERATOR_KEY_NAME", "tls1.key")

			Expect(CheckCertificate()).To(Equal(false))
		})
	})

	Context("When invalid certificate is provided", func() {
		It("returns false", func() {
			os.Setenv("APERTURE_OPERATOR_CERT_DIR", CertDir)
			os.Setenv("APERTURE_OPERATOR_CERT_NAME", "tls2.crt")
			os.Setenv("APERTURE_OPERATOR_KEY_NAME", "tls2.key")
			os.Setenv("APERTURE_OPERATOR_SERVICE_NAME", AppName)
			os.Setenv("APERTURE_OPERATOR_NAMESPACE", AppName)

			err := CheckAndGenerateCertForOperator()

			Expect(err).NotTo(HaveOccurred())

			certPath := fmt.Sprintf("%s/%s", os.Getenv("APERTURE_OPERATOR_CERT_DIR"), os.Getenv("APERTURE_OPERATOR_CERT_NAME"))
			serverCertPEM := new(bytes.Buffer)
			_ = pem.Encode(serverCertPEM, &pem.Block{
				Type:  "CERTIFICATE",
				Bytes: []byte(Test),
			})
			err = WriteFile(certPath, serverCertPEM)
			Expect(err).NotTo(HaveOccurred())
			Expect(CheckCertificate()).To(Equal(false))
		})
	})

	Context("When environment variables are provided", func() {
		It("uses default values", func() {
			os.Setenv("APERTURE_OPERATOR_CERT_DIR", "")
			os.Setenv("APERTURE_OPERATOR_CERT_NAME", "")
			os.Setenv("APERTURE_OPERATOR_SERVICE_NAME", AppName)
			os.Setenv("APERTURE_OPERATOR_KEY_NAME", "")

			CheckCertificate()

			Expect(os.Getenv("APERTURE_OPERATOR_CERT_DIR")).To(Equal("/tmp/k8s-webhook-server/serving-certs"))
			Expect(os.Getenv("APERTURE_OPERATOR_CERT_NAME")).To(Equal("tls.crt"))
			Expect(os.Getenv("APERTURE_OPERATOR_KEY_NAME")).To(Equal("tls.key"))
		})
	})
})

var _ = Describe("Tests for CheckAndGenerateCert", func() {
	Context("When service name is not provided in environment variable", func() {
		It("it should not create cert", func() {
			os.Setenv("APERTURE_OPERATOR_CERT_DIR", CertDir)
			os.Setenv("APERTURE_OPERATOR_CERT_NAME", "tls3.crt")
			os.Setenv("APERTURE_OPERATOR_KEY_NAME", "tls3.key")
			os.Setenv("APERTURE_OPERATOR_NAMESPACE", AppName)
			os.Setenv("APERTURE_OPERATOR_SERVICE_NAME", "")

			err := CheckAndGenerateCertForOperator()

			Expect(err).To(HaveOccurred())
		})
	})

	Context("When certificate is not provided", func() {
		It("it should create cert", func() {
			os.Setenv("APERTURE_OPERATOR_CERT_DIR", CertDir)
			os.Setenv("APERTURE_OPERATOR_CERT_NAME", "tls4.crt")
			os.Setenv("APERTURE_OPERATOR_KEY_NAME", "tls4.key")
			os.Setenv("APERTURE_OPERATOR_SERVICE_NAME", AppName)
			os.Setenv("APERTURE_OPERATOR_NAMESPACE", "")

			Expect(CheckCertificate()).To(Equal(false))
			Expect(CheckAndGenerateCertForOperator()).To(Succeed())
			Expect(CheckCertificate()).To(Equal(true))
		})
	})
})
