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
	"fmt"
	"net/http"
	"strings"

	"github.com/clarketm/json"
	"github.com/fluxninja/aperture/operator/api/v1alpha1"
	"github.com/fluxninja/aperture/pkg/config"
	policy "github.com/fluxninja/aperture/pkg/policies/controlplane"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// ControllerHooks injects the default spec of Aperture Controller in CR.
type ControllerHooks struct {
	decoder *admission.Decoder
}

// Handle receives incoming requests from MutatingWebhook for newly created Controllers, set defaults and validates them.
func (controllerHooks *ControllerHooks) Handle(ctx context.Context, req admission.Request) admission.Response {
	if strings.ToLower(req.AdmissionRequest.Kind.Kind) == "policy" {
		return controllerHooks.policyDefaults(ctx, req)
	}

	return controllerHooks.controllerDefaults(ctx, req)
}

// controllerDefaults validates and sets defaults for Aperture Controller Custom Resource.
func (controllerHooks *ControllerHooks) controllerDefaults(ctx context.Context, req admission.Request) admission.Response {
	controller := &v1alpha1.Controller{}

	err := config.UnmarshalYAML([]byte(req.Object.Raw), controller)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	if controller.Spec.Secrets.FluxNinjaPlugin.Create && controller.Spec.Secrets.FluxNinjaPlugin.Value == "" {
		return admission.Denied("The value for 'spec.secrets.fluxNinjaPlugin.value' can not be empty when 'spec.secrets.fluxNinjaPlugin.create' is set to true")
	}

	if controller.ObjectMeta.Annotations == nil {
		controller.ObjectMeta.Annotations = map[string]string{}
	}

	controller.ObjectMeta.Annotations[defaulterAnnotationKey] = "true"
	updatedController, err := json.Marshal(controller)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, updatedController)
}

// controllerDefaults validates and sets defaults for Aperture Controller Custom Resource.
func (controllerHooks *ControllerHooks) policyDefaults(ctx context.Context, req admission.Request) admission.Response {
	instance := &v1alpha1.Policy{}

	err := config.UnmarshalYAML(req.Object.Raw, instance)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	_, valid, msg, err := policy.ValidateAndCompile(ctx, "", instance.Spec.Raw)
	if err != nil || !valid {
		if err == nil {
			return admission.Errored(http.StatusBadRequest, fmt.Errorf(msg))
		}
		return admission.Errored(http.StatusBadRequest, err)
	}

	if instance.ObjectMeta.Annotations == nil {
		instance.ObjectMeta.Annotations = map[string]string{}
	}

	instance.ObjectMeta.Annotations[defaulterAnnotationKey] = "true"
	updatedPolicy, err := json.Marshal(instance)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, updatedPolicy)
}

// InjectDecoder injects the decoder.
func (controllerHooks *ControllerHooks) InjectDecoder(d *admission.Decoder) error {
	controllerHooks.decoder = d
	return nil
}
