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
	"reflect"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/utils/pointer"
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
			mutatingWebhookReconciler.AgentManager = true
			mutatingWebhookReconciler.ControllerManager = true

			instance = &admissionregistrationv1.MutatingWebhookConfiguration{
				ObjectMeta: v1.ObjectMeta{
					Name: appName,
				},
				Webhooks: []admissionregistrationv1.MutatingWebhook{
					{
						Name: "agent-defaulter.fluxninja.com",
						ClientConfig: admissionregistrationv1.WebhookClientConfig{
							CABundle: []byte(test),
							Service: &admissionregistrationv1.ServiceReference{
								Name:      "agent-defaulter.fluxninja.com",
								Namespace: test,
								Path:      pointer.String("/agent-defaulter"),
								Port:      pointer.Int32(443),
							},
						},
						NamespaceSelector: &v1.LabelSelector{},
						Rules: []admissionregistrationv1.RuleWithOperations{
							{
								Operations: []admissionregistrationv1.OperationType{"CREATE", "UPDATE"},
								Rule: admissionregistrationv1.Rule{
									APIGroups:   []string{"fluxninja.com"},
									APIVersions: []string{v1Alpha1Version},
									Resources:   []string{test},
									Scope:       &[]admissionregistrationv1.ScopeType{admissionregistrationv1.NamespacedScope}[0],
								},
							},
						},
						AdmissionReviewVersions: []string{v1Version},
						FailurePolicy:           &[]admissionregistrationv1.FailurePolicyType{admissionregistrationv1.Fail}[0],
						MatchPolicy:             &[]admissionregistrationv1.MatchPolicyType{admissionregistrationv1.Equivalent}[0],
						SideEffects:             &[]admissionregistrationv1.SideEffectClass{admissionregistrationv1.SideEffectClassNone}[0],
						TimeoutSeconds:          pointer.Int32Ptr(10),
					},
				},
			}
		})

		It("should not update resources when MutatingWebhookConfiguration is deleted", func() {
			res, err := mutatingWebhookReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name: appName,
				},
			})
			Expect(reflect.DeepEqual(res, ctrl.Result{})).To(Equal(true))
			Expect(err).ToNot(HaveOccurred())
		})

		It("should not update agent MutatingWebhookConfiguration when started for Controller manager", func() {
			copiedInstance := instance.DeepCopy()
			copiedInstance.Name = agentMutatingWebhookName
			Expect(k8sClient.Create(ctx, copiedInstance)).To(BeNil())

			mutatingWebhookReconciler.AgentManager = false

			_, err := mutatingWebhookReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name: agentMutatingWebhookName,
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

			Expect(k8sClient.Delete(ctx, copiedInstance)).To(BeNil())
		})

		It("should not update controller MutatingWebhookConfiguration when started for Agent manager", func() {
			copiedInstance := instance.DeepCopy()
			copiedInstance.Name = controllerMutatingWebhookName
			Expect(k8sClient.Create(ctx, copiedInstance)).To(BeNil())

			mutatingWebhookReconciler.ControllerManager = false

			_, err := mutatingWebhookReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name: controllerMutatingWebhookName,
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
			Expect(k8sClient.Delete(ctx, copiedInstance)).To(BeNil())
		})
	})

	Context("testing mwcEventFilters", func() {
		It("should allow only MutatingWebhookConfiguration Events for Agent and Controller in create", func() {
			mwcEventValid1 := event.CreateEvent{
				Object: &admissionregistrationv1.MutatingWebhookConfiguration{
					ObjectMeta: v1.ObjectMeta{
						Name: agentMutatingWebhookName,
					},
				},
			}

			mwcEventValid2 := event.CreateEvent{
				Object: &admissionregistrationv1.MutatingWebhookConfiguration{
					ObjectMeta: v1.ObjectMeta{
						Name: controllerMutatingWebhookName,
					},
				},
			}

			mwcEventInValid := event.CreateEvent{
				Object: &admissionregistrationv1.MutatingWebhookConfiguration{
					ObjectMeta: v1.ObjectMeta{
						Name: test,
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
						Name: agentMutatingWebhookName,
					},
				},
				ObjectNew: &admissionregistrationv1.MutatingWebhookConfiguration{
					ObjectMeta: v1.ObjectMeta{
						Name: agentMutatingWebhookName,
					},
					Webhooks: []admissionregistrationv1.MutatingWebhook{
						{
							Name: test,
						},
					},
				},
			}

			mwcEventValid2 := event.UpdateEvent{
				ObjectOld: &admissionregistrationv1.MutatingWebhookConfiguration{
					ObjectMeta: v1.ObjectMeta{
						Name: controllerMutatingWebhookName,
					},
				},
				ObjectNew: &admissionregistrationv1.MutatingWebhookConfiguration{
					ObjectMeta: v1.ObjectMeta{
						Name: controllerMutatingWebhookName,
					},
					Webhooks: []admissionregistrationv1.MutatingWebhook{
						{
							Name: test,
						},
					},
				},
			}

			mwcEventInValid1 := event.UpdateEvent{
				ObjectOld: &admissionregistrationv1.MutatingWebhookConfiguration{
					ObjectMeta: v1.ObjectMeta{
						Name: test,
					},
				},
				ObjectNew: &admissionregistrationv1.MutatingWebhookConfiguration{
					ObjectMeta: v1.ObjectMeta{
						Name: test,
					},
				},
			}

			mwcEventInValid2 := event.UpdateEvent{
				ObjectOld: &admissionregistrationv1.MutatingWebhookConfiguration{
					ObjectMeta: v1.ObjectMeta{
						Name: controllerMutatingWebhookName,
					},
				},
				ObjectNew: &admissionregistrationv1.MutatingWebhookConfiguration{
					ObjectMeta: v1.ObjectMeta{
						Name: controllerMutatingWebhookName,
					},
				},
			}

			mwcEventInValid3 := event.UpdateEvent{
				ObjectOld: &admissionregistrationv1.MutatingWebhookConfiguration{
					ObjectMeta: v1.ObjectMeta{
						Name: agentMutatingWebhookName,
					},
				},
				ObjectNew: &admissionregistrationv1.MutatingWebhookConfiguration{
					ObjectMeta: v1.ObjectMeta{
						Name: agentMutatingWebhookName,
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

			pred := eventFiltersForAgent()

			Expect(pred.Delete(mwcEventValid)).To(Equal(false))
			Expect(pred.Delete(mwcEventInValid)).To(Equal(true))
		})

	})
})
