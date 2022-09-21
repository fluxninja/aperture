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
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"

	"github.com/fluxninja/aperture/operator/api/v1alpha1"
	etcd "github.com/fluxninja/aperture/pkg/etcd/client"
	"github.com/fluxninja/aperture/pkg/prometheus"
)

var _ = Describe("Tests for containerSecurityContext", func() {
	Context("When ContainerSecurityContext is not enabled", func() {
		It("returns correct SecurityContext", func() {
			instance := &v1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.AgentSpec{
					CommonSpec: v1alpha1.CommonSpec{
						ContainerSecurityContext: v1alpha1.ContainerSecurityContext{
							Enabled: false,
						},
					},
				},
			}

			result := containerSecurityContext(instance.Spec.ContainerSecurityContext)
			Expect(result).To(Equal(&corev1.SecurityContext{}))
		})
	})

	Context("When ContainerSecurityContext is enabled", func() {
		It("returns correct SecurityContext", func() {
			instance := &v1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.AgentSpec{
					CommonSpec: v1alpha1.CommonSpec{
						ContainerSecurityContext: v1alpha1.ContainerSecurityContext{
							Enabled:                true,
							RunAsUser:              0,
							RunAsNonRootUser:       false,
							ReadOnlyRootFilesystem: false,
						},
					},
				},
			}

			expected := &corev1.SecurityContext{
				RunAsUser:              pointer.Int64Ptr(0),
				RunAsNonRoot:           pointer.BoolPtr(false),
				ReadOnlyRootFilesystem: pointer.BoolPtr(false),
			}

			result := containerSecurityContext(instance.Spec.ContainerSecurityContext)
			Expect(result).To(Equal(expected))
		})
	})
})

var _ = Describe("Tests for podSecurityContext", func() {
	Context("When PodSecurityContext is not enabled", func() {
		It("returns correct SecurityContext", func() {
			instance := &v1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.AgentSpec{
					CommonSpec: v1alpha1.CommonSpec{
						PodSecurityContext: v1alpha1.PodSecurityContext{
							Enabled: false,
						},
					},
				},
			}

			result := podSecurityContext(instance.Spec.PodSecurityContext)
			Expect(result).To(Equal(&corev1.PodSecurityContext{}))
		})
	})

	Context("When PodSecurityContext is enabled", func() {
		It("returns correct SecurityContext", func() {
			instance := &v1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.AgentSpec{
					CommonSpec: v1alpha1.CommonSpec{
						PodSecurityContext: v1alpha1.PodSecurityContext{
							Enabled: true,
							FsGroup: 1001,
						},
					},
				},
			}

			expected := &corev1.PodSecurityContext{
				FSGroup: pointer.Int64Ptr(1001),
			}

			result := podSecurityContext(instance.Spec.PodSecurityContext)
			Expect(result).To(Equal(expected))
		})
	})
})

var _ = Describe("Tests for imageString", func() {
	Context("When local image registry is provided", func() {
		It("returns correct image string", func() {
			instance := &v1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.AgentSpec{
					Image: v1alpha1.AgentImage{
						Image: v1alpha1.Image{
							Registry: test,
							Tag:      test,
						},
						Repository: test,
					},
				},
			}

			result := imageString(instance.Spec.Image.Image, instance.Spec.Image.Repository)
			Expect(result).To(Equal("test/test:test"))
		})
	})

	Context("When any image registry is not provided", func() {
		It("returns correct image string", func() {
			instance := &v1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.AgentSpec{
					Image: v1alpha1.AgentImage{
						Image: v1alpha1.Image{
							Tag: test,
						},
						Repository: test,
					},
				},
			}

			result := imageString(instance.Spec.Image.Image, instance.Spec.Image.Repository)
			Expect(result).To(Equal("test:test"))
		})
	})
})

