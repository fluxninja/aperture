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
	"context"
	"net/http"

	"github.com/fluxninja/aperture/v2/operator/controllers"

	"github.com/clarketm/json"
	controllerv1alpha1 "github.com/fluxninja/aperture/v2/operator/api/controller/v1alpha1"
	"github.com/fluxninja/aperture/v2/pkg/config"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// ControllerHooks injects the default spec of Aperture Controller in CR.
type ControllerHooks struct {
	decoder *admission.Decoder
}

// Handle receives incoming requests from MutatingWebhook for newly created Controllers, set defaults and validates them.
func (controllerHooks *ControllerHooks) Handle(ctx context.Context, req admission.Request) admission.Response {
	controller := &controllerv1alpha1.Controller{}

	err := config.UnmarshalYAML([]byte(req.Object.Raw), controller)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	if controller.Spec.Secrets.FluxNinjaExtension.Create && controller.Spec.Secrets.FluxNinjaExtension.Value == "" {
		return admission.Denied("The value for 'spec.secrets.fluxNinjaExtension.value' can not be empty when 'spec.secrets.fluxNinjaExtension.create' is set to true")
	}

	if controller.ObjectMeta.Annotations == nil {
		controller.ObjectMeta.Annotations = map[string]string{}
	}

	controller.ObjectMeta.Annotations[controllers.DefaulterAnnotationKey] = "true"
	updatedController, err := json.Marshal(controller)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, updatedController)
}

// InjectDecoder injects the decoder.
func (controllerHooks *ControllerHooks) InjectDecoder(d *admission.Decoder) error {
	controllerHooks.decoder = d
	return nil
}
