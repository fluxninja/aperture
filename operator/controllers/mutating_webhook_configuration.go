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
	"fmt"
	"os"

	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/fluxninja/aperture/operator/api/v1alpha1"
)

// podMutatingWebhookConfiguration prepares the MutatingWebhookConfiguration object for the Operator to mutate Pods, based on the provided parameter.
func podMutatingWebhookConfiguration(instance *v1alpha1.Agent) (*admissionregistrationv1.MutatingWebhookConfiguration, error) {
	certPath := fmt.Sprintf("%s/%s", os.Getenv("APERTURE_OPERATOR_CERT_DIR"), webhookClientCertName)
	cert, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}

	mutatingWebhookConfiguration := &admissionregistrationv1.MutatingWebhookConfiguration{
		ObjectMeta: v1.ObjectMeta{
			Name:        podMutatingWebhookName,
			Labels:      commonLabels(instance.Spec.Labels, instance.GetName(), operatorName),
			Annotations: getAgentAnnotationsWithOwnerRef(instance),
		},
		Webhooks: []admissionregistrationv1.MutatingWebhook{
			{
				Name: sidecarKey,
				ClientConfig: admissionregistrationv1.WebhookClientConfig{
					CABundle: cert,
					Service: &admissionregistrationv1.ServiceReference{
						Name:      os.Getenv("APERTURE_OPERATOR_SERVICE_NAME"),
						Namespace: os.Getenv("APERTURE_OPERATOR_NAMESPACE"),
						Path:      pointer.StringPtr(MutatingWebhookURI),
						Port:      pointer.Int32(443),
					},
				},
				NamespaceSelector: &v1.LabelSelector{
					MatchExpressions: []v1.LabelSelectorRequirement{
						{
							Key:      sidecarLabelKey,
							Operator: v1.LabelSelectorOpIn,
							Values:   []string{enabled},
						},
					},
				},
				Rules: []admissionregistrationv1.RuleWithOperations{
					{
						Operations: []admissionregistrationv1.OperationType{"CREATE"},
						Rule: admissionregistrationv1.Rule{
							APIGroups:   []string{""},
							APIVersions: []string{v1Version},
							Resources:   []string{"pods"},
							Scope:       &[]admissionregistrationv1.ScopeType{admissionregistrationv1.NamespacedScope}[0],
						},
					},
				},
				AdmissionReviewVersions: []string{v1Version},
				FailurePolicy:           &[]admissionregistrationv1.FailurePolicyType{admissionregistrationv1.Fail}[0],
				SideEffects:             &[]admissionregistrationv1.SideEffectClass{admissionregistrationv1.SideEffectClassNone}[0],
				TimeoutSeconds:          pointer.Int32Ptr(10),
			},
		},
	}

	return mutatingWebhookConfiguration, nil
}

// agentMutatingWebhookConfiguration prepares the MutatingWebhookConfiguration object for the Operator to mutate Agents, based on the provided parameter.
func agentMutatingWebhookConfiguration() (*admissionregistrationv1.MutatingWebhookConfiguration, error) {
	certPath := fmt.Sprintf("%s/%s", os.Getenv("APERTURE_OPERATOR_CERT_DIR"), webhookClientCertName)
	cert, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}

	mutatingWebhookConfiguration := &admissionregistrationv1.MutatingWebhookConfiguration{
		ObjectMeta: v1.ObjectMeta{
			Name: agentMutatingWebhookName,
			Labels: map[string]string{
				"app.kubernetes.io/name":       appName,
				"app.kubernetes.io/instance":   operatorName,
				"app.kubernetes.io/managed-by": operatorName,
				"app.kubernetes.io/component":  operatorName,
			},
		},
		Webhooks: []admissionregistrationv1.MutatingWebhook{
			{
				Name: fmt.Sprintf("%s.fluxninja.com", AgentMutatingWebhookURI),
				ClientConfig: admissionregistrationv1.WebhookClientConfig{
					CABundle: cert,
					Service: &admissionregistrationv1.ServiceReference{
						Name:      os.Getenv("APERTURE_OPERATOR_SERVICE_NAME"),
						Namespace: os.Getenv("APERTURE_OPERATOR_NAMESPACE"),
						Path:      pointer.StringPtr(fmt.Sprintf("/%s", AgentMutatingWebhookURI)),
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
							Resources:   []string{"agents"},
							Scope:       &[]admissionregistrationv1.ScopeType{admissionregistrationv1.NamespacedScope}[0],
						},
					},
				},
				AdmissionReviewVersions: []string{v1Version},
				FailurePolicy:           &[]admissionregistrationv1.FailurePolicyType{admissionregistrationv1.Fail}[0],
				SideEffects:             &[]admissionregistrationv1.SideEffectClass{admissionregistrationv1.SideEffectClassNone}[0],
				TimeoutSeconds:          pointer.Int32Ptr(10),
			},
		},
	}

	return mutatingWebhookConfiguration, nil
}

