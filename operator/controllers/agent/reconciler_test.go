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

package agent

import (
	"bytes"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
	"reflect"
	"time"

	. "github.com/fluxninja/aperture/v2/operator/controllers"
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

	agentv1alpha1 "github.com/fluxninja/aperture/v2/operator/api/agent/v1alpha1"
	"github.com/fluxninja/aperture/v2/operator/api/common"
	"github.com/fluxninja/aperture/v2/operator/controllers/mutatingwebhook"
)

var _ = Describe("Agent Reconcile", Ordered, func() {
	Context("testing Reconcile", func() {
		var instance *agentv1alpha1.Agent
		var reconciler *AgentReconciler

		BeforeEach(func() {
			instance = DefaultAgentInstance.DeepCopy()
			instance.Name = Test
			instance.Namespace = Test
			reconciler = &AgentReconciler{
				Client:           K8sClient,
				DynamicClient:    K8sDynamicClient,
				Scheme:           scheme.Scheme,
				Recorder:         K8sManager.GetEventRecorderFor(AppName),
				ApertureInjector: &mutatingwebhook.ApertureInjector{Client: K8sClient},
			}
		})

		It("should not create resources when Agent is deleted", func() {
			res, err := reconciler.Reconcile(Ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      Test,
					Namespace: Test,
				},
			})
			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).ToNot(HaveOccurred())
		})

		It("should create required resources when Agent is created with default parameters", func() {
			namespace := Test + "11"
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			Expect(K8sClient.Create(Ctx, ns)).To(Succeed())

			instance.Namespace = namespace
			Expect(K8sClient.Create(Ctx, instance)).To(Succeed())

			res, err := reconciler.Reconcile(Ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      Test,
					Namespace: namespace,
				},
			})

			createdAgentConfigMap := &corev1.ConfigMap{}
			agentConfigKey := types.NamespacedName{Name: AgentServiceName, Namespace: namespace}

			createdAgentService := &corev1.Service{}
			agentServiceKey := types.NamespacedName{Name: AgentServiceName, Namespace: namespace}

			createdClusterRole := &rbacv1.ClusterRole{}
			clusterRoleKey := types.NamespacedName{Name: AgentServiceName}

			createdClusterRoleBinding := &rbacv1.ClusterRoleBinding{}
			clusterRoleBindingKey := types.NamespacedName{Name: AgentServiceName}

			createdAgentServiceAccount := &corev1.ServiceAccount{}
			agentServiceAccountKey := types.NamespacedName{Name: AgentServiceName, Namespace: namespace}

			createdAgentDaemonset := &appsv1.DaemonSet{}
			agentDaemonsetKey := types.NamespacedName{Name: AgentServiceName, Namespace: namespace}

			createdMWC := &admissionregistrationv1.MutatingWebhookConfiguration{}
			mwcKey := types.NamespacedName{Name: PodMutatingWebhookName}

			createdAgentSecret := &corev1.Secret{}
			agentSecretKey := types.NamespacedName{Name: SecretName(Test, "agent", &instance.Spec.Secrets.FluxNinjaExtension), Namespace: namespace}

			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() bool {
				err1 := K8sClient.Get(Ctx, agentConfigKey, createdAgentConfigMap)
				err2 := K8sClient.Get(Ctx, agentServiceKey, createdAgentService)
				err3 := K8sClient.Get(Ctx, clusterRoleKey, createdClusterRole)
				err4 := K8sClient.Get(Ctx, clusterRoleBindingKey, createdClusterRoleBinding)
				err5 := K8sClient.Get(Ctx, agentServiceAccountKey, createdAgentServiceAccount)
				err6 := K8sClient.Get(Ctx, agentDaemonsetKey, createdAgentDaemonset)
				err7 := K8sClient.Get(Ctx, mwcKey, createdMWC)
				err8 := K8sClient.Get(Ctx, agentSecretKey, createdAgentSecret)
				return err1 == nil && err2 == nil && err3 == nil && err4 == nil &&
					err5 == nil && err6 == nil && err7 != nil && err8 != nil
			}, time.Second*10, time.Millisecond*250).Should(BeTrue())

			Expect(K8sClient.Get(Ctx, types.NamespacedName{Name: Test, Namespace: namespace}, instance)).To(Succeed())
			Expect(instance.Status.Resources).To(Equal("created"))

			Expect(K8sClient.Delete(Ctx, ns)).To(Succeed())
		})

		It("should create required resources when Agent is created with all parameters and without sidecar", func() {
			namespace := Test + "12"
			namespace1 := Test + "13"
			namespace2 := Test + "14"
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			Expect(K8sClient.Create(Ctx, ns)).To(Succeed())

			instance.Namespace = namespace
			instance.Spec.Sidecar.Enabled = false
			instance.Spec.Sidecar.EnableNamespaceByDefault = []string{namespace1}
			instance.Spec.Secrets.FluxNinjaExtension.Create = true
			instance.Spec.Secrets.FluxNinjaExtension.Value = Test
			Expect(K8sClient.Create(Ctx, instance)).To(Succeed())

			ns1 := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace1,
				},
			}
			Expect(K8sClient.Create(Ctx, ns1)).To(Succeed())

			ns2 := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace2,
					Labels: map[string]string{
						SidecarLabelKey: Enabled,
					},
				},
			}
			Expect(K8sClient.Create(Ctx, ns2)).To(Succeed())

			res, err := reconciler.Reconcile(Ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      Test,
					Namespace: namespace,
				},
			})

			createdAgentConfigMap := &corev1.ConfigMap{}
			agentConfigKey := types.NamespacedName{Name: AgentServiceName, Namespace: namespace}

			createdAgentService := &corev1.Service{}
			agentServiceKey := types.NamespacedName{Name: AgentServiceName, Namespace: namespace}

			createdClusterRole := &rbacv1.ClusterRole{}
			clusterRoleKey := types.NamespacedName{Name: AgentServiceName}

			createdClusterRoleBinding := &rbacv1.ClusterRoleBinding{}
			clusterRoleBindingKey := types.NamespacedName{Name: AgentServiceName}

			createdAgentServiceAccount := &corev1.ServiceAccount{}
			agentServiceAccountKey := types.NamespacedName{Name: AgentServiceName, Namespace: namespace}

			createdAgentDaemonset := &appsv1.DaemonSet{}
			agentDaemonsetKey := types.NamespacedName{Name: AgentServiceName, Namespace: namespace}

			createdMWC := &admissionregistrationv1.MutatingWebhookConfiguration{}
			mwcKey := types.NamespacedName{Name: PodMutatingWebhookName}

			createdAgentSecret := &corev1.Secret{}
			agentSecretKey := types.NamespacedName{Name: SecretName(Test, "agent", &instance.Spec.Secrets.FluxNinjaExtension), Namespace: namespace}

			createdAgentConfigMapNs1 := &corev1.ConfigMap{}
			agentConfigKeyNs1 := types.NamespacedName{Name: AgentServiceName, Namespace: namespace1}

			createdAgentConfigMapNs2 := &corev1.ConfigMap{}
			agentConfigKeyNs2 := types.NamespacedName{Name: AgentServiceName, Namespace: namespace2}

			createdAgentSecretNs1 := &corev1.Secret{}
			agentSecretKeyNs1 := types.NamespacedName{Name: SecretName(Test, "agent", &instance.Spec.Secrets.FluxNinjaExtension), Namespace: namespace1}

			createdAgentSecretNs2 := &corev1.Secret{}
			agentSecretKeyNs2 := types.NamespacedName{Name: SecretName(Test, "agent", &instance.Spec.Secrets.FluxNinjaExtension), Namespace: namespace2}

			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() bool {
				err1 := K8sClient.Get(Ctx, agentConfigKey, createdAgentConfigMap)
				err2 := K8sClient.Get(Ctx, agentServiceKey, createdAgentService)
				err3 := K8sClient.Get(Ctx, clusterRoleKey, createdClusterRole)
				err4 := K8sClient.Get(Ctx, clusterRoleBindingKey, createdClusterRoleBinding)
				err5 := K8sClient.Get(Ctx, agentServiceAccountKey, createdAgentServiceAccount)
				err6 := K8sClient.Get(Ctx, agentDaemonsetKey, createdAgentDaemonset)
				err7 := K8sClient.Get(Ctx, mwcKey, createdMWC)
				err8 := K8sClient.Get(Ctx, agentSecretKey, createdAgentSecret)
				err9 := K8sClient.Get(Ctx, agentConfigKeyNs1, createdAgentConfigMapNs1)
				err10 := K8sClient.Get(Ctx, agentConfigKeyNs2, createdAgentConfigMapNs2)
				err11 := K8sClient.Get(Ctx, agentSecretKeyNs1, createdAgentSecretNs1)
				err12 := K8sClient.Get(Ctx, agentSecretKeyNs2, createdAgentSecretNs2)
				return err1 == nil && err2 == nil && err3 == nil && err4 == nil &&
					err5 == nil && err6 == nil && err7 != nil && err8 == nil &&
					err9 != nil && err10 != nil && err11 != nil && err12 != nil
			}, time.Second*10, time.Millisecond*250).Should(BeTrue())

			Expect(K8sClient.Get(Ctx, types.NamespacedName{Name: Test, Namespace: namespace}, instance)).To(Succeed())
			Expect(instance.Status.Resources).To(Equal("created"))
			Expect(instance.Finalizers).To(Equal([]string{FinalizerName}))
			Expect(instance.Spec.Secrets.FluxNinjaExtension.Create).To(BeFalse())
			Expect(instance.Spec.Secrets.FluxNinjaExtension.Value).To(Equal(""))
			Expect(instance.Spec.ConfigSpec.FluxNinja.InstallationMode).To(Equal("KUBERNETES_DAEMONSET"))

			Expect(K8sClient.Delete(Ctx, ns)).To(Succeed())
			Expect(K8sClient.Delete(Ctx, ns1)).To(Succeed())
			Expect(K8sClient.Delete(Ctx, ns2)).To(Succeed())
		})

		It("should create required resources when Agent is created with all parameters and sidecar", func() {
			namespace := Test + "15"
			namespace1 := Test + "16"
			namespace2 := Test + "17"
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			Expect(K8sClient.Create(Ctx, ns)).To(Succeed())

			instance.Namespace = namespace
			instance.Spec.Sidecar.Enabled = true
			instance.Spec.Sidecar.EnableNamespaceByDefault = []string{namespace1}
			instance.Spec.CommonSpec.ServiceAccountSpec.Create = false
			encodedString := fmt.Sprintf("enc::%s::enc", base64.StdEncoding.EncodeToString([]byte(Test)))
			instance.Spec.Secrets.FluxNinjaExtension.Create = true
			instance.Spec.Secrets.FluxNinjaExtension.Value = Test
			Expect(K8sClient.Create(Ctx, instance)).To(Succeed())

			ns1 := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace1,
				},
			}
			Expect(K8sClient.Create(Ctx, ns1)).To(Succeed())

			ns2 := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace2,
					Labels: map[string]string{
						SidecarLabelKey: Enabled,
					},
				},
			}
			Expect(K8sClient.Create(Ctx, ns2)).To(Succeed())

			os.Setenv("APERTURE_OPERATOR_CERT_DIR", CertDir)
			os.Setenv("APERTURE_OPERATOR_CERT_NAME", "tls6.crt")
			os.Setenv("APERTURE_OPERATOR_NAMESPACE", AppName)
			os.Setenv("APERTURE_OPERATOR_SERVICE_NAME", AppName)
			certPath := fmt.Sprintf("%s/%s", os.Getenv("APERTURE_OPERATOR_CERT_DIR"), WebhookClientCertName)
			serverCertPEM := new(bytes.Buffer)
			_ = pem.Encode(serverCertPEM, &pem.Block{
				Type:  "CERTIFICATE",
				Bytes: []byte(Test),
			})
			err := WriteFile(certPath, serverCertPEM)
			Expect(err).NotTo(HaveOccurred())

			res, err := reconciler.Reconcile(Ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      Test,
					Namespace: namespace,
				},
			})

			createdAgentConfigMap := &corev1.ConfigMap{}
			agentConfigKey := types.NamespacedName{Name: AgentServiceName, Namespace: namespace}

			createdAgentService := &corev1.Service{}
			agentServiceKey := types.NamespacedName{Name: AgentServiceName, Namespace: namespace}

			createdClusterRole := &rbacv1.ClusterRole{}
			clusterRoleKey := types.NamespacedName{Name: AgentServiceName}

			createdClusterRoleBinding := &rbacv1.ClusterRoleBinding{}
			clusterRoleBindingKey := types.NamespacedName{Name: AgentServiceName}

			createdAgentServiceAccount := &corev1.ServiceAccount{}
			agentServiceAccountKey := types.NamespacedName{Name: AgentServiceName, Namespace: namespace}

			createdAgentDaemonset := &appsv1.DaemonSet{}
			agentDaemonsetKey := types.NamespacedName{Name: AgentServiceName, Namespace: namespace}

			createdMWC := &admissionregistrationv1.MutatingWebhookConfiguration{}
			mwcKey := types.NamespacedName{Name: PodMutatingWebhookName}

			createdAgentSecret := &corev1.Secret{}
			agentSecretKey := types.NamespacedName{Name: SecretName(Test, "agent", &instance.Spec.Secrets.FluxNinjaExtension), Namespace: namespace}

			createdAgentConfigMapNs1 := &corev1.ConfigMap{}
			agentConfigKeyNs1 := types.NamespacedName{Name: AgentServiceName, Namespace: namespace1}

			createdAgentConfigMapNs2 := &corev1.ConfigMap{}
			agentConfigKeyNs2 := types.NamespacedName{Name: AgentServiceName, Namespace: namespace2}

			createdAgentSecretNs1 := &corev1.Secret{}
			agentSecretKeyNs1 := types.NamespacedName{Name: SecretName(Test, "agent", &instance.Spec.Secrets.FluxNinjaExtension), Namespace: namespace1}

			createdAgentSecretNs2 := &corev1.Secret{}
			agentSecretKeyNs2 := types.NamespacedName{Name: SecretName(Test, "agent", &instance.Spec.Secrets.FluxNinjaExtension), Namespace: namespace2}

			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() bool {
				err1 := K8sClient.Get(Ctx, agentConfigKey, createdAgentConfigMap)
				err2 := K8sClient.Get(Ctx, agentServiceKey, createdAgentService)
				err3 := K8sClient.Get(Ctx, clusterRoleKey, createdClusterRole)
				err4 := K8sClient.Get(Ctx, clusterRoleBindingKey, createdClusterRoleBinding)
				err5 := K8sClient.Get(Ctx, agentServiceAccountKey, createdAgentServiceAccount)
				err6 := K8sClient.Get(Ctx, agentDaemonsetKey, createdAgentDaemonset)
				err7 := K8sClient.Get(Ctx, mwcKey, createdMWC)
				err8 := K8sClient.Get(Ctx, agentSecretKey, createdAgentSecret)
				err9 := K8sClient.Get(Ctx, agentConfigKeyNs1, createdAgentConfigMapNs1)
				err10 := K8sClient.Get(Ctx, agentConfigKeyNs2, createdAgentConfigMapNs2)
				err11 := K8sClient.Get(Ctx, agentSecretKeyNs1, createdAgentSecretNs1)
				err12 := K8sClient.Get(Ctx, agentSecretKeyNs2, createdAgentSecretNs2)
				return err1 != nil && err2 != nil && err3 == nil && err4 == nil &&
					err5 != nil && err6 != nil && err7 == nil && err8 != nil &&
					err9 == nil && err10 == nil && err11 == nil && err12 == nil
			}, time.Second*10, time.Millisecond*250).Should(BeTrue())

			Expect(K8sClient.Get(Ctx, types.NamespacedName{Name: Test, Namespace: namespace}, instance)).To(Succeed())
			Expect(instance.Status.Resources).To(Equal("created"))
			Expect(instance.Finalizers).To(Equal([]string{FinalizerName}))
			Expect(instance.Spec.Secrets.FluxNinjaExtension.Value).To(Equal(encodedString))
			Expect(instance.Spec.ConfigSpec.FluxNinja.InstallationMode).To(Equal("KUBERNETES_SIDECAR"))

			Expect(K8sClient.Delete(Ctx, ns)).To(Succeed())
			Expect(K8sClient.Delete(Ctx, ns1)).To(Succeed())
			Expect(K8sClient.Delete(Ctx, ns2)).To(Succeed())
		})

		It("should not create required resources when an Agent instance is already created", func() {
			namespace := Test + "18"
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			Expect(K8sClient.Create(Ctx, ns)).To(Succeed())

			instance.Namespace = namespace
			Expect(K8sClient.Create(Ctx, instance)).To(Succeed())

			res, err := reconciler.Reconcile(Ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      Test,
					Namespace: namespace,
				},
			})
			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).ToNot(HaveOccurred())

			createdAgentConfigMap := &corev1.ConfigMap{}
			agentConfigKey := types.NamespacedName{Name: AgentServiceName, Namespace: namespace}

			Eventually(func() bool {
				return K8sClient.Get(Ctx, agentConfigKey, createdAgentConfigMap) == nil
			}, time.Second*10, time.Millisecond*250).Should(BeTrue())

			Expect(K8sClient.Get(Ctx, types.NamespacedName{Name: Test, Namespace: namespace}, instance)).To(Succeed())
			Expect(instance.Status.Resources).To(Equal("created"))

			instanceNew := DefaultAgentInstance.DeepCopy()
			instanceNew.Name = TestTwo
			instanceNew.Namespace = namespace

			Expect(K8sClient.Create(Ctx, instanceNew)).To(Succeed())

			res, err = reconciler.Reconcile(Ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      TestTwo,
					Namespace: namespace,
				},
			})
			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).ToNot(HaveOccurred())

			Expect(K8sClient.Get(Ctx, types.NamespacedName{Name: TestTwo, Namespace: namespace}, instanceNew)).To(Succeed())
			Expect(instanceNew.Status.Resources).To(Equal("skipped"))

			Expect(K8sClient.Delete(Ctx, ns)).To(Succeed())
		})

		It("should delete required resources when an Agent instance is already deleted", func() {
			namespace := Test + "19"
			namespace1 := Test + "20"
			namespace2 := Test + "21"
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			Expect(K8sClient.Create(Ctx, ns)).To(Succeed())

			instance.Namespace = namespace
			instance.Spec.Sidecar.Enabled = true
			instance.Spec.Sidecar.EnableNamespaceByDefault = []string{namespace1}
			instance.Spec.CommonSpec.ServiceAccountSpec.Create = false
			instance.Spec.Secrets.FluxNinjaExtension.Create = true
			instance.Spec.Secrets.FluxNinjaExtension.Value = Test

			ns1 := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace1,
				},
			}
			Expect(K8sClient.Create(Ctx, ns1)).To(Succeed())

			ns2 := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace2,
					Labels: map[string]string{
						SidecarLabelKey: Enabled,
					},
				},
			}
			Expect(K8sClient.Create(Ctx, ns2)).To(Succeed())

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
			Expect(K8sClient.Create(Ctx, instance)).To(Succeed())

			res, err := reconciler.Reconcile(Ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      Test,
					Namespace: namespace,
				},
			})
			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).ToNot(HaveOccurred())

			createdClusterRoleBinding := &rbacv1.ClusterRoleBinding{}
			clusterRoleBindingKey := types.NamespacedName{Name: AgentServiceName}

			createdAgentConfigMapNs1 := &corev1.ConfigMap{}
			agentConfigKeyNs1 := types.NamespacedName{Name: AgentServiceName, Namespace: namespace1}

			createdAgentSecretNs2 := &corev1.Secret{}
			agentSecretKeyNs2 := types.NamespacedName{Name: SecretName(Test, "agent", &instance.Spec.Secrets.FluxNinjaExtension), Namespace: namespace2}

			Eventually(func() bool {
				return K8sClient.Get(Ctx, clusterRoleBindingKey, createdClusterRoleBinding) == nil &&
					K8sClient.Get(Ctx, agentConfigKeyNs1, createdAgentConfigMapNs1) == nil &&
					K8sClient.Get(Ctx, agentSecretKeyNs2, createdAgentSecretNs2) == nil
			}, time.Second*10, time.Millisecond*250).Should(BeTrue())

			Expect(K8sClient.Get(Ctx, types.NamespacedName{Name: Test, Namespace: namespace}, instance)).To(Succeed())
			Expect(instance.Status.Resources).To(Equal("created"))

			Expect(K8sClient.Delete(Ctx, instance)).To(Succeed())

			res, err = reconciler.Reconcile(Ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      Test,
					Namespace: namespace,
				},
			})
			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() bool {
				return K8sClient.Get(Ctx, clusterRoleBindingKey, createdClusterRoleBinding) != nil &&
					K8sClient.Get(Ctx, agentConfigKeyNs1, createdAgentConfigMapNs1) != nil &&
					K8sClient.Get(Ctx, agentSecretKeyNs2, createdAgentSecretNs2) != nil
			}, time.Second*10, time.Millisecond*250).Should(BeTrue())

			Expect(K8sClient.Delete(Ctx, ns)).To(Succeed())
		})

		It("should create required resources when Agent is updated to use sidecar mode", func() {
			namespace := Test + "22"
			namespace1 := Test + "23"
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			Expect(K8sClient.Create(Ctx, ns)).To(Succeed())

			instance.Namespace = namespace
			instance.Spec.Sidecar.Enabled = false
			instance.Spec.Secrets.FluxNinjaExtension.Create = true
			instance.Spec.Secrets.FluxNinjaExtension.Value = Test
			Expect(K8sClient.Create(Ctx, instance)).To(Succeed())

			ns1 := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace1,
				},
			}
			Expect(K8sClient.Create(Ctx, ns1)).To(Succeed())

			res, err := reconciler.Reconcile(Ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      Test,
					Namespace: namespace,
				},
			})

			createdAgentService := &corev1.Service{}
			agentServiceKey := types.NamespacedName{Name: AgentServiceName, Namespace: namespace}

			createdAgentDaemonset := &appsv1.DaemonSet{}
			agentDaemonsetKey := types.NamespacedName{Name: AgentServiceName, Namespace: namespace}

			createdMWC := &admissionregistrationv1.MutatingWebhookConfiguration{}
			mwcKey := types.NamespacedName{Name: PodMutatingWebhookName}

			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() bool {
				err1 := K8sClient.Get(Ctx, agentServiceKey, createdAgentService)
				err2 := K8sClient.Get(Ctx, agentDaemonsetKey, createdAgentDaemonset)
				err3 := K8sClient.Get(Ctx, mwcKey, createdMWC)
				return err1 == nil && err2 == nil && err3 != nil
			}, time.Second*10, time.Millisecond*250).Should(BeTrue())

			Expect(K8sClient.Get(Ctx, types.NamespacedName{Name: Test, Namespace: namespace}, instance)).To(Succeed())
			Expect(instance.Status.Resources).To(Equal("created"))

			instance.Spec.Sidecar.Enabled = true
			instance.Spec.Sidecar.EnableNamespaceByDefault = []string{namespace1}
			instance.Spec.CommonSpec.ServiceAccountSpec.Create = false
			encodedString := fmt.Sprintf("enc::%s::enc", base64.StdEncoding.EncodeToString([]byte(Test)))
			instance.Spec.Secrets.FluxNinjaExtension.Create = true
			instance.Spec.Secrets.FluxNinjaExtension.Value = encodedString
			instance.Annotations = map[string]string{}
			instance.ObjectMeta.Annotations[AgentModeChangeAnnotationKey] = "true"
			Expect(K8sClient.Update(Ctx, instance)).To(Succeed())

			os.Setenv("APERTURE_OPERATOR_CERT_DIR", CertDir)
			os.Setenv("APERTURE_OPERATOR_CERT_NAME", "tls7.crt")
			os.Setenv("APERTURE_OPERATOR_NAMESPACE", AppName)
			os.Setenv("APERTURE_OPERATOR_SERVICE_NAME", AppName)
			certPath := fmt.Sprintf("%s/%s", os.Getenv("APERTURE_OPERATOR_CERT_DIR"), WebhookClientCertName)
			serverCertPEM := new(bytes.Buffer)
			_ = pem.Encode(serverCertPEM, &pem.Block{
				Type:  "CERTIFICATE",
				Bytes: []byte(Test),
			})
			err = WriteFile(certPath, serverCertPEM)
			Expect(err).NotTo(HaveOccurred())

			res, err = reconciler.Reconcile(Ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      Test,
					Namespace: namespace,
				},
			})

			createdAgentService = &corev1.Service{}
			agentServiceKey = types.NamespacedName{Name: AgentServiceName, Namespace: namespace}

			createdAgentDaemonset = &appsv1.DaemonSet{}
			agentDaemonsetKey = types.NamespacedName{Name: AgentServiceName, Namespace: namespace}

			createdMWC = &admissionregistrationv1.MutatingWebhookConfiguration{}
			mwcKey = types.NamespacedName{Name: PodMutatingWebhookName}

			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() bool {
				err1 := K8sClient.Get(Ctx, agentServiceKey, createdAgentService)
				err2 := K8sClient.Get(Ctx, agentDaemonsetKey, createdAgentDaemonset)
				err3 := K8sClient.Get(Ctx, mwcKey, createdMWC)
				return err1 != nil && err2 != nil && err3 == nil
			}, time.Second*10, time.Millisecond*250).Should(BeTrue())

			Expect(K8sClient.Get(Ctx, types.NamespacedName{Name: Test, Namespace: namespace}, instance)).To(Succeed())
			Expect(instance.Status.Resources).To(Equal("created"))

			Expect(K8sClient.Delete(Ctx, ns)).To(Succeed())
			Expect(K8sClient.Delete(Ctx, ns1)).To(Succeed())
		})

		It("should create required resources when Agent is updated to use DaemonSet mode", func() {
			namespace := Test + "24"
			namespace1 := Test + "25"
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			Expect(K8sClient.Create(Ctx, ns)).To(Succeed())

			instance.Namespace = namespace
			instance.Spec.Sidecar.Enabled = true
			instance.Spec.Sidecar.EnableNamespaceByDefault = []string{namespace1}
			instance.Spec.Secrets.FluxNinjaExtension.Create = true
			instance.Spec.Secrets.FluxNinjaExtension.Value = Test
			Expect(K8sClient.Create(Ctx, instance)).To(Succeed())

			ns1 := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace1,
				},
			}
			Expect(K8sClient.Create(Ctx, ns1)).To(Succeed())

			res, err := reconciler.Reconcile(Ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      Test,
					Namespace: namespace,
				},
			})

			createdAgentService := &corev1.Service{}
			agentServiceKey := types.NamespacedName{Name: AgentServiceName, Namespace: namespace}

			createdAgentDaemonset := &appsv1.DaemonSet{}
			agentDaemonsetKey := types.NamespacedName{Name: AgentServiceName, Namespace: namespace}

			createdMWC := &admissionregistrationv1.MutatingWebhookConfiguration{}
			mwcKey := types.NamespacedName{Name: PodMutatingWebhookName}

			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() bool {
				err1 := K8sClient.Get(Ctx, agentServiceKey, createdAgentService)
				err2 := K8sClient.Get(Ctx, agentDaemonsetKey, createdAgentDaemonset)
				err3 := K8sClient.Get(Ctx, mwcKey, createdMWC)
				return err1 != nil && err2 != nil && err3 == nil
			}, time.Second*10, time.Millisecond*250).Should(BeTrue())

			Expect(K8sClient.Get(Ctx, types.NamespacedName{Name: Test, Namespace: namespace}, instance)).To(Succeed())
			Expect(instance.Status.Resources).To(Equal("created"))

			instance.Spec.Sidecar.Enabled = false
			instance.Spec.CommonSpec.ServiceAccountSpec.Create = false
			encodedString := fmt.Sprintf("enc::%s::enc", base64.StdEncoding.EncodeToString([]byte(Test)))
			instance.Spec.Secrets.FluxNinjaExtension.Create = true
			instance.Spec.Secrets.FluxNinjaExtension.Value = encodedString
			instance.Annotations = map[string]string{}
			instance.ObjectMeta.Annotations[AgentModeChangeAnnotationKey] = "true"
			Expect(K8sClient.Update(Ctx, instance)).To(Succeed())

			os.Setenv("APERTURE_OPERATOR_CERT_DIR", CertDir)
			os.Setenv("APERTURE_OPERATOR_CERT_NAME", "tls7.crt")
			os.Setenv("APERTURE_OPERATOR_NAMESPACE", AppName)
			os.Setenv("APERTURE_OPERATOR_SERVICE_NAME", AppName)
			certPath := fmt.Sprintf("%s/%s", os.Getenv("APERTURE_OPERATOR_CERT_DIR"), WebhookClientCertName)
			serverCertPEM := new(bytes.Buffer)
			_ = pem.Encode(serverCertPEM, &pem.Block{
				Type:  "CERTIFICATE",
				Bytes: []byte(Test),
			})
			err = WriteFile(certPath, serverCertPEM)
			Expect(err).NotTo(HaveOccurred())

			res, err = reconciler.Reconcile(Ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      Test,
					Namespace: namespace,
				},
			})

			createdAgentService = &corev1.Service{}
			agentServiceKey = types.NamespacedName{Name: AgentServiceName, Namespace: namespace}

			createdAgentDaemonset = &appsv1.DaemonSet{}
			agentDaemonsetKey = types.NamespacedName{Name: AgentServiceName, Namespace: namespace}

			createdMWC = &admissionregistrationv1.MutatingWebhookConfiguration{}
			mwcKey = types.NamespacedName{Name: PodMutatingWebhookName}

			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() bool {
				err1 := K8sClient.Get(Ctx, agentServiceKey, createdAgentService)
				err2 := K8sClient.Get(Ctx, agentDaemonsetKey, createdAgentDaemonset)
				err3 := K8sClient.Get(Ctx, mwcKey, createdMWC)
				return err1 == nil && err2 == nil && err3 != nil
			}, time.Second*10, time.Millisecond*250).Should(BeTrue())

			Expect(K8sClient.Get(Ctx, types.NamespacedName{Name: Test, Namespace: namespace}, instance)).To(Succeed())
			Expect(instance.Status.Resources).To(Equal("created"))

			Expect(K8sClient.Delete(Ctx, ns)).To(Succeed())
			Expect(K8sClient.Delete(Ctx, ns1)).To(Succeed())
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
		It("should allow only Agent Events in create", func() {
			AgentEvent := event.CreateEvent{
				Object: &agentv1alpha1.Agent{},
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
				ObjectOld: &agentv1alpha1.Agent{},
				ObjectNew: &agentv1alpha1.Agent{
					Spec: agentv1alpha1.AgentSpec{
						CommonSpec: common.CommonSpec{
							Command: TestArray,
						},
					},
				},
			}

			AgentEventValid2 := event.UpdateEvent{
				ObjectOld: &agentv1alpha1.Agent{},
				ObjectNew: &agentv1alpha1.Agent{
					ObjectMeta: metav1.ObjectMeta{
						DeletionTimestamp: &metav1.Time{Time: time.Now()},
					},
				},
			}

			AgentEventInvalid1 := event.UpdateEvent{
				ObjectOld: DefaultAgentInstance,
				ObjectNew: DefaultAgentInstance,
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
				ObjectOld: &agentv1alpha1.Agent{
					Spec: agentv1alpha1.AgentSpec{
						CommonSpec: common.CommonSpec{
							Secrets: common.Secrets{
								FluxNinjaExtension: common.APIKeySecret{
									Value: Test,
								},
							},
						},
					},
				},
				ObjectNew: &agentv1alpha1.Agent{
					Spec: agentv1alpha1.AgentSpec{
						CommonSpec: common.CommonSpec{
							Secrets: common.Secrets{
								FluxNinjaExtension: common.APIKeySecret{
									Value: "",
								},
							},
						},
					},
				},
			}

			AgentEventValid2 := event.UpdateEvent{
				ObjectOld: &agentv1alpha1.Agent{
					Spec: agentv1alpha1.AgentSpec{
						CommonSpec: common.CommonSpec{
							Secrets: common.Secrets{
								FluxNinjaExtension: common.APIKeySecret{
									Value: Test,
								},
							},
						},
					},
				},
				ObjectNew: &agentv1alpha1.Agent{
					Spec: agentv1alpha1.AgentSpec{
						CommonSpec: common.CommonSpec{
							Secrets: common.Secrets{
								FluxNinjaExtension: common.APIKeySecret{
									Value: fmt.Sprintf("enc::%s::enc", base64.StdEncoding.EncodeToString([]byte(Test))),
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
				Object:             &agentv1alpha1.Agent{},
				DeleteStateUnknown: true,
			}

			AgentEventInvalid := event.DeleteEvent{
				Object:             &agentv1alpha1.Agent{},
				DeleteStateUnknown: false,
			}

			pred := eventFiltersForAgent()

			Expect(pred.Delete(AgentEvent)).To(Equal(false))
			Expect(pred.Delete(AgentEventInvalid)).To(Equal(true))
		})
	})
})