var _ = Describe("Tests for imagePullSecrets", func() {
	Context("When only local image pullSecrets are provided", func() {
		It("returns correct imagePullSecrets", func() {
			instance := &v1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.AgentSpec{
					Image: v1alpha1.AgentImage{
						Image: v1alpha1.Image{
							PullSecrets: testArray,
						},
					},
				},
			}

			expected := []corev1.LocalObjectReference{
				{
					Name: test,
				},
			}

			result := imagePullSecrets(instance.Spec.Image.Image)
			Expect(result).To(Equal(expected))
		})
	})
})

var _ = Describe("Tests for containerEnvFrom", func() {
	Context("When only configMap is provided", func() {
		It("returns correct EnvFromSource", func() {
			instance := &v1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.AgentSpec{
					CommonSpec: v1alpha1.CommonSpec{
						ExtraEnvVarsCM: test,
					},
				},
			}

			expected := []corev1.EnvFromSource{
				{
					ConfigMapRef: &corev1.ConfigMapEnvSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: test,
						},
					},
				},
			}

			result := containerEnvFrom(instance.Spec.CommonSpec)
			Expect(result).To(Equal(expected))
		})
	})

	Context("When only secret is provided", func() {
		It("returns correct EnvFromSource", func() {
			instance := &v1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.AgentSpec{
					CommonSpec: v1alpha1.CommonSpec{
						ExtraEnvVarsSecret: test,
					},
				},
			}

			expected := []corev1.EnvFromSource{
				{
					SecretRef: &corev1.SecretEnvSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: test,
						},
					},
				},
			}

			result := containerEnvFrom(instance.Spec.CommonSpec)
			Expect(result).To(Equal(expected))
		})
	})

	Context("When both configMap and secret are provided", func() {
		It("returns correct EnvFromSource", func() {
			instance := &v1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.AgentSpec{
					CommonSpec: v1alpha1.CommonSpec{
						ExtraEnvVarsCM:     test,
						ExtraEnvVarsSecret: test,
					},
				},
			}

			expected := []corev1.EnvFromSource{
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
			}

			result := containerEnvFrom(instance.Spec.CommonSpec)
			Expect(result).To(Equal(expected))
		})
	})
})

