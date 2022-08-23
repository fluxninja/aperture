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
	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/utils/pointer"

	"aperture.tech/operators/aperture-operator/api/v1alpha1"
)

var _ = Describe("Service for Controller Webhook", func() {
	Context("Instance with default parameters", func() {
		It("returns correct Service", func() {
			instance := &v1alpha1.Aperture{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
			}

			expected := &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "agent-webhooks",
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
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{
							Name:       "https",
							Protocol:   corev1.Protocol("TCP"),
							Port:       int32(443),
							TargetPort: intstr.FromInt(8086),
						},
					},
					Selector: map[string]string{
						"app.kubernetes.io/name":       appName,
						"app.kubernetes.io/instance":   appName,
						"app.kubernetes.io/managed-by": operatorName,
						"app.kubernetes.io/component":  controllerServiceName,
					},
				},
			}

			result, _ := serviceForControllerWebhook(instance.DeepCopy(), logr.Logger{}, scheme.Scheme)
			Expect(result).To(Equal(expected))
		})
	})

	Context("Instance with all parameters", func() {
		It("returns correct Service", func() {
			instance := &v1alpha1.Aperture{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.ApertureSpec{
					Labels:      testMap,
					Annotations: testMap,
					Service: v1alpha1.ServiceSpec{
						Controller: v1alpha1.Service{
							Annotations: testMapTwo,
						},
					},
				},
			}

			expected := &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "agent-webhooks",
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
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{
							Name:       "https",
							Protocol:   corev1.Protocol("TCP"),
							Port:       int32(443),
							TargetPort: intstr.FromInt(8086),
						},
					},
					Selector: map[string]string{
						"app.kubernetes.io/name":       appName,
						"app.kubernetes.io/instance":   appName,
						"app.kubernetes.io/managed-by": operatorName,
						"app.kubernetes.io/component":  controllerServiceName,
					},
				},
			}

			result, _ := serviceForControllerWebhook(instance.DeepCopy(), logr.Logger{}, scheme.Scheme)
			Expect(result).To(Equal(expected))
		})
	})
})

var _ = Describe("Service for Controller", func() {
	Context("Instance with default parameters", func() {
		It("returns correct Service", func() {
			instance := &v1alpha1.Aperture{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
			}

			expected := &corev1.Service{
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
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{
							Name:       "grpc",
							Protocol:   corev1.Protocol("TCP"),
							Port:       int32(80),
							TargetPort: intstr.FromString("grpc"),
						},
					},
					Selector: map[string]string{
						"app.kubernetes.io/name":       appName,
						"app.kubernetes.io/instance":   appName,
						"app.kubernetes.io/managed-by": operatorName,
						"app.kubernetes.io/component":  controllerServiceName,
					},
				},
			}

			result, _ := serviceForController(instance.DeepCopy(), logr.Logger{}, scheme.Scheme)
			Expect(result).To(Equal(expected))
		})
	})

	Context("Instance with all parameters", func() {
		It("returns correct Service", func() {
			instance := &v1alpha1.Aperture{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.ApertureSpec{
					Labels:      testMap,
					Annotations: testMap,
					Service: v1alpha1.ServiceSpec{
						Controller: v1alpha1.Service{
							Annotations: testMapTwo,
						},
					},
				},
			}

			expected := &corev1.Service{
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
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{
							Name:       "grpc",
							Protocol:   corev1.Protocol("TCP"),
							Port:       int32(80),
							TargetPort: intstr.FromString("grpc"),
						},
					},
					Selector: map[string]string{
						"app.kubernetes.io/name":       appName,
						"app.kubernetes.io/instance":   appName,
						"app.kubernetes.io/managed-by": operatorName,
						"app.kubernetes.io/component":  controllerServiceName,
					},
				},
			}

			result, _ := serviceForController(instance.DeepCopy(), logr.Logger{}, scheme.Scheme)
			Expect(result).To(Equal(expected))
		})
	})
})

var _ = Describe("Service for Agent", func() {
	Context("Instance with default parameters", func() {
		It("returns correct Service", func() {
			instance := &v1alpha1.Aperture{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
			}

			expected := &corev1.Service{
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
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{
							Name:       "grpc",
							Protocol:   corev1.Protocol("TCP"),
							Port:       int32(80),
							TargetPort: intstr.FromString("grpc"),
						},
						{
							Name:       "grpc-otel",
							Protocol:   corev1.Protocol("TCP"),
							Port:       int32(4317),
							TargetPort: intstr.FromString("grpc-otel"),
						},
					},
					InternalTrafficPolicy: &[]corev1.ServiceInternalTrafficPolicyType{corev1.ServiceInternalTrafficPolicyLocal}[0],
					Selector: map[string]string{
						"app.kubernetes.io/name":       appName,
						"app.kubernetes.io/instance":   appName,
						"app.kubernetes.io/managed-by": operatorName,
						"app.kubernetes.io/component":  agentServiceName,
					},
				},
			}

			result, _ := serviceForAgent(instance.DeepCopy(), logr.Logger{}, scheme.Scheme)
			Expect(result).To(Equal(expected))
		})
	})

	Context("Instance with all parameters", func() {
		It("returns correct Service", func() {
			instance := &v1alpha1.Aperture{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.ApertureSpec{
					Labels:      testMap,
					Annotations: testMap,
					Service: v1alpha1.ServiceSpec{
						Agent: v1alpha1.Service{
							Annotations: testMapTwo,
						},
					},
				},
			}

			expected := &corev1.Service{
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
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{
							Name:       "grpc",
							Protocol:   corev1.Protocol("TCP"),
							Port:       int32(80),
							TargetPort: intstr.FromString("grpc"),
						},
						{
							Name:       "grpc-otel",
							Protocol:   corev1.Protocol("TCP"),
							Port:       int32(4317),
							TargetPort: intstr.FromString("grpc-otel"),
						},
					},
					InternalTrafficPolicy: &[]corev1.ServiceInternalTrafficPolicyType{corev1.ServiceInternalTrafficPolicyLocal}[0],
					Selector: map[string]string{
						"app.kubernetes.io/name":       appName,
						"app.kubernetes.io/instance":   appName,
						"app.kubernetes.io/managed-by": operatorName,
						"app.kubernetes.io/component":  agentServiceName,
					},
				},
			}

			result, _ := serviceForAgent(instance.DeepCopy(), logr.Logger{}, scheme.Scheme)
			Expect(result).To(Equal(expected))
		})
	})
})

var _ = Describe("Test Service Mutate", func() {
	It("Mutate should update required fields only", func() {
		expected := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{},
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Name: test,
						Port: 80,
					},
				},
				Selector: testMap,
			},
		}

		svc := &corev1.Service{}
		err := serviceMutate(svc, expected.Spec)()
		Expect(err).NotTo(HaveOccurred())
		Expect(svc).To(Equal(expected))
	})
})
