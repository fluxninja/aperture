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
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/fluxninja/aperture/operator/api/v1alpha1"
)

var _ = Describe("clusterRoleForAgent", func() {
	Context("Instance with default parameters", func() {
		It("returns correct ClusterRole", func() {
			instance := &v1alpha1.Agent{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Agent",
					APIVersion: "fluxninja.com/v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.AgentSpec{},
			}

			expected := &rbacv1.ClusterRole{
				ObjectMeta: metav1.ObjectMeta{
					Name: appName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       appName,
						"app.kubernetes.io/instance":   appName,
						"app.kubernetes.io/managed-by": operatorName,
						"app.kubernetes.io/component":  operatorName,
					},
					Annotations: map[string]string{
						"fluxninja.com/primary-resource-type": "Agent.fluxninja.com",
						"fluxninja.com/primary-resource":      fmt.Sprintf("%s/%s", appName, appName),
					},
				},
				Rules: []rbacv1.PolicyRule{
					{
						APIGroups: []string{""},
						Resources: []string{"services", "events", "endpoints", "pods", "nodes", "namespaces", "componentstatuses"},
						Verbs:     []string{"get", "list", "watch"},
					},
					{
						APIGroups: []string{"quota.openshift.io"},
						Resources: []string{"clusterresourcequotas"},
						Verbs:     []string{"get"},
					},
					{
						NonResourceURLs: []string{"/version", "/healthz"},
						Verbs:           []string{"get"},
					},
					{
						NonResourceURLs: []string{"/metrics"},
						Verbs:           []string{"get"},
					},
					{
						APIGroups: []string{""},
						Resources: []string{"nodes/metrics", "nodes/spec", "nodes/proxy", "nodes/stats"},
						Verbs:     []string{"get"},
					},
					{
						APIGroups:     []string{"policy"},
						Resources:     []string{"podsecuritypolicies"},
						Verbs:         []string{"use"},
						ResourceNames: []string{appName},
					},
					{
						APIGroups:     []string{"security.openshift.io"},
						Resources:     []string{"securitycontextconstraints"},
						Verbs:         []string{"use"},
						ResourceNames: []string{appName},
					},
					{
						APIGroups: []string{"coordination.k8s.io"},
						Resources: []string{"leases"},
						Verbs:     []string{"create", "delete", "get", "list", "patch", "update", "watch"},
					},
					{
						APIGroups: []string{"admissionregistration.k8s.io"},
						Resources: []string{"mutatingwebhookconfigurations"},
						Verbs:     []string{"create", "delete", "get", "list", "patch", "update", "watch"},
					},
					{
						APIGroups: []string{"fluxninja.com"},
						Resources: []string{"policies"},
						Verbs:     []string{"create", "delete", "get", "list", "patch", "update", "watch"},
					},
					{
						APIGroups: []string{""},
						Resources: []string{"events"},
						Verbs:     []string{"create", "patch"},
					},
					{
						APIGroups: []string{"fluxninja.com"},
						Resources: []string{"policies/finalizers"},
						Verbs:     []string{"update"},
					},
					{
						APIGroups: []string{"fluxninja.com"},
						Resources: []string{"policies/status"},
						Verbs:     []string{"get", "patch", "update"},
					},
				},
			}

			result := clusterRoleForAgent(instance.DeepCopy())
			Expect(result).To(Equal(expected))
		})
	})

	Context("Instance with all parameters", func() {
		It("returns correct ClusterRole", func() {
			instance := &v1alpha1.Agent{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Agent",
					APIVersion: "fluxninja.com/v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.AgentSpec{
					CommonSpec: v1alpha1.CommonSpec{
						Labels:      testMap,
						Annotations: testMap,
					},
				},
			}

			expected := &rbacv1.ClusterRole{
				ObjectMeta: metav1.ObjectMeta{
					Name: appName,
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
					},
				},
				Rules: []rbacv1.PolicyRule{
					{
						APIGroups: []string{""},
						Resources: []string{"services", "events", "endpoints", "pods", "nodes", "namespaces", "componentstatuses"},
						Verbs:     []string{"get", "list", "watch"},
					},
					{
						APIGroups: []string{"quota.openshift.io"},
						Resources: []string{"clusterresourcequotas"},
						Verbs:     []string{"get"},
					},
					{
						NonResourceURLs: []string{"/version", "/healthz"},
						Verbs:           []string{"get"},
					},
					{
						NonResourceURLs: []string{"/metrics"},
						Verbs:           []string{"get"},
					},
					{
						APIGroups: []string{""},
						Resources: []string{"nodes/metrics", "nodes/spec", "nodes/proxy", "nodes/stats"},
						Verbs:     []string{"get"},
					},
					{
						APIGroups:     []string{"policy"},
						Resources:     []string{"podsecuritypolicies"},
						Verbs:         []string{"use"},
						ResourceNames: []string{appName},
					},
					{
						APIGroups:     []string{"security.openshift.io"},
						Resources:     []string{"securitycontextconstraints"},
						Verbs:         []string{"use"},
						ResourceNames: []string{appName},
					},
					{
						APIGroups: []string{"coordination.k8s.io"},
						Resources: []string{"leases"},
						Verbs:     []string{"create", "delete", "get", "list", "patch", "update", "watch"},
					},
					{
						APIGroups: []string{"admissionregistration.k8s.io"},
						Resources: []string{"mutatingwebhookconfigurations"},
						Verbs:     []string{"create", "delete", "get", "list", "patch", "update", "watch"},
					},
					{
						APIGroups: []string{"fluxninja.com"},
						Resources: []string{"policies"},
						Verbs:     []string{"create", "delete", "get", "list", "patch", "update", "watch"},
					},
					{
						APIGroups: []string{""},
						Resources: []string{"events"},
						Verbs:     []string{"create", "patch"},
					},
					{
						APIGroups: []string{"fluxninja.com"},
						Resources: []string{"policies/finalizers"},
						Verbs:     []string{"update"},
					},
					{
						APIGroups: []string{"fluxninja.com"},
						Resources: []string{"policies/status"},
						Verbs:     []string{"get", "patch", "update"},
					},
				},
			}

			result := clusterRoleForAgent(instance.DeepCopy())
			Expect(result).To(Equal(expected))
		})
	})
})

