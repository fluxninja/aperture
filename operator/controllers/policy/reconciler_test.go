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

package policy

import (
	"os"
	"path/filepath"
	"reflect"
	"time"

	. "github.com/fluxninja/aperture/operator/controllers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/yaml"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	policyv1alpha1 "github.com/fluxninja/aperture/operator/api/policy/v1alpha1"
)

var _ = Describe("Policy Reconciler", Ordered, func() {
	Context("testing Reconcile", func() {
		var instance *policyv1alpha1.Policy
		policyTestReconciler := &PolicyReconciler{
			Client:   K8sClient,
			Scheme:   K8sManager.GetScheme(),
			Recorder: K8sManager.GetEventRecorderFor("aperture-policy"),
		}

		BeforeEach(func() {
			instance = DefaultPolicyInstance.DeepCopy()
			PolicyFilePath = PoliciesDir
		})

		It("should not create resources when Policy is deleted", func() {
			res, err := policyTestReconciler.Reconcile(Ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      Test,
					Namespace: Test,
				},
			})
			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).NotTo(HaveOccurred())
		})

		It("should create file when valid Policy is created", func() {
			namespace := Test + "101"
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			Expect(K8sClient.Create(Ctx, ns)).To(BeNil())

			instance.Namespace = namespace
			Expect(K8sClient.Create(Ctx, instance)).To(BeNil())

			res, err := policyTestReconciler.Reconcile(Ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      AppName,
					Namespace: namespace,
				},
			})
			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).NotTo(HaveOccurred())

			err = K8sClient.Get(Ctx, types.NamespacedName{Name: AppName, Namespace: namespace}, instance)
			Expect(err).NotTo(HaveOccurred())

			Expect(instance.Status.Status).To(Equal("uploaded"))
			fileContent, err := os.OpenFile(filepath.Join(PolicyFilePath, "aperture-test101.yaml"), os.O_RDONLY, 0o400)
			Expect(err).NotTo(HaveOccurred())
			yamlContent, err := yaml.JSONToYAML(instance.Spec.Raw)
			Expect(err).NotTo(HaveOccurred())
			content := make([]byte, len(yamlContent))
			fileContent.Read(content)
			Expect(string(content)).To(Equal(string(yamlContent)))
		})

		It("should not create file when invalid Policy is created", func() {
			namespace := Test + "102"
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			Expect(K8sClient.Create(Ctx, ns)).To(BeNil())

			instance.Namespace = namespace
			instance.Spec.Raw = []byte("{\"test\":\"test\"}")
			Expect(K8sClient.Create(Ctx, instance)).To(BeNil())

			res, err := policyTestReconciler.Reconcile(Ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      AppName,
					Namespace: namespace,
				},
			})
			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).NotTo(HaveOccurred())

			err = K8sClient.Get(Ctx, types.NamespacedName{Name: AppName, Namespace: namespace}, instance)
			Expect(err).NotTo(HaveOccurred())

			Expect(instance.Status.Status).To(Equal(FailedStatus))
			_, err = os.OpenFile(filepath.Join(PolicyFilePath, "aperture-test102.yaml"), os.O_RDONLY, 0o400)
			Expect(err).To(HaveOccurred())
		})

		It("should create file when invalid Policy is created with annotation", func() {
			namespace := Test + "103"
			os.Setenv("APERTURE_CONTROLLER_NAMESPACE", namespace)
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			Expect(K8sClient.Create(Ctx, ns)).To(BeNil())

			instance.Namespace = namespace
			instance.Spec.Raw = []byte("{\"test\":\"test\"}")
			annotations := map[string]string{
				DefaulterAnnotationKey: "true",
			}
			instance.Annotations = annotations
			Expect(K8sClient.Create(Ctx, instance)).To(BeNil())

			res, err := policyTestReconciler.Reconcile(Ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      AppName,
					Namespace: namespace,
				},
			})
			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).NotTo(HaveOccurred())

			err = K8sClient.Get(Ctx, types.NamespacedName{Name: AppName, Namespace: namespace}, instance)
			Expect(err).NotTo(HaveOccurred())

			Expect(instance.Status.Status).To(Equal("uploaded"))
			fileContent, err := os.OpenFile(filepath.Join(PolicyFilePath, "aperture.yaml"), os.O_RDONLY, 0o400)
			Expect(err).NotTo(HaveOccurred())
			yamlContent, err := yaml.JSONToYAML(instance.Spec.Raw)
			Expect(err).NotTo(HaveOccurred())
			content := make([]byte, len(yamlContent))
			fileContent.Read(content)
			Expect(string(content)).To(Equal(string(yamlContent)))
		})

		It("should delete file when valid Policy is deleted", func() {
			namespace := Test + "104"
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			Expect(K8sClient.Create(Ctx, ns)).To(BeNil())

			instance.Namespace = namespace
			Expect(K8sClient.Create(Ctx, instance)).To(BeNil())

			res, err := policyTestReconciler.Reconcile(Ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      AppName,
					Namespace: namespace,
				},
			})
			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).NotTo(HaveOccurred())

			err = K8sClient.Get(Ctx, types.NamespacedName{Name: AppName, Namespace: namespace}, instance)
			Expect(err).NotTo(HaveOccurred())

			Expect(instance.Status.Status).To(Equal("uploaded"))
			_, err = os.OpenFile(filepath.Join(PolicyFilePath, "aperture-test104.yaml"), os.O_RDONLY, 0o400)
			Expect(err).NotTo(HaveOccurred())

			Expect(K8sClient.Delete(Ctx, instance)).NotTo(HaveOccurred())
			res, err = policyTestReconciler.Reconcile(Ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      AppName,
					Namespace: namespace,
				},
			})
			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() bool {
				err := K8sClient.Get(Ctx, types.NamespacedName{Name: AppName, Namespace: namespace}, instance)
				return err != nil
			}, time.Second*10, time.Millisecond*250).Should(BeTrue())

			_, err = os.OpenFile(filepath.Join(PolicyFilePath, "aperture-test104.yaml"), os.O_RDONLY, 0o400)
			Expect(err).To(HaveOccurred())
		})

		AfterEach(func() {
			K8sClient.Delete(Ctx, instance)
		})
	})

	Context("testing eventFiltersForPolicy", func() {
		It("should allow only Policy Events in create", func() {
			policyEvent := event.CreateEvent{
				Object: &policyv1alpha1.Policy{},
			}
			cmEvent := event.CreateEvent{
				Object: &corev1.ConfigMap{},
			}
			secretEvent := event.CreateEvent{
				Object: &corev1.Secret{},
			}

			pred := eventFiltersForPolicy()

			Expect(pred.Create(policyEvent)).To(Equal(true))
			Expect(pred.Create(cmEvent)).To(Equal(false))
			Expect(pred.Create(secretEvent)).To(Equal(false))
		})

		It("should allow only Policy Events in update", func() {
			policyValidEvent := event.UpdateEvent{
				ObjectOld: &policyv1alpha1.Policy{
					Spec: runtime.RawExtension{
						Raw: []byte(Test),
					},
				},
				ObjectNew: &policyv1alpha1.Policy{
					Spec: runtime.RawExtension{
						Raw: []byte(TestTwo),
					},
				},
			}

			policyInValidEvent := event.UpdateEvent{
				ObjectOld: &policyv1alpha1.Policy{
					Spec: runtime.RawExtension{
						Raw: []byte(Test),
					},
				},
				ObjectNew: &policyv1alpha1.Policy{
					Spec: runtime.RawExtension{
						Raw: []byte(Test),
					},
				},
			}

			pred := eventFiltersForPolicy()

			Expect(pred.Update(policyValidEvent)).To(Equal(true))
			Expect(pred.Update(policyInValidEvent)).To(Equal(false))
		})

		It("should not allow policy Events in delete", func() {
			policyEvent := event.DeleteEvent{
				Object:             &policyv1alpha1.Policy{},
				DeleteStateUnknown: true,
			}

			policyEvent2 := event.DeleteEvent{
				Object:             &policyv1alpha1.Policy{},
				DeleteStateUnknown: false,
			}

			pred := eventFiltersForPolicy()

			Expect(pred.Delete(policyEvent)).To(Equal(false))
			Expect(pred.Delete(policyEvent2)).To(Equal(true))
		})

	})
})
