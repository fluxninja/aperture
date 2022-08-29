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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	"github.com/fluxninja/aperture/operator/api/v1alpha1"
)

var _ = Describe("ValidatingWebhookConfiguration for Controller", func() {
	Context("Instance with default parameters", func() {
		It("returns correct ValidatingWebhookConfiguration", func() {
			instance := &v1alpha1.Controller{
				TypeMeta: v1.TypeMeta{
					Kind:       "Controller",
					APIVersion: "fluxninja.com/v1alpha1",
				},
				ObjectMeta: v1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
			}

			expected := &admissionregistrationv1.ValidatingWebhookConfiguration{
				ObjectMeta: v1.ObjectMeta{
					Name: validatingWebhookServiceName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       appName,
						"app.kubernetes.io/instance":   appName,
						"app.kubernetes.io/managed-by": operatorName,
						"app.kubernetes.io/component":  controllerServiceName,
					},
					Annotations: map[string]string{
						"fluxninja.com/primary-resource-type": "Controller.fluxninja.com",
						"fluxninja.com/primary-resource":      fmt.Sprintf("%s/%s", appName, appName),
					},
				},
				Webhooks: []admissionregistrationv1.ValidatingWebhook{
					{
						Name: "cm-validator.fluxninja.com",
						ClientConfig: admissionregistrationv1.WebhookClientConfig{
							Service: &admissionregistrationv1.ServiceReference{
								Name:      validatingWebhookServiceName,
								Namespace: instance.GetNamespace(),
								Path:      pointer.StringPtr("/validate/configmap"),
								Port:      pointer.Int32(443),
							},
							CABundle: []byte(test),
						},
						NamespaceSelector: &v1.LabelSelector{
							MatchLabels: map[string]string{
								"kubernetes.io/metadata.name": instance.GetNamespace(),
							},
						},
						ObjectSelector: &v1.LabelSelector{
							MatchLabels: map[string]string{
								"fluxninja.com/validate": "true",
							},
						},
						Rules: []admissionregistrationv1.RuleWithOperations{
							{
								Operations: []admissionregistrationv1.OperationType{"CREATE", "UPDATE"},
								Rule: admissionregistrationv1.Rule{
									APIGroups:   []string{""},
									APIVersions: []string{v1Version},
									Resources:   []string{"configmaps"},
									Scope:       &[]admissionregistrationv1.ScopeType{admissionregistrationv1.NamespacedScope}[0],
								},
							},
						},
						AdmissionReviewVersions: []string{v1Version},
						FailurePolicy:           &[]admissionregistrationv1.FailurePolicyType{admissionregistrationv1.Fail}[0],
						SideEffects:             &[]admissionregistrationv1.SideEffectClass{admissionregistrationv1.SideEffectClassNone}[0],
						TimeoutSeconds:          pointer.Int32Ptr(5),
					},
				},
			}

			result := validatingWebhookConfiguration(instance.DeepCopy(), []byte(test))
			Expect(result).To(Equal(expected))
		})
	})

	Context("Instance with all parameters", func() {
		It("returns correct ValidatingWebhookConfiguration", func() {
			instance := &v1alpha1.Controller{
				TypeMeta: v1.TypeMeta{
					Kind:       "Controller",
					APIVersion: "fluxninja.com/v1alpha1",
				},
				ObjectMeta: v1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.ControllerSpec{
					CommonSpec: v1alpha1.CommonSpec{
						Labels:      testMap,
						Annotations: testMapTwo,
					},
				},
			}

			expected := &admissionregistrationv1.ValidatingWebhookConfiguration{
				ObjectMeta: v1.ObjectMeta{
					Name: validatingWebhookServiceName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       appName,
						"app.kubernetes.io/instance":   appName,
						"app.kubernetes.io/managed-by": operatorName,
						"app.kubernetes.io/component":  controllerServiceName,
						test:                           test,
					},
					Annotations: map[string]string{
						"fluxninja.com/primary-resource-type": "Controller.fluxninja.com",
						"fluxninja.com/primary-resource":      fmt.Sprintf("%s/%s", appName, appName),
						test:                                  test,
						testTwo:                               testTwo,
					},
				},
				Webhooks: []admissionregistrationv1.ValidatingWebhook{
					{
						Name: "cm-validator.fluxninja.com",
						ClientConfig: admissionregistrationv1.WebhookClientConfig{
							Service: &admissionregistrationv1.ServiceReference{
								Name:      validatingWebhookServiceName,
								Namespace: instance.GetNamespace(),
								Path:      pointer.StringPtr("/validate/configmap"),
								Port:      pointer.Int32(443),
							},
							CABundle: []byte(test),
						},
						NamespaceSelector: &v1.LabelSelector{
							MatchLabels: map[string]string{
								"kubernetes.io/metadata.name": instance.GetNamespace(),
							},
						},
						ObjectSelector: &v1.LabelSelector{
							MatchLabels: map[string]string{
								"fluxninja.com/validate": "true",
							},
						},
						Rules: []admissionregistrationv1.RuleWithOperations{
							{
								Operations: []admissionregistrationv1.OperationType{"CREATE", "UPDATE"},
								Rule: admissionregistrationv1.Rule{
									APIGroups:   []string{""},
									APIVersions: []string{v1Version},
									Resources:   []string{"configmaps"},
									Scope:       &[]admissionregistrationv1.ScopeType{admissionregistrationv1.NamespacedScope}[0],
								},
							},
						},
						AdmissionReviewVersions: []string{v1Version},
						FailurePolicy:           &[]admissionregistrationv1.FailurePolicyType{admissionregistrationv1.Fail}[0],
						SideEffects:             &[]admissionregistrationv1.SideEffectClass{admissionregistrationv1.SideEffectClassNone}[0],
						TimeoutSeconds:          pointer.Int32Ptr(5),
					},
				},
			}

			result := validatingWebhookConfiguration(instance.DeepCopy(), []byte(test))
			Expect(result).To(Equal(expected))
		})
	})
})

var _ = Describe("Test ValidatingWebhookConfiguration Mutate", func() {
	It("Mutate should update required fields only", func() {
		expected := &admissionregistrationv1.ValidatingWebhookConfiguration{
			ObjectMeta: v1.ObjectMeta{},
			Webhooks: []admissionregistrationv1.ValidatingWebhook{
				{
					Name:                    "cm-validator.fluxninja.com",
					AdmissionReviewVersions: testArray,
					ClientConfig: admissionregistrationv1.WebhookClientConfig{
						URL: &test,
					},
					NamespaceSelector: &v1.LabelSelector{
						MatchLabels: testMap,
					},
					ObjectSelector: &v1.LabelSelector{
						MatchLabels: testMap,
					},
					Rules: []admissionregistrationv1.RuleWithOperations{
						{
							Rule: admissionregistrationv1.Rule{
								APIGroups: testArray,
							},
						},
					},
					FailurePolicy:  &[]admissionregistrationv1.FailurePolicyType{admissionregistrationv1.Ignore}[0],
					SideEffects:    &[]admissionregistrationv1.SideEffectClass{admissionregistrationv1.SideEffectClassSome}[0],
					TimeoutSeconds: pointer.Int32Ptr(10),
				},
			},
		}

		vwc := &admissionregistrationv1.ValidatingWebhookConfiguration{
			Webhooks: []admissionregistrationv1.ValidatingWebhook{
				{
					Name: "cm-validator.fluxninja.com",
				},
			},
		}

		err := validatingWebhookConfigurationMutate(vwc, expected.Webhooks)()

		Expect(err).NotTo(HaveOccurred())
		Expect(vwc).To(Equal(expected))
	})
})
