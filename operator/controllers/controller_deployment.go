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
	"fmt"

	"github.com/go-logr/logr"
	"github.com/imdario/mergo"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/fluxninja/aperture/operator/api/v1alpha1"
)

// deploymentForAPIService prepares the Deployment object for the Controller.
func deploymentForController(instance *v1alpha1.Aperture, log logr.Logger, scheme *runtime.Scheme) (*appsv1.Deployment, error) {
	controllerSpec := instance.Spec.Controller

	podLabels := commonLabels(instance, controllerServiceName)
	if controllerSpec.PodLabels != nil {
		if err := mergo.Map(&podLabels, controllerSpec.PodLabels, mergo.WithOverride); err != nil {
			log.Info(fmt.Sprintf("failed to merge the Pod labels for Deployment. error: %s.", err.Error()))
		}
	}

	annotations := controllerSpec.PodAnnotations
	if annotations == nil {
		annotations = map[string]string{}
	}
	annotations[sidecarAnnotationKey] = "false"

	livenessProbe, readinessProbe := containerProbes(instance.Spec.Controller.CommonSpec)

	dep := &appsv1.Deployment{
		ObjectMeta: v1.ObjectMeta{
			Name:        controllerServiceName,
			Namespace:   instance.GetNamespace(),
			Labels:      commonLabels(instance, controllerServiceName),
			Annotations: instance.Spec.Annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &v1.LabelSelector{
				MatchLabels: selectorLabels(instance.GetName(), controllerServiceName),
			},
			// Controller does not support running multiple Replicas, hence it is hard-coded to 'Recreate'
			Strategy: appsv1.DeploymentStrategy{
				Type: "Recreate",
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: v1.ObjectMeta{
					Labels:      podLabels,
					Annotations: annotations,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName:            controllerServiceName,
					HostAliases:                   controllerSpec.HostAliases,
					ImagePullSecrets:              imagePullSecrets(instance.Spec.ImagePullSecrets, controllerSpec.Image),
					NodeSelector:                  controllerSpec.NodeSelector,
					Affinity:                      controllerSpec.Affinity,
					Tolerations:                   controllerSpec.Tolerations,
					SecurityContext:               podSecurityContext(controllerSpec.PodSecurityContext),
					TerminationGracePeriodSeconds: controllerSpec.TerminationGracePeriodSeconds,
					InitContainers:                controllerSpec.InitContainers,
					Containers: []corev1.Container{
						{
							Name:            controllerServiceName,
							Image:           imageString(instance.Spec.ImageRegistry, controllerSpec.Image),
							ImagePullPolicy: corev1.PullPolicy(controllerSpec.Image.PullPolicy),
							SecurityContext: containerSecurityContext(controllerSpec.ContainerSecurityContext),
							Command:         controllerSpec.Command,
							Args:            controllerSpec.Args,
							Env:             controllerEnv(instance),
							EnvFrom:         containerEnvFrom(controllerSpec.CommonSpec),
							Resources:       controllerSpec.Resources,
							Ports: []corev1.ContainerPort{
								{
									Name:          "grpc",
									ContainerPort: controllerSpec.ServerPort,
									Protocol:      "TCP",
								},
								{
									Name:          "webhooks-port",
									ContainerPort: 8086,
									Protocol:      "TCP",
								},
							},
							TerminationMessagePath:   "/dev/termination-log",
							TerminationMessagePolicy: corev1.TerminationMessageReadFile,
							LivenessProbe:            livenessProbe,
							ReadinessProbe:           readinessProbe,
							Lifecycle:                controllerSpec.LifecycleHooks,
							VolumeMounts:             controllerVolumeMounts(controllerSpec.CommonSpec),
						},
					},
					Volumes: controllerVolumes(instance),
				},
			},
		},
	}

	if controllerSpec.Sidecars != nil {
		dep.Spec.Template.Spec.Containers = append(dep.Spec.Template.Spec.Containers, controllerSpec.Sidecars...)
	}

	if err := ctrl.SetControllerReference(instance, dep, scheme); err != nil {
		return nil, err
	}
	return dep, nil
}

// deploymentMutate returns a mutate function that can be used to update the Deployment's spec.
func deploymentMutate(dep *appsv1.Deployment, spec appsv1.DeploymentSpec) controllerutil.MutateFn {
	return func() error {
		dep.Spec.Selector = spec.Selector
		dep.Spec.Strategy = spec.Strategy
		dep.Spec.Template.Annotations = spec.Template.Annotations
		dep.Spec.Template.Labels = spec.Template.Labels
		dep.Spec.Template.Spec.ServiceAccountName = spec.Template.Spec.ServiceAccountName
		dep.Spec.Template.Spec.HostAliases = spec.Template.Spec.HostAliases
		dep.Spec.Template.Spec.ImagePullSecrets = spec.Template.Spec.ImagePullSecrets
		dep.Spec.Template.Spec.HostAliases = spec.Template.Spec.HostAliases
		dep.Spec.Template.Spec.Affinity = spec.Template.Spec.Affinity
		dep.Spec.Template.Spec.NodeSelector = spec.Template.Spec.NodeSelector
		dep.Spec.Template.Spec.Tolerations = spec.Template.Spec.Tolerations
		dep.Spec.Template.Spec.PriorityClassName = spec.Template.Spec.PriorityClassName
		dep.Spec.Template.Spec.TopologySpreadConstraints = spec.Template.Spec.TopologySpreadConstraints
		dep.Spec.Template.Spec.SecurityContext = spec.Template.Spec.SecurityContext
		dep.Spec.Template.Spec.InitContainers = spec.Template.Spec.InitContainers
		dep.Spec.Template.Spec.Containers = spec.Template.Spec.Containers
		dep.Spec.Template.Spec.Volumes = spec.Template.Spec.Volumes
		return nil
	}
}
