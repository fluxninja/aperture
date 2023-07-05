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
	"fmt"

	. "github.com/fluxninja/aperture/v2/operator/controllers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/fluxninja/aperture/v2/operator/api/common"
	controllerv1alpha1 "github.com/fluxninja/aperture/v2/operator/api/controller/v1alpha1"
)

var _ = Describe("clusterRoleForController", func() {
	Context("Instance with default parameters", func() {
		It("returns correct ClusterRole", func() {
			instance := &controllerv1alpha1.Controller{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Controller",
					APIVersion: "fluxninja.com/v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      ControllerName,
					Namespace: AppName,
				},
				Spec: controllerv1alpha1.ControllerSpec{},
			}

			expected := &rbacv1.ClusterRole{
				ObjectMeta: metav1.ObjectMeta{
					Name: ControllerServiceName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       AppName,
						"app.kubernetes.io/instance":   ControllerName,
						"app.kubernetes.io/managed-by": OperatorName,
						"app.kubernetes.io/component":  ControllerServiceName,
					},
					Annotations: map[string]string{
						"fluxninja.com/primary-resource-type": "Controller.fluxninja.com",
						"fluxninja.com/primary-resource":      fmt.Sprintf("%s/%s", AppName, ControllerName),
					},
				},
				Rules: []rbacv1.PolicyRule{
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
			instance := &controllerv1alpha1.Controller{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Controller",
					APIVersion: "fluxninja.com/v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      ControllerName,
					Namespace: AppName,
				},
				Spec: controllerv1alpha1.ControllerSpec{
					CommonSpec: common.CommonSpec{
						Labels:      TestMap,
						Annotations: TestMap,
					},
				},
			}

			expected := &rbacv1.ClusterRole{
				ObjectMeta: metav1.ObjectMeta{
					Name: ControllerServiceName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       AppName,
						"app.kubernetes.io/instance":   ControllerName,
						"app.kubernetes.io/managed-by": OperatorName,
						"app.kubernetes.io/component":  ControllerServiceName,
						Test:                           Test,
					},
					Annotations: map[string]string{
						"fluxninja.com/primary-resource-type": "Controller.fluxninja.com",
						"fluxninja.com/primary-resource":      fmt.Sprintf("%s/%s", AppName, ControllerName),
						Test:                                  Test,
					},
				},
				Rules: []rbacv1.PolicyRule{
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

var _ = Describe("clusterRoleBindingForController", func() {
	It("returns correct ClusterRoleBinding when spec.ServiceAccountSpec.Create is false", func() {
		instance := &controllerv1alpha1.Controller{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Controller",
				APIVersion: "fluxninja.com/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      ControllerName,
				Namespace: AppName,
			},
			Spec: controllerv1alpha1.ControllerSpec{
				CommonSpec: common.CommonSpec{
					Labels:      TestMap,
					Annotations: TestMap,
				},
			},
		}

		expected := &rbacv1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name: ControllerServiceName,
				Labels: map[string]string{
					"app.kubernetes.io/name":       AppName,
					"app.kubernetes.io/instance":   ControllerName,
					"app.kubernetes.io/managed-by": OperatorName,
					"app.kubernetes.io/component":  ControllerServiceName,
					Test:                           Test,
				},
				Annotations: map[string]string{
					"fluxninja.com/primary-resource-type": "Controller.fluxninja.com",
					"fluxninja.com/primary-resource":      fmt.Sprintf("%s/%s", AppName, ControllerName),
					Test:                                  Test,
				},
			},
			RoleRef: rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     ControllerServiceName,
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      "",
					Namespace: instance.GetNamespace(),
				},
			},
		}

		result := clusterRoleBindingForController(instance.DeepCopy())
		Expect(result).To(Equal(expected))
	})

	It("returns correct ClusterRoleBinding when spec.ServiceAccountSpec.Create is true and spec.ServiceAccountSpec.Name is not provided", func() {
		instance := &controllerv1alpha1.Controller{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Controller",
				APIVersion: "fluxninja.com/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      ControllerName,
				Namespace: AppName,
			},
			Spec: controllerv1alpha1.ControllerSpec{
				CommonSpec: common.CommonSpec{
					Labels:      TestMap,
					Annotations: TestMap,
					ServiceAccountSpec: common.ServiceAccountSpec{
						Create: true,
					},
				},
			},
		}

		expected := &rbacv1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name: ControllerServiceName,
				Labels: map[string]string{
					"app.kubernetes.io/name":       AppName,
					"app.kubernetes.io/instance":   ControllerName,
					"app.kubernetes.io/managed-by": OperatorName,
					"app.kubernetes.io/component":  ControllerServiceName,
					Test:                           Test,
				},
				Annotations: map[string]string{
					"fluxninja.com/primary-resource-type": "Controller.fluxninja.com",
					"fluxninja.com/primary-resource":      fmt.Sprintf("%s/%s", AppName, ControllerName),
					Test:                                  Test,
				},
			},
			RoleRef: rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     ControllerServiceName,
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      ControllerServiceName,
					Namespace: instance.GetNamespace(),
				},
			},
		}

		result := clusterRoleBindingForController(instance.DeepCopy())
		Expect(result).To(Equal(expected))
	})

	It("returns correct ClusterRoleBinding when spec.ServiceAccountSpec.Create is true and spec.ServiceAccountSpec.Name is provided", func() {
		instance := &controllerv1alpha1.Controller{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Controller",
				APIVersion: "fluxninja.com/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      ControllerName,
				Namespace: AppName,
			},
			Spec: controllerv1alpha1.ControllerSpec{
				CommonSpec: common.CommonSpec{
					Labels:      TestMap,
					Annotations: TestMap,
					ServiceAccountSpec: common.ServiceAccountSpec{
						Create: true,
						Name:   Test,
					},
				},
			},
		}

		expected := &rbacv1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name: ControllerServiceName,
				Labels: map[string]string{
					"app.kubernetes.io/name":       AppName,
					"app.kubernetes.io/instance":   ControllerName,
					"app.kubernetes.io/managed-by": OperatorName,
					"app.kubernetes.io/component":  ControllerServiceName,
					Test:                           Test,
				},
				Annotations: map[string]string{
					"fluxninja.com/primary-resource-type": "Controller.fluxninja.com",
					"fluxninja.com/primary-resource":      fmt.Sprintf("%s/%s", AppName, ControllerName),
					Test:                                  Test,
				},
			},
			RoleRef: rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     ControllerServiceName,
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      Test,
					Namespace: instance.GetNamespace(),
				},
			},
		}

		result := clusterRoleBindingForController(instance.DeepCopy())
		Expect(result).To(Equal(expected))
	})
})
