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

package mutatingwebhook

import (
	"reflect"

	. "github.com/fluxninja/aperture/v2/operator/controllers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("MutatingWebhookConfiguration controller", Ordered, func() {
	Context("testing Reconcile", func() {

		var instance *admissionregistrationv1.MutatingWebhookConfiguration

		BeforeEach(func() {
			mutatingWebhookTestReconciler.AgentManager = true
			mutatingWebhookTestReconciler.ControllerManager = true

			instance = &admissionregistrationv1.MutatingWebhookConfiguration{
				ObjectMeta: v1.ObjectMeta{
					Name: AppName,
				},
				Webhooks: []admissionregistrationv1.MutatingWebhook{
					{
						Name: "agent-defaulter.fluxninja.com",
						ClientConfig: admissionregistrationv1.WebhookClientConfig{
							CABundle: []byte(Test),
							Service: &admissionregistrationv1.ServiceReference{
								Name:      "agent-defaulter.fluxninja.com",
								Namespace: Test,
								Path:      ptr.To("/agent-defaulter"),
								Port:      ptr.To[int32](443),
							},
						},
						NamespaceSelector: &v1.LabelSelector{},
						Rules: []admissionregistrationv1.RuleWithOperations{
							{
								Operations: []admissionregistrationv1.OperationType{"CREATE", "UPDATE"},
								Rule: admissionregistrationv1.Rule{
									APIGroups:   []string{"fluxninja.com"},
									APIVersions: []string{V1Alpha1Version},
									Resources:   []string{Test},
									Scope:       &[]admissionregistrationv1.ScopeType{admissionregistrationv1.NamespacedScope}[0],
								},
							},
						},
						AdmissionReviewVersions: []string{V1Version},
						FailurePolicy:           &[]admissionregistrationv1.FailurePolicyType{admissionregistrationv1.Fail}[0],
						MatchPolicy:             &[]admissionregistrationv1.MatchPolicyType{admissionregistrationv1.Equivalent}[0],
						SideEffects:             &[]admissionregistrationv1.SideEffectClass{admissionregistrationv1.SideEffectClassNone}[0],
						TimeoutSeconds:          ptr.To[int32](10),
					},
				},
			}
		})

		It("should not update resources when MutatingWebhookConfiguration is deleted", func() {
			res, err := mutatingWebhookTestReconciler.Reconcile(Ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name: AppName,
				},
			})
			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).ToNot(HaveOccurred())
		})

		It("should not update agent MutatingWebhookConfiguration when started for Controller manager", func() {
			copiedInstance := instance.DeepCopy()
			copiedInstance.Name = AgentMutatingWebhookName
			Expect(K8sClient.Create(Ctx, copiedInstance)).To(Succeed())

			mutatingWebhookTestReconciler.AgentManager = false

			_, err := mutatingWebhookTestReconciler.Reconcile(Ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name: AgentMutatingWebhookName,
				},
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(copiedInstance.Webhooks[0].AdmissionReviewVersions).To(Equal(instance.Webhooks[0].AdmissionReviewVersions))
			Expect(copiedInstance.Webhooks[0].ClientConfig).To(Equal(instance.Webhooks[0].ClientConfig))
			Expect(copiedInstance.Webhooks[0].NamespaceSelector).To(Equal(instance.Webhooks[0].NamespaceSelector))
			Expect(copiedInstance.Webhooks[0].Rules).To(Equal(instance.Webhooks[0].Rules))
			Expect(copiedInstance.Webhooks[0].FailurePolicy).To(Equal(instance.Webhooks[0].FailurePolicy))
			Expect(copiedInstance.Webhooks[0].SideEffects).To(Equal(instance.Webhooks[0].SideEffects))
			Expect(copiedInstance.Webhooks[0].TimeoutSeconds).To(Equal(instance.Webhooks[0].TimeoutSeconds))

			Expect(K8sClient.Delete(Ctx, copiedInstance)).To(Succeed())
		})

		It("should not update controller MutatingWebhookConfiguration when started for Agent manager", func() {
			copiedInstance := instance.DeepCopy()
			copiedInstance.Name = ControllerMutatingWebhookName
			Expect(K8sClient.Create(Ctx, copiedInstance)).To(Succeed())

			mutatingWebhookTestReconciler.ControllerManager = false

			_, err := mutatingWebhookTestReconciler.Reconcile(Ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name: ControllerMutatingWebhookName,
				},
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(copiedInstance.Webhooks[0].AdmissionReviewVersions).To(Equal(instance.Webhooks[0].AdmissionReviewVersions))
			Expect(copiedInstance.Webhooks[0].ClientConfig).To(Equal(instance.Webhooks[0].ClientConfig))
			Expect(copiedInstance.Webhooks[0].NamespaceSelector).To(Equal(instance.Webhooks[0].NamespaceSelector))
			Expect(copiedInstance.Webhooks[0].Rules).To(Equal(instance.Webhooks[0].Rules))
			Expect(copiedInstance.Webhooks[0].FailurePolicy).To(Equal(instance.Webhooks[0].FailurePolicy))
			Expect(copiedInstance.Webhooks[0].SideEffects).To(Equal(instance.Webhooks[0].SideEffects))
			Expect(copiedInstance.Webhooks[0].TimeoutSeconds).To(Equal(instance.Webhooks[0].TimeoutSeconds))
			Expect(K8sClient.Delete(Ctx, copiedInstance)).To(Succeed())
		})
	})

	Context("testing mwcEventFilters", func() {
		It("should allow only MutatingWebhookConfiguration Events for Agent and Controller in create", func() {
			mwcEventValid1 := event.CreateEvent{
				Object: &admissionregistrationv1.MutatingWebhookConfiguration{
					ObjectMeta: v1.ObjectMeta{
						Name: AgentMutatingWebhookName,
					},
				},
			}

			mwcEventValid2 := event.CreateEvent{
				Object: &admissionregistrationv1.MutatingWebhookConfiguration{
					ObjectMeta: v1.ObjectMeta{
						Name: ControllerMutatingWebhookName,
					},
				},
			}

			mwcEventInValid := event.CreateEvent{
				Object: &admissionregistrationv1.MutatingWebhookConfiguration{
					ObjectMeta: v1.ObjectMeta{
						Name: Test,
					},
				},
			}

			vwcEvent := event.CreateEvent{
				Object: &admissionregistrationv1.ValidatingWebhookConfiguration{},
			}

			pred := mwcEventFilters()

			Expect(pred.Create(mwcEventValid1)).To(Equal(true))
			Expect(pred.Create(mwcEventValid2)).To(Equal(true))
			Expect(pred.Create(mwcEventInValid)).To(Equal(false))
			Expect(pred.Create(vwcEvent)).To(Equal(false))
		})

		It("should allow only MutatingWebhookConfiguration Events for Agent or Controller in update", func() {
			mwcEventValid1 := event.UpdateEvent{
				ObjectOld: &admissionregistrationv1.MutatingWebhookConfiguration{
					ObjectMeta: v1.ObjectMeta{
						Name: AgentMutatingWebhookName,
					},
				},
				ObjectNew: &admissionregistrationv1.MutatingWebhookConfiguration{
					ObjectMeta: v1.ObjectMeta{
						Name: AgentMutatingWebhookName,
					},
					Webhooks: []admissionregistrationv1.MutatingWebhook{
						{
							Name: Test,
						},
					},
				},
			}

			mwcEventValid2 := event.UpdateEvent{
				ObjectOld: &admissionregistrationv1.MutatingWebhookConfiguration{
					ObjectMeta: v1.ObjectMeta{
						Name: ControllerMutatingWebhookName,
					},
				},
				ObjectNew: &admissionregistrationv1.MutatingWebhookConfiguration{
					ObjectMeta: v1.ObjectMeta{
						Name: ControllerMutatingWebhookName,
					},
					Webhooks: []admissionregistrationv1.MutatingWebhook{
						{
							Name: Test,
						},
					},
				},
			}

			mwcEventInValid1 := event.UpdateEvent{
				ObjectOld: &admissionregistrationv1.MutatingWebhookConfiguration{
					ObjectMeta: v1.ObjectMeta{
						Name: Test,
					},
				},
				ObjectNew: &admissionregistrationv1.MutatingWebhookConfiguration{
					ObjectMeta: v1.ObjectMeta{
						Name: Test,
					},
				},
			}

			mwcEventInValid2 := event.UpdateEvent{
				ObjectOld: &admissionregistrationv1.MutatingWebhookConfiguration{
					ObjectMeta: v1.ObjectMeta{
						Name: ControllerMutatingWebhookName,
					},
				},
				ObjectNew: &admissionregistrationv1.MutatingWebhookConfiguration{
					ObjectMeta: v1.ObjectMeta{
						Name: ControllerMutatingWebhookName,
					},
				},
			}

			mwcEventInValid3 := event.UpdateEvent{
				ObjectOld: &admissionregistrationv1.MutatingWebhookConfiguration{
					ObjectMeta: v1.ObjectMeta{
						Name: AgentMutatingWebhookName,
					},
				},
				ObjectNew: &admissionregistrationv1.MutatingWebhookConfiguration{
					ObjectMeta: v1.ObjectMeta{
						Name: AgentMutatingWebhookName,
					},
				},
			}

			vwcEvent := event.UpdateEvent{
				ObjectOld: &admissionregistrationv1.ValidatingWebhookConfiguration{},
				ObjectNew: &admissionregistrationv1.ValidatingWebhookConfiguration{},
			}

			pred := mwcEventFilters()

			Expect(pred.Update(mwcEventValid1)).To(Equal(true))
			Expect(pred.Update(mwcEventValid2)).To(Equal(true))
			Expect(pred.Update(mwcEventInValid1)).To(Equal(false))
			Expect(pred.Update(mwcEventInValid2)).To(Equal(false))
			Expect(pred.Update(mwcEventInValid3)).To(Equal(false))
			Expect(pred.Update(vwcEvent)).To(Equal(false))
		})

		It("should allow Events in delete when state is known", func() {
			mwcEventValid := event.DeleteEvent{
				Object:             &admissionregistrationv1.MutatingWebhookConfiguration{},
				DeleteStateUnknown: true,
			}

			mwcEventInValid := event.DeleteEvent{
				Object:             &admissionregistrationv1.MutatingWebhookConfiguration{},
				DeleteStateUnknown: false,
			}

			pred := mwcEventFilters()

			Expect(pred.Delete(mwcEventValid)).To(Equal(false))
			Expect(pred.Delete(mwcEventInValid)).To(Equal(true))
		})

	})
})
