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
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/fluxninja/aperture/operator/api/v1alpha1"
)

// daemonsetForAgent prepares the Daemonset object for the Agent.
func daemonsetForAgent(instance *v1alpha1.Agent, log logr.Logger, scheme *runtime.Scheme) (*appsv1.DaemonSet, error) {
	spec := instance.Spec

	podLabels := commonLabels(spec.Labels, instance.GetName(), agentServiceName)
	if spec.PodLabels != nil {
		if err := mergo.Map(&podLabels, spec.PodLabels, mergo.WithOverride); err != nil {
			log.Info(fmt.Sprintf("failed to merge the Pod labels for Deployment. error: %s.", err.Error()))
		}
	}

	probeScheme := corev1.URISchemeHTTP
	if instance.Spec.ConfigSpec.Server.TLS.Enabled {
		probeScheme = corev1.URISchemeHTTPS
	}

	livenessProbe, readinessProbe := containerProbes(spec.CommonSpec, probeScheme)

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

	daemonset := &appsv1.DaemonSet{
		ObjectMeta: v1.ObjectMeta{
			Name:        agentServiceName,
			Namespace:   instance.GetNamespace(),
			Labels:      commonLabels(spec.Labels, instance.GetName(), agentServiceName),
			Annotations: spec.Annotations,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &v1.LabelSelector{
				MatchLabels: selectorLabels(instance.GetName(), agentServiceName),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: v1.ObjectMeta{
					Labels:      podLabels,
					Annotations: spec.PodAnnotations,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName:            agentServiceName,
					ImagePullSecrets:              imagePullSecrets(spec.Image.Image),
					NodeSelector:                  spec.NodeSelector,
					Affinity:                      spec.Affinity,
					Tolerations:                   spec.Tolerations,
					SecurityContext:               podSecurityContext(spec.PodSecurityContext),
					TerminationGracePeriodSeconds: pointer.Int64(spec.TerminationGracePeriodSeconds),
					InitContainers:                spec.InitContainers,
					Containers: []corev1.Container{
						{
							Name:            agentServiceName,
							Image:           imageString(spec.Image.Image, spec.Image.Repository),
							ImagePullPolicy: corev1.PullPolicy(spec.Image.PullPolicy),
							SecurityContext: containerSecurityContext(spec.ContainerSecurityContext),
							Command:         spec.Command,
							Args:            spec.Args,
							Env:             agentEnv(instance, ""),
							EnvFrom:         containerEnvFrom(spec.CommonSpec),
							Resources:       spec.Resources,
							Ports: []corev1.ContainerPort{
								{
									Name:          server,
									ContainerPort: serverPort,
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          grpcOtel,
									ContainerPort: otelGRPCPort,
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          httpOtel,
									ContainerPort: otelHTTPPort,
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          distCache,
									ContainerPort: distCachePort,
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          memberList,
									ContainerPort: memberListPort,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							TerminationMessagePath:   "/dev/termination-log",
							TerminationMessagePolicy: corev1.TerminationMessageReadFile,
							LivenessProbe:            livenessProbe,
							ReadinessProbe:           readinessProbe,
							Lifecycle:                spec.LifecycleHooks,
							VolumeMounts:             agentVolumeMounts(spec),
						},
					},
					Volumes: agentVolumes(spec),
				},
			},
		},
	}

	if spec.Sidecars != nil {
		daemonset.Spec.Template.Spec.Containers = append(daemonset.Spec.Template.Spec.Containers, spec.Sidecars...)
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
