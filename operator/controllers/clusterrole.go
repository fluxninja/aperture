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
	"context"
	"fmt"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/fluxninja/aperture/operator/api/v1alpha1"
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
	}

	roleRef = rbacv1.RoleRef{
		APIGroup: "rbac.authorization.k8s.io",
		Kind:     "ClusterRole",
		Name:     appName,
	}
)

// clusterRoleForAgent prepares the ClusterRole object for the Agent based on the provided parameter.
func clusterRoleForAgent(instance *v1alpha1.Agent) *rbacv1.ClusterRole {
	clusterRole := &rbacv1.ClusterRole{
		ObjectMeta: v1.ObjectMeta{
			Name:        appName,
			Labels:      commonLabels(instance.Spec.Labels, instance.GetName(), operatorName),
			Annotations: getAgentAnnotationsWithOwnerRef(instance),
		},
		Rules: rules,
	}

	return clusterRole
}

// clusterRoleForController prepares the ClusterRole object for the Controller based on the provided parameter.
func clusterRoleForController(instance *v1alpha1.Controller) *rbacv1.ClusterRole {
	clusterRole := &rbacv1.ClusterRole{
		ObjectMeta: v1.ObjectMeta{
			Name:        appName,
			Labels:      commonLabels(instance.Spec.Labels, instance.GetName(), operatorName),
			Annotations: getControllerAnnotationsWithOwnerRef(instance),
		},
		Rules: rules,
	}

	return clusterRole
}

// clusterRoleBindingForAgent prepares the ClusterRoleBinding object for the Agent based on the provided parameter.
func clusterRoleBindingForAgent(instance *v1alpha1.Agent) *rbacv1.ClusterRoleBinding {
	clusterRoleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: v1.ObjectMeta{
			Name:        agentServiceName,
			Labels:      commonLabels(instance.Spec.Labels, instance.GetName(), agentServiceName),
			Annotations: getAgentAnnotationsWithOwnerRef(instance),
		},
		RoleRef: roleRef,
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      agentServiceName,
				Namespace: instance.GetNamespace(),
			},
		},
	}

	return clusterRoleBinding
}

// clusterRoleBindingForController prepares the ClusterRoleBinding object for the Controller based on the provided parameter.
func clusterRoleBindingForController(instance *v1alpha1.Controller) *rbacv1.ClusterRoleBinding {
	clusterRoleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: v1.ObjectMeta{
			Name:        controllerServiceName,
			Labels:      commonLabels(instance.Spec.Labels, instance.GetName(), controllerServiceName),
			Annotations: getControllerAnnotationsWithOwnerRef(instance),
		},
		RoleRef: roleRef,
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      controllerServiceName,
				Namespace: instance.GetNamespace(),
			},
		},
	}

	return clusterRoleBinding
}

// clusterRoleMutate returns a mutate function that can be used to update the ClusterRole's spec.
func clusterRoleMutate(cr *rbacv1.ClusterRole, rules []rbacv1.PolicyRule) controllerutil.MutateFn {
	return func() error {
		cr.Rules = rules
		return nil
	}
}

// clusterRoleBindingMutate returns a mutate function that can be used to update the Service's service spec.
func clusterRoleBindingMutate(crb *rbacv1.ClusterRoleBinding, roleRef rbacv1.RoleRef, subjects []rbacv1.Subject) controllerutil.MutateFn {
	return func() error {
		crb.RoleRef = roleRef
		return nil
	}
}

// updateClusterRoleBinding appends the Serviaccount in the ClusterRoleBinding if not exists.
func updateClusterRoleBinding(client client.Client, subject rbacv1.Subject, ctx context.Context, namespace string) error {
	crb := &rbacv1.ClusterRoleBinding{}
	err := client.Get(ctx, types.NamespacedName{Name: agentServiceName, Namespace: namespace}, crb)
	if err != nil {
		return fmt.Errorf("failed to Get the ClusterRoleBinding. Error: %+v", err)
	}

	for _, sub := range crb.Subjects {
		if sub.Name == subject.Name && sub.Namespace == subject.Namespace {
			return nil
		}
	}

	crb.Subjects = append(crb.Subjects, subject)
	err = client.Update(ctx, crb)
	if err != nil {
		if errors.IsConflict(err) {
			return updateClusterRoleBinding(client, subject, ctx, namespace)
		}

		return fmt.Errorf("failed to Update the ClusterRoleBinding. Error: %+v", err)
	}

	return nil
}
