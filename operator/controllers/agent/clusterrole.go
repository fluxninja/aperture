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

func genRules(instance *agentv1alpha1.Agent) []rbacv1.PolicyRule {
	rules := []rbacv1.PolicyRule{
		{
			APIGroups: []string{""},
			Resources: []string{"pods", "nodes"},
			Verbs:     []string{"get", "list", "watch"},
		},
	}

	if !instance.Spec.Sidecar.Enabled {
		newRules := []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"services", "events", "endpoints", "namespaces", "componentstatuses"},
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
		}
		rules = append(rules, newRules...)
	}

	if instance.Spec.ConfigSpec.ServiceDiscoverySpec.KubernetesDiscoveryConfig.AutoscaleEnabled {
		newRules := []rbacv1.PolicyRule{
			{
				APIGroups: []string{"*"},
				Resources: []string{"*"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{"*"},
				Resources: []string{"*/scale"},
				Verbs:     []string{"get", "update", "patch"},
			},
		}

		rules = append(rules, newRules...)
	}

	return rules
}

// clusterRoleForAgent prepares the ClusterRole object for the Agent based on the provided parameter.
func clusterRoleForAgent(instance *agentv1alpha1.Agent) *rbacv1.ClusterRole {
	clusterRole := &rbacv1.ClusterRole{
		ObjectMeta: v1.ObjectMeta{
			Name:        controllers.AgentServiceName,
			Labels:      controllers.CommonLabels(instance.Spec.Labels, instance.GetName(), controllers.OperatorName),
			Annotations: controllers.AgentAnnotationsWithOwnerRef(instance),
		},
		Rules: genRules(instance),
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
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     controllers.AgentServiceName,
		},
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
