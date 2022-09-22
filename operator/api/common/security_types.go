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

package common

// ServiceAccountSpec defines the the configuration for Service account for Agent and Controller.
type ServiceAccountSpec struct {
	// Specifies whether a ServiceAccount should be created
	Create bool `json:"create" default:"true"`

	// Additional Service Account annotations
	//+kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// Automount service account token for the server service account
	//+kubebuilder:validation:Optional
	AutomountServiceAccountToken bool `json:"automountServiceAccountToken" default:"true"`
}
