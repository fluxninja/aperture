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
	"encoding/base64"
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

var _ = Describe("Agent Reconcile", Ordered, func() {
	Context("testing Reconcile", func() {
		var instance *v1alpha1.Agent
		var reconciler *AgentReconciler

		BeforeEach(func() {
			instance = defaultAgentInstance.DeepCopy()
			instance.Name = test
			instance.Namespace = test
			reconciler = &AgentReconciler{
				Client:           k8sClient,
				Scheme:           scheme.Scheme,
				Recorder:         k8sManager.GetEventRecorderFor(appName),
				ApertureInjector: &ApertureInjector{Client: k8sClient},
			}
		})

		It("should not create resources when Agent is deleted", func() {
			res, err := reconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      test,
					Namespace: test,
				},
			})
			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).ToNot(HaveOccurred())
		})

		It("should create required resources when Agent is created with default parameters", func() {
			namespace := test + "11"
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

			createdAgentConfigMap := &corev1.ConfigMap{}
			agentConfigKey := types.NamespacedName{Name: agentServiceName, Namespace: namespace}

			createdAgentService := &corev1.Service{}
			agentServiceKey := types.NamespacedName{Name: agentServiceName, Namespace: namespace}

			createdClusterRole := &rbacv1.ClusterRole{}
			clusterRoleKey := types.NamespacedName{Name: appName}

			createdClusterRoleBinding := &rbacv1.ClusterRoleBinding{}
			clusterRoleBindingKey := types.NamespacedName{Name: agentServiceName}

			createdAgentServiceAccount := &corev1.ServiceAccount{}
			agentServiceAccountKey := types.NamespacedName{Name: agentServiceName, Namespace: namespace}

			createdAgentDaemonset := &appsv1.DaemonSet{}
			agentDaemonsetKey := types.NamespacedName{Name: agentServiceName, Namespace: namespace}

			createdMWC := &admissionregistrationv1.MutatingWebhookConfiguration{}
			mwcKey := types.NamespacedName{Name: mutatingWebhookName}

			createdAgentSecret := &corev1.Secret{}
			agentSecretKey := types.NamespacedName{Name: secretName(test, "agent", &instance.Spec.FluxNinjaPlugin.APIKeySecret), Namespace: namespace}

			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() bool {
				err1 := k8sClient.Get(ctx, agentConfigKey, createdAgentConfigMap)
				err2 := k8sClient.Get(ctx, agentServiceKey, createdAgentService)
				err3 := k8sClient.Get(ctx, clusterRoleKey, createdClusterRole)
				err4 := k8sClient.Get(ctx, clusterRoleBindingKey, createdClusterRoleBinding)
				err5 := k8sClient.Get(ctx, agentServiceAccountKey, createdAgentServiceAccount)
				err6 := k8sClient.Get(ctx, agentDaemonsetKey, createdAgentDaemonset)
				err7 := k8sClient.Get(ctx, mwcKey, createdMWC)
				err8 := k8sClient.Get(ctx, agentSecretKey, createdAgentSecret)
				return err1 == nil && err2 == nil && err3 == nil && err4 == nil &&
					err5 == nil && err6 == nil && err7 != nil && err8 != nil
			}, time.Second*10, time.Millisecond*250).Should(BeTrue())

			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: test, Namespace: namespace}, instance)).To(BeNil())
			Expect(instance.Status.Resources).To(Equal("created"))

			Expect(k8sClient.Delete(ctx, ns)).To(BeNil())
		})

		It("should create required resources when Agent is created with all parameters and without sidecar", func() {
			namespace := test + "12"
			namespace1 := test + "13"
			namespace2 := test + "14"
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			Expect(k8sClient.Create(ctx, ns)).To(BeNil())

			instance.Namespace = namespace
			instance.Spec.Sidecar.Enabled = false
			instance.Spec.Sidecar.EnableNamespaceByDefault = []string{namespace1}
			instance.Spec.FluxNinjaPlugin.Enabled = true
			instance.Spec.FluxNinjaPlugin.APIKeySecret.Create = true
			instance.Spec.FluxNinjaPlugin.APIKeySecret.Value = test
			Expect(k8sClient.Create(ctx, instance)).To(BeNil())

			ns1 := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace1,
				},
			}
			Expect(k8sClient.Create(ctx, ns1)).To(BeNil())

			ns2 := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace2,
					Labels: map[string]string{
						sidecarLabelKey: enabled,
					},
				},
			}
			Expect(k8sClient.Create(ctx, ns2)).To(BeNil())

			res, err := reconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      test,
					Namespace: namespace,
				},
			})

			createdAgentConfigMap := &corev1.ConfigMap{}
			agentConfigKey := types.NamespacedName{Name: agentServiceName, Namespace: namespace}

			createdAgentService := &corev1.Service{}
			agentServiceKey := types.NamespacedName{Name: agentServiceName, Namespace: namespace}

			createdClusterRole := &rbacv1.ClusterRole{}
			clusterRoleKey := types.NamespacedName{Name: appName}

			createdClusterRoleBinding := &rbacv1.ClusterRoleBinding{}
			clusterRoleBindingKey := types.NamespacedName{Name: agentServiceName}

			createdAgentServiceAccount := &corev1.ServiceAccount{}
			agentServiceAccountKey := types.NamespacedName{Name: agentServiceName, Namespace: namespace}

			createdAgentDaemonset := &appsv1.DaemonSet{}
			agentDaemonsetKey := types.NamespacedName{Name: agentServiceName, Namespace: namespace}

			createdMWC := &admissionregistrationv1.MutatingWebhookConfiguration{}
			mwcKey := types.NamespacedName{Name: mutatingWebhookName}

			createdAgentSecret := &corev1.Secret{}
			agentSecretKey := types.NamespacedName{Name: secretName(test, "agent", &instance.Spec.FluxNinjaPlugin.APIKeySecret), Namespace: namespace}

			createdAgentConfigMapNs1 := &corev1.ConfigMap{}
			agentConfigKeyNs1 := types.NamespacedName{Name: agentServiceName, Namespace: namespace1}

			createdAgentConfigMapNs2 := &corev1.ConfigMap{}
			agentConfigKeyNs2 := types.NamespacedName{Name: agentServiceName, Namespace: namespace2}

			createdAgentSecretNs1 := &corev1.Secret{}
			agentSecretKeyNs1 := types.NamespacedName{Name: secretName(test, "agent", &instance.Spec.FluxNinjaPlugin.APIKeySecret), Namespace: namespace1}

			createdAgentSecretNs2 := &corev1.Secret{}
			agentSecretKeyNs2 := types.NamespacedName{Name: secretName(test, "agent", &instance.Spec.FluxNinjaPlugin.APIKeySecret), Namespace: namespace2}

			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() bool {
				err1 := k8sClient.Get(ctx, agentConfigKey, createdAgentConfigMap)
				err2 := k8sClient.Get(ctx, agentServiceKey, createdAgentService)
				err3 := k8sClient.Get(ctx, clusterRoleKey, createdClusterRole)
				err4 := k8sClient.Get(ctx, clusterRoleBindingKey, createdClusterRoleBinding)
				err5 := k8sClient.Get(ctx, agentServiceAccountKey, createdAgentServiceAccount)
				err6 := k8sClient.Get(ctx, agentDaemonsetKey, createdAgentDaemonset)
				err7 := k8sClient.Get(ctx, mwcKey, createdMWC)
				err8 := k8sClient.Get(ctx, agentSecretKey, createdAgentSecret)
				err9 := k8sClient.Get(ctx, agentConfigKeyNs1, createdAgentConfigMapNs1)
				err10 := k8sClient.Get(ctx, agentConfigKeyNs2, createdAgentConfigMapNs2)
				err11 := k8sClient.Get(ctx, agentSecretKeyNs1, createdAgentSecretNs1)
				err12 := k8sClient.Get(ctx, agentSecretKeyNs2, createdAgentSecretNs2)
				return err1 == nil && err2 == nil && err3 == nil && err4 == nil &&
					err5 == nil && err6 == nil && err7 != nil && err8 == nil &&
					err9 != nil && err10 != nil && err11 != nil && err12 != nil
			}, time.Second*10, time.Millisecond*250).Should(BeTrue())

			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: test, Namespace: namespace}, instance)).To(BeNil())
			Expect(instance.Status.Resources).To(Equal("created"))
			Expect(instance.Finalizers).To(Equal([]string{finalizerName}))
			Expect(instance.Spec.FluxNinjaPlugin.APIKeySecret.Create).To(BeFalse())
			Expect(instance.Spec.FluxNinjaPlugin.APIKeySecret.Value).To(Equal(""))

			Expect(k8sClient.Delete(ctx, ns)).To(BeNil())
			Expect(k8sClient.Delete(ctx, ns1)).To(BeNil())
			Expect(k8sClient.Delete(ctx, ns2)).To(BeNil())
		})

		It("should create required resources when Agent is created with all parameters and sidecar", func() {
			namespace := test + "15"
			namespace1 := test + "16"
			namespace2 := test + "17"
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			Expect(k8sClient.Create(ctx, ns)).To(BeNil())

			instance.Namespace = namespace
			instance.Spec.Sidecar.Enabled = true
			instance.Spec.Sidecar.EnableNamespaceByDefault = []string{namespace1}
			instance.Spec.FluxNinjaPlugin.Enabled = true
			instance.Spec.CommonSpec.ServiceAccountSpec.Create = false
			encodedString := fmt.Sprintf("enc::%s::enc", base64.StdEncoding.EncodeToString([]byte(test)))
			instance.Spec.FluxNinjaPlugin.APIKeySecret.Create = true
			instance.Spec.FluxNinjaPlugin.APIKeySecret.Value = test
			Expect(k8sClient.Create(ctx, instance)).To(BeNil())

			ns1 := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace1,
				},
			}
			Expect(k8sClient.Create(ctx, ns1)).To(BeNil())

			ns2 := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace2,
					Labels: map[string]string{
						sidecarLabelKey: enabled,
					},
				},
			}
			Expect(k8sClient.Create(ctx, ns2)).To(BeNil())

			os.Setenv("APERTURE_OPERATOR_CERT_DIR", certDir)
			os.Setenv("APERTURE_OPERATOR_CERT_NAME", "tls6.crt")
			os.Setenv("APERTURE_OPERATOR_NAMESPACE", appName)
			os.Setenv("APERTURE_OPERATOR_SERVICE_NAME", appName)
			certPath := fmt.Sprintf("%s/%s", os.Getenv("APERTURE_OPERATOR_CERT_DIR"), webhookClientCertName)
			serverCertPEM := new(bytes.Buffer)
			_ = pem.Encode(serverCertPEM, &pem.Block{
				Type:  "CERTIFICATE",
				Bytes: []byte(test),
			})
			err := writeFile(certPath, serverCertPEM)
			Expect(err).NotTo(HaveOccurred())

			res, err := reconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      test,
					Namespace: namespace,
				},
			})

			createdAgentConfigMap := &corev1.ConfigMap{}
			agentConfigKey := types.NamespacedName{Name: agentServiceName, Namespace: namespace}

			createdAgentService := &corev1.Service{}
			agentServiceKey := types.NamespacedName{Name: agentServiceName, Namespace: namespace}

			createdClusterRole := &rbacv1.ClusterRole{}
			clusterRoleKey := types.NamespacedName{Name: appName}

			createdClusterRoleBinding := &rbacv1.ClusterRoleBinding{}
			clusterRoleBindingKey := types.NamespacedName{Name: agentServiceName}

			createdAgentServiceAccount := &corev1.ServiceAccount{}
			agentServiceAccountKey := types.NamespacedName{Name: agentServiceName, Namespace: namespace}

			createdAgentDaemonset := &appsv1.DaemonSet{}
			agentDaemonsetKey := types.NamespacedName{Name: agentServiceName, Namespace: namespace}

			createdMWC := &admissionregistrationv1.MutatingWebhookConfiguration{}
			mwcKey := types.NamespacedName{Name: mutatingWebhookName}

			createdAgentSecret := &corev1.Secret{}
			agentSecretKey := types.NamespacedName{Name: secretName(test, "agent", &instance.Spec.FluxNinjaPlugin.APIKeySecret), Namespace: namespace}

			createdAgentConfigMapNs1 := &corev1.ConfigMap{}
			agentConfigKeyNs1 := types.NamespacedName{Name: agentServiceName, Namespace: namespace1}

			createdAgentConfigMapNs2 := &corev1.ConfigMap{}
			agentConfigKeyNs2 := types.NamespacedName{Name: agentServiceName, Namespace: namespace2}

			createdAgentSecretNs1 := &corev1.Secret{}
			agentSecretKeyNs1 := types.NamespacedName{Name: secretName(test, "agent", &instance.Spec.FluxNinjaPlugin.APIKeySecret), Namespace: namespace1}

			createdAgentSecretNs2 := &corev1.Secret{}
			agentSecretKeyNs2 := types.NamespacedName{Name: secretName(test, "agent", &instance.Spec.FluxNinjaPlugin.APIKeySecret), Namespace: namespace2}

			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() bool {
				err1 := k8sClient.Get(ctx, agentConfigKey, createdAgentConfigMap)
				err2 := k8sClient.Get(ctx, agentServiceKey, createdAgentService)
				err3 := k8sClient.Get(ctx, clusterRoleKey, createdClusterRole)
				err4 := k8sClient.Get(ctx, clusterRoleBindingKey, createdClusterRoleBinding)
				err5 := k8sClient.Get(ctx, agentServiceAccountKey, createdAgentServiceAccount)
				err6 := k8sClient.Get(ctx, agentDaemonsetKey, createdAgentDaemonset)
				err7 := k8sClient.Get(ctx, mwcKey, createdMWC)
				err8 := k8sClient.Get(ctx, agentSecretKey, createdAgentSecret)
				err9 := k8sClient.Get(ctx, agentConfigKeyNs1, createdAgentConfigMapNs1)
				err10 := k8sClient.Get(ctx, agentConfigKeyNs2, createdAgentConfigMapNs2)
				err11 := k8sClient.Get(ctx, agentSecretKeyNs1, createdAgentSecretNs1)
				err12 := k8sClient.Get(ctx, agentSecretKeyNs2, createdAgentSecretNs2)
				return err1 != nil && err2 != nil && err3 == nil && err4 == nil &&
					err5 != nil && err6 != nil && err7 == nil && err8 != nil &&
					err9 == nil && err10 == nil && err11 == nil && err12 == nil
			}, time.Second*10, time.Millisecond*250).Should(BeTrue())

			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: test, Namespace: namespace}, instance)).To(BeNil())
			Expect(instance.Status.Resources).To(Equal("created"))
			Expect(instance.Finalizers).To(Equal([]string{finalizerName}))
			Expect(instance.Spec.FluxNinjaPlugin.APIKeySecret.Value).To(Equal(encodedString))

			Expect(k8sClient.Delete(ctx, ns)).To(BeNil())
			Expect(k8sClient.Delete(ctx, ns1)).To(BeNil())
			Expect(k8sClient.Delete(ctx, ns2)).To(BeNil())
		})

		It("should not create required resources when an Agent instance is already created", func() {
			namespace := test + "18"
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

			createdAgentConfigMap := &corev1.ConfigMap{}
			agentConfigKey := types.NamespacedName{Name: agentServiceName, Namespace: namespace}

			Eventually(func() bool {
				return k8sClient.Get(ctx, agentConfigKey, createdAgentConfigMap) == nil
			}, time.Second*10, time.Millisecond*250).Should(BeTrue())

			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: test, Namespace: namespace}, instance)).To(BeNil())
			Expect(instance.Status.Resources).To(Equal("created"))

			instanceNew := defaultAgentInstance.DeepCopy()
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

		It("should delete required resources when an Agent instance is already deleted", func() {
			namespace := test + "19"
			namespace1 := test + "20"
			namespace2 := test + "21"
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			Expect(k8sClient.Create(ctx, ns)).To(BeNil())

			instance.Namespace = namespace
			instance.Spec.Sidecar.Enabled = true
			instance.Spec.Sidecar.EnableNamespaceByDefault = []string{namespace1}
			instance.Spec.FluxNinjaPlugin.Enabled = true
			instance.Spec.CommonSpec.ServiceAccountSpec.Create = false
			instance.Spec.FluxNinjaPlugin.APIKeySecret.Create = true
			instance.Spec.FluxNinjaPlugin.APIKeySecret.Value = test

			ns1 := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace1,
				},
			}
			Expect(k8sClient.Create(ctx, ns1)).To(BeNil())

			ns2 := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace2,
					Labels: map[string]string{
						sidecarLabelKey: enabled,
					},
				},
			}
			Expect(k8sClient.Create(ctx, ns2)).To(BeNil())

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
			clusterRoleBindingKey := types.NamespacedName{Name: agentServiceName}

			createdAgentConfigMapNs1 := &corev1.ConfigMap{}
			agentConfigKeyNs1 := types.NamespacedName{Name: agentServiceName, Namespace: namespace1}

			createdAgentSecretNs2 := &corev1.Secret{}
			agentSecretKeyNs2 := types.NamespacedName{Name: secretName(test, "agent", &instance.Spec.FluxNinjaPlugin.APIKeySecret), Namespace: namespace2}

			Eventually(func() bool {
				return k8sClient.Get(ctx, clusterRoleBindingKey, createdClusterRoleBinding) == nil &&
					k8sClient.Get(ctx, agentConfigKeyNs1, createdAgentConfigMapNs1) == nil &&
					k8sClient.Get(ctx, agentSecretKeyNs2, createdAgentSecretNs2) == nil
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
				return k8sClient.Get(ctx, clusterRoleBindingKey, createdClusterRoleBinding) != nil &&
					k8sClient.Get(ctx, agentConfigKeyNs1, createdAgentConfigMapNs1) != nil &&
					k8sClient.Get(ctx, agentSecretKeyNs2, createdAgentSecretNs2) != nil
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
		It("should allow only Agent Events in create", func() {
			AgentEvent := event.CreateEvent{
				Object: &v1alpha1.Agent{},
			}
			cmEvent := event.CreateEvent{
				Object: &corev1.ConfigMap{},
			}

			pred := eventFiltersForAgent()

			Expect(pred.Create(AgentEvent)).To(Equal(true))
			Expect(pred.Create(cmEvent)).To(Equal(false))
		})

		It("should allow only Agent Events in update when Spec is changed or instance is deleted", func() {
			AgentEventValid1 := event.UpdateEvent{
				ObjectOld: &v1alpha1.Agent{},
				ObjectNew: &v1alpha1.Agent{
					Spec: v1alpha1.AgentSpec{
						CommonSpec: v1alpha1.CommonSpec{
							FluxNinjaPlugin: v1alpha1.FluxNinjaPluginSpec{
								Enabled: true,
							},
						},
					},
				},
			}

			AgentEventValid2 := event.UpdateEvent{
				ObjectOld: &v1alpha1.Agent{},
				ObjectNew: &v1alpha1.Agent{
					ObjectMeta: metav1.ObjectMeta{
						DeletionTimestamp: &metav1.Time{Time: time.Now()},
					},
				},
			}

			AgentEventInvalid1 := event.UpdateEvent{
				ObjectOld: defaultAgentInstance,
				ObjectNew: defaultAgentInstance,
			}

			AgentEventInvalid2 := event.UpdateEvent{
				ObjectOld: &corev1.ConfigMap{},
				ObjectNew: &corev1.ConfigMap{},
			}

			pred := eventFiltersForAgent()

			Expect(pred.Update(AgentEventValid1)).To(Equal(true))
			Expect(pred.Update(AgentEventValid2)).To(Equal(true))
			Expect(pred.Update(AgentEventInvalid1)).To(Equal(false))
			Expect(pred.Update(AgentEventInvalid2)).To(Equal(true))
		})

		It("should allow only Agent Events in update when secret is changed", func() {
			AgentEventValid1 := event.UpdateEvent{
				ObjectOld: &v1alpha1.Agent{
					Spec: v1alpha1.AgentSpec{
						CommonSpec: v1alpha1.CommonSpec{
							FluxNinjaPlugin: v1alpha1.FluxNinjaPluginSpec{
								APIKeySecret: v1alpha1.APIKeySecret{
									Value: test,
								},
							},
						},
					},
				},
				ObjectNew: &v1alpha1.Agent{
					Spec: v1alpha1.AgentSpec{
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

			AgentEventValid2 := event.UpdateEvent{
				ObjectOld: &v1alpha1.Agent{
					Spec: v1alpha1.AgentSpec{
						CommonSpec: v1alpha1.CommonSpec{
							FluxNinjaPlugin: v1alpha1.FluxNinjaPluginSpec{
								APIKeySecret: v1alpha1.APIKeySecret{
									Value: test,
								},
							},
						},
					},
				},
				ObjectNew: &v1alpha1.Agent{
					Spec: v1alpha1.AgentSpec{
						CommonSpec: v1alpha1.CommonSpec{
							FluxNinjaPlugin: v1alpha1.FluxNinjaPluginSpec{
								APIKeySecret: v1alpha1.APIKeySecret{
									Value: fmt.Sprintf("enc::%s::enc", base64.StdEncoding.EncodeToString([]byte(test))),
								},
							},
						},
					},
				},
			}

			pred := eventFiltersForAgent()

			Expect(pred.Update(AgentEventValid1)).To(Equal(false))
			Expect(pred.Update(AgentEventValid2)).To(Equal(false))
		})

		It("should allow Events in delete when state is known", func() {
			AgentEvent := event.DeleteEvent{
				Object:             &v1alpha1.Agent{},
				DeleteStateUnknown: true,
			}

			AgentEventInvalid := event.DeleteEvent{
				Object:             &v1alpha1.Agent{},
				DeleteStateUnknown: false,
			}

			pred := eventFiltersForAgent()

			Expect(pred.Delete(AgentEvent)).To(Equal(false))
			Expect(pred.Delete(AgentEventInvalid)).To(Equal(true))
		})
	})
})
