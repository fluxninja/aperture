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
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/fluxninja/aperture/v2/operator/controllers"

	"github.com/clarketm/json"
	"github.com/fluxninja/aperture/v2/operator/api/agent/v1alpha1"
	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/utils"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// AgentHooks injects the default spec of Aperture Agent in CR.
type AgentHooks struct{}

// Handle receives incoming requests from MutatingWebhook for newly created Agents, set defaults and validates them.
func (agentHooks *AgentHooks) Handle(ctx context.Context, req admission.Request) admission.Response {
	newAgent := &v1alpha1.Agent{}
	oldAgent := &v1alpha1.Agent{}

	err := config.UnmarshalYAML([]byte(req.Object.Raw), newAgent)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	if !newAgent.Spec.ConfigSpec.FluxNinja.EnableCloudController {
		if len(newAgent.Spec.ConfigSpec.Etcd.Endpoints) == 0 {
			return admission.Denied("At least one etcd endpoint must be provided under spec.config.etcd.endpoints.")
		}

		if newAgent.Spec.ConfigSpec.Prometheus.Address == "" {
			return admission.Denied("The address for Prometheus must be provided under spec.config.prometheus.address.")
		}
	} else if newAgent.Spec.ConfigSpec.FluxNinja.Endpoint == "" {
		return admission.Denied("The endpoint for Flux Ninja must be provided under spec.config.fluxNinja.endpoint.")
	}

	if (newAgent.Spec.Image.Digest == "" && newAgent.Spec.Image.Tag == "") || (newAgent.Spec.Image.Digest != "" && newAgent.Spec.Image.Tag != "") {
		return admission.Denied("Either 'spec.image.digest' or 'spec.image.tag' should be provided.")
	}

	if newAgent.Spec.Secrets.FluxNinjaExtension.Create && newAgent.Spec.Secrets.FluxNinjaExtension.Value == "" {
		return admission.Denied("The value for 'spec.secrets.fluxNinjaExtension.value' can not be empty when 'spec.secrets.fluxNinjaExtension.create' is set to true")
	}

	if newAgent.Spec.ConfigSpec.FluxNinja.InstallationMode != utils.InstallationModeCloudAgent && newAgent.Spec.ConfigSpec.AgentInfo.AgentGroup == utils.ApertureCloudAgentGroup {
		return admission.Denied(
			fmt.Sprintf("'%s' is a reserved group name for FluxNinja Cloud Agents. Please use a different agent group name", utils.ApertureCloudAgentGroup))
	}

	if newAgent.ObjectMeta.Annotations == nil {
		newAgent.ObjectMeta.Annotations = map[string]string{}
	}

	newAgent.ObjectMeta.Annotations[controllers.DefaulterAnnotationKey] = "true"

	if len(req.OldObject.Raw) != 0 {
		err = config.UnmarshalYAML([]byte(req.OldObject.Raw), oldAgent)
		if err != nil {
			return admission.Errored(http.StatusBadRequest, err)
		}

		if oldAgent.Spec.Sidecar.Enabled != newAgent.Spec.Sidecar.Enabled {
			newAgent.ObjectMeta.Annotations[controllers.AgentModeChangeAnnotationKey] = "true"
		}
	}

	if newAgent.Spec.Sidecar.Enabled {
		newAgent.Spec.ConfigSpec.FluxNinja.InstallationMode = utils.InstallationModeKubernetesSidecar
	} else if strings.ToLower(newAgent.Spec.DeploymentConfigSpec.Type) != "deployment" {
		newAgent.Spec.ConfigSpec.FluxNinja.InstallationMode = utils.InstallationModeKubernetesDaemonSet
	}

	updatedAgent, err := json.Marshal(newAgent)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, updatedAgent)
}
