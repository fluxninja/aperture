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
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	"github.com/fluxninja/aperture/operator/api/v1alpha1"
)

var _ = Describe("MutatingWebhookConfiguration for Controller", func() {
	Context("Instance with all parameters", func() {
		It("returns correct MutatingWebhookConfiguration", func() {
			instance := &v1alpha1.Agent{
				TypeMeta: v1.TypeMeta{
					Kind:       "Agent",
					APIVersion: "fluxninja.com/v1alpha1",
				},
				ObjectMeta: v1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.AgentSpec{
					CommonSpec: v1alpha1.CommonSpec{
						Labels:      testMap,
						Annotations: testMapTwo,
					},
				},
			}

			os.Setenv("APERTURE_OPERATOR_CERT_DIR", certDir)
			os.Setenv("APERTURE_OPERATOR_CERT_NAME", "tls5.crt")
			os.Setenv("APERTURE_OPERATOR_NAMESPACE", appName)
			os.Setenv("APERTURE_OPERATOR_SERVICE_NAME", appName)
			certPath := fmt.Sprintf("%s/%s", os.Getenv("APERTURE_OPERATOR_CERT_DIR"), webhookClientCertName)
			serverCertPEM := new(bytes.Buffer)
			_ = pem.Encode(serverCertPEM, &pem.Block{
				Type:  "CERTIFICATE",
				Bytes: []byte(test),
			})
			err := writeFile(certPath, serverCertPEM)
			Expect(err).NotTo(HaveOccurred())

			expected := &admissionregistrationv1.MutatingWebhookConfiguration{
				ObjectMeta: v1.ObjectMeta{
					Name: podMutatingWebhookName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       appName,
						"app.kubernetes.io/instance":   appName,
						"app.kubernetes.io/managed-by": operatorName,
						"app.kubernetes.io/component":  operatorName,
						test:                           test,
					},
					Annotations: map[string]string{
						"fluxninja.com/primary-resource-type": "Agent.fluxninja.com",
						"fluxninja.com/primary-resource":      fmt.Sprintf("%s/%s", appName, appName),
						test:                                  test,
						testTwo:                               testTwo,
					},
				},
				Webhooks: []admissionregistrationv1.MutatingWebhook{
					{
						Name: sidecarKey,
						ClientConfig: admissionregistrationv1.WebhookClientConfig{
							Service: &admissionregistrationv1.ServiceReference{
								Name:      appName,
								Namespace: appName,
								Path:      pointer.StringPtr(MutatingWebhookURI),
								Port:      pointer.Int32(443),
							},
							CABundle: serverCertPEM.Bytes(),
						},
						NamespaceSelector: &v1.LabelSelector{
							MatchExpressions: []v1.LabelSelectorRequirement{
								{
									Key:      sidecarLabelKey,
									Operator: v1.LabelSelectorOpIn,
									Values:   []string{enabled},
								},
							},
						},
						Rules: []admissionregistrationv1.RuleWithOperations{
							{
								Operations: []admissionregistrationv1.OperationType{"CREATE"},
								Rule: admissionregistrationv1.Rule{
									APIGroups:   []string{""},
									APIVersions: []string{v1Version},
									Resources:   []string{"pods"},
									Scope:       &[]admissionregistrationv1.ScopeType{admissionregistrationv1.NamespacedScope}[0],
								},
							},
						},
						AdmissionReviewVersions: []string{v1Version},
						FailurePolicy:           &[]admissionregistrationv1.FailurePolicyType{admissionregistrationv1.Fail}[0],
						SideEffects:             &[]admissionregistrationv1.SideEffectClass{admissionregistrationv1.SideEffectClassNone}[0],
						TimeoutSeconds:          pointer.Int32Ptr(10),
					},
				},
			}

			result, err := podMutatingWebhookConfiguration(instance.DeepCopy())

			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(expected))

			os.Remove(certPath)
		})
	})
})

var _ = Describe("MutatingWebhookConfiguration for Agent CR", func() {
	Context("Instance with all parameters", func() {
		It("returns correct MutatingWebhookConfiguration", func() {
			os.Setenv("APERTURE_OPERATOR_CERT_DIR", certDir)
			os.Setenv("APERTURE_OPERATOR_CERT_NAME", "tls10.crt")
			os.Setenv("APERTURE_OPERATOR_NAMESPACE", appName)
			os.Setenv("APERTURE_OPERATOR_SERVICE_NAME", appName)
			certPath := fmt.Sprintf("%s/%s", os.Getenv("APERTURE_OPERATOR_CERT_DIR"), webhookClientCertName)
			serverCertPEM := new(bytes.Buffer)
			_ = pem.Encode(serverCertPEM, &pem.Block{
				Type:  "CERTIFICATE",
				Bytes: []byte(test),
			})
			err := writeFile(certPath, serverCertPEM)
			Expect(err).NotTo(HaveOccurred())

			expected := &admissionregistrationv1.MutatingWebhookConfiguration{
				ObjectMeta: v1.ObjectMeta{
					Name: agentMutatingWebhookName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       appName,
						"app.kubernetes.io/instance":   operatorName,
						"app.kubernetes.io/managed-by": operatorName,
						"app.kubernetes.io/component":  operatorName,
					},
				},
				Webhooks: []admissionregistrationv1.MutatingWebhook{
					{
						Name: "agent-defaulter.fluxninja.com",
						ClientConfig: admissionregistrationv1.WebhookClientConfig{
							Service: &admissionregistrationv1.ServiceReference{
								Name:      appName,
								Namespace: appName,
								Path:      pointer.StringPtr("/agent-defaulter"),
								Port:      pointer.Int32(443),
							},
							CABundle: serverCertPEM.Bytes(),
						},
						NamespaceSelector: &v1.LabelSelector{},
						Rules: []admissionregistrationv1.RuleWithOperations{
							{
								Operations: []admissionregistrationv1.OperationType{"CREATE", "UPDATE"},
								Rule: admissionregistrationv1.Rule{
									APIGroups:   []string{"fluxninja.com"},
									APIVersions: []string{v1Alpha1Version},
									Resources:   []string{"agents"},
									Scope:       &[]admissionregistrationv1.ScopeType{admissionregistrationv1.NamespacedScope}[0],
								},
							},
						},
						AdmissionReviewVersions: []string{v1Version},
						FailurePolicy:           &[]admissionregistrationv1.FailurePolicyType{admissionregistrationv1.Fail}[0],
						SideEffects:             &[]admissionregistrationv1.SideEffectClass{admissionregistrationv1.SideEffectClassNone}[0],
						TimeoutSeconds:          pointer.Int32Ptr(10),
					},
				},
			}

			result, err := agentMutatingWebhookConfiguration()

			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(expected))

			os.Remove(certPath)
		})
	})
})