var _ = Describe("clusterRoleForController", func() {
	Context("Instance with default parameters", func() {
		It("returns correct ClusterRole", func() {
			instance := &v1alpha1.Controller{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Controller",
					APIVersion: "fluxninja.com/v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.ControllerSpec{},
			}

			expected := &rbacv1.ClusterRole{
				ObjectMeta: metav1.ObjectMeta{
					Name: appName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       appName,
						"app.kubernetes.io/instance":   appName,
						"app.kubernetes.io/managed-by": operatorName,
						"app.kubernetes.io/component":  operatorName,
					},
					Annotations: map[string]string{
						"fluxninja.com/primary-resource-type": "Controller.fluxninja.com",
						"fluxninja.com/primary-resource":      fmt.Sprintf("%s/%s", appName, appName),
					},
				},
				Rules: []rbacv1.PolicyRule{
					{
						APIGroups: []string{""},
						Resources: []string{"services", "events", "endpoints", "pods", "nodes", "namespaces", "componentstatuses"},
						Verbs:     []string{"get", "list", "watch"},
					},
					{
						APIGroups: []string{"quota.openshift.io"},
						Resources: []string{"clusterresourcequotas"},
						Verbs:     []string{"get"},
					},
					{
						NonResourceURLs: []string{"/version", "/healthz"},
						Verbs:           []string{"get"},
					},
					{
						NonResourceURLs: []string{"/metrics"},
						Verbs:           []string{"get"},
					},
					{
						APIGroups: []string{""},
						Resources: []string{"nodes/metrics", "nodes/spec", "nodes/proxy", "nodes/stats"},
						Verbs:     []string{"get"},
					},
					{
						APIGroups:     []string{"policy"},
						Resources:     []string{"podsecuritypolicies"},
						Verbs:         []string{"use"},
						ResourceNames: []string{appName},
					},
					{
						APIGroups:     []string{"security.openshift.io"},
						Resources:     []string{"securitycontextconstraints"},
						Verbs:         []string{"use"},
						ResourceNames: []string{appName},
					},
					{
						APIGroups: []string{"coordination.k8s.io"},
						Resources: []string{"leases"},
						Verbs:     []string{"create", "delete", "get", "list", "patch", "update", "watch"},
					},
					{
						APIGroups: []string{"admissionregistration.k8s.io"},
						Resources: []string{"mutatingwebhookconfigurations"},
						Verbs:     []string{"create", "delete", "get", "list", "patch", "update", "watch"},
					},
					{
						APIGroups: []string{"fluxninja.com"},
						Resources: []string{"policies"},
						Verbs:     []string{"create", "delete", "get", "list", "patch", "update", "watch"},
					},
					{
						APIGroups: []string{""},
						Resources: []string{"events"},
						Verbs:     []string{"create", "patch"},
					},
					{
						APIGroups: []string{"fluxninja.com"},
						Resources: []string{"policies/finalizers"},
						Verbs:     []string{"update"},
					},
					{
						APIGroups: []string{"fluxninja.com"},
						Resources: []string{"policies/status"},
						Verbs:     []string{"get", "patch", "update"},
					},
				},
			}

			result := clusterRoleForController(instance.DeepCopy())
			Expect(result).To(Equal(expected))
		})
	})

	Context("Instance with all parameters", func() {
		It("returns correct ClusterRole", func() {
			instance := &v1alpha1.Controller{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Controller",
					APIVersion: "fluxninja.com/v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.ControllerSpec{
					CommonSpec: v1alpha1.CommonSpec{
						Labels:      testMap,
						Annotations: testMap,
					},
				},
			}

			expected := &rbacv1.ClusterRole{
				ObjectMeta: metav1.ObjectMeta{
					Name: appName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       appName,
						"app.kubernetes.io/instance":   appName,
						"app.kubernetes.io/managed-by": operatorName,
						"app.kubernetes.io/component":  operatorName,
						test:                           test,
					},
					Annotations: map[string]string{
						"fluxninja.com/primary-resource-type": "Controller.fluxninja.com",
						"fluxninja.com/primary-resource":      fmt.Sprintf("%s/%s", appName, appName),
						test:                                  test,
					},
				},
				Rules: []rbacv1.PolicyRule{
					{
						APIGroups: []string{""},
						Resources: []string{"services", "events", "endpoints", "pods", "nodes", "namespaces", "componentstatuses"},
						Verbs:     []string{"get", "list", "watch"},
					},
					{
						APIGroups: []string{"quota.openshift.io"},
						Resources: []string{"clusterresourcequotas"},
						Verbs:     []string{"get"},
					},
					{
						NonResourceURLs: []string{"/version", "/healthz"},
						Verbs:           []string{"get"},
					},
					{
						NonResourceURLs: []string{"/metrics"},
						Verbs:           []string{"get"},
					},
					{
						APIGroups: []string{""},
						Resources: []string{"nodes/metrics", "nodes/spec", "nodes/proxy", "nodes/stats"},
						Verbs:     []string{"get"},
					},
					{
						APIGroups:     []string{"policy"},
						Resources:     []string{"podsecuritypolicies"},
						Verbs:         []string{"use"},
						ResourceNames: []string{appName},
					},
					{
						APIGroups:     []string{"security.openshift.io"},
						Resources:     []string{"securitycontextconstraints"},
						Verbs:         []string{"use"},
						ResourceNames: []string{appName},
					},
					{
						APIGroups: []string{"coordination.k8s.io"},
						Resources: []string{"leases"},
						Verbs:     []string{"create", "delete", "get", "list", "patch", "update", "watch"},
					},
					{
						APIGroups: []string{"admissionregistration.k8s.io"},
						Resources: []string{"mutatingwebhookconfigurations"},
						Verbs:     []string{"create", "delete", "get", "list", "patch", "update", "watch"},
					},
					{
						APIGroups: []string{"fluxninja.com"},
						Resources: []string{"policies"},
						Verbs:     []string{"create", "delete", "get", "list", "patch", "update", "watch"},
					},
					{
						APIGroups: []string{""},
						Resources: []string{"events"},
						Verbs:     []string{"create", "patch"},
					},
					{
						APIGroups: []string{"fluxninja.com"},
						Resources: []string{"policies/finalizers"},
						Verbs:     []string{"update"},
					},
					{
						APIGroups: []string{"fluxninja.com"},
						Resources: []string{"policies/status"},
						Verbs:     []string{"get", "patch", "update"},
					},
				},
			}

			result := clusterRoleForController(instance.DeepCopy())
			Expect(result).To(Equal(expected))
		})
	})
})

