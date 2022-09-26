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
	"fmt"
	"os"

	"github.com/fluxninja/aperture/operator/controllers"

	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

// agentMutatingWebhookConfiguration prepares the MutatingWebhookConfiguration object for the Operator to mutate Agents, based on the provided parameter.
func agentMutatingWebhookConfiguration() (*admissionregistrationv1.MutatingWebhookConfiguration, error) {
	certPath := fmt.Sprintf("%s/%s", os.Getenv("APERTURE_OPERATOR_CERT_DIR"), controllers.WebhookClientCertName)
	cert, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}

	mutatingWebhookConfiguration := &admissionregistrationv1.MutatingWebhookConfiguration{
		ObjectMeta: v1.ObjectMeta{
			Name: controllers.AgentMutatingWebhookName,
			Labels: map[string]string{
				"app.kubernetes.io/name":       controllers.AppName,
				"app.kubernetes.io/instance":   controllers.OperatorName,
				"app.kubernetes.io/managed-by": controllers.OperatorName,
				"app.kubernetes.io/component":  controllers.OperatorName,
			},
		},
		Webhooks: []admissionregistrationv1.MutatingWebhook{
			{
				Name: fmt.Sprintf("%s.fluxninja.com", controllers.AgentMutatingWebhookURI),
				ClientConfig: admissionregistrationv1.WebhookClientConfig{
					CABundle: cert,
					Service: &admissionregistrationv1.ServiceReference{
						Name:      os.Getenv("APERTURE_OPERATOR_SERVICE_NAME"),
						Namespace: os.Getenv("APERTURE_OPERATOR_NAMESPACE"),
						Path:      pointer.StringPtr(fmt.Sprintf("/%s", controllers.AgentMutatingWebhookURI)),
						Port:      pointer.Int32(443),
					},
				},
				NamespaceSelector: &v1.LabelSelector{},
				Rules: []admissionregistrationv1.RuleWithOperations{
					{
						Operations: []admissionregistrationv1.OperationType{"CREATE", "UPDATE"},
						Rule: admissionregistrationv1.Rule{
							APIGroups:   []string{"fluxninja.com"},
							APIVersions: []string{controllers.V1Alpha1Version},
							Resources:   []string{"agents"},
							Scope:       &[]admissionregistrationv1.ScopeType{admissionregistrationv1.NamespacedScope}[0],
						},
					},
				},
				AdmissionReviewVersions: []string{controllers.V1Version},
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
	certPath := fmt.Sprintf("%s/%s", os.Getenv("APERTURE_OPERATOR_CERT_DIR"), controllers.WebhookClientCertName)
	cert, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}

	mutatingWebhookConfiguration := &admissionregistrationv1.MutatingWebhookConfiguration{
		ObjectMeta: v1.ObjectMeta{
			Name: controllers.ControllerMutatingWebhookName,
			Labels: map[string]string{
				"app.kubernetes.io/name":       controllers.AppName,
				"app.kubernetes.io/instance":   controllers.OperatorName,
				"app.kubernetes.io/managed-by": controllers.OperatorName,
				"app.kubernetes.io/component":  controllers.OperatorName,
			},
		},
		Webhooks: []admissionregistrationv1.MutatingWebhook{
			{
				Name: fmt.Sprintf("%s.fluxninja.com", controllers.ControllerMutatingWebhookURI),
				ClientConfig: admissionregistrationv1.WebhookClientConfig{
					CABundle: cert,
					Service: &admissionregistrationv1.ServiceReference{
						Name:      os.Getenv("APERTURE_OPERATOR_SERVICE_NAME"),
						Namespace: os.Getenv("APERTURE_OPERATOR_NAMESPACE"),
						Path:      pointer.StringPtr(fmt.Sprintf("/%s", controllers.ControllerMutatingWebhookURI)),
						Port:      pointer.Int32(443),
					},
				},
				NamespaceSelector: &v1.LabelSelector{},
				Rules: []admissionregistrationv1.RuleWithOperations{
					{
						Operations: []admissionregistrationv1.OperationType{"CREATE", "UPDATE"},
						Rule: admissionregistrationv1.Rule{
							APIGroups:   []string{"fluxninja.com"},
							APIVersions: []string{controllers.V1Alpha1Version},
							Resources:   []string{"controllers"},
							Scope:       &[]admissionregistrationv1.ScopeType{admissionregistrationv1.NamespacedScope}[0],
						},
					},
				},
				AdmissionReviewVersions: []string{controllers.V1Version},
				FailurePolicy:           &[]admissionregistrationv1.FailurePolicyType{admissionregistrationv1.Fail}[0],
				SideEffects:             &[]admissionregistrationv1.SideEffectClass{admissionregistrationv1.SideEffectClassNone}[0],
				TimeoutSeconds:          pointer.Int32Ptr(10),
			},
		},
	}

	return mutatingWebhookConfiguration, nil
}
