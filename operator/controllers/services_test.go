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

	"github.com/fluxninja/aperture/operator/api/v1alpha1"
	"github.com/fluxninja/aperture/pkg/distcache"
	"github.com/fluxninja/aperture/pkg/net/listener"
	"github.com/fluxninja/aperture/pkg/otel"
)

var _ = Describe("Service for Controller", func() {
	Context("Instance with default parameters", func() {
		It("returns correct Service", func() {
			instance := &v1alpha1.Controller{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.ControllerSpec{
					ConfigSpec: v1alpha1.ControllerConfigSpec{
						CommonConfigSpec: v1alpha1.CommonConfigSpec{
							Server: v1alpha1.ServerConfigSpec{
								ListenerConfig: listener.ListenerConfig{
									Addr: ":8080",
								},
							},
							Otel: otel.OtelConfig{
								GRPCAddr: ":4317",
								HTTPAddr: ":4318",
							},
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
					},
					Annotations: nil,
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion:         "fluxninja.com/v1alpha1",
							Name:               instance.GetName(),
							Kind:               "Controller",
							Controller:         pointer.BoolPtr(true),
							BlockOwnerDeletion: pointer.BoolPtr(true),
						},
					},
				},
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{
							Name:       server,
							Protocol:   corev1.Protocol(tcp),
							Port:       int32(8080),
							TargetPort: intstr.FromString(server),
						},
						{
							Name:       grpcOtel,
							Protocol:   corev1.Protocol(tcp),
							Port:       int32(4317),
							TargetPort: intstr.FromString(grpcOtel),
						},
						{
							Name:       httpOtel,
							Protocol:   corev1.Protocol(tcp),
							Port:       int32(4318),
							TargetPort: intstr.FromString(httpOtel),
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

			result, err := serviceForController(instance.DeepCopy(), logr.Logger{}, scheme.Scheme)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(expected))
		})
	})

	Context("Instance with all parameters", func() {
		It("returns correct Service", func() {
			instance := &v1alpha1.Controller{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.ControllerSpec{
					CommonSpec: v1alpha1.CommonSpec{
						Labels:      testMap,
						Annotations: testMap,
						Service: v1alpha1.Service{
							Annotations: testMapTwo,
						},
					},
					ConfigSpec: v1alpha1.ControllerConfigSpec{
						CommonConfigSpec: v1alpha1.CommonConfigSpec{
							Server: v1alpha1.ServerConfigSpec{
								ListenerConfig: listener.ListenerConfig{
									Addr: ":8080",
								},
							},
							Otel: otel.OtelConfig{
								GRPCAddr: ":4317",
								HTTPAddr: ":4318",
							},
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
							APIVersion:         "fluxninja.com/v1alpha1",
							Name:               instance.GetName(),
							Kind:               "Controller",
							Controller:         pointer.BoolPtr(true),
							BlockOwnerDeletion: pointer.BoolPtr(true),
						},
					},
				},
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{
							Name:       server,
							Protocol:   corev1.Protocol(tcp),
							Port:       int32(8080),
							TargetPort: intstr.FromString(server),
						},
						{
							Name:       grpcOtel,
							Protocol:   corev1.Protocol(tcp),
							Port:       int32(4317),
							TargetPort: intstr.FromString(grpcOtel),
						},
						{
							Name:       httpOtel,
							Protocol:   corev1.Protocol(tcp),
							Port:       int32(4318),
							TargetPort: intstr.FromString(httpOtel),
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

			result, err := serviceForController(instance.DeepCopy(), logr.Logger{}, scheme.Scheme)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(expected))
		})
	})
})

var _ = Describe("Service for Agent", func() {
	Context("Instance with default parameters", func() {
		It("returns correct Service", func() {
			instance := &v1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.AgentSpec{
					ConfigSpec: v1alpha1.AgentConfigSpec{
						CommonConfigSpec: v1alpha1.CommonConfigSpec{
							Server: v1alpha1.ServerConfigSpec{
								ListenerConfig: listener.ListenerConfig{
									Addr: ":8080",
								},
							},
							Otel: otel.OtelConfig{
								GRPCAddr: ":4317",
								HTTPAddr: ":4318",
							},
						},
						DistCache: distcache.DistCacheConfig{
							BindAddr:           ":3320",
							MemberlistBindAddr: ":3322",
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
					},
					Annotations: nil,
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion:         "fluxninja.com/v1alpha1",
							Name:               instance.GetName(),
							Kind:               "Agent",
							Controller:         pointer.BoolPtr(true),
							BlockOwnerDeletion: pointer.BoolPtr(true),
						},
					},
				},
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{
							Name:       server,
							Protocol:   corev1.Protocol(tcp),
							Port:       int32(8080),
							TargetPort: intstr.FromString(server),
						},
						{
							Name:       grpcOtel,
							Protocol:   corev1.Protocol(tcp),
							Port:       int32(4317),
							TargetPort: intstr.FromString(grpcOtel),
						},
						{
							Name:       httpOtel,
							Protocol:   corev1.Protocol(tcp),
							Port:       int32(4318),
							TargetPort: intstr.FromString(httpOtel),
						},
						{
							Name:       distCache,
							Protocol:   corev1.Protocol(tcp),
							Port:       int32(3320),
							TargetPort: intstr.FromString(distCache),
						},
						{
							Name:       memberList,
							Protocol:   corev1.Protocol(tcp),
							Port:       int32(3322),
							TargetPort: intstr.FromString(memberList),
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

			result, err := serviceForAgent(instance.DeepCopy(), logr.Logger{}, scheme.Scheme)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(expected))
		})
	})

	Context("Instance with all parameters", func() {
		It("returns correct Service", func() {
			instance := &v1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.AgentSpec{
					CommonSpec: v1alpha1.CommonSpec{
						Labels:      testMap,
						Annotations: testMap,
						Service: v1alpha1.Service{
							Annotations: testMapTwo,
						},
					},
					ConfigSpec: v1alpha1.AgentConfigSpec{
						CommonConfigSpec: v1alpha1.CommonConfigSpec{
							Server: v1alpha1.ServerConfigSpec{
								ListenerConfig: listener.ListenerConfig{
									Addr: ":8080",
								},
							},
							Otel: otel.OtelConfig{
								GRPCAddr: ":4317",
								HTTPAddr: ":4318",
							},
						},
						DistCache: distcache.DistCacheConfig{
							BindAddr:           ":3320",
							MemberlistBindAddr: ":3322",
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
							APIVersion:         "fluxninja.com/v1alpha1",
							Name:               instance.GetName(),
							Kind:               "Agent",
							Controller:         pointer.BoolPtr(true),
							BlockOwnerDeletion: pointer.BoolPtr(true),
						},
					},
				},
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{
							Name:       server,
							Protocol:   corev1.Protocol(tcp),
							Port:       int32(8080),
							TargetPort: intstr.FromString(server),
						},
						{
							Name:       grpcOtel,
							Protocol:   corev1.Protocol(tcp),
							Port:       int32(4317),
							TargetPort: intstr.FromString(grpcOtel),
						},
						{
							Name:       httpOtel,
							Protocol:   corev1.Protocol(tcp),
							Port:       int32(4318),
							TargetPort: intstr.FromString(httpOtel),
						},
						{
							Name:       distCache,
							Protocol:   corev1.Protocol(tcp),
							Port:       int32(3320),
							TargetPort: intstr.FromString(distCache),
						},
						{
							Name:       memberList,
							Protocol:   corev1.Protocol(tcp),
							Port:       int32(3322),
							TargetPort: intstr.FromString(memberList),
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

			result, err := serviceForAgent(instance.DeepCopy(), logr.Logger{}, scheme.Scheme)

			Expect(err).NotTo(HaveOccurred())
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
