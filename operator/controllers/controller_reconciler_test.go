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
	"bytes"
	"encoding/pem"
	"fmt"
	"os"
	"reflect"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/fluxninja/aperture/operator/api/v1alpha1"
)

var _ = Describe("Controller Reconciler", Ordered, func() {
	Context("testing Reconcile", func() {
		var instance *v1alpha1.Controller
		var reconciler *ControllerReconciler

		BeforeEach(func() {
			instance = defaultControllerInstance.DeepCopy()
			instance.Name = test
			instance.Namespace = test
			reconciler = &ControllerReconciler{
				Client:   k8sClient,
				Scheme:   scheme.Scheme,
				Recorder: k8sManager.GetEventRecorderFor(appName),
			}
		})

		It("should not create resources when Controller is deleted", func() {
			res, err := reconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      test,
					Namespace: test,
				},
			})
			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).ToNot(HaveOccurred())
		})

		It("should create required resources when Controller is created with default parameters", func() {
			namespace := test + "31"
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			Expect(k8sClient.Create(ctx, ns)).To(BeNil())

			instance.Namespace = namespace
			Expect(k8sClient.Create(ctx, instance)).To(BeNil())

			res, err := reconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      test,
					Namespace: namespace,
				},
			})

			createdControllerConfigMap := &corev1.ConfigMap{}
			controllerConfigKey := types.NamespacedName{Name: controllerServiceName, Namespace: namespace}

			createdControllerService := &corev1.Service{}
			controllerServiceKey := types.NamespacedName{Name: controllerServiceName, Namespace: namespace}

			createdControllerWebhookService := &corev1.Service{}
			controllerWebhookServiceKey := types.NamespacedName{Name: validatingWebhookServiceName, Namespace: namespace}

			createdClusterRole := &rbacv1.ClusterRole{}
			clusterRoleKey := types.NamespacedName{Name: appName}

			createdClusterRoleBinding := &rbacv1.ClusterRoleBinding{}
			clusterRoleBindingKey := types.NamespacedName{Name: controllerServiceName}

			createdControllerServiceAccount := &corev1.ServiceAccount{}
			controllerServiceAccountKey := types.NamespacedName{Name: controllerServiceName, Namespace: namespace}

			createdControllerDeployment := &appsv1.Deployment{}
			controllerDeploymentKey := types.NamespacedName{Name: controllerServiceName, Namespace: namespace}

			createdVWC := &admissionregistrationv1.ValidatingWebhookConfiguration{}
			vwcKey := types.NamespacedName{Name: validatingWebhookServiceName}

			createdControllerSecret := &corev1.Secret{}
			controllerSecretKey := types.NamespacedName{Name: secretName(test, "controller", &instance.Spec.FluxNinjaPlugin.APIKeySecret), Namespace: namespace}

			createdControllerCertSecret := &corev1.Secret{}
			controllerCertSecretKey := types.NamespacedName{Name: fmt.Sprintf("%s-controller-cert", instance.GetName()), Namespace: namespace}

			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() bool {
				err1 := k8sClient.Get(ctx, controllerConfigKey, createdControllerConfigMap)
				err2 := k8sClient.Get(ctx, controllerServiceKey, createdControllerService)
				err3 := k8sClient.Get(ctx, controllerWebhookServiceKey, createdControllerWebhookService)
				err4 := k8sClient.Get(ctx, clusterRoleKey, createdClusterRole)
				err5 := k8sClient.Get(ctx, clusterRoleBindingKey, createdClusterRoleBinding)
				err6 := k8sClient.Get(ctx, controllerServiceAccountKey, createdControllerServiceAccount)
				err7 := k8sClient.Get(ctx, controllerDeploymentKey, createdControllerDeployment)
				err8 := k8sClient.Get(ctx, vwcKey, createdVWC)
				err9 := k8sClient.Get(ctx, controllerSecretKey, createdControllerSecret)
				err10 := k8sClient.Get(ctx, controllerCertSecretKey, createdControllerCertSecret)
				return err1 == nil && err2 == nil && err3 == nil && err4 == nil && err5 == nil &&
					err6 == nil && err7 == nil && err8 == nil && err9 != nil && err10 == nil
			}, time.Second*10, time.Millisecond*250).Should(BeTrue())

			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: test, Namespace: namespace}, instance)).To(BeNil())
			Expect(instance.Status.Resources).To(Equal("created"))

			Expect(k8sClient.Delete(ctx, ns)).To(BeNil())
		})

		It("should create required resources when Controller is created with all parameters", func() {
			namespace := test + "32"
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			Expect(k8sClient.Create(ctx, ns)).To(BeNil())

			instance.Namespace = namespace
			instance.Spec.FluxNinjaPlugin.Enabled = true
			instance.Spec.FluxNinjaPlugin.APIKeySecret.Create = true
			instance.Spec.FluxNinjaPlugin.APIKeySecret.Value = test
			Expect(k8sClient.Create(ctx, instance)).To(BeNil())

			res, err := reconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      test,
					Namespace: namespace,
				},
			})

			createdControllerConfigMap := &corev1.ConfigMap{}
			controllerConfigKey := types.NamespacedName{Name: controllerServiceName, Namespace: namespace}

			createdControllerService := &corev1.Service{}
			controllerServiceKey := types.NamespacedName{Name: controllerServiceName, Namespace: namespace}

			createdControllerWebhookService := &corev1.Service{}
			controllerWebhookServiceKey := types.NamespacedName{Name: validatingWebhookServiceName, Namespace: namespace}

			createdClusterRole := &rbacv1.ClusterRole{}
			clusterRoleKey := types.NamespacedName{Name: appName}

			createdClusterRoleBinding := &rbacv1.ClusterRoleBinding{}
			clusterRoleBindingKey := types.NamespacedName{Name: controllerServiceName}

			createdControllerServiceAccount := &corev1.ServiceAccount{}
			controllerServiceAccountKey := types.NamespacedName{Name: controllerServiceName, Namespace: namespace}

			createdControllerDeployment := &appsv1.Deployment{}
			controllerDeploymentKey := types.NamespacedName{Name: controllerServiceName, Namespace: namespace}

			createdVWC := &admissionregistrationv1.ValidatingWebhookConfiguration{}
			vwcKey := types.NamespacedName{Name: validatingWebhookServiceName}

			createdControllerSecret := &corev1.Secret{}
			controllerSecretKey := types.NamespacedName{Name: secretName(test, "controller", &instance.Spec.FluxNinjaPlugin.APIKeySecret), Namespace: namespace}

			createdControllerCertSecret := &corev1.Secret{}
			controllerCertSecretKey := types.NamespacedName{Name: fmt.Sprintf("%s-controller-cert", instance.GetName()), Namespace: namespace}

			Eventually(func() bool {
				err1 := k8sClient.Get(ctx, controllerConfigKey, createdControllerConfigMap)
				err2 := k8sClient.Get(ctx, controllerServiceKey, createdControllerService)
				err3 := k8sClient.Get(ctx, controllerWebhookServiceKey, createdControllerWebhookService)
				err4 := k8sClient.Get(ctx, clusterRoleKey, createdClusterRole)
				err5 := k8sClient.Get(ctx, clusterRoleBindingKey, createdClusterRoleBinding)
				err6 := k8sClient.Get(ctx, controllerServiceAccountKey, createdControllerServiceAccount)
				err7 := k8sClient.Get(ctx, controllerDeploymentKey, createdControllerDeployment)
				err8 := k8sClient.Get(ctx, vwcKey, createdVWC)
				err9 := k8sClient.Get(ctx, controllerSecretKey, createdControllerSecret)
				err10 := k8sClient.Get(ctx, controllerCertSecretKey, createdControllerCertSecret)
				return err1 == nil && err2 == nil && err3 == nil && err4 == nil && err5 == nil &&
					err6 == nil && err7 == nil && err8 == nil && err9 == nil && err10 == nil
			}, time.Second*10, time.Millisecond*250).Should(BeTrue())

			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).ToNot(HaveOccurred())

			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: test, Namespace: namespace}, instance)).To(BeNil())
			Expect(instance.Status.Resources).To(Equal("created"))
			Expect(instance.Finalizers).To(Equal([]string{finalizerName}))
			Expect(instance.Spec.FluxNinjaPlugin.APIKeySecret.Create).To(BeFalse())
			Expect(instance.Spec.FluxNinjaPlugin.APIKeySecret.Value).To(Equal(""))

			Expect(k8sClient.Delete(ctx, ns)).To(BeNil())
		})

		It("should not create required resources when an Controller instance is already created", func() {
			namespace := test + "33"
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			Expect(k8sClient.Create(ctx, ns)).To(BeNil())

			instance.Namespace = namespace
			Expect(k8sClient.Create(ctx, instance)).To(BeNil())

			res, err := reconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      test,
					Namespace: namespace,
				},
			})
			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).ToNot(HaveOccurred())

			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: test, Namespace: namespace}, instance)).To(BeNil())
			Expect(instance.Status.Resources).To(Equal("created"))

			instanceNew := defaultControllerInstance.DeepCopy()
			instanceNew.Name = testTwo
			instanceNew.Namespace = namespace

			Expect(k8sClient.Create(ctx, instanceNew)).To(BeNil())

			res, err = reconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      testTwo,
					Namespace: namespace,
				},
			})
			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).ToNot(HaveOccurred())

			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: testTwo, Namespace: namespace}, instanceNew)).To(BeNil())
			Expect(instanceNew.Status.Resources).To(Equal("skipped"))

			Expect(k8sClient.Delete(ctx, ns)).To(BeNil())
		})

		It("should delete required resources when an Controller instance is already deleted", func() {
			namespace := test + "34"
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			Expect(k8sClient.Create(ctx, ns)).To(BeNil())

			instance.Namespace = namespace
			instance.Spec.FluxNinjaPlugin.Enabled = true
			instance.Spec.CommonSpec.ServiceAccountSpec.Create = false
			instance.Spec.FluxNinjaPlugin.APIKeySecret.Create = true
			instance.Spec.FluxNinjaPlugin.APIKeySecret.Value = test

			os.Setenv("APERTURE_OPERATOR_CERT_DIR", certDir)
			os.Setenv("APERTURE_OPERATOR_CERT_NAME", "tls6.crt")
			os.Setenv("APERTURE_OPERATOR_NAMESPACE", appName)
			os.Setenv("APERTURE_OPERATOR_SERVICE_NAME", appName)
			certPath := fmt.Sprintf("%s/%s", os.Getenv("APERTURE_OPERATOR_CERT_DIR"), os.Getenv("APERTURE_OPERATOR_CERT_NAME"))
			serverCertPEM := new(bytes.Buffer)
			err := pem.Encode(serverCertPEM, &pem.Block{
				Type:  "CERTIFICATE",
				Bytes: []byte(test),
			})
			Expect(err).NotTo(HaveOccurred())

			err = writeFile(certPath, serverCertPEM)
			Expect(err).NotTo(HaveOccurred())
			Expect(k8sClient.Create(ctx, instance)).To(BeNil())

			res, err := reconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      test,
					Namespace: namespace,
				},
			})
			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).ToNot(HaveOccurred())

			createdClusterRoleBinding := &rbacv1.ClusterRoleBinding{}
			clusterRoleBindingKey := types.NamespacedName{Name: controllerServiceName}

			Eventually(func() bool {
				return k8sClient.Get(ctx, clusterRoleBindingKey, createdClusterRoleBinding) == nil
			}, time.Second*10, time.Millisecond*250).Should(BeTrue())

			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: test, Namespace: namespace}, instance)).To(BeNil())
			Expect(instance.Status.Resources).To(Equal("created"))

			Expect(k8sClient.Delete(ctx, instance)).To(BeNil())

			res, err = reconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      test,
					Namespace: namespace,
				},
			})
			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() bool {
				return k8sClient.Get(ctx, clusterRoleBindingKey, createdClusterRoleBinding) != nil
			}, time.Second*10, time.Millisecond*250).Should(BeTrue())

			Expect(k8sClient.Delete(ctx, ns)).To(BeNil())
		})

		AfterEach(func() {
			_ = k8sClient.Delete(ctx, instance)
			_, err := reconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      test,
					Namespace: test,
				},
			})
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("testing eventFilters", func() {
		It("should allow only Controller Events in create", func() {
			ControllerEvent := event.CreateEvent{
				Object: &v1alpha1.Controller{},
			}
			cmEvent := event.CreateEvent{
				Object: &corev1.ConfigMap{},
			}

			pred := eventFiltersForController()

			Expect(pred.Create(ControllerEvent)).To(Equal(true))
			Expect(pred.Create(cmEvent)).To(Equal(false))
		})

		It("should allow only Controller Events in update when Spec is changed or instance is deleted", func() {
			ControllerEventValid1 := event.UpdateEvent{
				ObjectOld: &v1alpha1.Controller{},
				ObjectNew: &v1alpha1.Controller{
					Spec: v1alpha1.ControllerSpec{
						CommonSpec: v1alpha1.CommonSpec{
							FluxNinjaPlugin: v1alpha1.FluxNinjaPluginSpec{
								Enabled: true,
							},
						},
					},
				},
			}

			ControllerEventValid2 := event.UpdateEvent{
				ObjectOld: &v1alpha1.Controller{},
				ObjectNew: &v1alpha1.Controller{
					ObjectMeta: metav1.ObjectMeta{
						DeletionTimestamp: &metav1.Time{Time: time.Now()},
					},
				},
			}

			ControllerEventInvalid1 := event.UpdateEvent{
				ObjectOld: defaultControllerInstance,
				ObjectNew: defaultControllerInstance,
			}

			ControllerEventInvalid2 := event.UpdateEvent{
				ObjectOld: &corev1.ConfigMap{},
				ObjectNew: &corev1.ConfigMap{},
			}

			pred := eventFiltersForController()

			Expect(pred.Update(ControllerEventValid1)).To(Equal(true))
			Expect(pred.Update(ControllerEventValid2)).To(Equal(true))
			Expect(pred.Update(ControllerEventInvalid1)).To(Equal(false))
			Expect(pred.Update(ControllerEventInvalid2)).To(Equal(true))
		})

		It("should allow only Controller Events in update when secret is changed", func() {
			ControllerEventValid1 := event.UpdateEvent{
				ObjectOld: &v1alpha1.Controller{
					Spec: v1alpha1.ControllerSpec{
						CommonSpec: v1alpha1.CommonSpec{
							FluxNinjaPlugin: v1alpha1.FluxNinjaPluginSpec{
								APIKeySecret: v1alpha1.APIKeySecret{
									Value: test,
								},
							},
						},
					},
				},
				ObjectNew: &v1alpha1.Controller{
					Spec: v1alpha1.ControllerSpec{
						CommonSpec: v1alpha1.CommonSpec{
							FluxNinjaPlugin: v1alpha1.FluxNinjaPluginSpec{
								APIKeySecret: v1alpha1.APIKeySecret{
									Value: "",
								},
							},
						},
					},
				},
			}

			pred := eventFiltersForController()

			Expect(pred.Update(ControllerEventValid1)).To(Equal(false))
		})

		It("should allow Events in delete when state is known", func() {
			ControllerEvent := event.DeleteEvent{
				Object:             &v1alpha1.Controller{},
				DeleteStateUnknown: true,
			}

			ControllerEventInvalid := event.DeleteEvent{
				Object:             &v1alpha1.Controller{},
				DeleteStateUnknown: false,
			}

			pred := eventFiltersForController()

			Expect(pred.Delete(ControllerEvent)).To(Equal(false))
			Expect(pred.Delete(ControllerEventInvalid)).To(Equal(true))
		})
	})
})
