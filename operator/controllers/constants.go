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

const (
	// MutatingWebhookURI defines the URI for the Mutating Webhook for Pods.
	MutatingWebhookURI = "/mutate-pod"
	// AgentMutatingWebhookURI defines the URI for the Mutating Webhook for Agents.
	AgentMutatingWebhookURI = "agent-defaulter"
	// ControllerMutatingWebhookURI defines the URI for the Mutating Webhook for Controllers.
	ControllerMutatingWebhookURI = "controller-defaulter"

	secretKey                     = "apiKey"
	appName                       = "aperture"
	operatorName                  = appName + "-operator"
	controllerServiceName         = appName + "-controller"
	agentServiceName              = appName + "-agent"
	podMutatingWebhookName        = appName + "-injector"
	agentMutatingWebhookName      = appName + "-" + AgentMutatingWebhookURI
	controllerMutatingWebhookName = appName + "-" + ControllerMutatingWebhookURI
	validatingWebhookServiceName  = controllerServiceName + "-webhook"
	finalizerName                 = "fluxninja.com/finalizer"
	sidecarKey                    = "sidecar.fluxninja.com"
	sidecarAnnotationKey          = sidecarKey + "/injection"
	sidecarLabelKey               = appName + "-injection"
	agentGroupKey                 = sidecarKey + "/agent-group"
	v1Version                     = "v1"
	v1Alpha1Version               = "v1alpha1"
	enabled                       = "enabled"
	validatingWebhookSvcName      = validatingWebhookServiceName
	webhookClientCertName         = "client.pem"
	controllerCertKeyName         = "key.pem"
	controllerCertName            = "crt.pem"
	controllerCertPath            = "/etc/aperture/aperture-controller/certs"
	server                        = "server"
	grpcOtel                      = "grpc-otel"
	httpOtel                      = "http-otel"
	tcp                           = "TCP"
	distCache                     = "dist-cache"
	memberList                    = "memberlist"
	apertureFluxNinjaPlugin       = "aperture-plugin-fluxninja"
	defaulterAnnotationKey        = "fluxninja.com/set-defaults"
	failedStatus                  = "failed"
)
