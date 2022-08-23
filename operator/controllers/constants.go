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
	// MutatingWebhookURI defines the URI for the Mutating Webhook.
	MutatingWebhookURI = "/mutate-pod"

	secretKey                = "apiKey"
	appName                  = "aperture"
	operatorName             = appName + "-operator"
	controllerServiceName    = appName + "-controller"
	agentServiceName         = appName + "-agent"
	mutatingWebhookName      = appName + "-injector"
	finalizerName            = "aperture.tech/finalizer"
	sidecarKey               = "sidecar.aperture.tech"
	sidecarAnnotationKey     = sidecarKey + "/injection"
	sidecarLabelKey          = appName + "-injection"
	agentGroupKey            = sidecarKey + "/agent-group"
	v1Version                = "v1"
	enabled                  = "enabled"
	validatingWebhookSvcName = "agent-webhooks"
	webhookClientCertName    = "client.pem"
	controllerCertKeyName    = "key.pem"
	controllerCertName       = "crt.pem"
	controllerCertPath       = "/etc/aperture/aperture-controller/certs"
)
