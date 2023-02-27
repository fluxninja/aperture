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
	_ "embed"
	"fmt"
	"strings"

	"github.com/fluxninja/aperture/operator/controllers"

	"github.com/clarketm/json"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/yaml"

	agentv1alpha1 "github.com/fluxninja/aperture/operator/api/agent/v1alpha1"
)

// configMapForAgentConfig prepares the ConfigMap object for the Agent.
func configMapForAgentConfig(
	ctx context.Context,
	client_ client.Client,
	instance *agentv1alpha1.Agent,
	scheme *runtime.Scheme,
) (*corev1.ConfigMap, error) {
	// HACK: auto-inject controller's cert to agent if we detect it connects to local controller.
	// FIXME: Mounting optional configmap with controller CA would be a better
	// solution, but right now we don't have a configmap with CA, just a secret
	// that contains also the key.
	// Also, we should figure out the way so that both CA file from config and
	// controller's CA can be used simultaneously.
	// FIXME: Even better solution would be for local controller to be signed
	// using https://kubernetes.io/docs/tasks/tls/managing-tls-in-a-cluster/
	var localControllerCert []byte
	for _, endpoint := range instance.Spec.ConfigSpec.AgentFunctions.Endpoints {
		controllerNS := localControllerNamespaceFromEndpoint(endpoint)
		if controllerNS == "" {
			continue
		}

		caFile := &instance.Spec.ConfigSpec.AgentFunctions.ClientConfig.GRPCClient.ClientTLSConfig.CAFile
		if *caFile != "" {
			continue
		}

		if client_ == nil {
			continue
		}

		var secrets corev1.SecretList
		_ = client_.List(ctx, &secrets, &client.ListOptions{Namespace: controllerNS})

		for _, secret := range secrets.Items {
			if !strings.HasSuffix(secret.Name, "-controller-cert") {
				continue
			}
			localControllerCert = append(localControllerCert, secret.Data[controllers.ControllerCertName]...)
			if localControllerCert != nil {
				*caFile = "/etc/aperture/aperture-agent/config/controller-ca.pem"
			}
		}
	}

	jsonConfig, err := json.Marshal(instance.Spec.ConfigSpec)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Agent config to JSON. Error: '%s'", err.Error())
	}

	config, err := yaml.JSONToYAML(jsonConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Agent config to YAML. Error: '%s'", err.Error())
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{
			Name:        controllers.AgentServiceName,
			Namespace:   instance.GetNamespace(),
			Labels:      controllers.CommonLabels(instance.Spec.Labels, instance.GetName(), controllers.AgentServiceName),
			Annotations: instance.Spec.Annotations,
		},
		Data: map[string]string{
			"aperture-agent.yaml": string(config),
		},
	}

	if localControllerCert != nil {
		cm.Data["controller-ca.pem"] = string(localControllerCert)
	}

	if scheme != nil {
		if err := ctrl.SetControllerReference(instance, cm, scheme); err != nil {
			return nil, err
		}
	}

	return cm, nil
}

// CreateConfigMapForAgent calls the Kubernetes API to create the provided Agent ConfigMap resource.
func CreateConfigMapForAgent(
	client client.Client, recorder record.EventRecorder, configMap *corev1.ConfigMap, ctx context.Context, instance *agentv1alpha1.Agent) (
	controllerutil.OperationResult, error,
) {
	res, err := controllerutil.CreateOrUpdate(ctx, client, configMap, controllers.ConfigMapMutate(configMap, configMap.Data))
	if err != nil {
		if errors.IsConflict(err) {
			return CreateConfigMapForAgent(client, recorder, configMap, ctx, instance)
		}

		msg := fmt.Sprintf("failed to create ConfigMap '%s' for Instance '%s' in Namespace '%s'. Response='%v', Error='%s'",
			configMap.GetName(), instance.GetName(), instance.GetNamespace(), res, err.Error())
		if recorder != nil {
			recorder.Event(instance, corev1.EventTypeNormal, "ConfigMapCreationFailed", msg)
		}
		return controllerutil.OperationResultNone, fmt.Errorf(msg)
	}

	if recorder != nil {
		switch res {
		case controllerutil.OperationResultCreated:
			recorder.Eventf(instance, corev1.EventTypeNormal, "ConfigMapCreationSuccessful",
				"Created ConfigMap '%s' in Namespace '%s'", configMap.GetName(), configMap.GetNamespace())
		case controllerutil.OperationResultUpdated:
			recorder.Eventf(instance, corev1.EventTypeNormal, "ConfigMapUpdationSuccessful",
				"Updated ConfigMap '%s' in Namespace '%s'", configMap.GetName(), configMap.GetNamespace())
		case controllerutil.OperationResultNone:
		default:
		}
	}

	return res, nil
}

// CreateAgentConfigMapInNamespace creates the Agent ConfigMap in the given namespace instead of the default one.
func CreateAgentConfigMapInNamespace(
	ctx context.Context,
	client client.Client,
	instance *agentv1alpha1.Agent,
	namespace string,
) *corev1.ConfigMap {
	configMap, _ := configMapForAgentConfig(ctx, client, instance, nil)
	configMap.Namespace = namespace
	configMap.Annotations = controllers.AgentAnnotationsWithOwnerRef(instance)

	return configMap
}

func localControllerNamespaceFromEndpoint(endpoint string) string {
	addr, port, ok := strings.Cut(endpoint, ":")
	if !ok {
		return ""
	}

	if port != "8080" {
		return ""
	}

	subdomains := strings.Split(addr, ".")
	if len(subdomains) < 2 {
		return ""
	}

	if subdomains[0] != "aperture-controller" {
		return ""
	}

	tail := strings.Join(subdomains[2:], ".") + "."
	if !strings.HasPrefix("svc.cluster.local.", tail) { //nolint:gocritic
		return ""
	}

	return subdomains[1]
}
