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
								ListenerConfig: listener.ListenerConfig{
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
						Name: "cm-validator.fluxninja.com",
						ClientConfig: admissionregistrationv1.WebhookClientConfig{
							Service: &admissionregistrationv1.ServiceReference{
								Name:      ControllerServiceName,
								Namespace: instance.GetNamespace(),
								Path:      pointer.StringPtr("/validate/configmap"),
								Port:      pointer.Int32(8080),
							},
							CABundle: []byte(Test),
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
									APIVersions: []string{V1Version},
									Resources:   []string{"configmaps"},
									Scope:       &[]admissionregistrationv1.ScopeType{admissionregistrationv1.NamespacedScope}[0],
								},
							},
						},
						AdmissionReviewVersions: []string{V1Version},
						FailurePolicy:           &[]admissionregistrationv1.FailurePolicyType{admissionregistrationv1.Fail}[0],
						SideEffects:             &[]admissionregistrationv1.SideEffectClass{admissionregistrationv1.SideEffectClassNone}[0],
						TimeoutSeconds:          pointer.Int32Ptr(5),
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
								ListenerConfig: listener.ListenerConfig{
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
						Name: "cm-validator.fluxninja.com",
						ClientConfig: admissionregistrationv1.WebhookClientConfig{
							Service: &admissionregistrationv1.ServiceReference{
								Name:      ControllerServiceName,
								Namespace: instance.GetNamespace(),
								Path:      pointer.StringPtr("/validate/configmap"),
								Port:      pointer.Int32(80),
							},
							CABundle: []byte(Test),
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
									APIVersions: []string{V1Version},
									Resources:   []string{"configmaps"},
									Scope:       &[]admissionregistrationv1.ScopeType{admissionregistrationv1.NamespacedScope}[0],
								},
							},
						},
						AdmissionReviewVersions: []string{V1Version},
						FailurePolicy:           &[]admissionregistrationv1.FailurePolicyType{admissionregistrationv1.Fail}[0],
						SideEffects:             &[]admissionregistrationv1.SideEffectClass{admissionregistrationv1.SideEffectClassNone}[0],
						TimeoutSeconds:          pointer.Int32Ptr(5),
					},
				},
			}

			result := validatingWebhookConfiguration(instance.DeepCopy(), []byte(Test))
			Expect(result).To(Equal(expected))
		})
	})
})
