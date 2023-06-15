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
	"fmt"

	"github.com/fluxninja/aperture/v2/operator/controllers"

	"github.com/go-logr/logr"
	"github.com/imdario/mergo"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	controllerv1alpha1 "github.com/fluxninja/aperture/v2/operator/api/controller/v1alpha1"
)

// serviceForController prepares an object of Service for Controller based on the provided parameters.
func serviceForController(instance *controllerv1alpha1.Controller, log logr.Logger, scheme *runtime.Scheme) (*corev1.Service, error) {
	spec := instance.Spec
	controllerServiceSpec := spec.Service
	annotations := spec.Annotations
	if controllerServiceSpec.Annotations != nil {
		if annotations == nil {
			annotations = map[string]string{}
		}
		if err := mergo.Map(&annotations, controllerServiceSpec.Annotations, mergo.WithOverride); err != nil {
			log.Info(fmt.Sprintf("failed to merge the Controller Service annotations. error: %s.", err.Error()))
		}
	}

	serverPort, err := controllers.GetPort(spec.ConfigSpec.Server.Listener.Addr)
	if err != nil {
		return nil, err
	}

	svc := &corev1.Service{
		ObjectMeta: v1.ObjectMeta{
			Name:        controllers.ControllerServiceName,
			Namespace:   instance.GetNamespace(),
			Labels:      controllers.CommonLabels(spec.Labels, instance.GetName(), controllers.ControllerServiceName),
			Annotations: spec.Annotations,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       controllers.Server,
					Protocol:   corev1.Protocol(controllers.TCP),
					Port:       int32(serverPort),
					TargetPort: intstr.FromString(controllers.Server),
				},
			},
			Selector: controllers.SelectorLabels(instance.GetName(), controllers.ControllerServiceName),
		},
	}

	if err := ctrl.SetControllerReference(instance, svc, scheme); err != nil {
		return nil, err
	}

	return svc, nil
}

// createService calls the Kubernetes API to create the provided Controller Service resource.
func (r *ControllerReconciler) createService(service *corev1.Service, ctx context.Context, instance *controllerv1alpha1.Controller) error {
	res, err := controllerutil.CreateOrUpdate(ctx, r.Client, service, controllers.ServiceMutate(service, service.Spec))
	if err != nil {
		if errors.IsConflict(err) {
			return r.createService(service, ctx, instance)
		}

		msg := fmt.Sprintf("failed to create Service '%s' for Instance '%s' in Namespace '%s'. Response='%v', Error='%s'",
			service.GetName(), instance.GetName(), instance.GetNamespace(), res, err.Error())
		r.Recorder.Event(instance, corev1.EventTypeNormal, "ServiceCreationFailed", msg)
		return fmt.Errorf(msg)
	}

	switch res {
	case controllerutil.OperationResultCreated:
		r.Recorder.Eventf(instance, corev1.EventTypeNormal, "ServiceCreationSuccessful", "Created Service '%s'", service.GetName())
	case controllerutil.OperationResultUpdated:
		r.Recorder.Eventf(instance, corev1.EventTypeNormal, "ServiceUpdationSuccessful", "Updated Service '%s'", service.GetName())
	case controllerutil.OperationResultNone:
	default:
	}

	return nil
}