var _ = Describe("clusterRoleBindingForAgent", func() {
	It("returns correct ClusterRoleBinding", func() {
		instance := &v1alpha1.Agent{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Agent",
				APIVersion: "fluxninja.com/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      appName,
				Namespace: appName,
			},
			Spec: v1alpha1.AgentSpec{
				CommonSpec: v1alpha1.CommonSpec{
					Labels:      testMap,
					Annotations: testMap,
				},
			},
		}

		expected := &rbacv1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name: agentServiceName,
				Labels: map[string]string{
					"app.kubernetes.io/name":       appName,
					"app.kubernetes.io/instance":   appName,
					"app.kubernetes.io/managed-by": operatorName,
					"app.kubernetes.io/component":  agentServiceName,
					test:                           test,
				},
				Annotations: map[string]string{
					"fluxninja.com/primary-resource-type": "Agent.fluxninja.com",
					"fluxninja.com/primary-resource":      fmt.Sprintf("%s/%s", appName, appName),
					test:                                  test,
				},
			},
			RoleRef: rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     appName,
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      agentServiceName,
					Namespace: instance.GetNamespace(),
				},
			},
		}

		result := clusterRoleBindingForAgent(instance.DeepCopy())
		Expect(result).To(Equal(expected))
	})
})

var _ = Describe("clusterRoleBindingForController", func() {
	It("returns correct ClusterRoleBinding", func() {
		instance := &v1alpha1.Controller{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Controller",
				APIVersion: "fluxninja.com/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      appName,
				Namespace: appName,
			},
			Spec: v1alpha1.ControllerSpec{
				CommonSpec: v1alpha1.CommonSpec{
					Labels:      testMap,
					Annotations: testMap,
				},
			},
		}

		expected := &rbacv1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name: controllerServiceName,
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
				},
			},
			RoleRef: rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     appName,
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      controllerServiceName,
					Namespace: instance.GetNamespace(),
				},
			},
		}

		result := clusterRoleBindingForController(instance.DeepCopy())
		Expect(result).To(Equal(expected))
	})
})

