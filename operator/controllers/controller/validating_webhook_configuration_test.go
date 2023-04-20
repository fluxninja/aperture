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

	. "github.com/fluxninja/aperture/operator/controllers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	"github.com/fluxninja/aperture/operator/api/common"
	controllerv1alpha1 "github.com/fluxninja/aperture/operator/api/controller/v1alpha1"
	"github.com/fluxninja/aperture/pkg/net/listener"
)

var _ = Describe("ValidatingWebhookConfiguration for Controller", func() {
	Context("Instance with default parameters", func() {
		It("returns correct ValidatingWebhookConfiguration", func() {
			instance := &controllerv1alpha1.Controller{
				TypeMeta: v1.TypeMeta{
					Kind:       "Controller",
					APIVersion: "fluxninja.com/v1alpha1",
				},
				ObjectMeta: v1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: controllerv1alpha1.ControllerSpec{
					ConfigSpec: controllerv1alpha1.ControllerConfigSpec{
						CommonConfigSpec: common.CommonConfigSpec{
							Server: common.ServerConfigSpec{
								Listener: listener.ListenerConfig{
									Addr: ":8080",
								},
							},
						},
					},
				},
			}

			expected := &admissionregistrationv1.ValidatingWebhookConfiguration{
				ObjectMeta: v1.ObjectMeta{
					Name: ControllerServiceName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       AppName,
						"app.kubernetes.io/instance":   AppName,
						"app.kubernetes.io/managed-by": OperatorName,
						"app.kubernetes.io/component":  ControllerServiceName,
					},
					Annotations: map[string]string{
						"fluxninja.com/primary-resource-type": "Controller.fluxninja.com",
						"fluxninja.com/primary-resource":      fmt.Sprintf("%s/%s", AppName, AppName),
					},
				},
				Webhooks: []admissionregistrationv1.ValidatingWebhook{
					{
						Name: PolicyValidatingWebhookName,
						ClientConfig: admissionregistrationv1.WebhookClientConfig{
							Service: &admissionregistrationv1.ServiceReference{
								Name:      ControllerServiceName,
								Namespace: instance.GetNamespace(),
								Path:      pointer.String(PolicyValidatingWebhookURI),
								Port:      pointer.Int32(8080),
							},
							CABundle: []byte(Test),
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
									APIGroups:   []string{"fluxninja.com"},
									APIVersions: []string{V1Alpha1Version},
									Resources:   []string{"policies"},
									Scope:       &[]admissionregistrationv1.ScopeType{admissionregistrationv1.NamespacedScope}[0],
								},
							},
						},
						AdmissionReviewVersions: []string{V1Version},
						FailurePolicy:           &[]admissionregistrationv1.FailurePolicyType{admissionregistrationv1.Fail}[0],
						SideEffects:             &[]admissionregistrationv1.SideEffectClass{admissionregistrationv1.SideEffectClassNone}[0],
						TimeoutSeconds:          pointer.Int32(10),
					},
				},
			}

			result := validatingWebhookConfiguration(instance.DeepCopy(), []byte(Test))
			Expect(result).To(Equal(expected))
		})
	})

	Context("Instance with all parameters", func() {
		It("returns correct ValidatingWebhookConfiguration", func() {
			instance := &controllerv1alpha1.Controller{
				TypeMeta: v1.TypeMeta{
					Kind:       "Controller",
					APIVersion: "fluxninja.com/v1alpha1",
				},
				ObjectMeta: v1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: controllerv1alpha1.ControllerSpec{
					CommonSpec: common.CommonSpec{
						Labels:      TestMap,
						Annotations: TestMapTwo,
					},
					ConfigSpec: controllerv1alpha1.ControllerConfigSpec{
						CommonConfigSpec: common.CommonConfigSpec{
							Server: common.ServerConfigSpec{
								Listener: listener.ListenerConfig{
									Addr: ":80",
								},
							},
						},
					},
				},
			}

			expected := &admissionregistrationv1.ValidatingWebhookConfiguration{
				ObjectMeta: v1.ObjectMeta{
					Name: ControllerServiceName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       AppName,
						"app.kubernetes.io/instance":   AppName,
						"app.kubernetes.io/managed-by": OperatorName,
						"app.kubernetes.io/component":  ControllerServiceName,
						Test:                           Test,
					},
					Annotations: map[string]string{
						"fluxninja.com/primary-resource-type": "Controller.fluxninja.com",
						"fluxninja.com/primary-resource":      fmt.Sprintf("%s/%s", AppName, AppName),
						Test:                                  Test,
						TestTwo:                               TestTwo,
					},
				},
				Webhooks: []admissionregistrationv1.ValidatingWebhook{
					{
						Name: PolicyValidatingWebhookName,
						ClientConfig: admissionregistrationv1.WebhookClientConfig{
							Service: &admissionregistrationv1.ServiceReference{
								Name:      ControllerServiceName,
								Namespace: instance.GetNamespace(),
								Path:      pointer.String(PolicyValidatingWebhookURI),
								Port:      pointer.Int32(80),
							},
							CABundle: []byte(Test),
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
									APIGroups:   []string{"fluxninja.com"},
									APIVersions: []string{V1Alpha1Version},
									Resources:   []string{"policies"},
									Scope:       &[]admissionregistrationv1.ScopeType{admissionregistrationv1.NamespacedScope}[0],
								},
							},
						},
						AdmissionReviewVersions: []string{V1Version},
						FailurePolicy:           &[]admissionregistrationv1.FailurePolicyType{admissionregistrationv1.Fail}[0],
						SideEffects:             &[]admissionregistrationv1.SideEffectClass{admissionregistrationv1.SideEffectClassNone}[0],
						TimeoutSeconds:          pointer.Int32(10),
					},
				},
			}

			result := validatingWebhookConfiguration(instance.DeepCopy(), []byte(Test))
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
					Name:                    PolicyValidatingWebhookName,
					AdmissionReviewVersions: TestArray,
					ClientConfig: admissionregistrationv1.WebhookClientConfig{
						URL: &Test,
					},
					ObjectSelector: &v1.LabelSelector{
						MatchLabels: TestMap,
					},
					Rules: []admissionregistrationv1.RuleWithOperations{
						{
							Rule: admissionregistrationv1.Rule{
								APIGroups: TestArray,
							},
						},
					},
					FailurePolicy:  &[]admissionregistrationv1.FailurePolicyType{admissionregistrationv1.Ignore}[0],
					SideEffects:    &[]admissionregistrationv1.SideEffectClass{admissionregistrationv1.SideEffectClassSome}[0],
					TimeoutSeconds: pointer.Int32(10),
				},
			},
		}

		vwc := &admissionregistrationv1.ValidatingWebhookConfiguration{
			Webhooks: []admissionregistrationv1.ValidatingWebhook{
				{
					Name: PolicyValidatingWebhookName,
				},
			},
		}

		err := ValidatingWebhookConfigurationMutate(vwc, expected.Webhooks)()

		Expect(err).NotTo(HaveOccurred())
		Expect(vwc).To(Equal(expected))
	})
})