// controllerMutatingWebhookConfiguration prepares the MutatingWebhookConfiguration object for the Operator to mutate Controllers, based on the provided parameter.
func controllerMutatingWebhookConfiguration() (*admissionregistrationv1.MutatingWebhookConfiguration, error) {
	certPath := fmt.Sprintf("%s/%s", os.Getenv("APERTURE_OPERATOR_CERT_DIR"), webhookClientCertName)
	cert, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}

	mutatingWebhookConfiguration := &admissionregistrationv1.MutatingWebhookConfiguration{
		ObjectMeta: v1.ObjectMeta{
			Name: controllerMutatingWebhookName,
			Labels: map[string]string{
				"app.kubernetes.io/name":       appName,
				"app.kubernetes.io/instance":   operatorName,
				"app.kubernetes.io/managed-by": operatorName,
				"app.kubernetes.io/component":  operatorName,
			},
		},
		Webhooks: []admissionregistrationv1.MutatingWebhook{
			{
				Name: fmt.Sprintf("%s.fluxninja.com", ControllerMutatingWebhookURI),
				ClientConfig: admissionregistrationv1.WebhookClientConfig{
					CABundle: cert,
					Service: &admissionregistrationv1.ServiceReference{
						Name:      os.Getenv("APERTURE_OPERATOR_SERVICE_NAME"),
						Namespace: os.Getenv("APERTURE_OPERATOR_NAMESPACE"),
						Path:      pointer.StringPtr(fmt.Sprintf("/%s", ControllerMutatingWebhookURI)),
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
							Resources:   []string{"controllers"},
							Scope:       &[]admissionregistrationv1.ScopeType{admissionregistrationv1.NamespacedScope}[0],
						},
					},
				},
				AdmissionReviewVersions: []string{v1Version},
				FailurePolicy:           &[]admissionregistrationv1.FailurePolicyType{admissionregistrationv1.Fail}[0],
				SideEffects:             &[]admissionregistrationv1.SideEffectClass{admissionregistrationv1.SideEffectClassNone}[0],
				TimeoutSeconds:          pointer.Int32Ptr(10),
			},
		},
	}

	return mutatingWebhookConfiguration, nil
}

// mutatingWebhookConfigurationMutate returns a mutate function that can be used to update the MutatingWebhookConfiguration's spec.
func mutatingWebhookConfigurationMutate(mwc *admissionregistrationv1.MutatingWebhookConfiguration, webhooks []admissionregistrationv1.MutatingWebhook) controllerutil.MutateFn {
	return func() error {
		mwc.Webhooks[0].AdmissionReviewVersions = webhooks[0].AdmissionReviewVersions
		mwc.Webhooks[0].ClientConfig = webhooks[0].ClientConfig
		mwc.Webhooks[0].NamespaceSelector = webhooks[0].NamespaceSelector
		mwc.Webhooks[0].Rules = webhooks[0].Rules
		mwc.Webhooks[0].FailurePolicy = webhooks[0].FailurePolicy
		mwc.Webhooks[0].SideEffects = webhooks[0].SideEffects
		mwc.Webhooks[0].TimeoutSeconds = webhooks[0].TimeoutSeconds
		return nil
	}
}
