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

package agent

import (
	"bytes"
	"encoding/pem"
	"fmt"
	"os"

	. "github.com/fluxninja/aperture/operator/controllers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	agentv1alpha1 "github.com/fluxninja/aperture/operator/api/agent/v1alpha1"
	"github.com/fluxninja/aperture/operator/api/common"
)

var _ = Describe("MutatingWebhookConfiguration for Pods", func() {
	Context("Instance with all parameters", func() {
		It("returns correct MutatingWebhookConfiguration", func() {
			instance := &agentv1alpha1.Agent{
				TypeMeta: v1.TypeMeta{
					Kind:       "Agent",
					APIVersion: "fluxninja.com/v1alpha1",
				},
				ObjectMeta: v1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: agentv1alpha1.AgentSpec{
					CommonSpec: common.CommonSpec{
						Labels:      TestMap,
						Annotations: TestMapTwo,
					},
				},
			}

			os.Setenv("APERTURE_OPERATOR_CERT_DIR", CertDir)
			os.Setenv("APERTURE_OPERATOR_CERT_NAME", "tls5.crt")
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
					Name: PodMutatingWebhookName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       AppName,
						"app.kubernetes.io/instance":   AppName,
						"app.kubernetes.io/managed-by": OperatorName,
						"app.kubernetes.io/component":  OperatorName,
						Test:                           Test,
					},
					Annotations: map[string]string{
						"fluxninja.com/primary-resource-type": "Agent.fluxninja.com",
						"fluxninja.com/primary-resource":      fmt.Sprintf("%s/%s", AppName, AppName),
						Test:                                  Test,
						TestTwo:                               TestTwo,
					},
				},
				Webhooks: []admissionregistrationv1.MutatingWebhook{
					{
						Name: SidecarKey,
						ClientConfig: admissionregistrationv1.WebhookClientConfig{
							Service: &admissionregistrationv1.ServiceReference{
								Name:      AppName,
								Namespace: AppName,
								Path:      pointer.StringPtr(MutatingWebhookURI),
								Port:      pointer.Int32(443),
							},
							CABundle: serverCertPEM.Bytes(),
						},
						NamespaceSelector: &v1.LabelSelector{
							MatchExpressions: []v1.LabelSelectorRequirement{
								{
									Key:      SidecarLabelKey,
									Operator: v1.LabelSelectorOpIn,
									Values:   []string{Enabled},
								},
							},
						},
						Rules: []admissionregistrationv1.RuleWithOperations{
							{
								Operations: []admissionregistrationv1.OperationType{"CREATE"},
								Rule: admissionregistrationv1.Rule{
									APIGroups:   []string{""},
									APIVersions: []string{V1Version},
									Resources:   []string{"pods"},
									Scope:       &[]admissionregistrationv1.ScopeType{admissionregistrationv1.NamespacedScope}[0],
								},
							},
						},
						AdmissionReviewVersions: []string{V1Version},
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
