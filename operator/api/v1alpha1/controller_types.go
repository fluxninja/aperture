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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
)

// ControllerSpec defines the desired state for the Controller.
type ControllerSpec struct {
	// CommonSpec defines the common state between Agent and Controller
	CommonSpec `json:",inline"`

	// Image configuration
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:={tag:"latest",pullPolicy:"IfNotPresent",registry:"docker.io/fluxninja",repository:"aperture-controller"}
	Image Image `json:"image"`

	// Pod's host aliases
	//+kubebuilder:validation:Optional
	HostAliases []corev1.HostAlias `json:"hostAliases"`
}
