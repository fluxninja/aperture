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
	"encoding/base64"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

var _ = Describe("Test ClusterRole Mutate", func() {
	It("Mutate should update required fields only", func() {
		expected := &rbacv1.ClusterRole{
			ObjectMeta: metav1.ObjectMeta{},
			Rules: []rbacv1.PolicyRule{
				{
					APIGroups: TestArray,
					Resources: TestArray,
					Verbs:     TestArray,
				},
			},
		}

		cr := &rbacv1.ClusterRole{}
		err := ClusterRoleMutate(cr, expected.Rules)()
		Expect(err).NotTo(HaveOccurred())
		Expect(cr).To(Equal(expected))
	})
})

var _ = Describe("Test ClusterRoleBinding Mutate", func() {
	It("Mutate should update required fields only", func() {
		expected := &rbacv1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{},
			RoleRef: rbacv1.RoleRef{
				APIGroup: Test,
				Kind:     Test,
				Name:     Test,
			},
		}

		crb := &rbacv1.ClusterRoleBinding{}
		err := ClusterRoleBindingMutate(crb, expected.RoleRef, expected.Subjects)()
		Expect(err).NotTo(HaveOccurred())
		Expect(crb).To(Equal(expected))
	})
})

var _ = Describe("Test ConfigMap Mutate", func() {
	It("Mutate should update required fields only", func() {
		expected := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{},
			Data:       TestMap,
		}

		cm := &corev1.ConfigMap{}
		err := ConfigMapMutate(cm, expected.Data)()

		Expect(err).NotTo(HaveOccurred())
		Expect(cm).To(Equal(expected))
	})
})

var _ = Describe("Test MutatingWebhookConfiguration Mutate", func() {
	It("Mutate should update required fields only", func() {
		expected := &admissionregistrationv1.MutatingWebhookConfiguration{
			ObjectMeta: metav1.ObjectMeta{},
			Webhooks: []admissionregistrationv1.MutatingWebhook{
				{
					Name:                    "cm-validator.fluxninja.com",
					AdmissionReviewVersions: TestArray,
					ClientConfig: admissionregistrationv1.WebhookClientConfig{
						URL: &Test,
					},
					NamespaceSelector: &metav1.LabelSelector{
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

		mwc := &admissionregistrationv1.MutatingWebhookConfiguration{
			Webhooks: []admissionregistrationv1.MutatingWebhook{
				{
					Name: "cm-validator.fluxninja.com",
				},
			},
		}
		err := MutatingWebhookConfigurationMutate(mwc, expected.Webhooks)()

		Expect(err).NotTo(HaveOccurred())
		Expect(mwc).To(Equal(expected))
	})
})

var _ = Describe("Test Secret Mutate", func() {
	It("Mutate should update required fields only", func() {
		expected := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{},
			Data: map[string][]byte{
				Test: []byte(base64.StdEncoding.EncodeToString([]byte(Test))),
			},
		}

		secret := &corev1.Secret{}
		err := SecretMutate(secret, expected.Data)()

		Expect(err).NotTo(HaveOccurred())
		Expect(secret).To(Equal(expected))
	})
})

var _ = Describe("Test ServiceAccount Mutate", func() {
	It("Mutate should update required fields only", func() {
		expected := &corev1.ServiceAccount{
			ObjectMeta:                   metav1.ObjectMeta{},
			AutomountServiceAccountToken: pointer.Bool(true),
		}

		sa := &corev1.ServiceAccount{}
		err := ServiceAccountMutate(sa, expected.AutomountServiceAccountToken)()

		Expect(err).NotTo(HaveOccurred())
		Expect(sa).To(Equal(expected))
	})
})

var _ = Describe("Test Service Mutate", func() {
	It("Mutate should update required fields only", func() {
		expected := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{},
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Name: Test,
						Port: 80,
					},
				},
				Selector: TestMap,
			},
		}

		svc := &corev1.Service{}
		err := ServiceMutate(svc, expected.Spec)()

		Expect(err).NotTo(HaveOccurred())
		Expect(svc).To(Equal(expected))
	})
})

var _ = Describe("Test ValidatingWebhookConfiguration Mutate", func() {
	It("Mutate should update required fields only", func() {
		expected := &admissionregistrationv1.ValidatingWebhookConfiguration{
			ObjectMeta: metav1.ObjectMeta{},
			Webhooks: []admissionregistrationv1.ValidatingWebhook{
				{
					Name:                    "cm-validator.fluxninja.com",
					AdmissionReviewVersions: TestArray,
					ClientConfig: admissionregistrationv1.WebhookClientConfig{
						URL: &Test,
					},
					NamespaceSelector: &metav1.LabelSelector{
						MatchLabels: TestMap,
					},
					ObjectSelector: &metav1.LabelSelector{
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
					Name: "cm-validator.fluxninja.com",
				},
			},
		}

		err := ValidatingWebhookConfigurationMutate(vwc, expected.Webhooks)()

		Expect(err).NotTo(HaveOccurred())
		Expect(vwc).To(Equal(expected))
	})
})
