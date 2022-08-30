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
	"encoding/base64"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/utils/pointer"

	"github.com/fluxninja/aperture/operator/api/v1alpha1"
)

var _ = Describe("Secret for Agent", func() {
	Context("Instance with default parameters", func() {
		It("returns correct Secret", func() {
			instance := &v1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.AgentSpec{
					CommonSpec: v1alpha1.CommonSpec{
						FluxNinjaPlugin: v1alpha1.FluxNinjaPluginSpec{
							APIKeySecret: v1alpha1.APIKeySecret{
								Value: test,
							},
						},
					},
				},
			}

			expected := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-agent-apikey", appName),
					Namespace: appName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       appName,
						"app.kubernetes.io/instance":   appName,
						"app.kubernetes.io/managed-by": operatorName,
						"app.kubernetes.io/component":  agentServiceName,
					},
					Annotations: nil,
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion:         "fluxninja.com/v1alpha1",
							Name:               instance.GetName(),
							Kind:               "Agent",
							Controller:         pointer.BoolPtr(true),
							BlockOwnerDeletion: pointer.BoolPtr(true),
						},
					},
				},
				Data: map[string][]byte{
					"apiKey": []byte(test),
				},
			}

			result, _ := secretForAgentAPIKey(instance.DeepCopy(), scheme.Scheme)
			Expect(result).To(Equal(expected))
		})
	})

	Context("Instance with all parameters", func() {
		It("returns correct Secret", func() {
			instance := &v1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.AgentSpec{
					CommonSpec: v1alpha1.CommonSpec{
						Annotations: testMap,
						FluxNinjaPlugin: v1alpha1.FluxNinjaPluginSpec{
							APIKeySecret: v1alpha1.APIKeySecret{
								SecretKeyRef: v1alpha1.SecretKeyRef{
									Name: test,
									Key:  test,
								},
								Value: test,
							},
						},
					},
				},
			}

			expected := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      test,
					Namespace: appName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       appName,
						"app.kubernetes.io/instance":   appName,
						"app.kubernetes.io/managed-by": operatorName,
						"app.kubernetes.io/component":  agentServiceName,
					},
					Annotations: testMap,
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion:         "fluxninja.com/v1alpha1",
							Name:               instance.GetName(),
							Kind:               "Agent",
							Controller:         pointer.BoolPtr(true),
							BlockOwnerDeletion: pointer.BoolPtr(true),
						},
					},
				},
				Data: map[string][]byte{
					test: []byte(test),
				},
			}

			result, _ := secretForAgentAPIKey(instance.DeepCopy(), scheme.Scheme)
			Expect(result).To(Equal(expected))
		})
	})
})

var _ = Describe("Secret for Controller", func() {
	Context("Instance with default parameters", func() {
		It("returns correct Secret", func() {
			instance := &v1alpha1.Controller{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.ControllerSpec{
					CommonSpec: v1alpha1.CommonSpec{
						FluxNinjaPlugin: v1alpha1.FluxNinjaPluginSpec{
							APIKeySecret: v1alpha1.APIKeySecret{
								Value: test,
							},
						},
					},
				},
			}

			expected := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-controller-apikey", appName),
					Namespace: appName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       appName,
						"app.kubernetes.io/instance":   appName,
						"app.kubernetes.io/managed-by": operatorName,
						"app.kubernetes.io/component":  controllerServiceName,
					},
					Annotations: nil,
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion:         "fluxninja.com/v1alpha1",
							Name:               instance.GetName(),
							Kind:               "Controller",
							Controller:         pointer.BoolPtr(true),
							BlockOwnerDeletion: pointer.BoolPtr(true),
						},
					},
				},
				Data: map[string][]byte{
					"apiKey": []byte(test),
				},
			}

			result, _ := secretForControllerAPIKey(instance.DeepCopy(), scheme.Scheme)
			Expect(result).To(Equal(expected))
		})
	})

	Context("Instance with all parameters", func() {
		It("returns correct Secret", func() {
			instance := &v1alpha1.Controller{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.ControllerSpec{
					CommonSpec: v1alpha1.CommonSpec{
						Annotations: testMap,
						FluxNinjaPlugin: v1alpha1.FluxNinjaPluginSpec{
							APIKeySecret: v1alpha1.APIKeySecret{
								SecretKeyRef: v1alpha1.SecretKeyRef{
									Name: test,
									Key:  test,
								},
								Value: test,
							},
						},
					},
				},
			}

			expected := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      test,
					Namespace: appName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       appName,
						"app.kubernetes.io/instance":   appName,
						"app.kubernetes.io/managed-by": operatorName,
						"app.kubernetes.io/component":  controllerServiceName,
					},
					Annotations: testMap,
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion:         "fluxninja.com/v1alpha1",
							Name:               instance.GetName(),
							Kind:               "Controller",
							Controller:         pointer.BoolPtr(true),
							BlockOwnerDeletion: pointer.BoolPtr(true),
						},
					},
				},
				Data: map[string][]byte{
					test: []byte(test),
				},
			}

			result, _ := secretForControllerAPIKey(instance.DeepCopy(), scheme.Scheme)
			Expect(result).To(Equal(expected))
		})
	})
})

