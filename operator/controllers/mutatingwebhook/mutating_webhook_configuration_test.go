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

package mutatingwebhook

import (
	"bytes"
	"encoding/pem"
	"fmt"
	"os"

	. "github.com/fluxninja/aperture/v2/operator/controllers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

var _ = Describe("MutatingWebhookConfiguration for Agent CR", func() {
	Context("Instance with all parameters", func() {
		It("returns correct MutatingWebhookConfiguration", func() {
			os.Setenv("APERTURE_OPERATOR_CERT_DIR", CertDir)
			os.Setenv("APERTURE_OPERATOR_CERT_NAME", "tls10.crt")
			os.Setenv("APERTURE_OPERATOR_NAMESPACE", AppName)
			os.Setenv("APERTURE_OPERATOR_SERVICE_NAME", AppName)
			certPath := fmt.Sprintf("%s/%s", os.Getenv("APERTURE_OPERATOR_CERT_DIR"), WebhookClientCertName)
			serverCertPEM := new(bytes.Buffer)
			_ = pem.Encode(serverCertPEM, &pem.Block{
				Type:  "CERTIFICATE",
				Bytes: []byte(Test),
			})
			err := WriteFile(certPath, serverCertPEM)
			Expect(err).NotTo(HaveOccurred())

			expected := &admissionregistrationv1.MutatingWebhookConfiguration{
				ObjectMeta: v1.ObjectMeta{
					Name: AgentMutatingWebhookName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       AppName,
						"app.kubernetes.io/instance":   OperatorName,
						"app.kubernetes.io/managed-by": OperatorName,
						"app.kubernetes.io/component":  OperatorName,
					},
				},
				Webhooks: []admissionregistrationv1.MutatingWebhook{
					{
						Name: "agent-defaulter.fluxninja.com",
						ClientConfig: admissionregistrationv1.WebhookClientConfig{
							Service: &admissionregistrationv1.ServiceReference{
								Name:      AppName,
								Namespace: AppName,
								Path:      pointer.String("/agent-defaulter"),
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
									APIVersions: []string{V1Alpha1Version},
									Resources:   []string{"agents"},
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
			os.Setenv("APERTURE_OPERATOR_CERT_DIR", CertDir)
			os.Setenv("APERTURE_OPERATOR_CERT_NAME", "tls11.crt")
			os.Setenv("APERTURE_OPERATOR_NAMESPACE", AppName)
			os.Setenv("APERTURE_OPERATOR_SERVICE_NAME", AppName)
			certPath := fmt.Sprintf("%s/%s", os.Getenv("APERTURE_OPERATOR_CERT_DIR"), WebhookClientCertName)
			serverCertPEM := new(bytes.Buffer)
			_ = pem.Encode(serverCertPEM, &pem.Block{
				Type:  "CERTIFICATE",
				Bytes: []byte(Test),
			})
			err := WriteFile(certPath, serverCertPEM)
			Expect(err).NotTo(HaveOccurred())

			expected := &admissionregistrationv1.MutatingWebhookConfiguration{
				ObjectMeta: v1.ObjectMeta{
					Name: ControllerMutatingWebhookName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       AppName,
						"app.kubernetes.io/instance":   OperatorName,
						"app.kubernetes.io/managed-by": OperatorName,
						"app.kubernetes.io/component":  OperatorName,
					},
				},
				Webhooks: []admissionregistrationv1.MutatingWebhook{
					{
						Name: "controller-defaulter.fluxninja.com",
						ClientConfig: admissionregistrationv1.WebhookClientConfig{
							Service: &admissionregistrationv1.ServiceReference{
								Name:      AppName,
								Namespace: AppName,
								Path:      pointer.String("/controller-defaulter"),
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
									APIVersions: []string{V1Alpha1Version},
									Resources:   []string{"controllers"},
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

			result, err := controllerMutatingWebhookConfiguration()

			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(expected))

			os.Remove(certPath)
		})
	})
})