var _ = Describe("Test ClusterRole Mutate", func() {
	It("Mutate should update required fields only", func() {
		expected := &rbacv1.ClusterRole{
			ObjectMeta: metav1.ObjectMeta{},
			Rules: []rbacv1.PolicyRule{
				{
					APIGroups: testArray,
					Resources: testArray,
					Verbs:     testArray,
				},
			},
		}

		cr := &rbacv1.ClusterRole{}
		err := clusterRoleMutate(cr, expected.Rules)()
		Expect(err).NotTo(HaveOccurred())
		Expect(cr).To(Equal(expected))
	})
})

var _ = Describe("Test ClusterRoleBinding Mutate", func() {
	It("Mutate should update required fields only", func() {
		expected := &rbacv1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{},
			RoleRef: rbacv1.RoleRef{
				APIGroup: test,
				Kind:     test,
				Name:     test,
			},
		}

		crb := &rbacv1.ClusterRoleBinding{}
		err := clusterRoleBindingMutate(crb, expected.RoleRef, expected.Subjects)()
		Expect(err).NotTo(HaveOccurred())
		Expect(crb).To(Equal(expected))
	})
})

var _ = Describe("Test updateClusterRoleBinding", func() {
	It("should update the ClusterRoleBinding when valid Subject is provided", func() {
		crb := &rbacv1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      agentServiceName,
				Namespace: test,
			},
			RoleRef: rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     appName,
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      test,
					Namespace: test,
				},
			},
		}
		err := k8sClient.Create(ctx, crb)
		Expect(err).NotTo(HaveOccurred())

		subject := &rbacv1.Subject{
			Kind:      "ServiceAccount",
			Name:      testTwo,
			Namespace: test,
		}
		Expect(updateClusterRoleBinding(k8sClient, *subject, ctx, test)).To(BeNil())

		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: agentServiceName, Namespace: test}, crb)).To(BeNil())
		Expect(len(crb.Subjects)).To(Equal(2))
	})

	It("should not update the ClusterRoleBinding when valid Subject is not provided", func() {
		crb := &rbacv1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      agentServiceName,
				Namespace: test,
			},
			RoleRef: rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     appName,
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      test,
					Namespace: test,
				},
			},
		}
		k8sClient.Create(ctx, crb)

		subject := &rbacv1.Subject{
			Kind: "ServiceAccount",
			Name: testTwo,
		}
		Expect(updateClusterRoleBinding(k8sClient, *subject, ctx, test) != nil).To(BeTrue())
	})
})