var _ = Describe("Secret for Controller Cert", func() {

	Context("Instance with default parameters", func() {
		It("returns correct Secret", func() {
			instance := &v1alpha1.Controller{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.ControllerSpec{},
			}

			expected := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-controller-cert", appName),
					Namespace: appName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       appName,
						"app.kubernetes.io/instance":   appName,
						"app.kubernetes.io/managed-by": operatorName,
						"app.kubernetes.io/component":  controllerServiceName,
					},
					Annotations: nil,
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion:         "fluxninja.com/v1alpha1",
							Name:               instance.GetName(),
							Kind:               "Controller",
							Controller:         pointer.BoolPtr(true),
							BlockOwnerDeletion: pointer.BoolPtr(true),
						},
					},
				},
				Data: map[string][]byte{
					controllerCertName:    []byte(test),
					controllerCertKeyName: []byte(test),
				},
			}

			result, _ := secretForControllerCert(instance.DeepCopy(), scheme.Scheme, bytes.NewBuffer([]byte(test)), bytes.NewBuffer([]byte(test)))
			Expect(result).To(Equal(expected))
		})
	})

	Context("Instance with all parameters", func() {
		It("returns correct Secret", func() {
			instance := &v1alpha1.Controller{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.ControllerSpec{
					CommonSpec: v1alpha1.CommonSpec{
						Annotations: testMap,
					},
				},
			}

			expected := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-controller-cert", appName),
					Namespace: appName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       appName,
						"app.kubernetes.io/instance":   appName,
						"app.kubernetes.io/managed-by": operatorName,
						"app.kubernetes.io/component":  controllerServiceName,
					},
					Annotations: testMap,
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion:         "fluxninja.com/v1alpha1",
							Name:               instance.GetName(),
							Kind:               "Controller",
							Controller:         pointer.BoolPtr(true),
							BlockOwnerDeletion: pointer.BoolPtr(true),
						},
					},
				},
				Data: map[string][]byte{
					controllerCertName:    []byte(test),
					controllerCertKeyName: []byte(test),
				},
			}

			result, _ := secretForControllerCert(instance.DeepCopy(), scheme.Scheme, bytes.NewBuffer([]byte(test)), bytes.NewBuffer([]byte(test)))
			Expect(result).To(Equal(expected))
		})
	})
})

var _ = Describe("Test Secret Mutate", func() {
	It("Mutate should update required fields only", func() {
		expected := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{},
			Data: map[string][]byte{
				test: []byte(base64.StdEncoding.EncodeToString([]byte(test))),
			},
		}

		secret := &corev1.Secret{}
		err := secretMutate(secret, expected.Data)()
		Expect(err).NotTo(HaveOccurred())
		Expect(secret).To(Equal(expected))
	})
})
