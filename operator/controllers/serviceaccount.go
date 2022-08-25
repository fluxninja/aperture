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

	"github.com/imdario/mergo"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/fluxninja/aperture/operator/api/v1alpha1"
)

// serviceAccountForAgent prepares the ServiceAccount object for the Agent based on the provided parameter.
func serviceAccountForAgent(instance *v1alpha1.Aperture, scheme *runtime.Scheme) (*corev1.ServiceAccount, error) {
	saSpec := instance.Spec.Agent.ServiceAccountSpec

	annotations := instance.Spec.Annotations
	if annotations == nil {
		annotations = saSpec.Annotations
	} else if saSpec.Annotations != nil {
		if err := mergo.Map(&annotations, saSpec.Annotations, mergo.WithOverride); err != nil {
			return nil, fmt.Errorf(fmt.Sprintf("failed to merge the annotations for ServiceAccount of Agent. error: %s.", err.Error()))
		}
	}

	sa := &corev1.ServiceAccount{
		ObjectMeta: v1.ObjectMeta{
			Name:        agentServiceName,
			Namespace:   instance.GetNamespace(),
			Labels:      commonLabels(instance, agentServiceName),
			Annotations: annotations,
		},
		AutomountServiceAccountToken: &saSpec.AutomountServiceAccountToken,
	}

	if err := ctrl.SetControllerReference(instance, sa, scheme); err != nil {
		return nil, err
	}

	return sa, nil
}

// serviceAccountForController prepares the ServiceAccount object for the Controller based on the provided parameter.
func serviceAccountForController(instance *v1alpha1.Aperture, scheme *runtime.Scheme) (*corev1.ServiceAccount, error) {
	saSpec := instance.Spec.Controller.ServiceAccountSpec

	annotations := instance.Spec.Annotations
	if annotations == nil {
		annotations = saSpec.Annotations
	} else if saSpec.Annotations != nil {
		if err := mergo.Map(&annotations, saSpec.Annotations, mergo.WithOverride); err != nil {
			return nil, fmt.Errorf(fmt.Sprintf("failed to merge the annotations for ServiceAccount of Controller. error: %s.", err.Error()))
		}
	}

	sa := &corev1.ServiceAccount{
		ObjectMeta: v1.ObjectMeta{
			Name:        controllerServiceName,
			Namespace:   instance.GetNamespace(),
			Labels:      commonLabels(instance, controllerServiceName),
			Annotations: annotations,
		},
		AutomountServiceAccountToken: &saSpec.AutomountServiceAccountToken,
	}

	if err := ctrl.SetControllerReference(instance, sa, scheme); err != nil {
		return nil, err
	}

	return sa, nil
}

// serviceAccountMutate returns a mutate function that can be used to update the ClusterRole's spec.
func serviceAccountMutate(sa *corev1.ServiceAccount, automountServiceAccountToken *bool) controllerutil.MutateFn {
	return func() error {
		sa.AutomountServiceAccountToken = automountServiceAccountToken
		return nil
	}
}

// createServiceAccount calls the Kubernetes API to create the provided ServiceAccount resource.
func (r *ApertureReconciler) createServiceAccount(sa *corev1.ServiceAccount, ctx context.Context, instance *v1alpha1.Aperture) error {
	res, err := controllerutil.CreateOrUpdate(ctx, r.Client, sa, serviceAccountMutate(sa, sa.AutomountServiceAccountToken))
	if err != nil {
		if errors.IsConflict(err) {
			return r.createServiceAccount(sa, ctx, instance)
		}

		msg := fmt.Sprintf("failed to create ServiceAccount '%s' for Instance '%s' in Namespace '%s'. Response='%v', Error='%s'",
			sa.GetName(), instance.GetName(), instance.GetNamespace(), res, err.Error())
		r.Recorder.Event(instance, corev1.EventTypeNormal, "ServiceAccountCreationFailed", msg)
		return fmt.Errorf(msg)
	}

	switch res {
	case controllerutil.OperationResultCreated:
		r.Recorder.Eventf(instance, corev1.EventTypeNormal, "ServiceAccountCreationSuccessful", "Created ServiceAccount '%s'", sa.GetName())
	case controllerutil.OperationResultUpdated:
		r.Recorder.Eventf(instance, corev1.EventTypeNormal, "ServiceAccountUpdationSuccessful", "Updated ServiceAccount '%s'", sa.GetName())
	case controllerutil.OperationResultNone:
	default:
	}

	return nil
}
