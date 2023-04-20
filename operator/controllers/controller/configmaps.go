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
	"bytes"
	"context"
	_ "embed"
	"fmt"

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

	controllerv1alpha1 "github.com/fluxninja/aperture/operator/api/controller/v1alpha1"
)

// configMapForAgentConfig prepares the ConfigMap object for the Controller.
func configMapForControllerConfig(instance *controllerv1alpha1.Controller, scheme *runtime.Scheme) (*corev1.ConfigMap, error) {
	jsonConfig, err := json.Marshal(instance.Spec.ConfigSpec)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Controller config to JSON. Error: '%s'", err.Error())
	}

	config, err := yaml.JSONToYAML(jsonConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Controller config to YAML. Error: '%s'", err.Error())
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{
			Name:        controllers.ControllerServiceName,
			Namespace:   instance.GetNamespace(),
			Labels:      controllers.CommonLabels(instance.Spec.Labels, instance.GetName(), controllers.ControllerServiceName),
			Annotations: instance.Spec.Annotations,
		},
		Data: map[string]string{
			"aperture-controller.yaml": string(config),
		},
	}

	if err := ctrl.SetControllerReference(instance, cm, scheme); err != nil {
		return nil, err
	}

	return cm, nil
}

// configMapForControllerClientCert prepares the ConfigMap object for the Aperture Controller client Certificate.
func configMapForControllerClientCert(instance *controllerv1alpha1.Controller, scheme *runtime.Scheme, clientCert *bytes.Buffer) (*corev1.ConfigMap, error) {
	cm := &corev1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{
			Name:        fmt.Sprintf("%s-controller-client-cert", instance.GetName()),
			Namespace:   instance.GetNamespace(),
			Labels:      controllers.CommonLabels(instance.Spec.Labels, instance.GetName(), controllers.AppName),
			Annotations: instance.Spec.Annotations,
		},
		Data: map[string]string{
			controllers.ControllerClientCertKey: clientCert.String(),
		},
	}

	return cm, nil
}

// createConfigMap calls the Kubernetes API to create the provided Controller ConfigMap resource.
func createConfigMapForController(
	client client.Client, recorder record.EventRecorder, configMap *corev1.ConfigMap, ctx context.Context, instance *controllerv1alpha1.Controller) (
	controllerutil.OperationResult, error,
) {
	res, err := controllerutil.CreateOrUpdate(ctx, client, configMap, controllers.ConfigMapMutate(configMap, configMap.Data))
	if err != nil {
		if errors.IsConflict(err) {
			return createConfigMapForController(client, recorder, configMap, ctx, instance)
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
