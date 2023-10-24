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

	controllerv1alpha1 "github.com/fluxninja/aperture/v2/operator/api/controller/v1alpha1"
)

// deploymentForAPIService prepares the Deployment object for the Controller.
func deploymentForController(instance *controllerv1alpha1.Controller, log logr.Logger, scheme *runtime.Scheme) (*appsv1.Deployment, error) {
	spec := instance.Spec

	podLabels := controllers.CommonLabels(spec.Labels, instance.GetName(), controllers.ControllerServiceName)
	if spec.PodLabels != nil {
		if err := mergo.Map(&podLabels, spec.PodLabels, mergo.WithOverride); err != nil {
			log.Info(fmt.Sprintf("failed to merge the Pod labels for Deployment. error: %s.", err.Error()))
		}
	}

	annotations := spec.PodAnnotations
	if annotations == nil {
		annotations = map[string]string{}
	}
	annotations[controllers.SidecarAnnotationKey] = "false"

	probeScheme := corev1.URISchemeHTTP
	if instance.Spec.ConfigSpec.Server.TLS.Enabled {
		probeScheme = corev1.URISchemeHTTPS
	}

	livenessProbe, readinessProbe := controllers.ContainerProbes(spec.CommonSpec, probeScheme)

	serverPort, err := controllers.GetPort(spec.ConfigSpec.Server.Listener.Addr)
	if err != nil {
		return nil, err
	}

	dep := &appsv1.Deployment{
		ObjectMeta: v1.ObjectMeta{
			Name:        controllers.ControllerResourcesName(instance),
			Namespace:   instance.GetNamespace(),
			Labels:      controllers.CommonLabels(spec.Labels, instance.GetName(), controllers.ControllerServiceName),
			Annotations: spec.Annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &v1.LabelSelector{
				MatchLabels: controllers.SelectorLabels(instance.GetName(), controllers.ControllerServiceName),
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
					ServiceAccountName:            controllers.ControllerServiceAccountName(instance),
					HostAliases:                   spec.HostAliases,
					ImagePullSecrets:              controllers.ImagePullSecrets(spec.Image.Image),
					NodeSelector:                  spec.NodeSelector,
					Affinity:                      spec.Affinity,
					Tolerations:                   spec.Tolerations,
					SecurityContext:               controllers.PodSecurityContext(spec.PodSecurityContext),
					TerminationGracePeriodSeconds: pointer.Int64(spec.TerminationGracePeriodSeconds),
					InitContainers:                spec.InitContainers,
					Containers: []corev1.Container{
						{
							Name:            controllers.ControllerServiceName,
							Image:           controllers.ImageString(spec.Image.Image, spec.Image.Repository),
							ImagePullPolicy: corev1.PullPolicy(spec.Image.PullPolicy),
							SecurityContext: controllers.ContainerSecurityContext(spec.ContainerSecurityContext),
							Command:         spec.Command,
							Args:            spec.Args,
							Env:             controllers.ControllerEnv(instance),
							EnvFrom:         controllers.ContainerEnvFrom(spec.CommonSpec),
							Resources:       spec.Resources,
							Ports: []corev1.ContainerPort{
								{
									Name:          controllers.Server,
									ContainerPort: serverPort,
									Protocol:      controllers.TCP,
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
							VolumeMounts:             controllers.ControllerVolumeMounts(spec.CommonSpec),
						},
					},
					Volumes: controllers.ControllerVolumes(instance),
				},
			},
		},
	}

	if spec.Sidecars != nil {
		dep.Spec.Template.Spec.Containers = append(dep.Spec.Template.Spec.Containers, spec.Sidecars...)
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
		dep.Spec.Template.Spec.Affinity = spec.Template.Spec.Affinity
		dep.Spec.Template.Spec.NodeSelector = spec.Template.Spec.NodeSelector
		dep.Spec.Template.Spec.Tolerations = spec.Template.Spec.Tolerations
		dep.Spec.Template.Spec.TopologySpreadConstraints = spec.Template.Spec.TopologySpreadConstraints
		dep.Spec.Template.Spec.SecurityContext = spec.Template.Spec.SecurityContext
		dep.Spec.Template.Spec.InitContainers = spec.Template.Spec.InitContainers
		dep.Spec.Template.Spec.Containers = spec.Template.Spec.Containers
		dep.Spec.Template.Spec.Volumes = spec.Template.Spec.Volumes
		return nil
	}
}
