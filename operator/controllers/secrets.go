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
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/fluxninja/aperture/operator/api/v1alpha1"
)

// secretForAgentAPIKey prepares the Secret object for the ApiKey of Agent.
func secretForAgentAPIKey(instance *v1alpha1.Agent, scheme *runtime.Scheme) (*corev1.Secret, error) {
	spec := &instance.Spec.Secrets.FluxNinjaPlugin

	if spec.Value == "" {
		return nil, fmt.Errorf("value for the ApiKey of Agent cannot be empty")
	}

	secret := &corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name:        secretName(instance.GetName(), "agent", spec),
			Namespace:   instance.GetNamespace(),
			Labels:      commonLabels(instance.Spec.Labels, instance.GetName(), agentServiceName),
			Annotations: instance.Spec.Annotations,
		},
		Data: map[string][]byte{
			secretDataKey(&spec.SecretKeyRef): []byte(spec.Value),
		},
	}

	if scheme != nil {
		if err := ctrl.SetControllerReference(instance, secret, scheme); err != nil {
			return nil, err
		}
	}

	return secret, nil
}

// secretForControllerAPIKey prepares the Secret object for the ApiKey of Agent.
func secretForControllerAPIKey(instance *v1alpha1.Controller, scheme *runtime.Scheme) (*corev1.Secret, error) {
	spec := &instance.Spec.Secrets.FluxNinjaPlugin

	if spec.Value == "" {
		return nil, fmt.Errorf("value for the ApiKey of Controller cannot be empty")
	}

	secret := &corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name:        secretName(instance.GetName(), "controller", spec),
			Namespace:   instance.GetNamespace(),
			Labels:      commonLabels(instance.Spec.Labels, instance.GetName(), controllerServiceName),
			Annotations: instance.Spec.Annotations,
		},
		Data: map[string][]byte{
			secretDataKey(&spec.SecretKeyRef): []byte(spec.Value),
		},
	}

	if err := ctrl.SetControllerReference(instance, secret, scheme); err != nil {
		return nil, err
	}

	return secret, nil
}

// secretForControllerApiKey prepares the Secret object for the ApiKey of Agent.
func secretForControllerCert(instance *v1alpha1.Controller, scheme *runtime.Scheme, serverCert, serverKey *bytes.Buffer) (*corev1.Secret, error) {
	secret := &corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name:        fmt.Sprintf("%s-controller-cert", instance.GetName()),
			Namespace:   instance.GetNamespace(),
			Labels:      commonLabels(instance.Spec.Labels, instance.GetName(), controllerServiceName),
			Annotations: instance.Spec.Annotations,
		},
		Data: map[string][]byte{
			controllerCertName:    serverCert.Bytes(),
			controllerCertKeyName: serverKey.Bytes(),
		},
	}

	if err := ctrl.SetControllerReference(instance, secret, scheme); err != nil {
		return nil, err
	}

	return secret, nil
}

// secretMutate returns a mutate function that can be used to update the Secret's data.
func secretMutate(secret *corev1.Secret, data map[string][]byte) controllerutil.MutateFn {
	return func() error {
		secret.Data = data
		return nil
	}
}

// createSecret calls the Kubernetes API to create the provided Agent Secret resource.
func createSecretForAgent(
	client client.Client, recorder record.EventRecorder, secret *corev1.Secret, ctx context.Context, instance *v1alpha1.Agent) (
	controllerutil.OperationResult, error,
) {
	res, err := controllerutil.CreateOrUpdate(ctx, client, secret, secretMutate(secret, secret.Data))
	if err != nil {
		if errors.IsConflict(err) {
			return createSecretForAgent(client, recorder, secret, ctx, instance)
		}

		msg := fmt.Sprintf("failed to create Secret '%s' for Instance '%s' in Namespace '%s'. Response='%v', Error='%s'",
			secret.GetName(), instance.GetName(), instance.GetNamespace(), res, err.Error())
		if recorder != nil {
			recorder.Event(instance, corev1.EventTypeNormal, "SecretCreationFailed", msg)
		}
		return controllerutil.OperationResultNone, fmt.Errorf(msg)
	}

	if recorder != nil {
		switch res {
		case controllerutil.OperationResultCreated:
			recorder.Eventf(instance, corev1.EventTypeNormal, "SecretCreationSuccessful",
				"Created Secret '%s' in Namespace '%s'", secret.GetName(), secret.GetNamespace())
		case controllerutil.OperationResultUpdated:
			recorder.Eventf(instance, corev1.EventTypeNormal, "SecretUpdationSuccessful",
				"Updated Secret '%s' in Namespace '%s'", secret.GetName(), secret.GetNamespace())
		case controllerutil.OperationResultNone:
		default:
		}
	}

	return res, nil
}

// createSecret calls the Kubernetes API to create the provided Controller Secret resource.
func createSecretForController(
	client client.Client, recorder record.EventRecorder, secret *corev1.Secret, ctx context.Context, instance *v1alpha1.Controller) (
	controllerutil.OperationResult, error,
) {
	res, err := controllerutil.CreateOrUpdate(ctx, client, secret, secretMutate(secret, secret.Data))
	if err != nil {
		if errors.IsConflict(err) {
			return createSecretForController(client, recorder, secret, ctx, instance)
		}

		msg := fmt.Sprintf("failed to create Secret '%s' for Instance '%s' in Namespace '%s'. Response='%v', Error='%s'",
			secret.GetName(), instance.GetName(), instance.GetNamespace(), res, err.Error())
		if recorder != nil {
			recorder.Event(instance, corev1.EventTypeNormal, "SecretCreationFailed", msg)
		}
		return controllerutil.OperationResultNone, fmt.Errorf(msg)
	}

	if recorder != nil {
		switch res {
		case controllerutil.OperationResultCreated:
			recorder.Eventf(instance, corev1.EventTypeNormal, "SecretCreationSuccessful",
				"Created Secret '%s' in Namespace '%s'", secret.GetName(), secret.GetNamespace())
		case controllerutil.OperationResultUpdated:
			recorder.Eventf(instance, corev1.EventTypeNormal, "SecretUpdationSuccessful",
				"Updated Secret '%s' in Namespace '%s'", secret.GetName(), secret.GetNamespace())
		case controllerutil.OperationResultNone:
		default:
		}
	}

	return res, nil
}

// createAgentSecretInNamespace creates the Agent Secret for ApiKey in the given namespace instead of the default one.
func createAgentSecretInNamespace(instance *v1alpha1.Agent, namespace string) (*corev1.Secret, error) {
	copiedInstance := instance.DeepCopy()
	value := copiedInstance.Spec.Secrets.FluxNinjaPlugin.Value
	value = strings.TrimPrefix(value, "enc::")
	value = strings.TrimSuffix(value, "::enc")
	decodedValue, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, err
	}
	copiedInstance.Spec.Secrets.FluxNinjaPlugin.Value = string(decodedValue)
	secret, err := secretForAgentAPIKey(copiedInstance, nil)
	if err != nil {
		return nil, err
	}

	secret.Namespace = namespace
	secret.OwnerReferences = []v1.OwnerReference{}
	secret.Annotations = getAgentAnnotationsWithOwnerRef(copiedInstance)

	return secret, nil
}
