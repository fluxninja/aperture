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

package mutatingwebhook

import (
	. "github.com/fluxninja/aperture/v2/operator/controllers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Test updateClusterRoleBinding", func() {
	It("should update the ClusterRoleBinding when valid Subject is provided", func() {
		crb := &rbacv1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      AgentServiceName,
				Namespace: Test,
			},
			RoleRef: rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     AppName,
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      Test,
					Namespace: Test,
				},
			},
		}
		err := K8sClient.Create(Ctx, crb)
		Expect(err).NotTo(HaveOccurred())

		subject := &rbacv1.Subject{
			Kind:      "ServiceAccount",
			Name:      TestTwo,
			Namespace: Test,
		}
		Expect(updateClusterRoleBinding(K8sClient, *subject, Ctx, Test)).To(Succeed())

		Expect(K8sClient.Get(Ctx, types.NamespacedName{Name: AgentServiceName, Namespace: Test}, crb)).To(Succeed())
		Expect(len(crb.Subjects)).To(Equal(2))
	})

	It("should not update the ClusterRoleBinding when valid Subject is not provided", func() {
		crb := &rbacv1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      AgentServiceName,
				Namespace: Test,
			},
			RoleRef: rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     AppName,
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      Test,
					Namespace: Test,
				},
			},
		}
		K8sClient.Create(Ctx, crb)

		subject := &rbacv1.Subject{
			Kind: "ServiceAccount",
			Name: TestTwo,
		}
		Expect(updateClusterRoleBinding(K8sClient, *subject, Ctx, Test) != nil).To(BeTrue())
	})
})