var _ = Describe("MutatingWebhookConfiguration for Controller CR", func() {
	Context("Instance with all parameters", func() {
		It("returns correct MutatingWebhookConfiguration", func() {
			os.Setenv("APERTURE_OPERATOR_CERT_DIR", certDir)
			os.Setenv("APERTURE_OPERATOR_CERT_NAME", "tls11.crt")
			os.Setenv("APERTURE_OPERATOR_NAMESPACE", appName)
			os.Setenv("APERTURE_OPERATOR_SERVICE_NAME", appName)
			certPath := fmt.Sprintf("%s/%s", os.Getenv("APERTURE_OPERATOR_CERT_DIR"), webhookClientCertName)
			serverCertPEM := new(bytes.Buffer)
			_ = pem.Encode(serverCertPEM, &pem.Block{
				Type:  "CERTIFICATE",
				Bytes: []byte(test),
			})
			err := writeFile(certPath, serverCertPEM)
			Expect(err).NotTo(HaveOccurred())

			expected := &admissionregistrationv1.MutatingWebhookConfiguration{
				ObjectMeta: v1.ObjectMeta{
					Name: controllerMutatingWebhookName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       appName,
						"app.kubernetes.io/instance":   operatorName,
						"app.kubernetes.io/managed-by": operatorName,
						"app.kubernetes.io/component":  operatorName,
					},
				},
				Webhooks: []admissionregistrationv1.MutatingWebhook{
					{
						Name: "controller-defaulter.fluxninja.com",
						ClientConfig: admissionregistrationv1.WebhookClientConfig{
							Service: &admissionregistrationv1.ServiceReference{
								Name:      appName,
								Namespace: appName,
								Path:      pointer.StringPtr("/controller-defaulter"),
								Port:      pointer.Int32(443),
							},
							CABundle: serverCertPEM.Bytes(),
						},
						NamespaceSelector: &v1.LabelSelector{},
						Rules: []admissionregistrationv1.RuleWithOperations{
							{
								Operations: []admissionregistrationv1.OperationType{"CREATE", "UPDATE"},
								Rule: admissionregistrationv1.Rule{
									APIGroups:   []string{"fluxninja.com"},
									APIVersions: []string{v1Alpha1Version},
									Resources:   []string{"controllers", "policies"},
									Scope:       &[]admissionregistrationv1.ScopeType{admissionregistrationv1.NamespacedScope}[0],
								},
							},
						},
						AdmissionReviewVersions: []string{v1Version},
						FailurePolicy:           &[]admissionregistrationv1.FailurePolicyType{admissionregistrationv1.Fail}[0],
						SideEffects:             &[]admissionregistrationv1.SideEffectClass{admissionregistrationv1.SideEffectClassNone}[0],
						TimeoutSeconds:          pointer.Int32Ptr(10),
					},
				},
			}

			result, err := controllerMutatingWebhookConfiguration()

			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(expected))

			os.Remove(certPath)
		})
	})
})

var _ = Describe("Test MutatingWebhookConfiguration Mutate", func() {
	It("Mutate should update required fields only", func() {
		expected := &admissionregistrationv1.MutatingWebhookConfiguration{
			ObjectMeta: v1.ObjectMeta{},
			Webhooks: []admissionregistrationv1.MutatingWebhook{
				{
					Name:                    "cm-validator.fluxninja.com",
					AdmissionReviewVersions: testArray,
					ClientConfig: admissionregistrationv1.WebhookClientConfig{
						URL: &test,
					},
					NamespaceSelector: &v1.LabelSelector{
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

		mwc := &admissionregistrationv1.MutatingWebhookConfiguration{
			Webhooks: []admissionregistrationv1.MutatingWebhook{
				{
					Name: "cm-validator.fluxninja.com",
				},
			},
		}
		err := mutatingWebhookConfigurationMutate(mwc, expected.Webhooks)()

		Expect(err).NotTo(HaveOccurred())
		Expect(mwc).To(Equal(expected))
	})
})
