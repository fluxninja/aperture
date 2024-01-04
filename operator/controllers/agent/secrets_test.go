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
	"fmt"

	. "github.com/fluxninja/aperture/v2/operator/controllers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/utils/ptr"

	agentv1alpha1 "github.com/fluxninja/aperture/v2/operator/api/agent/v1alpha1"
	"github.com/fluxninja/aperture/v2/operator/api/common"
)

var _ = Describe("Secret for Agent", func() {
	Context("Instance with default parameters", func() {
		It("returns correct Secret", func() {
			instance := &agentv1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: agentv1alpha1.AgentSpec{
					CommonSpec: common.CommonSpec{
						Secrets: common.Secrets{
							FluxNinjaExtension: common.APIKeySecret{
								Value: Test,
							},
						},
					},
				},
			}

			expected := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-agent-apikey", AppName),
					Namespace: AppName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       AppName,
						"app.kubernetes.io/instance":   AppName,
						"app.kubernetes.io/managed-by": OperatorName,
						"app.kubernetes.io/component":  AgentServiceName,
					},
					Annotations: nil,
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion:         "fluxninja.com/v1alpha1",
							Name:               instance.GetName(),
							Kind:               "Agent",
							Controller:         ptr.To(true),
							BlockOwnerDeletion: ptr.To(true),
						},
					},
				},
				Data: map[string][]byte{
					"apiKey": []byte(Test),
				},
			}

			result, err := secretForAgentAPIKey(instance.DeepCopy(), scheme.Scheme)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(expected))
		})
	})

	Context("Instance with all parameters", func() {
		It("returns correct Secret", func() {
			instance := &agentv1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: agentv1alpha1.AgentSpec{
					CommonSpec: common.CommonSpec{
						Annotations: TestMap,
						Secrets: common.Secrets{
							FluxNinjaExtension: common.APIKeySecret{
								SecretKeyRef: common.SecretKeyRef{
									Name: Test,
									Key:  Test,
								},
								Value: Test,
							},
						},
					},
				},
			}

			expected := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      Test,
					Namespace: AppName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       AppName,
						"app.kubernetes.io/instance":   AppName,
						"app.kubernetes.io/managed-by": OperatorName,
						"app.kubernetes.io/component":  AgentServiceName,
					},
					Annotations: TestMap,
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion:         "fluxninja.com/v1alpha1",
							Name:               instance.GetName(),
							Kind:               "Agent",
							Controller:         ptr.To(true),
							BlockOwnerDeletion: ptr.To(true),
						},
					},
				},
				Data: map[string][]byte{
					Test: []byte(Test),
				},
			}

			result, err := secretForAgentAPIKey(instance.DeepCopy(), scheme.Scheme)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(expected))
		})
	})
})
