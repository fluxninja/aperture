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

package namespace

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"time"

	. "github.com/fluxninja/aperture/operator/controllers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	agentv1alpha1 "github.com/fluxninja/aperture/operator/api/agent/v1alpha1"
)

var _ = Describe("Namespace controller", func() {
	Context("testing Reconcile", func() {
		var instance *agentv1alpha1.Agent

		BeforeEach(func() {
			instance = DefaultAgentInstance.DeepCopy()
			instance.Status.Resources = "created"
		})

		It("should not create resources when namespace is deleted", func() {
			res, err := namespaceTestReconciler.Reconcile(Ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name: Test,
				},
			})
			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).ToNot(HaveOccurred())
		})

		It("should not create resources when creation time is greater than 5 seconds", func() {
			namespace := Test + "1"
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			Expect(K8sClient.Create(Ctx, ns)).To(BeNil())
			time.Sleep(5 * time.Second)

			res, err := namespaceTestReconciler.Reconcile(Ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name: namespace,
				},
			})
			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).ToNot(HaveOccurred())
		})

		It("should not create resources when sidecar injection is disabled", func() {
			namespace := Test + "2"
			Expect(K8sClient.Create(Ctx, instance)).To(BeNil())

			createdConfigMap := &corev1.ConfigMap{}
			configKey := types.NamespacedName{Name: AgentServiceName, Namespace: namespace}

			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			Expect(K8sClient.Create(Ctx, ns)).To(BeNil())

			_, err := namespaceTestReconciler.Reconcile(Ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name: namespace,
				},
			})
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() bool {
				err := K8sClient.Get(Ctx, configKey, createdConfigMap)
				return errors.IsNotFound(err)
			}, time.Second*10, time.Millisecond*250).Should(BeTrue())

			Expect(K8sClient.Delete(Ctx, ns)).To(BeNil())
		})

		It("should not create resources when sidecar injection is enabled but label is not present", func() {
			namespace := Test + "3"
			instance.Spec.Sidecar.Enabled = true
			Expect(K8sClient.Create(Ctx, instance)).To(BeNil())
			instance.Status.Resources = "created"
			Expect(K8sClient.Status().Update(Ctx, instance)).To(BeNil())

			createdConfigMap := &corev1.ConfigMap{}
			configKey := types.NamespacedName{Name: AgentServiceName, Namespace: namespace}

			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			Expect(K8sClient.Create(Ctx, ns)).To(BeNil())

			_, err := namespaceTestReconciler.Reconcile(Ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name: namespace,
				},
			})
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() bool {
				err := K8sClient.Get(Ctx, configKey, createdConfigMap)
				return errors.IsNotFound(err)
			}, time.Second*10, time.Millisecond*250).Should(BeTrue())

			Expect(K8sClient.Delete(Ctx, ns)).To(BeNil())
		})

		It("should not create resources when sidecar injection is enabled but instance is not in created state", func() {
			namespace := Test + "4"
			instance.Spec.Sidecar.Enabled = true
			Expect(K8sClient.Create(Ctx, instance)).To(BeNil())
			instance.Status.Resources = "skipped"
			Expect(K8sClient.Status().Update(Ctx, instance)).To(BeNil())

			createdConfigMap := &corev1.ConfigMap{}
			configKey := types.NamespacedName{Name: AgentServiceName, Namespace: namespace}

			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			Expect(K8sClient.Create(Ctx, ns)).To(BeNil())

			_, err := namespaceTestReconciler.Reconcile(Ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name: namespace,
				},
			})
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() bool {
				err := K8sClient.Get(Ctx, configKey, createdConfigMap)
				return errors.IsNotFound(err)
			}, time.Second*10, time.Millisecond*250).Should(BeTrue())

			Expect(K8sClient.Delete(Ctx, ns)).To(BeNil())
		})

		It("should create ConfigMap when sidecar injection is enabled", func() {
			namespace := Test + "5"
			instance.Spec.Sidecar.Enabled = true
			Expect(K8sClient.Create(Ctx, instance)).To(BeNil())
			instance.Status.Resources = "created"
			Expect(K8sClient.Status().Update(Ctx, instance)).To(BeNil())

			createdConfigMap := &corev1.ConfigMap{}
			configKey := types.NamespacedName{Name: AgentServiceName, Namespace: namespace}

			createdSecret := &corev1.Secret{}
			secretKey := types.NamespacedName{Name: AgentServiceName, Namespace: namespace}

			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
					Labels: map[string]string{
						SidecarLabelKey: Enabled,
					},
				},
			}
			Expect(K8sClient.Create(Ctx, ns)).To(BeNil())

			_, err := namespaceTestReconciler.Reconcile(Ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name: namespace,
				},
			})
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() bool {
				err1 := K8sClient.Get(Ctx, configKey, createdConfigMap)
				err2 := K8sClient.Get(Ctx, secretKey, createdSecret)
				return err1 == nil && errors.IsNotFound(err2)
			}, time.Second*10, time.Millisecond*250).Should(BeTrue())

			Expect(K8sClient.Delete(Ctx, ns)).To(BeNil())
		})

		It("should create ConfigMap and Secret when sidecar injection and FluxNinja plugin are enabled", func() {
			namespace := Test + "6"

			instance.Spec.Sidecar.Enabled = true
			instance.Spec.Sidecar.EnableNamespaceByDefault = []string{namespace}
			instance.Spec.Secrets.FluxNinjaPlugin.Create = true
			instance.Spec.Secrets.FluxNinjaPlugin.Value = fmt.Sprintf("enc::%s::enc", base64.StdEncoding.EncodeToString([]byte(Test)))
			Expect(K8sClient.Create(Ctx, instance)).To(BeNil())
			instance.Status.Resources = "created"
			Expect(K8sClient.Status().Update(Ctx, instance)).To(BeNil())

			createdConfigMap := &corev1.ConfigMap{}
			configKey := types.NamespacedName{Name: AgentServiceName, Namespace: namespace}

			createdSecret := &corev1.Secret{}
			secretKey := types.NamespacedName{Name: SecretName(AppName, "agent", &instance.Spec.Secrets.FluxNinjaPlugin), Namespace: namespace}

			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			Expect(K8sClient.Create(Ctx, ns)).To(BeNil())

			_, err := namespaceTestReconciler.Reconcile(Ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name: namespace,
				},
			})
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() bool {
				err1 := K8sClient.Get(Ctx, configKey, createdConfigMap)
				err2 := K8sClient.Get(Ctx, secretKey, createdSecret)
				return err1 == nil && err2 == nil
			}, time.Second*10, time.Millisecond*250).Should(BeTrue())
			Expect(createdSecret.Data["apiKey"]).To(Equal([]byte(Test)))
			Expect(K8sClient.Get(Ctx, types.NamespacedName{Name: namespace}, ns)).To(BeNil())
			Expect(ns.Labels[SidecarLabelKey]).To(Equal(Enabled))

			Expect(K8sClient.Delete(Ctx, ns)).To(BeNil())
		})

		AfterEach(func() {
			K8sClient.Delete(Ctx, DefaultAgentInstance.DeepCopy())
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
							SidecarLabelKey: Enabled,
						},
					},
				},
			}

			nsEventInvalid1 := event.UpdateEvent{
				ObjectOld: &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							SidecarLabelKey: Enabled,
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
							"app.kubernetes.io/component": AgentServiceName,
						},
					},
					Data: TestMap,
				},
				ObjectNew: &corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app.kubernetes.io/component": AgentServiceName,
						},
					},
					Data: TestMapTwo,
				},
			}

			cmEventInvalid1 := event.UpdateEvent{
				ObjectOld: &corev1.ConfigMap{
					Data: TestMap,
				},
				ObjectNew: &corev1.ConfigMap{
					Data: TestMapTwo,
				},
			}

			cmEventInvalid2 := event.UpdateEvent{
				ObjectOld: &corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app.kubernetes.io/component": AgentServiceName,
						},
					},
					Data: TestMap,
				},
				ObjectNew: &corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app.kubernetes.io/component": AgentServiceName,
						},
					},
					Data: TestMap,
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
							"app.kubernetes.io/component": AgentServiceName,
						},
					},
					Data: map[string][]byte{
						Test: []byte(Test),
					},
				},
				ObjectNew: &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app.kubernetes.io/component": AgentServiceName,
						},
					},
					Data: map[string][]byte{
						Test: []byte(TestTwo),
					},
				},
			}

			secretEventInvalid1 := event.UpdateEvent{
				ObjectOld: &corev1.Secret{
					Data: map[string][]byte{
						Test: []byte(Test),
					},
				},
				ObjectNew: &corev1.Secret{
					Data: map[string][]byte{
						Test: []byte(TestTwo),
					},
				},
			}

			secretEventInvalid2 := event.UpdateEvent{
				ObjectOld: &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app.kubernetes.io/component": AgentServiceName,
						},
					},
					Data: map[string][]byte{
						Test: []byte(Test),
					},
				},
				ObjectNew: &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app.kubernetes.io/component": AgentServiceName,
						},
					},
					Data: map[string][]byte{
						Test: []byte(Test),
					},
				},
			}

			pred := namespaceEventFilters()

			Expect(pred.Update(secretEventValid)).To(Equal(true))
			Expect(pred.Update(secretEventInvalid1)).To(Equal(false))
			Expect(pred.Update(secretEventInvalid2)).To(Equal(false))
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
							"app.kubernetes.io/component": AgentServiceName,
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
							"app.kubernetes.io/component": ControllerServiceName,
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
							"app.kubernetes.io/component": AgentServiceName,
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
							"app.kubernetes.io/component": ControllerServiceName,
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
})
