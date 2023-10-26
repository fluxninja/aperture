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
	"fmt"

	"github.com/fluxninja/aperture/v2/operator/controllers"

	"github.com/go-logr/logr"
	"github.com/imdario/mergo"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	agentv1alpha1 "github.com/fluxninja/aperture/v2/operator/api/agent/v1alpha1"
)

// daemonsetForAgent prepares the DaemonSet object for the Agent.
func daemonsetForAgent(instance *agentv1alpha1.Agent, log logr.Logger, scheme *runtime.Scheme) (*appsv1.DaemonSet, error) {
	spec := instance.Spec

	podLabels := controllers.CommonLabels(spec.Labels, instance.GetName(), controllers.AgentServiceName)
	if spec.PodLabels != nil {
		if err := mergo.Map(&podLabels, spec.PodLabels, mergo.WithOverride); err != nil {
			log.Info(fmt.Sprintf("failed to merge the Pod labels for Deployment. error: %s.", err.Error()))
		}
	}

	probeScheme := corev1.URISchemeHTTP
	if instance.Spec.ConfigSpec.Server.TLS.Enabled {
		probeScheme = corev1.URISchemeHTTPS
	}

	livenessProbe, readinessProbe := controllers.ContainerProbes(spec.CommonSpec, probeScheme)

	serverPort, err := controllers.GetPort(spec.ConfigSpec.Server.Listener.Addr)
	if err != nil {
		return nil, err
	}

	distCachePort, err := controllers.GetPort(spec.ConfigSpec.DistCache.BindAddr)
	if err != nil {
		return nil, err
	}

	memberListPort, err := controllers.GetPort(spec.ConfigSpec.DistCache.MemberlistBindAddr)
	if err != nil {
		return nil, err
	}

	daemonset := &appsv1.DaemonSet{
		ObjectMeta: v1.ObjectMeta{
			Name:        controllers.AgentResourceName(instance),
			Namespace:   instance.GetNamespace(),
			Labels:      controllers.CommonLabels(spec.Labels, instance.GetName(), controllers.AgentServiceName),
			Annotations: spec.Annotations,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &v1.LabelSelector{
				MatchLabels: controllers.SelectorLabels(instance.GetName(), controllers.AgentServiceName),
			},
			MinReadySeconds: spec.MinReadySeconds,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: v1.ObjectMeta{
					Labels:      podLabels,
					Annotations: spec.PodAnnotations,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName:            controllers.AgentServiceAccountName(instance),
					ImagePullSecrets:              controllers.ImagePullSecrets(spec.Image.Image),
					NodeSelector:                  spec.NodeSelector,
					Affinity:                      spec.Affinity,
					Tolerations:                   spec.Tolerations,
					SecurityContext:               controllers.PodSecurityContext(spec.PodSecurityContext),
					TerminationGracePeriodSeconds: pointer.Int64(spec.TerminationGracePeriodSeconds),
					InitContainers:                spec.InitContainers,
					Containers: []corev1.Container{
						{
							Name:            controllers.AgentServiceName,
							Image:           controllers.ImageString(spec.Image.Image, spec.Image.Repository),
							ImagePullPolicy: corev1.PullPolicy(spec.Image.PullPolicy),
							SecurityContext: controllers.ContainerSecurityContext(spec.ContainerSecurityContext),
							Command:         spec.Command,
							Args:            spec.Args,
							Env:             controllers.AgentEnv(instance, ""),
							EnvFrom:         controllers.ContainerEnvFrom(spec.CommonSpec),
							Resources:       spec.Resources,
							Ports: []corev1.ContainerPort{
								{
									Name:          controllers.Server,
									ContainerPort: serverPort,
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          controllers.DistCache,
									ContainerPort: distCachePort,
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          controllers.MemberList,
									ContainerPort: memberListPort,
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          controllers.OtelDebugPort,
									ContainerPort: int32(spec.ConfigSpec.OTel.Ports.DebugPort),
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          controllers.OtelHealthcheckPort,
									ContainerPort: int32(spec.ConfigSpec.OTel.Ports.HealthCheckPort),
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          controllers.OtelPprofPort,
									ContainerPort: int32(spec.ConfigSpec.OTel.Ports.PprofPort),
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          controllers.OtelZpagesPort,
									ContainerPort: int32(spec.ConfigSpec.OTel.Ports.ZpagesPort),
									Protocol:      corev1.ProtocolTCP,
								},
							},
							TerminationMessagePath:   "/dev/termination-log",
							TerminationMessagePolicy: corev1.TerminationMessageReadFile,
							LivenessProbe:            livenessProbe,
							ReadinessProbe:           readinessProbe,
							Lifecycle:                spec.LifecycleHooks,
							VolumeMounts:             controllers.AgentVolumeMounts(spec),
						},
					},
					Volumes: controllers.AgentVolumes(instance),
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
