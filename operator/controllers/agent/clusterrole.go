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
	"github.com/fluxninja/aperture/operator/controllers"

	rbacv1 "k8s.io/api/rbac/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	agentv1alpha1 "github.com/fluxninja/aperture/operator/api/agent/v1alpha1"
)

var (
	rules = []rbacv1.PolicyRule{
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
			ResourceNames: []string{controllers.AppName},
		},
		{
			APIGroups:     []string{"security.openshift.io"},
			Resources:     []string{"securitycontextconstraints"},
			Verbs:         []string{"use"},
			ResourceNames: []string{controllers.AppName},
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
	}

	roleRef = rbacv1.RoleRef{
		APIGroup: "rbac.authorization.k8s.io",
		Kind:     "ClusterRole",
		Name:     controllers.AppName,
	}
)

// clusterRoleForAgent prepares the ClusterRole object for the Agent based on the provided parameter.
func clusterRoleForAgent(instance *agentv1alpha1.Agent) *rbacv1.ClusterRole {
	clusterRole := &rbacv1.ClusterRole{
		ObjectMeta: v1.ObjectMeta{
			Name:        controllers.AppName,
			Labels:      controllers.CommonLabels(instance.Spec.Labels, instance.GetName(), controllers.OperatorName),
			Annotations: controllers.AgentAnnotationsWithOwnerRef(instance),
		},
		Rules: rules,
	}

	return clusterRole
}

// clusterRoleBindingForAgent prepares the ClusterRoleBinding object for the Agent based on the provided parameter.
func clusterRoleBindingForAgent(instance *agentv1alpha1.Agent) *rbacv1.ClusterRoleBinding {
	clusterRoleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: v1.ObjectMeta{
			Name:        controllers.AgentServiceName,
			Labels:      controllers.CommonLabels(instance.Spec.Labels, instance.GetName(), controllers.AgentServiceName),
			Annotations: controllers.AgentAnnotationsWithOwnerRef(instance),
		},
		RoleRef: roleRef,
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      controllers.AgentServiceName,
				Namespace: instance.GetNamespace(),
			},
		},
	}

	return clusterRoleBinding
}
