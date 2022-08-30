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
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/fluxninja/aperture/operator/api/v1alpha1"
)

// validatingWebhookConfiguration prepares the ValidatingWebhookConfiguration object for the Controller, based on the provided parameter.
func validatingWebhookConfiguration(instance *v1alpha1.Controller, cert []byte) *admissionregistrationv1.ValidatingWebhookConfiguration {
	validatingWebhookConfiguration := &admissionregistrationv1.ValidatingWebhookConfiguration{
		ObjectMeta: v1.ObjectMeta{
			Name:        validatingWebhookServiceName,
			Labels:      commonLabels(instance.Spec.Labels, instance.GetName(), controllerServiceName),
			Annotations: getControllerAnnotationsWithOwnerRef(instance),
		},
		Webhooks: []admissionregistrationv1.ValidatingWebhook{
			{
				Name: "cm-validator.fluxninja.com",
				ClientConfig: admissionregistrationv1.WebhookClientConfig{
					Service: &admissionregistrationv1.ServiceReference{
						Name:      validatingWebhookServiceName,
						Namespace: instance.GetNamespace(),
						Path:      pointer.StringPtr("/validate/configmap"),
						Port:      pointer.Int32(443),
					},
					CABundle: cert,
				},
				NamespaceSelector: &v1.LabelSelector{
					MatchLabels: map[string]string{
						"kubernetes.io/metadata.name": instance.GetNamespace(),
					},
				},
				ObjectSelector: &v1.LabelSelector{
					MatchLabels: map[string]string{
						"fluxninja.com/validate": "true",
					},
				},
				Rules: []admissionregistrationv1.RuleWithOperations{
					{
						Operations: []admissionregistrationv1.OperationType{"CREATE", "UPDATE"},
						Rule: admissionregistrationv1.Rule{
							APIGroups:   []string{""},
							APIVersions: []string{v1Version},
							Resources:   []string{"configmaps"},
							Scope:       &[]admissionregistrationv1.ScopeType{admissionregistrationv1.NamespacedScope}[0],
						},
					},
				},
				AdmissionReviewVersions: []string{v1Version},
				FailurePolicy:           &[]admissionregistrationv1.FailurePolicyType{admissionregistrationv1.Fail}[0],
				SideEffects:             &[]admissionregistrationv1.SideEffectClass{admissionregistrationv1.SideEffectClassNone}[0],
				TimeoutSeconds:          pointer.Int32Ptr(5),
			},
		},
	}

	return validatingWebhookConfiguration
}

// validatingWebhookConfigurationMutate returns a mutate function that can be used to update the ValidatingWebhookConfiguration's spec.
func validatingWebhookConfigurationMutate(vwc *admissionregistrationv1.ValidatingWebhookConfiguration, webhooks []admissionregistrationv1.ValidatingWebhook) controllerutil.MutateFn {
	return func() error {
		vwc.Webhooks[0].AdmissionReviewVersions = webhooks[0].AdmissionReviewVersions
		vwc.Webhooks[0].ClientConfig = webhooks[0].ClientConfig
		vwc.Webhooks[0].NamespaceSelector = webhooks[0].NamespaceSelector
		vwc.Webhooks[0].ObjectSelector = webhooks[0].ObjectSelector
		vwc.Webhooks[0].Rules = webhooks[0].Rules
		vwc.Webhooks[0].FailurePolicy = webhooks[0].FailurePolicy
		vwc.Webhooks[0].SideEffects = webhooks[0].SideEffects
		vwc.Webhooks[0].TimeoutSeconds = webhooks[0].TimeoutSeconds
		return nil
	}
}
