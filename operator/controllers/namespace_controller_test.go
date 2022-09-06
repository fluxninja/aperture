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
	"encoding/base64"
	"fmt"
	"reflect"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/fluxninja/aperture/operator/api/v1alpha1"
)

var _ = Describe("Namespace controller", func() {
	Context("testing Reconcile", func() {
		var instance *v1alpha1.Agent

		BeforeEach(func() {
			instance = defaultAgentInstance.DeepCopy()
			instance.Status.Resources = "created"
		})

		It("should not create resources when namespace is deleted", func() {
			res, err := namespaceReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name: test,
				},
			})
			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).ToNot(HaveOccurred())
		})

		It("should not create resources when creation time is greater than 5 seconds", func() {
			namespace := test + "1"
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			Expect(k8sClient.Create(ctx, ns)).To(BeNil())
			time.Sleep(5 * time.Second)

			res, err := namespaceReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name: namespace,
				},
			})
			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).ToNot(HaveOccurred())
		})

		It("should not create resources when sidecar injection is disabled", func() {
			namespace := test + "2"
			Expect(k8sClient.Create(ctx, instance)).To(BeNil())

			createdConfigMap := &corev1.ConfigMap{}
			configKey := types.NamespacedName{Name: agentServiceName, Namespace: namespace}

			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			Expect(k8sClient.Create(ctx, ns)).To(BeNil())

			_, err := namespaceReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name: namespace,
				},
			})
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() bool {
				err := k8sClient.Get(ctx, configKey, createdConfigMap)
				return errors.IsNotFound(err)
			}, time.Second*10, time.Millisecond*250).Should(BeTrue())

			Expect(k8sClient.Delete(ctx, ns)).To(BeNil())
		})

		It("should not create resources when sidecar injection is enabled but label is not present", func() {
			namespace := test + "3"
			instance.Spec.Sidecar.Enabled = true
			Expect(k8sClient.Create(ctx, instance)).To(BeNil())
			instance.Status.Resources = "created"
			Expect(k8sClient.Status().Update(ctx, instance)).To(BeNil())

			createdConfigMap := &corev1.ConfigMap{}
			configKey := types.NamespacedName{Name: agentServiceName, Namespace: namespace}

			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			Expect(k8sClient.Create(ctx, ns)).To(BeNil())

			_, err := namespaceReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name: namespace,
				},
			})
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() bool {
				err := k8sClient.Get(ctx, configKey, createdConfigMap)
				return errors.IsNotFound(err)
			}, time.Second*10, time.Millisecond*250).Should(BeTrue())

			Expect(k8sClient.Delete(ctx, ns)).To(BeNil())
		})

		It("should not create resources when sidecar injection is enabled but instance is not in created state", func() {
			namespace := test + "4"
			instance.Spec.Sidecar.Enabled = true
			Expect(k8sClient.Create(ctx, instance)).To(BeNil())
			instance.Status.Resources = "skipped"
			Expect(k8sClient.Status().Update(ctx, instance)).To(BeNil())

			createdConfigMap := &corev1.ConfigMap{}
			configKey := types.NamespacedName{Name: agentServiceName, Namespace: namespace}

			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			Expect(k8sClient.Create(ctx, ns)).To(BeNil())

			_, err := namespaceReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name: namespace,
				},
			})
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() bool {
				err := k8sClient.Get(ctx, configKey, createdConfigMap)
				return errors.IsNotFound(err)
			}, time.Second*10, time.Millisecond*250).Should(BeTrue())

			Expect(k8sClient.Delete(ctx, ns)).To(BeNil())
		})

		It("should create ConfigMap when sidecar injection is enabled", func() {
			namespace := test + "5"
			instance.Spec.Sidecar.Enabled = true
			Expect(k8sClient.Create(ctx, instance)).To(BeNil())
			instance.Status.Resources = "created"
			Expect(k8sClient.Status().Update(ctx, instance)).To(BeNil())

			createdConfigMap := &corev1.ConfigMap{}
			configKey := types.NamespacedName{Name: agentServiceName, Namespace: namespace}

			createdSecret := &corev1.Secret{}
			secretKey := types.NamespacedName{Name: agentServiceName, Namespace: namespace}

			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
					Labels: map[string]string{
						sidecarLabelKey: enabled,
					},
				},
			}
			Expect(k8sClient.Create(ctx, ns)).To(BeNil())

			_, err := namespaceReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name: namespace,
				},
			})
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() bool {
				err1 := k8sClient.Get(ctx, configKey, createdConfigMap)
				err2 := k8sClient.Get(ctx, secretKey, createdSecret)
				return err1 == nil && errors.IsNotFound(err2)
			}, time.Second*10, time.Millisecond*250).Should(BeTrue())

			Expect(k8sClient.Delete(ctx, ns)).To(BeNil())
		})

		It("should create ConfigMap and Secret when sidecar injection and FluxNinja plugin are enabled", func() {
			namespace := test + "6"

			instance.Spec.Sidecar.Enabled = true
			instance.Spec.Sidecar.EnableNamespaceByDefault = []string{namespace}
			instance.Spec.Secrets.FluxNinjaPlugin.Create = true
			instance.Spec.Secrets.FluxNinjaPlugin.Value = fmt.Sprintf("enc::%s::enc", base64.StdEncoding.EncodeToString([]byte(test)))
			Expect(k8sClient.Create(ctx, instance)).To(BeNil())
			instance.Status.Resources = "created"
			Expect(k8sClient.Status().Update(ctx, instance)).To(BeNil())

			createdConfigMap := &corev1.ConfigMap{}
			configKey := types.NamespacedName{Name: agentServiceName, Namespace: namespace}

			createdSecret := &corev1.Secret{}
			secretKey := types.NamespacedName{Name: secretName(appName, "agent", &instance.Spec.Secrets.FluxNinjaPlugin), Namespace: namespace}

			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			Expect(k8sClient.Create(ctx, ns)).To(BeNil())

			_, err := namespaceReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name: namespace,
				},
			})
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() bool {
				err1 := k8sClient.Get(ctx, configKey, createdConfigMap)
				err2 := k8sClient.Get(ctx, secretKey, createdSecret)
				return err1 == nil && err2 == nil
			}, time.Second*10, time.Millisecond*250).Should(BeTrue())
			Expect(createdSecret.Data["apiKey"]).To(Equal([]byte(test)))
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: namespace}, ns)).To(BeNil())
			Expect(ns.Labels[sidecarLabelKey]).To(Equal(enabled))

			Expect(k8sClient.Delete(ctx, ns)).To(BeNil())
		})

		AfterEach(func() {
			k8sClient.Delete(ctx, defaultAgentInstance.DeepCopy())
		})
	})

	Context("testing namespaceEventFilters", func() {
		It("should allow only namespace Events in create", func() {
			nsEvent := event.CreateEvent{
				Object: &corev1.Namespace{},
			}
			cmEvent := event.CreateEvent{
				Object: &corev1.ConfigMap{},
			}
			secretEvent := event.CreateEvent{
				Object: &corev1.Secret{},
			}

			pred := namespaceEventFilters()

			Expect(pred.Create(nsEvent)).To(Equal(true))
			Expect(pred.Create(cmEvent)).To(Equal(false))
			Expect(pred.Create(secretEvent)).To(Equal(false))
		})

		It("should allow only namespace Events in update when label is added", func() {
			nsEventValid := event.UpdateEvent{
				ObjectOld: &corev1.Namespace{},
				ObjectNew: &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							sidecarLabelKey: enabled,
						},
					},
				},
			}

			nsEventInvalid1 := event.UpdateEvent{
				ObjectOld: &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							sidecarLabelKey: enabled,
						},
					},
				},
				ObjectNew: &corev1.Namespace{},
			}

			nsEventInvalid2 := event.UpdateEvent{
				ObjectOld: &corev1.Namespace{},
				ObjectNew: &corev1.Namespace{},
			}

			pred := namespaceEventFilters()

			Expect(pred.Update(nsEventValid)).To(Equal(true))
			Expect(pred.Update(nsEventInvalid1)).To(Equal(false))
			Expect(pred.Update(nsEventInvalid2)).To(Equal(false))
		})

		It("should allow only ConfigMap Events in update when data is changed", func() {
			cmEventValid := event.UpdateEvent{
				ObjectOld: &corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app.kubernetes.io/component": agentServiceName,
						},
					},
					Data: testMap,
				},
				ObjectNew: &corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app.kubernetes.io/component": agentServiceName,
						},
					},
					Data: testMapTwo,
				},
			}

			cmEventInvalid1 := event.UpdateEvent{
				ObjectOld: &corev1.ConfigMap{
					Data: testMap,
				},
				ObjectNew: &corev1.ConfigMap{
					Data: testMapTwo,
				},
			}

			cmEventInvalid2 := event.UpdateEvent{
				ObjectOld: &corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app.kubernetes.io/component": agentServiceName,
						},
					},
					Data: testMap,
				},
				ObjectNew: &corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app.kubernetes.io/component": agentServiceName,
						},
					},
					Data: testMap,
				},
			}

			pred := namespaceEventFilters()

			Expect(pred.Update(cmEventValid)).To(Equal(true))
			Expect(pred.Update(cmEventInvalid1)).To(Equal(false))
			Expect(pred.Update(cmEventInvalid2)).To(Equal(false))
		})

		It("should allow only Secret Events in update when data is changed", func() {
			secretEventValid := event.UpdateEvent{
				ObjectOld: &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app.kubernetes.io/component": agentServiceName,
						},
					},
					Data: map[string][]byte{
						test: []byte(test),
					},
				},
				ObjectNew: &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app.kubernetes.io/component": agentServiceName,
						},
					},
					Data: map[string][]byte{
						test: []byte(testTwo),
					},
				},
			}

			secretEventInvalid1 := event.UpdateEvent{
				ObjectOld: &corev1.Secret{
					Data: map[string][]byte{
						test: []byte(test),
					},
				},
				ObjectNew: &corev1.Secret{
					Data: map[string][]byte{
						test: []byte(testTwo),
					},
				},
			}

			secretEventInvalid2 := event.UpdateEvent{
				ObjectOld: &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app.kubernetes.io/component": agentServiceName,
						},
					},
					Data: map[string][]byte{
						test: []byte(test),
					},
				},
				ObjectNew: &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app.kubernetes.io/component": agentServiceName,
						},
					},
					Data: map[string][]byte{
						test: []byte(test),
					},
				},
			}

			pred := namespaceEventFilters()

			Expect(pred.Update(secretEventValid)).To(Equal(true))
			Expect(pred.Update(secretEventInvalid1)).To(Equal(false))
			Expect(pred.Update(secretEventInvalid2)).To(Equal(false))
		})
	})

	It("should not allow namespace Events in delete", func() {
		nsEvent := event.DeleteEvent{
			Object: &corev1.Namespace{},
		}

		pred := namespaceEventFilters()

		Expect(pred.Delete(nsEvent)).To(Equal(false))
	})

	It("should allow ConfigMap Events in delete when proper labels are present", func() {
		cmEventValid := event.DeleteEvent{
			Object: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app.kubernetes.io/component": agentServiceName,
					},
				},
			},
		}

		cmEventInvalid1 := event.DeleteEvent{
			Object: &corev1.ConfigMap{},
		}

		cmEventInvalid2 := event.DeleteEvent{
			Object: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app.kubernetes.io/component": controllerServiceName,
					},
				},
			},
		}

		pred := namespaceEventFilters()

		Expect(pred.Delete(cmEventValid)).To(Equal(true))
		Expect(pred.Delete(cmEventInvalid1)).To(Equal(false))
		Expect(pred.Delete(cmEventInvalid2)).To(Equal(false))
	})

	It("should allow Secret Events in delete when proper labels are present", func() {
		secretEventValid := event.DeleteEvent{
			Object: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app.kubernetes.io/component": agentServiceName,
					},
				},
			},
		}

		secretEventInvalid1 := event.DeleteEvent{
			Object: &corev1.Secret{},
		}

		secretEventInvalid2 := event.DeleteEvent{
			Object: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app.kubernetes.io/component": controllerServiceName,
					},
				},
			},
		}

		pred := namespaceEventFilters()

		Expect(pred.Delete(secretEventValid)).To(Equal(true))
		Expect(pred.Delete(secretEventInvalid1)).To(Equal(false))
		Expect(pred.Delete(secretEventInvalid2)).To(Equal(false))
	})
})
