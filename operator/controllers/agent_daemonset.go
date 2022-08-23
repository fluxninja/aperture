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

	"aperture.tech/operators/aperture-operator/api/v1alpha1"
)

// daemonsetForAgent prepares the Daemonset object for the Agent.
func daemonsetForAgent(instance *v1alpha1.Aperture, log logr.Logger, scheme *runtime.Scheme) (*appsv1.DaemonSet, error) {
	agentSpec := instance.Spec.Agent

	podLabels := commonLabels(instance, agentServiceName)
	if agentSpec.PodLabels != nil {
		if err := mergo.Map(&podLabels, agentSpec.PodLabels, mergo.WithOverride); err != nil {
			log.Info(fmt.Sprintf("failed to merge the Pod labels for Deployment. error: %s.", err.Error()))
		}
	}

	livenessProbe, readinessProbe := containerProbes(instance.Spec.Agent.CommonSpec)

	daemonset := &appsv1.DaemonSet{
		ObjectMeta: v1.ObjectMeta{
			Name:        agentServiceName,
			Namespace:   instance.GetNamespace(),
			Labels:      commonLabels(instance, agentServiceName),
			Annotations: instance.Spec.Annotations,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &v1.LabelSelector{
				MatchLabels: selectorLabels(instance.GetName(), agentServiceName),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: v1.ObjectMeta{
					Labels:      podLabels,
					Annotations: agentSpec.PodAnnotations,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName:            agentServiceName,
					ImagePullSecrets:              imagePullSecrets(instance.Spec.ImagePullSecrets, agentSpec.Image),
					NodeSelector:                  agentSpec.NodeSelector,
					Affinity:                      agentSpec.Affinity,
					Tolerations:                   agentSpec.Tolerations,
					SecurityContext:               podSecurityContext(agentSpec.PodSecurityContext),
					TerminationGracePeriodSeconds: agentSpec.TerminationGracePeriodSeconds,
					InitContainers:                agentSpec.InitContainers,
					Containers: []corev1.Container{
						{
							Name:            agentServiceName,
							Image:           imageString(instance.Spec.ImageRegistry, agentSpec.Image),
							ImagePullPolicy: corev1.PullPolicy(agentSpec.Image.PullPolicy),
							SecurityContext: containerSecurityContext(agentSpec.ContainerSecurityContext),
							Command:         agentSpec.Command,
							Args:            agentSpec.Args,
							Env:             agentEnv(instance, ""),
							EnvFrom:         containerEnvFrom(agentSpec.CommonSpec),
							Resources:       agentSpec.Resources,
							Ports: []corev1.ContainerPort{
								{
									Name:          "grpc",
									ContainerPort: agentSpec.ServerPort,
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          "grpc-otel",
									ContainerPort: 4317,
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          "dist-cache",
									ContainerPort: agentSpec.DistributedCachePort,
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          "memberlist",
									ContainerPort: agentSpec.MemberListPort,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							TerminationMessagePath:   "/dev/termination-log",
							TerminationMessagePolicy: corev1.TerminationMessageReadFile,
							LivenessProbe:            livenessProbe,
							ReadinessProbe:           readinessProbe,
							Lifecycle:                agentSpec.LifecycleHooks,
							VolumeMounts:             agentVolumeMounts(agentSpec),
						},
					},
					Volumes: agentVolumes(agentSpec),
				},
			},
		},
	}

	if agentSpec.Sidecars != nil {
		daemonset.Spec.Template.Spec.Containers = append(daemonset.Spec.Template.Spec.Containers, agentSpec.Sidecars...)
	}

	if err := ctrl.SetControllerReference(instance, daemonset, scheme); err != nil {
		return nil, err
	}
	return daemonset, nil
}

// daemonsetMutate returns a mutate function that can be used to update the DaemonSet's spec.
func daemonsetMutate(dms *appsv1.DaemonSet, spec appsv1.DaemonSetSpec) controllerutil.MutateFn {
	return func() error {
		dms.Spec.Selector = spec.Selector
		dms.Spec.Template.Annotations = spec.Template.Annotations
		dms.Spec.Template.Labels = spec.Template.Labels
		dms.Spec.Template.Spec.ServiceAccountName = spec.Template.Spec.ServiceAccountName
		dms.Spec.Template.Spec.ImagePullSecrets = spec.Template.Spec.ImagePullSecrets
		dms.Spec.Template.Spec.NodeSelector = spec.Template.Spec.NodeSelector
		dms.Spec.Template.Spec.Tolerations = spec.Template.Spec.Tolerations
		dms.Spec.Template.Spec.SecurityContext = spec.Template.Spec.SecurityContext
		dms.Spec.Template.Spec.InitContainers = spec.Template.Spec.InitContainers
		dms.Spec.Template.Spec.Containers = spec.Template.Spec.Containers
		dms.Spec.Template.Spec.Volumes = spec.Template.Spec.Volumes
		return nil
	}
}
