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
	"fmt"
	"os"

	"github.com/fluxninja/aperture/v2/operator/controllers"

	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	agentv1alpha1 "github.com/fluxninja/aperture/v2/operator/api/agent/v1alpha1"
)

// podMutatingWebhookConfiguration prepares the MutatingWebhookConfiguration object for the Operator to mutate Pods, based on the provided parameter.
func podMutatingWebhookConfiguration(instance *agentv1alpha1.Agent) (*admissionregistrationv1.MutatingWebhookConfiguration, error) {
	certPath := fmt.Sprintf("%s/%s", os.Getenv("APERTURE_OPERATOR_CERT_DIR"), controllers.WebhookClientCertName)
	cert, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}

	mutatingWebhookConfiguration := &admissionregistrationv1.MutatingWebhookConfiguration{
		ObjectMeta: v1.ObjectMeta{
			Name:        controllers.PodMutatingWebhookName,
			Labels:      controllers.CommonLabels(instance.Spec.Labels, instance.GetName(), controllers.OperatorName),
			Annotations: controllers.AgentAnnotationsWithOwnerRef(instance),
		},
		Webhooks: []admissionregistrationv1.MutatingWebhook{
			{
				Name: controllers.SidecarKey,
				ClientConfig: admissionregistrationv1.WebhookClientConfig{
					CABundle: cert,
					Service: &admissionregistrationv1.ServiceReference{
						Name:      os.Getenv("APERTURE_OPERATOR_SERVICE_NAME"),
						Namespace: os.Getenv("APERTURE_OPERATOR_NAMESPACE"),
						Path:      pointer.String(controllers.MutatingWebhookURI),
						Port:      pointer.Int32(443),
					},
				},
				NamespaceSelector: &v1.LabelSelector{
					MatchExpressions: []v1.LabelSelectorRequirement{
						{
							Key:      controllers.SidecarLabelKey,
							Operator: v1.LabelSelectorOpIn,
							Values:   []string{controllers.Enabled},
						},
					},
				},
				Rules: []admissionregistrationv1.RuleWithOperations{
					{
						Operations: []admissionregistrationv1.OperationType{"CREATE"},
						Rule: admissionregistrationv1.Rule{
							APIGroups:   []string{""},
							APIVersions: []string{controllers.V1Version},
							Resources:   []string{"pods"},
							Scope:       &[]admissionregistrationv1.ScopeType{admissionregistrationv1.NamespacedScope}[0],
						},
					},
				},
				AdmissionReviewVersions: []string{controllers.V1Version},
				FailurePolicy:           &[]admissionregistrationv1.FailurePolicyType{admissionregistrationv1.Fail}[0],
				SideEffects:             &[]admissionregistrationv1.SideEffectClass{admissionregistrationv1.SideEffectClassNone}[0],
				TimeoutSeconds:          pointer.Int32(10),
			},
		},
	}

	return mutatingWebhookConfiguration, nil
}
