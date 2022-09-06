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

	"github.com/go-logr/logr"
	"github.com/imdario/mergo"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/fluxninja/aperture/operator/api/v1alpha1"
)

// serviceForAgent prepares an object of Service for Agent based on the provided parameters.
func serviceForAgent(instance *v1alpha1.Agent, log logr.Logger, scheme *runtime.Scheme) (*corev1.Service, error) {
	spec := instance.Spec
	agentServiceSpec := spec.Service
	annotations := spec.Annotations
	if agentServiceSpec.Annotations != nil {
		if annotations == nil {
			annotations = map[string]string{}
		}
		if err := mergo.Map(&annotations, agentServiceSpec.Annotations, mergo.WithOverride); err != nil {
			log.Info(fmt.Sprintf("failed to merge the Agent Service annotations. error: %s.", err.Error()))
		}
	}

	serverPort, err := getPort(spec.ConfigSpec.Server.Addr)
	if err != nil {
		return nil, err
	}

	otelGRPCPort, err := getPort(spec.ConfigSpec.Otel.GRPCAddr)
	if err != nil {
		return nil, err
	}

	otelHTTPPort, err := getPort(spec.ConfigSpec.Otel.HTTPAddr)
	if err != nil {
		return nil, err
	}

	distCachePort, err := getPort(spec.ConfigSpec.DistCache.BindAddr)
	if err != nil {
		return nil, err
	}

	memberListPort, err := getPort(spec.ConfigSpec.DistCache.MemberlistBindAddr)
	if err != nil {
		return nil, err
	}

	svc := &corev1.Service{
		ObjectMeta: v1.ObjectMeta{
			Name:        agentServiceName,
			Namespace:   instance.GetNamespace(),
			Labels:      commonLabels(spec.Labels, instance.GetName(), agentServiceName),
			Annotations: annotations,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       server,
					Protocol:   corev1.Protocol(tcp),
					Port:       int32(serverPort),
					TargetPort: intstr.FromString(server),
				},
				{
					Name:       grpcOtel,
					Protocol:   corev1.Protocol(tcp),
					Port:       int32(otelGRPCPort),
					TargetPort: intstr.FromString(grpcOtel),
				},
				{
					Name:       httpOtel,
					Protocol:   corev1.Protocol(tcp),
					Port:       int32(otelHTTPPort),
					TargetPort: intstr.FromString(httpOtel),
				},
				{
					Name:       distCache,
					Protocol:   corev1.Protocol(tcp),
					Port:       int32(distCachePort),
					TargetPort: intstr.FromString(distCache),
				},
				{
					Name:       memberList,
					Protocol:   corev1.Protocol(tcp),
					Port:       int32(memberListPort),
					TargetPort: intstr.FromString(memberList),
				},
			},
			InternalTrafficPolicy: &[]corev1.ServiceInternalTrafficPolicyType{corev1.ServiceInternalTrafficPolicyLocal}[0],
			Selector:              selectorLabels(instance.GetName(), agentServiceName),
		},
	}

	if err := ctrl.SetControllerReference(instance, svc, scheme); err != nil {
		return nil, err
	}

	return svc, nil
}

// serviceForController prepares an object of Service for Controller based on the provided parameters.
func serviceForController(instance *v1alpha1.Controller, log logr.Logger, scheme *runtime.Scheme) (*corev1.Service, error) {
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

	serverPort, err := getPort(spec.ConfigSpec.Server.Addr)
	if err != nil {
		return nil, err
	}

	otelGRPCPort, err := getPort(spec.ConfigSpec.Otel.GRPCAddr)
	if err != nil {
		return nil, err
	}

	otelHTTPPort, err := getPort(spec.ConfigSpec.Otel.HTTPAddr)
	if err != nil {
		return nil, err
	}

	svc := &corev1.Service{
		ObjectMeta: v1.ObjectMeta{
			Name:        controllerServiceName,
			Namespace:   instance.GetNamespace(),
			Labels:      commonLabels(spec.Labels, instance.GetName(), controllerServiceName),
			Annotations: spec.Annotations,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       server,
					Protocol:   corev1.Protocol(tcp),
					Port:       int32(serverPort),
					TargetPort: intstr.FromString(server),
				},
				{
					Name:       grpcOtel,
					Protocol:   corev1.Protocol(tcp),
					Port:       int32(otelGRPCPort),
					TargetPort: intstr.FromString(grpcOtel),
				},
				{
					Name:       httpOtel,
					Protocol:   corev1.Protocol(tcp),
					Port:       int32(otelHTTPPort),
					TargetPort: intstr.FromString(httpOtel),
				},
			},
			Selector: selectorLabels(instance.GetName(), controllerServiceName),
		},
	}

	if err := ctrl.SetControllerReference(instance, svc, scheme); err != nil {
		return nil, err
	}

	return svc, nil
}

// serviceMutate returns a mutate function that can be used to update the Service's spec.
func serviceMutate(svc *corev1.Service, spec corev1.ServiceSpec) controllerutil.MutateFn {
	return func() error {
		svc.Spec.Ports = spec.Ports
		svc.Spec.Selector = spec.Selector
		return nil
	}
}

// createService calls the Kubernetes API to create the provided Agent Service resource.
func (r *AgentReconciler) createService(service *corev1.Service, ctx context.Context, instance *v1alpha1.Agent) error {
	res, err := controllerutil.CreateOrUpdate(ctx, r.Client, service, serviceMutate(service, service.Spec))
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

// createService calls the Kubernetes API to create the provided Controller Service resource.
func (r *ControllerReconciler) createService(service *corev1.Service, ctx context.Context, instance *v1alpha1.Controller) error {
	res, err := controllerutil.CreateOrUpdate(ctx, r.Client, service, serviceMutate(service, service.Spec))
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
