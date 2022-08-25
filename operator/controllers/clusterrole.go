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

// clusterRoleForAgent prepares the ClusterRole object for the Agent based on the provided parameter.
func clusterRoleForAgent(instance *v1alpha1.Aperture) *rbacv1.ClusterRole {
	clusterRole := &rbacv1.ClusterRole{
		ObjectMeta: v1.ObjectMeta{
			Name:        appName,
			Labels:      commonLabels(instance, agentServiceName),
			Annotations: getAnnotationsWithOwnerRef(instance),
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
				Verbs:     []string{"get"},
			},
		},
	}

	return clusterRole
}

// clusterRoleBindingForAgent prepares the ClusterRoleBinding object for the Agent based on the provided parameter.
func clusterRoleBindingForAgent(instance *v1alpha1.Aperture) *rbacv1.ClusterRoleBinding {
	clusterRoleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: v1.ObjectMeta{
			Name:        appName,
			Labels:      commonLabels(instance, agentServiceName),
			Annotations: getAnnotationsWithOwnerRef(instance),
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

	if !instance.Spec.Sidecar.Enabled {
		clusterRoleBinding.Subjects = append(clusterRoleBinding.Subjects, rbacv1.Subject{
			Kind:      "ServiceAccount",
			Name:      agentServiceName,
			Namespace: instance.GetNamespace(),
		})
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
	err := client.Get(ctx, types.NamespacedName{Name: appName, Namespace: namespace}, crb)
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