var _ = Describe("Tests for containerProbes", func() {
	Context("When only livenessProbe is provided", func() {
		It("returns correct Probe", func() {
			instance := &v1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.AgentSpec{
					CommonSpec: v1alpha1.CommonSpec{
						LivenessProbe: v1alpha1.Probe{
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
						Path:   "/v1/status/liveness",
						Port:   intstr.FromString(server),
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

			liveness, readiness := containerProbes(instance.Spec.CommonSpec, corev1.URISchemeHTTP)
			Expect(liveness).To(Equal(expectedLiveness))
			Expect(readiness).To(Equal(expectedReadiness))
		})
	})

	Context("When only custom livenessProbe is provided", func() {
		It("returns correct Probe", func() {
			probe := &corev1.Probe{
				ProbeHandler: corev1.ProbeHandler{
					HTTPGet: &corev1.HTTPGetAction{
						Path: "/v1/status/liveness",
						Port: intstr.FromString(server),
					},
				},
				InitialDelaySeconds: 10,
				TimeoutSeconds:      10,
				PeriodSeconds:       10,
				FailureThreshold:    5,
				SuccessThreshold:    1,
			}

			instance := &v1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.AgentSpec{
					CommonSpec: v1alpha1.CommonSpec{
						CustomLivenessProbe: probe,
					},
				},
			}

			expectedLiveness := probe

			var expectedReadiness *corev1.Probe

			liveness, readiness := containerProbes(instance.Spec.CommonSpec, corev1.URISchemeHTTP)
			Expect(liveness).To(Equal(expectedLiveness))
			Expect(readiness).To(Equal(expectedReadiness))
		})
	})

	Context("When only readinessProbe is provided", func() {
		It("returns correct Probe", func() {
			instance := &v1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.AgentSpec{
					CommonSpec: v1alpha1.CommonSpec{
						ReadinessProbe: v1alpha1.Probe{
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
						Path:   "/v1/status/readiness",
						Port:   intstr.FromString(server),
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

			liveness, readiness := containerProbes(instance.Spec.CommonSpec, corev1.URISchemeHTTP)
			Expect(liveness).To(Equal(expectedLiveness))
			Expect(readiness).To(Equal(expectedReadiness))
		})
	})

	Context("When only custom readinessProbe is provided", func() {
		It("returns correct Probe", func() {
			probe := &corev1.Probe{
				ProbeHandler: corev1.ProbeHandler{
					HTTPGet: &corev1.HTTPGetAction{
						Path: "/v1/status/readiness",
						Port: intstr.FromString(server),
					},
				},
				InitialDelaySeconds: 10,
				TimeoutSeconds:      10,
				PeriodSeconds:       10,
				FailureThreshold:    5,
				SuccessThreshold:    1,
			}

			instance := &v1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.AgentSpec{
					CommonSpec: v1alpha1.CommonSpec{
						CustomReadinessProbe: probe,
					},
				},
			}

			expectedReadiness := probe

			var expectedLiveness *corev1.Probe

			liveness, readiness := containerProbes(instance.Spec.CommonSpec, corev1.URISchemeHTTP)
			Expect(liveness).To(Equal(expectedLiveness))
			Expect(readiness).To(Equal(expectedReadiness))
		})
	})

	Context("When both livenessProbe and readinessProbe are provided", func() {
		It("returns correct Probe", func() {
			instance := &v1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.AgentSpec{
					CommonSpec: v1alpha1.CommonSpec{
						LivenessProbe: v1alpha1.Probe{
							Enabled:             true,
							InitialDelaySeconds: 10,
							TimeoutSeconds:      10,
							PeriodSeconds:       10,
							FailureThreshold:    1,
							SuccessThreshold:    1,
						},
						ReadinessProbe: v1alpha1.Probe{
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
						Path:   "/v1/status/readiness",
						Port:   intstr.FromString(server),
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
						Path:   "/v1/status/liveness",
						Port:   intstr.FromString(server),
						Scheme: corev1.URISchemeHTTP,
					},
				},
				InitialDelaySeconds: 10,
				TimeoutSeconds:      10,
				PeriodSeconds:       10,
				FailureThreshold:    1,
				SuccessThreshold:    1,
			}

			liveness, readiness := containerProbes(instance.Spec.CommonSpec, corev1.URISchemeHTTP)
			Expect(liveness).To(Equal(expectedLiveness))
			Expect(readiness).To(Equal(expectedReadiness))
		})
	})
})

var _ = Describe("Tests for agentEnv", func() {
	Context("When extra Env are not provided", func() {
		It("returns correct EnvVarSource", func() {
			instance := &v1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.AgentSpec{
					CommonSpec: v1alpha1.CommonSpec{
						Secrets: v1alpha1.Secrets{
							FluxNinjaPlugin: v1alpha1.APIKeySecret{
								Create: true,
								SecretKeyRef: v1alpha1.SecretKeyRef{
									Name: test,
									Key:  test,
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
				{
					Name: "APERTURE_AGENT_FLUXNINJA_PLUGIN_API_KEY",
					ValueFrom: &corev1.EnvVarSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: test,
							},
							Key:      test,
							Optional: pointer.BoolPtr(false),
						},
					},
				},
			}

			result := agentEnv(instance, "")
			Expect(result).To(Equal(expected))
		})
	})

	Context("When extra Env are provided", func() {
		It("returns correct EnvVarSource", func() {
			instance := &v1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.AgentSpec{
					CommonSpec: v1alpha1.CommonSpec{
						ExtraEnvVars: []corev1.EnvVar{
							{
								Name:  test,
								Value: test,
							},
						},
					},
				},
			}

			expected := []corev1.EnvVar{
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
			}

			result := agentEnv(instance, "")
			Expect(result).To(Equal(expected))
		})
	})
})

var _ = Describe("Tests for agentVolumeMounts", func() {
	Context("When extra VolumeMounts are not provided", func() {
		It("returns correct VolumeMount", func() {
			instance := &v1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.AgentSpec{},
			}

			expected := []corev1.VolumeMount{
				{
					Name:      "aperture-agent-config",
					MountPath: "/etc/aperture/aperture-agent/config",
				},
			}

			result := agentVolumeMounts(instance.Spec)
			Expect(result).To(Equal(expected))
		})
	})

	Context("When extra VolumeMounts are provided", func() {
		It("returns correct VolumeMount", func() {
			instance := &v1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.AgentSpec{
					CommonSpec: v1alpha1.CommonSpec{
						ExtraVolumeMounts: []corev1.VolumeMount{
							{
								Name:      test,
								MountPath: test,
							},
						},
					},
				},
			}

			expected := []corev1.VolumeMount{
				{
					Name:      test,
					MountPath: test,
				},
				{
					Name:      "aperture-agent-config",
					MountPath: "/etc/aperture/aperture-agent/config",
				},
			}

			result := agentVolumeMounts(instance.Spec)
			Expect(result).To(Equal(expected))
		})
	})
})

var _ = Describe("Tests for agentVolumes", func() {
	Context("When extra Volumes are not provided", func() {
		It("returns correct Volume", func() {
			instance := &v1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.AgentSpec{},
			}

			expected := []corev1.Volume{
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
			}

			result := agentVolumes(instance.Spec)
			Expect(result).To(Equal(expected))
		})
	})

	Context("When extra Volumes are provided", func() {
		It("returns correct Volume", func() {
			instance := &v1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.AgentSpec{
					CommonSpec: v1alpha1.CommonSpec{
						ExtraVolumes: []corev1.Volume{
							{
								Name: test,
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
			}

			result := agentVolumes(instance.Spec)
			Expect(result).To(Equal(expected))
		})
	})
})

var _ = Describe("Tests for controllerEnv", func() {
	Context("When extra Env are not provided", func() {
		It("returns correct EnvVarSource", func() {
			instance := &v1alpha1.Controller{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.ControllerSpec{
					CommonSpec: v1alpha1.CommonSpec{
						Secrets: v1alpha1.Secrets{
							FluxNinjaPlugin: v1alpha1.APIKeySecret{
								Create: true,
								SecretKeyRef: v1alpha1.SecretKeyRef{
									Name: test,
									Key:  test,
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
							APIVersion: v1Version,
							FieldPath:  "spec.nodeName",
						},
					},
				},
				{
					Name: "APERTURE_CONTROLLER_FLUXNINJA_PLUGIN_API_KEY",
					ValueFrom: &corev1.EnvVarSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: test,
							},
							Key:      test,
							Optional: pointer.BoolPtr(false),
						},
					},
				},
			}

			result := controllerEnv(instance)
			Expect(result).To(Equal(expected))
		})
	})

	Context("When extra Env are provided", func() {
		It("returns correct EnvVarSource", func() {
			instance := &v1alpha1.Controller{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.ControllerSpec{
					CommonSpec: v1alpha1.CommonSpec{
						ExtraEnvVars: []corev1.EnvVar{
							{
								Name:  test,
								Value: test,
							},
						},
					},
				},
			}

			expected := []corev1.EnvVar{
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
			}

			result := controllerEnv(instance)
			Expect(result).To(Equal(expected))
		})
	})
})

var _ = Describe("Tests for controllerVolumeMounts", func() {
	Context("When extra VolumeMounts are not provided", func() {
		It("returns correct VolumeMount", func() {
			instance := &v1alpha1.Controller{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.ControllerSpec{},
			}

			expected := []corev1.VolumeMount{
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
			}

			result := controllerVolumeMounts(instance.Spec.CommonSpec)
			Expect(result).To(Equal(expected))
		})
	})

	Context("When extra VolumeMounts are provided", func() {
		It("returns correct VolumeMount", func() {
			instance := &v1alpha1.Controller{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.ControllerSpec{
					CommonSpec: v1alpha1.CommonSpec{
						ExtraVolumeMounts: []corev1.VolumeMount{
							{
								Name:      test,
								MountPath: test,
							},
						},
					},
				},
			}

			expected := []corev1.VolumeMount{
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
			}

			result := controllerVolumeMounts(instance.Spec.CommonSpec)
			Expect(result).To(Equal(expected))
		})
	})
})

var _ = Describe("Tests for controllerVolumes", func() {
	Context("When extra Volumes are not provided", func() {
		It("returns correct Volume", func() {
			instance := &v1alpha1.Controller{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.ControllerSpec{},
			}

			expected := []corev1.Volume{
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
			}

			result := controllerVolumes(instance.DeepCopy())
			Expect(result).To(Equal(expected))
		})
	})

	Context("When extra Volumes are provided", func() {
		It("returns correct Volume", func() {
			instance := &v1alpha1.Controller{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.ControllerSpec{
					CommonSpec: v1alpha1.CommonSpec{
						ExtraVolumes: []corev1.Volume{
							{
								Name: test,
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
			}

			result := controllerVolumes(instance.DeepCopy())
			Expect(result).To(Equal(expected))
		})
	})
})

var _ = Describe("Tests for commonLabels", func() {
	Context("When global labels are not provided", func() {
		It("returns correct labels", func() {
			instance := &v1alpha1.Controller{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.ControllerSpec{},
			}

			expected := map[string]string{
				"app.kubernetes.io/name":       appName,
				"app.kubernetes.io/instance":   instance.GetName(),
				"app.kubernetes.io/managed-by": operatorName,
				"app.kubernetes.io/component":  test,
			}

			result := commonLabels(instance.Spec.Labels, instance.GetName(), test)
			Expect(result).To(Equal(expected))
		})
	})

	Context("When global labels are provided", func() {
		It("returns correct labels", func() {
			instance := &v1alpha1.Controller{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.ControllerSpec{
					CommonSpec: v1alpha1.CommonSpec{
						Labels: testMap,
					},
				},
			}

			expected := map[string]string{
				"app.kubernetes.io/name":       appName,
				"app.kubernetes.io/instance":   instance.GetName(),
				"app.kubernetes.io/managed-by": operatorName,
				"app.kubernetes.io/component":  test,
				test:                           test,
			}

			result := commonLabels(instance.Spec.Labels, instance.GetName(), test)
			Expect(result).To(Equal(expected))
		})
	})
})

var _ = Describe("Tests for checkEtcdEndpoints", func() {
	Context("When Etcd endpoints are not provided", func() {
		It("returns correct etcd config", func() {
			instance := &v1alpha1.Controller{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.ControllerSpec{},
			}

			expected := etcd.EtcdConfig{
				Endpoints: []string{
					fmt.Sprintf("http://%s-etcd.%s:2379", appName, appName),
				},
			}

			result := checkEtcdEndpoints(instance.Spec.ConfigSpec.Etcd, instance.Name, instance.Namespace)
			Expect(result).To(Equal(expected))
		})
	})

	Context("When Etcd endpoints are provided", func() {
		It("returns correct etcd config", func() {
			instance := &v1alpha1.Controller{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.ControllerSpec{
					ConfigSpec: v1alpha1.ControllerConfigSpec{
						CommonConfigSpec: v1alpha1.CommonConfigSpec{
							Etcd: etcd.EtcdConfig{
								Endpoints: testArray,
							},
						},
					},
				},
			}

			expected := etcd.EtcdConfig{
				Endpoints: testArray,
			}

			result := checkEtcdEndpoints(instance.Spec.ConfigSpec.Etcd, instance.Name, instance.Namespace)
			Expect(result).To(Equal(expected))
		})
	})

	Context("When Etcd endpoints are provided with empty string", func() {
		It("returns correct etcd config", func() {
			instance := &v1alpha1.Controller{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.ControllerSpec{
					ConfigSpec: v1alpha1.ControllerConfigSpec{
						CommonConfigSpec: v1alpha1.CommonConfigSpec{
							Etcd: etcd.EtcdConfig{
								Endpoints: []string{""},
							},
						},
					},
				},
			}

			expected := etcd.EtcdConfig{
				Endpoints: []string{
					fmt.Sprintf("http://%s-etcd.%s:2379", appName, appName),
				},
			}

			result := checkEtcdEndpoints(instance.Spec.ConfigSpec.Etcd, instance.Name, instance.Namespace)
			Expect(result).To(Equal(expected))
		})
	})
})

var _ = Describe("Tests for checkPrometheusAddress", func() {
	Context("When prometheus address is not provided", func() {
		It("returns correct prometheus address", func() {
			instance := &v1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.AgentSpec{},
			}

			expected := fmt.Sprintf("http://%s-prometheus-server.%s:80", appName, appName)

			result := checkPrometheusAddress(instance.Spec.ConfigSpec.Prometheus.Address, instance.Name, instance.Namespace)
			Expect(result).To(Equal(expected))
		})
	})

	Context("When prometheus address is provided", func() {
		It("returns correct prometheus address", func() {
			instance := &v1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.AgentSpec{
					ConfigSpec: v1alpha1.AgentConfigSpec{
						CommonConfigSpec: v1alpha1.CommonConfigSpec{
							Prometheus: prometheus.PrometheusConfig{
								Address: test,
							},
						},
					},
				},
			}

			result := checkPrometheusAddress(instance.Spec.ConfigSpec.Prometheus.Address, instance.Name, instance.Namespace)
			Expect(result).To(Equal(test))
		})
	})
})

var _ = Describe("Tests for secretName", func() {
	Context("When secret name is provided", func() {
		It("returns correct secret name", func() {
			instance := &v1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.AgentSpec{
					CommonSpec: v1alpha1.CommonSpec{
						Secrets: v1alpha1.Secrets{
							FluxNinjaPlugin: v1alpha1.APIKeySecret{
								SecretKeyRef: v1alpha1.SecretKeyRef{
									Name: test,
								},
							},
						},
					},
				},
			}

			result := secretName(appName, "agent", &instance.Spec.Secrets.FluxNinjaPlugin)
			Expect(result).To(Equal(test))
		})
	})

	Context("When secret name is not provided for Agent", func() {
		It("returns correct secret name for Agent", func() {
			instance := &v1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.AgentSpec{
					CommonSpec: v1alpha1.CommonSpec{
						Secrets: v1alpha1.Secrets{
							FluxNinjaPlugin: v1alpha1.APIKeySecret{
								SecretKeyRef: v1alpha1.SecretKeyRef{},
							},
						},
					},
				},
			}

			expected := fmt.Sprintf("%s-agent-apikey", appName)

			result := secretName(appName, "agent", &instance.Spec.Secrets.FluxNinjaPlugin)
			Expect(result).To(Equal(expected))
		})
	})

	Context("When secret name is not provided for Controller", func() {
		It("returns correct secret name for controller", func() {
			instance := &v1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.AgentSpec{
					CommonSpec: v1alpha1.CommonSpec{
						Secrets: v1alpha1.Secrets{
							FluxNinjaPlugin: v1alpha1.APIKeySecret{
								SecretKeyRef: v1alpha1.SecretKeyRef{},
							},
						},
					},
				},
			}

			expected := fmt.Sprintf("%s-controller-apikey", appName)

			result := secretName(appName, "controller", &instance.Spec.Secrets.FluxNinjaPlugin)
			Expect(result).To(Equal(expected))
		})
	})
})

var _ = Describe("Tests for secretDataKey", func() {
	Context("When secret key is provided", func() {
		It("returns correct secret key", func() {
			instance := &v1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.AgentSpec{
					CommonSpec: v1alpha1.CommonSpec{
						Secrets: v1alpha1.Secrets{
							FluxNinjaPlugin: v1alpha1.APIKeySecret{
								SecretKeyRef: v1alpha1.SecretKeyRef{
									Key: test,
								},
							},
						},
					},
				},
			}

			result := secretDataKey(&instance.Spec.Secrets.FluxNinjaPlugin.SecretKeyRef)
			Expect(result).To(Equal(test))
		})
	})

	Context("When secret key is not provided", func() {
		It("returns correct secret key", func() {
			instance := &v1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.AgentSpec{
					CommonSpec: v1alpha1.CommonSpec{
						Secrets: v1alpha1.Secrets{
							FluxNinjaPlugin: v1alpha1.APIKeySecret{
								SecretKeyRef: v1alpha1.SecretKeyRef{},
							},
						},
					},
				},
			}

			result := secretDataKey(&instance.Spec.Secrets.FluxNinjaPlugin.SecretKeyRef)
			Expect(result).To(Equal(secretKey))
		})
	})
})

