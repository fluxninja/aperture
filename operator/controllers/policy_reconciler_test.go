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
	"os"
	"path/filepath"
	"reflect"
	"time"

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

	"github.com/fluxninja/aperture/operator/api/v1alpha1"
)

var _ = Describe("Policy Reconciler", Ordered, func() {
	Context("testing Reconcile", func() {
		var instance *v1alpha1.Policy

		BeforeEach(func() {
			instance = defaultPolicyInstance.DeepCopy()
			policyFilePath = policiesDir
		})

		It("should not create resources when Policy is deleted", func() {
			res, err := policyReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      test,
					Namespace: test,
				},
			})
			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).NotTo(HaveOccurred())
		})

		It("should create file when valid Policy is created", func() {
			namespace := test + "101"
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			Expect(k8sClient.Create(ctx, ns)).To(BeNil())

			instance.Namespace = namespace
			Expect(k8sClient.Create(ctx, instance)).To(BeNil())

			res, err := policyReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      appName,
					Namespace: namespace,
				},
			})
			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).NotTo(HaveOccurred())

			err = k8sClient.Get(ctx, types.NamespacedName{Name: appName, Namespace: namespace}, instance)
			Expect(err).NotTo(HaveOccurred())

			Expect(instance.Status.Status).To(Equal("uploaded"))
			fileContent, err := os.OpenFile(filepath.Join(policyFilePath, "aperture-test101.yaml"), os.O_RDONLY, 0o400)
			Expect(err).NotTo(HaveOccurred())
			yamlContent, err := yaml.JSONToYAML(instance.Spec.Raw)
			Expect(err).NotTo(HaveOccurred())
			content := make([]byte, len(yamlContent))
			fileContent.Read(content)
			Expect(string(content)).To(Equal(string(yamlContent)))
		})

		It("should not create file when invalid Policy is created", func() {
			namespace := test + "102"
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			Expect(k8sClient.Create(ctx, ns)).To(BeNil())

			instance.Namespace = namespace
			instance.Spec.Raw = []byte("{\"test\":\"test\"}")
			Expect(k8sClient.Create(ctx, instance)).To(BeNil())

			res, err := policyReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      appName,
					Namespace: namespace,
				},
			})
			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).NotTo(HaveOccurred())

			err = k8sClient.Get(ctx, types.NamespacedName{Name: appName, Namespace: namespace}, instance)
			Expect(err).NotTo(HaveOccurred())

			Expect(instance.Status.Status).To(Equal(failedStatus))
			_, err = os.OpenFile(filepath.Join(policyFilePath, "aperture-test102.yaml"), os.O_RDONLY, 0o400)
			Expect(err).To(HaveOccurred())
		})

		It("should create file when invalid Policy is created with annotation", func() {
			namespace := test + "103"
			os.Setenv("APERTURE_CONTROLLER_NAMESPACE", namespace)
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			Expect(k8sClient.Create(ctx, ns)).To(BeNil())

			instance.Namespace = namespace
			instance.Spec.Raw = []byte("{\"test\":\"test\"}")
			annotations := map[string]string{
				defaulterAnnotationKey: "true",
			}
			instance.Annotations = annotations
			Expect(k8sClient.Create(ctx, instance)).To(BeNil())

			res, err := policyReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      appName,
					Namespace: namespace,
				},
			})
			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).NotTo(HaveOccurred())

			err = k8sClient.Get(ctx, types.NamespacedName{Name: appName, Namespace: namespace}, instance)
			Expect(err).NotTo(HaveOccurred())

			Expect(instance.Status.Status).To(Equal("uploaded"))
			fileContent, err := os.OpenFile(filepath.Join(policyFilePath, "aperture.yaml"), os.O_RDONLY, 0o400)
			Expect(err).NotTo(HaveOccurred())
			yamlContent, err := yaml.JSONToYAML(instance.Spec.Raw)
			Expect(err).NotTo(HaveOccurred())
			content := make([]byte, len(yamlContent))
			fileContent.Read(content)
			Expect(string(content)).To(Equal(string(yamlContent)))
		})

		It("should delete file when valid Policy is deleted", func() {
			namespace := test + "104"
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			Expect(k8sClient.Create(ctx, ns)).To(BeNil())

			instance.Namespace = namespace
			Expect(k8sClient.Create(ctx, instance)).To(BeNil())

			res, err := policyReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      appName,
					Namespace: namespace,
				},
			})
			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).NotTo(HaveOccurred())

			err = k8sClient.Get(ctx, types.NamespacedName{Name: appName, Namespace: namespace}, instance)
			Expect(err).NotTo(HaveOccurred())

			Expect(instance.Status.Status).To(Equal("uploaded"))
			_, err = os.OpenFile(filepath.Join(policyFilePath, "aperture-test104.yaml"), os.O_RDONLY, 0o400)
			Expect(err).NotTo(HaveOccurred())

			Expect(k8sClient.Delete(ctx, instance)).NotTo(HaveOccurred())
			res, err = policyReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      appName,
					Namespace: namespace,
				},
			})
			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: appName, Namespace: namespace}, instance)
				return err != nil
			}, time.Second*10, time.Millisecond*250).Should(BeTrue())

			_, err = os.OpenFile(filepath.Join(policyFilePath, "aperture-test104.yaml"), os.O_RDONLY, 0o400)
			Expect(err).To(HaveOccurred())
		})

		AfterEach(func() {
			k8sClient.Delete(ctx, instance)
		})
	})

	Context("testing eventFiltersForPolicy", func() {
		It("should allow only Policy Events in create", func() {
			policyEvent := event.CreateEvent{
				Object: &v1alpha1.Policy{},
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
				ObjectOld: &v1alpha1.Policy{
					Spec: runtime.RawExtension{
						Raw: []byte(test),
					},
				},
				ObjectNew: &v1alpha1.Policy{
					Spec: runtime.RawExtension{
						Raw: []byte(testTwo),
					},
				},
			}

			policyInValidEvent := event.UpdateEvent{
				ObjectOld: &v1alpha1.Policy{
					Spec: runtime.RawExtension{
						Raw: []byte(test),
					},
				},
				ObjectNew: &v1alpha1.Policy{
					Spec: runtime.RawExtension{
						Raw: []byte(test),
					},
				},
			}

			pred := eventFiltersForPolicy()

			Expect(pred.Update(policyValidEvent)).To(Equal(true))
			Expect(pred.Update(policyInValidEvent)).To(Equal(false))
		})

		It("should not allow policy Events in delete", func() {
			policyEvent := event.DeleteEvent{
				Object:             &v1alpha1.Policy{},
				DeleteStateUnknown: true,
			}

			policyEvent2 := event.DeleteEvent{
				Object:             &v1alpha1.Policy{},
				DeleteStateUnknown: false,
			}

			pred := eventFiltersForPolicy()

			Expect(pred.Delete(policyEvent)).To(Equal(false))
			Expect(pred.Delete(policyEvent2)).To(Equal(true))
		})

	})
})
