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
	"context"
	"encoding/json"
	"net/http"

	"github.com/fluxninja/aperture/operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// ApertureInjector injects the sidecar container of Aperture Agent in Pods.
type ApertureInjector struct {
	Client   client.Client
	decoder  *admission.Decoder
	Instance *v1alpha1.Aperture
}

// Handle receives incomming requests from MutatingWebhook for newly created Pods and injects Agent container in them.
func (apertureInjector *ApertureInjector) Handle(ctx context.Context, req admission.Request) admission.Response {
	pod := &corev1.Pod{}

	err := apertureInjector.decoder.Decode(req, pod)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	if pod.Annotations != nil && pod.Annotations[sidecarAnnotationKey] == "false" {
		return admission.Allowed("")
	}

	updatedPod := pod.DeepCopy()
	agentPod(apertureInjector.Instance, updatedPod)
	marshaledPod, err := json.Marshal(updatedPod)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	serviceAccount := pod.Spec.ServiceAccountName
	if serviceAccount == "" {
		serviceAccount = "default"
	}

	subject := rbacv1.Subject{
		Kind:      "ServiceAccount",
		Name:      serviceAccount,
		Namespace: req.Namespace,
	}

	if err := updateClusterRoleBinding(apertureInjector.Client, subject, ctx, apertureInjector.Instance.GetNamespace()); err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)
}

// InjectDecoder injects the decoder.
func (apertureInjector *ApertureInjector) InjectDecoder(d *admission.Decoder) error {
	apertureInjector.decoder = d
	return nil
}