var _ = Describe("Tests for checkCertificate", func() {
	Context("When certificate is provided", func() {
		It("returns true", func() {
			os.Setenv("APERTURE_OPERATOR_CERT_DIR", certDir)
			os.Setenv("APERTURE_OPERATOR_CERT_NAME", "tls.crt")
			os.Setenv("APERTURE_OPERATOR_KEY_NAME", "tls.key")
			os.Setenv("APERTURE_OPERATOR_SERVICE_NAME", appName)
			os.Setenv("APERTURE_OPERATOR_NAMESPACE", appName)

			err := CheckAndGenerateCertForOperator()

			Expect(err).NotTo(HaveOccurred())
			Expect(checkCertificate()).To(Equal(true))
		})
	})

	Context("When certificate is not provided", func() {
		It("returns false", func() {
			os.Setenv("APERTURE_OPERATOR_CERT_DIR", certDir)
			os.Setenv("APERTURE_OPERATOR_CERT_NAME", "tls1.crt")
			os.Setenv("APERTURE_OPERATOR_SERVICE_NAME", appName)
			os.Setenv("APERTURE_OPERATOR_KEY_NAME", "tls1.key")

			Expect(checkCertificate()).To(Equal(false))
		})
	})

	Context("When invalid certificate is provided", func() {
		It("returns false", func() {
			os.Setenv("APERTURE_OPERATOR_CERT_DIR", certDir)
			os.Setenv("APERTURE_OPERATOR_CERT_NAME", "tls2.crt")
			os.Setenv("APERTURE_OPERATOR_KEY_NAME", "tls2.key")
			os.Setenv("APERTURE_OPERATOR_SERVICE_NAME", appName)
			os.Setenv("APERTURE_OPERATOR_NAMESPACE", appName)

			err := CheckAndGenerateCertForOperator()

			Expect(err).NotTo(HaveOccurred())

			certPath := fmt.Sprintf("%s/%s", os.Getenv("APERTURE_OPERATOR_CERT_DIR"), os.Getenv("APERTURE_OPERATOR_CERT_NAME"))
			serverCertPEM := new(bytes.Buffer)
			_ = pem.Encode(serverCertPEM, &pem.Block{
				Type:  "CERTIFICATE",
				Bytes: []byte(test),
			})
			err = writeFile(certPath, serverCertPEM)
			Expect(err).NotTo(HaveOccurred())
			Expect(checkCertificate()).To(Equal(false))
		})
	})

	Context("When environment variables are provided", func() {
		It("uses default values", func() {
			os.Setenv("APERTURE_OPERATOR_CERT_DIR", "")
			os.Setenv("APERTURE_OPERATOR_CERT_NAME", "")
			os.Setenv("APERTURE_OPERATOR_SERVICE_NAME", appName)
			os.Setenv("APERTURE_OPERATOR_KEY_NAME", "")

			checkCertificate()

			Expect(os.Getenv("APERTURE_OPERATOR_CERT_DIR")).To(Equal("/tmp/k8s-webhook-server/serving-certs"))
			Expect(os.Getenv("APERTURE_OPERATOR_CERT_NAME")).To(Equal("tls.crt"))
			Expect(os.Getenv("APERTURE_OPERATOR_KEY_NAME")).To(Equal("tls.key"))
		})
	})
})

