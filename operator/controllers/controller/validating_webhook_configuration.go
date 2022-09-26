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

package controller

import (
	"github.com/fluxninja/aperture/operator/controllers"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	controllerv1alpha1 "github.com/fluxninja/aperture/operator/api/controller/v1alpha1"
)

// validatingWebhookConfiguration prepares the ValidatingWebhookConfiguration object for the Controller, based on the provided parameter.
func validatingWebhookConfiguration(instance *controllerv1alpha1.Controller, cert []byte) *admissionregistrationv1.ValidatingWebhookConfiguration {
	serverPort, _ := controllers.GetPort(instance.Spec.ConfigSpec.Server.Addr)

	validatingWebhookConfiguration := &admissionregistrationv1.ValidatingWebhookConfiguration{
		ObjectMeta: v1.ObjectMeta{
			Name:        controllers.ControllerServiceName,
			Labels:      controllers.CommonLabels(instance.Spec.Labels, instance.GetName(), controllers.ControllerServiceName),
			Annotations: controllers.ControllerAnnotationsWithOwnerRef(instance),
		},
		Webhooks: []admissionregistrationv1.ValidatingWebhook{
			{
				Name: controllers.PolicyValidatingWebhookName,
				ClientConfig: admissionregistrationv1.WebhookClientConfig{
					Service: &admissionregistrationv1.ServiceReference{
						Name:      controllers.ControllerServiceName,
						Namespace: instance.GetNamespace(),
						Path:      pointer.StringPtr(controllers.PolicyValidatingWebhookURI),
						Port:      pointer.Int32(serverPort),
					},
					CABundle: cert,
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
							APIGroups:   []string{"fluxninja.com"},
							APIVersions: []string{controllers.V1Alpha1Version},
							Resources:   []string{"policies"},
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

	return validatingWebhookConfiguration
}
