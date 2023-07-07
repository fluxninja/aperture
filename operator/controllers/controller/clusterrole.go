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

package controller

import (
	"github.com/fluxninja/aperture/v2/operator/controllers"

	rbacv1 "k8s.io/api/rbac/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	controllerv1alpha1 "github.com/fluxninja/aperture/v2/operator/api/controller/v1alpha1"
)

var rules = []rbacv1.PolicyRule{
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

// clusterRoleForController prepares the ClusterRole object for the Controller based on the provided parameter.
func clusterRoleForController(instance *controllerv1alpha1.Controller) *rbacv1.ClusterRole {
	clusterRole := &rbacv1.ClusterRole{
		ObjectMeta: v1.ObjectMeta{
			Name:        controllers.ControllerResourcesNamespacedName(instance),
			Labels:      controllers.CommonLabels(instance.Spec.Labels, instance.GetName(), controllers.ControllerServiceName),
			Annotations: controllers.ControllerAnnotationsWithOwnerRef(instance),
		},
		Rules: rules,
	}

	return clusterRole
}

// clusterRoleBindingForController prepares the ClusterRoleBinding object for the Controller based on the provided parameter.
func clusterRoleBindingForController(instance *controllerv1alpha1.Controller) *rbacv1.ClusterRoleBinding {
	clusterRoleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: v1.ObjectMeta{
			Name:        controllers.ControllerResourcesNamespacedName(instance),
			Labels:      controllers.CommonLabels(instance.Spec.Labels, instance.GetName(), controllers.ControllerServiceName),
			Annotations: controllers.ControllerAnnotationsWithOwnerRef(instance),
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     controllers.ControllerResourcesNamespacedName(instance),
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      controllers.ServiceAccountName(instance),
				Namespace: instance.GetNamespace(),
			},
		},
	}

	return clusterRoleBinding
}