var _ = Describe("Tests for CheckAndGenerateCert", func() {
	Context("When service name is not provided in environment variable", func() {
		It("it should not create cert", func() {
			os.Setenv("APERTURE_OPERATOR_CERT_DIR", certDir)
			os.Setenv("APERTURE_OPERATOR_CERT_NAME", "tls3.crt")
			os.Setenv("APERTURE_OPERATOR_KEY_NAME", "tls3.key")
			os.Setenv("APERTURE_OPERATOR_NAMESPACE", appName)
			os.Setenv("APERTURE_OPERATOR_SERVICE_NAME", "")

			err := CheckAndGenerateCertForOperator()

			Expect(err).To(HaveOccurred())
		})
	})

	Context("When certificate is not provided", func() {
		It("it should create cert", func() {
			os.Setenv("APERTURE_OPERATOR_CERT_DIR", certDir)
			os.Setenv("APERTURE_OPERATOR_CERT_NAME", "tls4.crt")
			os.Setenv("APERTURE_OPERATOR_KEY_NAME", "tls4.key")
			os.Setenv("APERTURE_OPERATOR_SERVICE_NAME", appName)
			os.Setenv("APERTURE_OPERATOR_NAMESPACE", "")

			Expect(checkCertificate()).To(Equal(false))
			Expect(CheckAndGenerateCertForOperator()).To(BeNil())
			Expect(checkCertificate()).To(Equal(true))
		})
	})
})

