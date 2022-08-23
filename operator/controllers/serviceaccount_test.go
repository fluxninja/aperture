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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/utils/pointer"

	"aperture.tech/operators/aperture-operator/api/v1alpha1"
)

var _ = Describe("ServiceAccount for Controller", func() {
	Context("Instance with default parameters", func() {
		It("returns correct ServiceAccount", func() {
			instance := &v1alpha1.Aperture{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.ApertureSpec{
					Controller: v1alpha1.ControllerSpec{
						CommonSpec: v1alpha1.CommonSpec{
							ServiceAccountSpec: v1alpha1.ServiceAccountSpec{
								Create:                       true,
								AutomountServiceAccountToken: true,
							},
						},
					},
				},
			}

			expected := &corev1.ServiceAccount{
				ObjectMeta: metav1.ObjectMeta{
					Name:      controllerServiceName,
					Namespace: appName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       appName,
						"app.kubernetes.io/instance":   appName,
						"app.kubernetes.io/managed-by": operatorName,
						"app.kubernetes.io/component":  controllerServiceName,
					},
					Annotations: nil,
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion:         "aperture.tech/v1alpha1",
							Name:               instance.GetName(),
							Kind:               "Aperture",
							Controller:         pointer.BoolPtr(true),
							BlockOwnerDeletion: pointer.BoolPtr(true),
						},
					},
				},
				AutomountServiceAccountToken: pointer.BoolPtr(true),
			}

			result, _ := serviceAccountForController(instance.DeepCopy(), scheme.Scheme)
			Expect(result).To(Equal(expected))
		})
	})

	Context("Instance with all parameters", func() {
		It("returns correct ServiceAccount", func() {
			instance := &v1alpha1.Aperture{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.ApertureSpec{
					Labels:      testMap,
					Annotations: testMap,
					Controller: v1alpha1.ControllerSpec{
						CommonSpec: v1alpha1.CommonSpec{
							ServiceAccountSpec: v1alpha1.ServiceAccountSpec{
								Create:                       true,
								Annotations:                  testMapTwo,
								AutomountServiceAccountToken: false,
							},
						},
					},
				},
			}

			expected := &corev1.ServiceAccount{
				ObjectMeta: metav1.ObjectMeta{
					Name:      controllerServiceName,
					Namespace: appName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       appName,
						"app.kubernetes.io/instance":   appName,
						"app.kubernetes.io/managed-by": operatorName,
						"app.kubernetes.io/component":  controllerServiceName,
						test:                           test,
					},
					Annotations: map[string]string{
						test:    test,
						testTwo: testTwo,
					},
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion:         "aperture.tech/v1alpha1",
							Name:               instance.GetName(),
							Kind:               "Aperture",
							Controller:         pointer.BoolPtr(true),
							BlockOwnerDeletion: pointer.BoolPtr(true),
						},
					},
				},
				AutomountServiceAccountToken: pointer.BoolPtr(false),
			}

			result, _ := serviceAccountForController(instance.DeepCopy(), scheme.Scheme)
			Expect(result).To(Equal(expected))
		})
	})
})

var _ = Describe("ServiceAccount for Agent", func() {
	Context("Instance with default parameters", func() {
		It("returns correct ServiceAccount", func() {
			instance := &v1alpha1.Aperture{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.ApertureSpec{
					Agent: v1alpha1.AgentSpec{
						CommonSpec: v1alpha1.CommonSpec{
							ServiceAccountSpec: v1alpha1.ServiceAccountSpec{
								Create:                       true,
								AutomountServiceAccountToken: true,
							},
						},
					},
				},
			}

			expected := &corev1.ServiceAccount{
				ObjectMeta: metav1.ObjectMeta{
					Name:      agentServiceName,
					Namespace: appName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       appName,
						"app.kubernetes.io/instance":   appName,
						"app.kubernetes.io/managed-by": operatorName,
						"app.kubernetes.io/component":  agentServiceName,
					},
					Annotations: nil,
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion:         "aperture.tech/v1alpha1",
							Name:               instance.GetName(),
							Kind:               "Aperture",
							Controller:         pointer.BoolPtr(true),
							BlockOwnerDeletion: pointer.BoolPtr(true),
						},
					},
				},
				AutomountServiceAccountToken: pointer.BoolPtr(true),
			}

			result, _ := serviceAccountForAgent(instance.DeepCopy(), scheme.Scheme)
			Expect(result).To(Equal(expected))
		})
	})

	Context("Instance with all parameters", func() {
		It("returns correct ServiceAccount", func() {
			instance := &v1alpha1.Aperture{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.ApertureSpec{
					Labels:      testMap,
					Annotations: testMap,
					Agent: v1alpha1.AgentSpec{
						CommonSpec: v1alpha1.CommonSpec{
							ServiceAccountSpec: v1alpha1.ServiceAccountSpec{
								Create:                       true,
								Annotations:                  testMapTwo,
								AutomountServiceAccountToken: false,
							},
						},
					},
				},
			}

			expected := &corev1.ServiceAccount{
				ObjectMeta: metav1.ObjectMeta{
					Name:      agentServiceName,
					Namespace: appName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       appName,
						"app.kubernetes.io/instance":   appName,
						"app.kubernetes.io/managed-by": operatorName,
						"app.kubernetes.io/component":  agentServiceName,
						test:                           test,
					},
					Annotations: map[string]string{
						test:    test,
						testTwo: testTwo,
					},
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion:         "aperture.tech/v1alpha1",
							Name:               instance.GetName(),
							Kind:               "Aperture",
							Controller:         pointer.BoolPtr(true),
							BlockOwnerDeletion: pointer.BoolPtr(true),
						},
					},
				},
				AutomountServiceAccountToken: pointer.BoolPtr(false),
			}

			result, _ := serviceAccountForAgent(instance.DeepCopy(), scheme.Scheme)
			Expect(result).To(Equal(expected))
		})
	})
})

var _ = Describe("Test ServiceAccount Mutate", func() {
	It("Mutate should update required fields only", func() {
		expected := &corev1.ServiceAccount{
			ObjectMeta:                   metav1.ObjectMeta{},
			AutomountServiceAccountToken: pointer.BoolPtr(true),
		}

		sa := &corev1.ServiceAccount{}
		err := serviceAccountMutate(sa, expected.AutomountServiceAccountToken)()
		Expect(err).NotTo(HaveOccurred())
		Expect(sa).To(Equal(expected))
	})
})
