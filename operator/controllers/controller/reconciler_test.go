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
	"bytes"
	"encoding/pem"
	"fmt"
	"os"
	"reflect"
	"time"

	. "github.com/fluxninja/aperture/operator/controllers"
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

	"github.com/fluxninja/aperture/operator/api/common"
	controllerv1alpha1 "github.com/fluxninja/aperture/operator/api/controller/v1alpha1"
)

var _ = Describe("Controller Reconciler", Ordered, func() {
	Context("testing Reconcile", func() {
		var instance *controllerv1alpha1.Controller
		var reconciler *ControllerReconciler

		BeforeEach(func() {
			instance = DefaultControllerInstance.DeepCopy()
			instance.Name = Test
			instance.Namespace = Test
			reconciler = &ControllerReconciler{
				Client:        K8sClient,
				DynamicClient: K8sDynamicClient,
				Scheme:        scheme.Scheme,
				Recorder:      K8sManager.GetEventRecorderFor(AppName),
			}
		})

		It("should not create resources when Controller is deleted", func() {
			res, err := reconciler.Reconcile(Ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      Test,
					Namespace: Test,
				},
			})
			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).ToNot(HaveOccurred())
		})

		It("should create required resources when Controller is created with default parameters", func() {
			namespace := Test + "31"
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			Expect(K8sClient.Create(Ctx, ns)).To(BeNil())

			instance.Namespace = namespace
			Expect(K8sClient.Create(Ctx, instance)).To(BeNil())

			res, err := reconciler.Reconcile(Ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      Test,
					Namespace: namespace,
				},
			})

			createdControllerConfigMap := &corev1.ConfigMap{}
			controllerConfigKey := types.NamespacedName{Name: ControllerServiceName, Namespace: namespace}

			createdControllerService := &corev1.Service{}
			controllerServiceKey := types.NamespacedName{Name: ControllerServiceName, Namespace: namespace}

			createdClusterRole := &rbacv1.ClusterRole{}
			clusterRoleKey := types.NamespacedName{Name: AppName}

			createdClusterRoleBinding := &rbacv1.ClusterRoleBinding{}
			clusterRoleBindingKey := types.NamespacedName{Name: ControllerServiceName}

			createdControllerServiceAccount := &corev1.ServiceAccount{}
			controllerServiceAccountKey := types.NamespacedName{Name: ControllerServiceName, Namespace: namespace}

			createdControllerDeployment := &appsv1.Deployment{}
			controllerDeploymentKey := types.NamespacedName{Name: ControllerServiceName, Namespace: namespace}

			createdVWC := &admissionregistrationv1.ValidatingWebhookConfiguration{}
			vwcKey := types.NamespacedName{Name: ControllerServiceName}

			createdControllerSecret := &corev1.Secret{}
			controllerSecretKey := types.NamespacedName{Name: SecretName(Test, "controller", &instance.Spec.Secrets.FluxNinjaPlugin), Namespace: namespace}

			createdControllerCertSecret := &corev1.Secret{}
			controllerCertSecretKey := types.NamespacedName{Name: fmt.Sprintf("%s-controller-cert", instance.GetName()), Namespace: namespace}

			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() bool {
				err1 := K8sClient.Get(Ctx, controllerConfigKey, createdControllerConfigMap)
				err2 := K8sClient.Get(Ctx, controllerServiceKey, createdControllerService)
				err3 := K8sClient.Get(Ctx, clusterRoleKey, createdClusterRole)
				err4 := K8sClient.Get(Ctx, clusterRoleBindingKey, createdClusterRoleBinding)
				err5 := K8sClient.Get(Ctx, controllerServiceAccountKey, createdControllerServiceAccount)
				err6 := K8sClient.Get(Ctx, controllerDeploymentKey, createdControllerDeployment)
				err7 := K8sClient.Get(Ctx, vwcKey, createdVWC)
				err8 := K8sClient.Get(Ctx, controllerSecretKey, createdControllerSecret)
				err9 := K8sClient.Get(Ctx, controllerCertSecretKey, createdControllerCertSecret)
				return err1 == nil && err2 == nil && err3 == nil && err4 == nil &&
					err5 == nil && err6 == nil && err7 == nil && err8 != nil && err9 == nil
			}, time.Second*10, time.Millisecond*250).Should(BeTrue())

			Expect(K8sClient.Get(Ctx, types.NamespacedName{Name: Test, Namespace: namespace}, instance)).To(BeNil())
			Expect(instance.Status.Resources).To(Equal("created"))

			Expect(K8sClient.Delete(Ctx, ns)).To(BeNil())
		})

		It("should create required resources when Controller is created with all parameters", func() {
			namespace := Test + "32"
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			Expect(K8sClient.Create(Ctx, ns)).To(BeNil())

			instance.Namespace = namespace
			instance.Spec.Secrets.FluxNinjaPlugin.Create = true
			instance.Spec.Secrets.FluxNinjaPlugin.Value = Test
			Expect(K8sClient.Create(Ctx, instance)).To(BeNil())

			res, err := reconciler.Reconcile(Ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      Test,
					Namespace: namespace,
				},
			})

			createdControllerConfigMap := &corev1.ConfigMap{}
			controllerConfigKey := types.NamespacedName{Name: ControllerServiceName, Namespace: namespace}

			createdControllerService := &corev1.Service{}
			controllerServiceKey := types.NamespacedName{Name: ControllerServiceName, Namespace: namespace}

			createdClusterRole := &rbacv1.ClusterRole{}
			clusterRoleKey := types.NamespacedName{Name: AppName}

			createdClusterRoleBinding := &rbacv1.ClusterRoleBinding{}
			clusterRoleBindingKey := types.NamespacedName{Name: ControllerServiceName}

			createdControllerServiceAccount := &corev1.ServiceAccount{}
			controllerServiceAccountKey := types.NamespacedName{Name: ControllerServiceName, Namespace: namespace}

			createdControllerDeployment := &appsv1.Deployment{}
			controllerDeploymentKey := types.NamespacedName{Name: ControllerServiceName, Namespace: namespace}

			createdVWC := &admissionregistrationv1.ValidatingWebhookConfiguration{}
			vwcKey := types.NamespacedName{Name: ControllerServiceName}

			createdControllerSecret := &corev1.Secret{}
			controllerSecretKey := types.NamespacedName{Name: SecretName(Test, "controller", &instance.Spec.Secrets.FluxNinjaPlugin), Namespace: namespace}

			createdControllerCertSecret := &corev1.Secret{}
			controllerCertSecretKey := types.NamespacedName{Name: fmt.Sprintf("%s-controller-cert", instance.GetName()), Namespace: namespace}

			Eventually(func() bool {
				err1 := K8sClient.Get(Ctx, controllerConfigKey, createdControllerConfigMap)
				err2 := K8sClient.Get(Ctx, controllerServiceKey, createdControllerService)
				err3 := K8sClient.Get(Ctx, clusterRoleKey, createdClusterRole)
				err4 := K8sClient.Get(Ctx, clusterRoleBindingKey, createdClusterRoleBinding)
				err5 := K8sClient.Get(Ctx, controllerServiceAccountKey, createdControllerServiceAccount)
				err6 := K8sClient.Get(Ctx, controllerDeploymentKey, createdControllerDeployment)
				err7 := K8sClient.Get(Ctx, vwcKey, createdVWC)
				err8 := K8sClient.Get(Ctx, controllerSecretKey, createdControllerSecret)
				err9 := K8sClient.Get(Ctx, controllerCertSecretKey, createdControllerCertSecret)
				return err1 == nil && err2 == nil && err3 == nil && err4 == nil &&
					err5 == nil && err6 == nil && err7 == nil && err8 == nil && err9 == nil
			}, time.Second*10, time.Millisecond*250).Should(BeTrue())

			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).ToNot(HaveOccurred())

			Expect(K8sClient.Get(Ctx, types.NamespacedName{Name: Test, Namespace: namespace}, instance)).To(BeNil())
			Expect(instance.Status.Resources).To(Equal("created"))
			Expect(instance.Finalizers).To(Equal([]string{FinalizerName}))
			Expect(instance.Spec.Secrets.FluxNinjaPlugin.Create).To(BeFalse())
			Expect(instance.Spec.Secrets.FluxNinjaPlugin.Value).To(Equal(""))

			Expect(K8sClient.Delete(Ctx, ns)).To(BeNil())
		})

		It("should not create required resources when an Controller instance is already created", func() {
			namespace := Test + "33"
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			Expect(K8sClient.Create(Ctx, ns)).To(BeNil())

			instance.Namespace = namespace
			Expect(K8sClient.Create(Ctx, instance)).To(BeNil())

			res, err := reconciler.Reconcile(Ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      Test,
					Namespace: namespace,
				},
			})
			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).ToNot(HaveOccurred())

			Expect(K8sClient.Get(Ctx, types.NamespacedName{Name: Test, Namespace: namespace}, instance)).To(BeNil())
			Expect(instance.Status.Resources).To(Equal("created"))

			instanceNew := DefaultControllerInstance.DeepCopy()
			instanceNew.Name = TestTwo
			instanceNew.Namespace = namespace

			Expect(K8sClient.Create(Ctx, instanceNew)).To(BeNil())

			res, err = reconciler.Reconcile(Ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      TestTwo,
					Namespace: namespace,
				},
			})
			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).ToNot(HaveOccurred())

			Expect(K8sClient.Get(Ctx, types.NamespacedName{Name: TestTwo, Namespace: namespace}, instanceNew)).To(BeNil())
			Expect(instanceNew.Status.Resources).To(Equal("skipped"))

			Expect(K8sClient.Delete(Ctx, ns)).To(BeNil())
		})

		It("should delete required resources when an Controller instance is already deleted", func() {
			namespace := Test + "34"
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			Expect(K8sClient.Create(Ctx, ns)).To(BeNil())

			instance.Namespace = namespace
			instance.Spec.CommonSpec.ServiceAccountSpec.Create = false
			instance.Spec.Secrets.FluxNinjaPlugin.Create = true
			instance.Spec.Secrets.FluxNinjaPlugin.Value = Test

			os.Setenv("APERTURE_OPERATOR_CERT_DIR", CertDir)
			os.Setenv("APERTURE_OPERATOR_CERT_NAME", "tls6.crt")
			os.Setenv("APERTURE_OPERATOR_NAMESPACE", AppName)
			os.Setenv("APERTURE_OPERATOR_SERVICE_NAME", AppName)
			certPath := fmt.Sprintf("%s/%s", os.Getenv("APERTURE_OPERATOR_CERT_DIR"), os.Getenv("APERTURE_OPERATOR_CERT_NAME"))
			serverCertPEM := new(bytes.Buffer)
			err := pem.Encode(serverCertPEM, &pem.Block{
				Type:  "CERTIFICATE",
				Bytes: []byte(Test),
			})
			Expect(err).NotTo(HaveOccurred())

			err = WriteFile(certPath, serverCertPEM)
			Expect(err).NotTo(HaveOccurred())
			Expect(K8sClient.Create(Ctx, instance)).To(BeNil())

			res, err := reconciler.Reconcile(Ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      Test,
					Namespace: namespace,
				},
			})
			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).ToNot(HaveOccurred())

			createdClusterRoleBinding := &rbacv1.ClusterRoleBinding{}
			clusterRoleBindingKey := types.NamespacedName{Name: ControllerServiceName}

			Eventually(func() bool {
				return K8sClient.Get(Ctx, clusterRoleBindingKey, createdClusterRoleBinding) == nil
			}, time.Second*10, time.Millisecond*250).Should(BeTrue())

			Expect(K8sClient.Get(Ctx, types.NamespacedName{Name: Test, Namespace: namespace}, instance)).To(BeNil())
			Expect(instance.Status.Resources).To(Equal("created"))

			Expect(K8sClient.Delete(Ctx, instance)).To(BeNil())

			res, err = reconciler.Reconcile(Ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      Test,
					Namespace: namespace,
				},
			})
			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() bool {
				return K8sClient.Get(Ctx, clusterRoleBindingKey, createdClusterRoleBinding) != nil
			}, time.Second*10, time.Millisecond*250).Should(BeTrue())

			Expect(K8sClient.Delete(Ctx, ns)).To(BeNil())
		})

		AfterEach(func() {
			_ = K8sClient.Delete(Ctx, instance)
			_, err := reconciler.Reconcile(Ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      Test,
					Namespace: Test,
				},
			})
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("testing eventFilters", func() {
		It("should allow only Controller Events in create", func() {
			ControllerEvent := event.CreateEvent{
				Object: &controllerv1alpha1.Controller{},
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
				ObjectOld: &controllerv1alpha1.Controller{},
				ObjectNew: &controllerv1alpha1.Controller{
					Spec: controllerv1alpha1.ControllerSpec{
						CommonSpec: common.CommonSpec{
							Command: TestArray,
						},
					},
				},
			}

			ControllerEventValid2 := event.UpdateEvent{
				ObjectOld: &controllerv1alpha1.Controller{},
				ObjectNew: &controllerv1alpha1.Controller{
					ObjectMeta: metav1.ObjectMeta{
						DeletionTimestamp: &metav1.Time{Time: time.Now()},
					},
				},
			}

			ControllerEventInvalid1 := event.UpdateEvent{
				ObjectOld: DefaultControllerInstance,
				ObjectNew: DefaultControllerInstance,
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
				ObjectOld: &controllerv1alpha1.Controller{
					Spec: controllerv1alpha1.ControllerSpec{
						CommonSpec: common.CommonSpec{
							Secrets: common.Secrets{
								FluxNinjaPlugin: common.APIKeySecret{
									Value: Test,
								},
							},
						},
					},
				},
				ObjectNew: &controllerv1alpha1.Controller{
					Spec: controllerv1alpha1.ControllerSpec{
						CommonSpec: common.CommonSpec{
							Secrets: common.Secrets{
								FluxNinjaPlugin: common.APIKeySecret{
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
				Object:             &controllerv1alpha1.Controller{},
				DeleteStateUnknown: true,
			}

			ControllerEventInvalid := event.DeleteEvent{
				Object:             &controllerv1alpha1.Controller{},
				DeleteStateUnknown: false,
			}

			pred := eventFiltersForController()

			Expect(pred.Delete(ControllerEvent)).To(Equal(false))
			Expect(pred.Delete(ControllerEventInvalid)).To(Equal(true))
		})
	})
})