var _ = Describe("Tests for getPolicyFileName", func() {
	Context("When controller namespace is same as policy namespace", func() {
		It("it should use only name for filename", func() {
			os.Setenv("APERTURE_CONTROLLER_NAMESPACE", appName)

			filename := getPolicyFileName(test, appName)

			Expect(filename).To(Equal("test.yaml"))
		})
	})

	Context("When controller namespace is not same as policy namespace", func() {
		It("it should use name and namespace for filename", func() {
			filename := getPolicyFileName(test, test)

			Expect(filename).To(Equal("test-test.yaml"))
		})
	})
})

// checkEtcdEndpoints generates endpoints list based on the release name if that is not provided else returns the provided values.
func checkEtcdEndpoints(etcd etcd.EtcdConfig, name, namespace string) etcd.EtcdConfig {
	endpoints := []string{}
	if etcd.Endpoints != nil {
		for _, endpoint := range etcd.Endpoints {
			if endpoint != "" {
				endpoints = append(endpoints, endpoint)
			}
		}
	}

	if len(endpoints) == 0 {
		endpoints = append(endpoints, fmt.Sprintf("http://%s-etcd.%s:2379", name, namespace))
	}

	etcd.Endpoints = endpoints
	return etcd
}

// checkPrometheusAddress generates prometheus address based on the release name if that is not provided else returns the provided value.
func checkPrometheusAddress(address, name, namespace string) string {
	if address == "" {
		address = fmt.Sprintf("http://%s-prometheus-server.%s:80", name, namespace)
	}
	return strings.TrimRight(address, "/")
}
