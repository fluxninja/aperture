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

	. "github.com/fluxninja/aperture/operator/controllers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	agentv1alpha1 "github.com/fluxninja/aperture/operator/api/agent/v1alpha1"
	"github.com/fluxninja/aperture/operator/api/common"
)

var _ = Describe("clusterRoleForAgent", func() {
	Context("Instance with default parameters", func() {
		It("returns correct ClusterRole", func() {
			instance := &agentv1alpha1.Agent{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Agent",
					APIVersion: "fluxninja.com/v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: agentv1alpha1.AgentSpec{},
			}

			expected := &rbacv1.ClusterRole{
				ObjectMeta: metav1.ObjectMeta{
					Name: AgentServiceName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       AppName,
						"app.kubernetes.io/instance":   AppName,
						"app.kubernetes.io/managed-by": OperatorName,
						"app.kubernetes.io/component":  OperatorName,
					},
					Annotations: map[string]string{
						"fluxninja.com/primary-resource-type": "Agent.fluxninja.com",
						"fluxninja.com/primary-resource":      fmt.Sprintf("%s/%s", AppName, AppName),
					},
				},
				Rules: []rbacv1.PolicyRule{
					{
						APIGroups: []string{""},
						Resources: []string{"services", "events", "endpoints", "pods", "nodes", "namespaces", "componentstatuses"},
						Verbs:     []string{"get", "list", "watch"},
					},
					{
						NonResourceURLs: []string{"/version", "/healthz", "/metrics"},
						Verbs:           []string{"get"},
					},
					{
						APIGroups: []string{""},
						Resources: []string{"nodes/metrics", "nodes/spec", "nodes/proxy", "nodes/stats"},
						Verbs:     []string{"get"},
					},
					{
						APIGroups: []string{""},
						Resources: []string{"events"},
						Verbs:     []string{"create", "patch"},
					},
				},
			}

			result := clusterRoleForAgent(instance.DeepCopy())
			Expect(result).To(Equal(expected))
		})
	})

	Context("Instance with all parameters", func() {
		It("returns correct ClusterRole", func() {
			instance := &agentv1alpha1.Agent{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Agent",
					APIVersion: "fluxninja.com/v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: agentv1alpha1.AgentSpec{
					CommonSpec: common.CommonSpec{
						Labels:      TestMap,
						Annotations: TestMap,
					},
				},
			}

			expected := &rbacv1.ClusterRole{
				ObjectMeta: metav1.ObjectMeta{
					Name: AgentServiceName,
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
					},
				},
				Rules: []rbacv1.PolicyRule{
					{
						APIGroups: []string{""},
						Resources: []string{"services", "events", "endpoints", "pods", "nodes", "namespaces", "componentstatuses"},
						Verbs:     []string{"get", "list", "watch"},
					},
					{
						NonResourceURLs: []string{"/version", "/healthz", "/metrics"},
						Verbs:           []string{"get"},
					},
					{
						APIGroups: []string{""},
						Resources: []string{"nodes/metrics", "nodes/spec", "nodes/proxy", "nodes/stats"},
						Verbs:     []string{"get"},
					},
					{
						APIGroups: []string{""},
						Resources: []string{"events"},
						Verbs:     []string{"create", "patch"},
					},
				},
			}

			result := clusterRoleForAgent(instance.DeepCopy())
			Expect(result).To(Equal(expected))
		})
	})
})

var _ = Describe("clusterRoleBindingForAgent", func() {
	It("returns correct ClusterRoleBinding", func() {
		instance := &agentv1alpha1.Agent{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Agent",
				APIVersion: "fluxninja.com/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      AppName,
				Namespace: AppName,
			},
			Spec: agentv1alpha1.AgentSpec{
				CommonSpec: common.CommonSpec{
					Labels:      TestMap,
					Annotations: TestMap,
				},
			},
		}

		expected := &rbacv1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name: AgentServiceName,
				Labels: map[string]string{
					"app.kubernetes.io/name":       AppName,
					"app.kubernetes.io/instance":   AppName,
					"app.kubernetes.io/managed-by": OperatorName,
					"app.kubernetes.io/component":  AgentServiceName,
					Test:                           Test,
				},
				Annotations: map[string]string{
					"fluxninja.com/primary-resource-type": "Agent.fluxninja.com",
					"fluxninja.com/primary-resource":      fmt.Sprintf("%s/%s", AppName, AppName),
					Test:                                  Test,
				},
			},
			RoleRef: rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     AgentServiceName,
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      AgentServiceName,
					Namespace: instance.GetNamespace(),
				},
			},
		}

		result := clusterRoleBindingForAgent(instance.DeepCopy())
		Expect(result).To(Equal(expected))
	})
})
