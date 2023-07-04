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
	. "github.com/fluxninja/aperture/v2/operator/controllers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/utils/pointer"

	"github.com/fluxninja/aperture/v2/operator/api/common"
	controllerv1alpha1 "github.com/fluxninja/aperture/v2/operator/api/controller/v1alpha1"
)

var _ = Describe("ServiceAccount for Controller", func() {
	Context("Instance with default parameters", func() {
		It("returns correct ServiceAccount", func() {
			instance := &controllerv1alpha1.Controller{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: controllerv1alpha1.ControllerSpec{
					CommonSpec: common.CommonSpec{
						ServiceAccountSpec: common.ServiceAccountSpec{
							Create:                       true,
							AutomountServiceAccountToken: true,
						},
					},
				},
			}

			expected := &corev1.ServiceAccount{
				ObjectMeta: metav1.ObjectMeta{
					Name:      ControllerServiceName,
					Namespace: AppName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       AppName,
						"app.kubernetes.io/instance":   AppName,
						"app.kubernetes.io/managed-by": OperatorName,
						"app.kubernetes.io/component":  ControllerServiceName,
					},
				},
				AutomountServiceAccountToken: pointer.Bool(true),
			}

			result, err := serviceAccountForController(instance.DeepCopy(), scheme.Scheme)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(expected))
		})
	})

	Context("Instance with all parameters", func() {
		It("returns correct ServiceAccount", func() {
			instance := &controllerv1alpha1.Controller{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: controllerv1alpha1.ControllerSpec{
					CommonSpec: common.CommonSpec{
						Labels:      TestMap,
						Annotations: TestMap,
						ServiceAccountSpec: common.ServiceAccountSpec{
							Create:                       true,
							Annotations:                  TestMapTwo,
							AutomountServiceAccountToken: false,
						},
					},
				},
			}

			expected := &corev1.ServiceAccount{
				ObjectMeta: metav1.ObjectMeta{
					Name:      ControllerServiceName,
					Namespace: AppName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       AppName,
						"app.kubernetes.io/instance":   AppName,
						"app.kubernetes.io/managed-by": OperatorName,
						"app.kubernetes.io/component":  ControllerServiceName,
						Test:                           Test,
					},
					Annotations: map[string]string{
						Test:    Test,
						TestTwo: TestTwo,
					},
				},
				AutomountServiceAccountToken: pointer.Bool(false),
			}

			result, err := serviceAccountForController(instance.DeepCopy(), scheme.Scheme)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(expected))
		})
	})
})
