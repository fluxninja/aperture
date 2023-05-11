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
	"context"
	_ "embed"
	"net/http"
	"strings"

	"github.com/fluxninja/aperture/v2/operator/controllers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

//go:embed old_hook_test.tpl
var oldAgentSampleYAML string

//go:embed hook_test.tpl
var agentSampleYAML string

//go:embed hook_sidecar_test.tpl
var agentSidecarSampleYAML string

//go:embed hook_invalid_test.tpl
var agentInvalidSampleYAML string

var _ = Describe("Agent Hook Tests", Ordered, func() {
	Context("testing Handle", func() {
		It("should add defaults in spec when valid instance is provided. Sidecar disabled", func() {
			agentHook := AgentHooks{}

			res := agentHook.Handle(context.Background(), admission.Request{
				AdmissionRequest: v1.AdmissionRequest{
					Object: runtime.RawExtension{Raw: []byte(agentSampleYAML)},
				},
			})

			Expect(res.Allowed).To(Equal(true))
			Expect(len(res.Patches) > 0).To(Equal(true))
			By("overriding installation mode to KUBERNETES_DAEMONSET")
			patchFound := false
			for _, patch := range res.Patches {
				if patch.Operation == "add" && patch.Path == "/spec/config/fluxninja" {
					changes := patch.Value.(map[string]interface{})
					patchFound = changes["installation_mode"] == "KUBERNETES_DAEMONSET"
				}
			}
			Expect(patchFound).To(Equal(true))
		})

		It("should add defaults in spec when valid instance is provided. Sidecar enabled", func() {
			agentHook := AgentHooks{}

			res := agentHook.Handle(context.Background(), admission.Request{
				AdmissionRequest: v1.AdmissionRequest{
					Object: runtime.RawExtension{Raw: []byte(agentSidecarSampleYAML)},
				},
			})

			Expect(res.Allowed).To(Equal(true))
			Expect(len(res.Patches) > 0).To(Equal(true))
			By("overriding installation mode to KUBERNETES_SIDECAR")
			patchFound := false
			for _, patch := range res.Patches {
				if patch.Operation == "add" && patch.Path == "/spec/config/fluxninja" {
					changes := patch.Value.(map[string]interface{})
					patchFound = changes["installation_mode"] == "KUBERNETES_SIDECAR"
				}
			}
			Expect(patchFound).To(Equal(true))
		})

		It("should add annotation when valid instance is provided and installation mode is changed", func() {
			agentHook := AgentHooks{}

			res := agentHook.Handle(context.Background(), admission.Request{
				AdmissionRequest: v1.AdmissionRequest{
					Object:    runtime.RawExtension{Raw: []byte(agentSampleYAML)},
					OldObject: runtime.RawExtension{Raw: []byte(oldAgentSampleYAML)},
				},
			})

			Expect(res.Allowed).To(Equal(true))
			Expect(len(res.Patches) > 0).To(Equal(true))
			patchFound := false
			for _, patch := range res.Patches {
				if patch.Operation == "add" && patch.Path == "/metadata/annotations" {
					changes := patch.Value.(map[string]interface{})
					patchFound = changes[controllers.AgentModeChangeAnnotationKey] == "true"
				}
			}

			Expect(patchFound).To(Equal(true))
		})

		It("should not add defaults in spec when invalid instance is provided", func() {
			agentHook := AgentHooks{}

			res := agentHook.Handle(context.Background(), admission.Request{
				AdmissionRequest: v1.AdmissionRequest{
					Object: runtime.RawExtension{Raw: []byte(agentInvalidSampleYAML)},
				},
			})

			Expect(res.Allowed).To(Equal(false))
			Expect(int(res.Result.Code)).To(Equal(http.StatusBadRequest))
			Expect(strings.Contains(res.Result.Message, "PullPolicy")).To(Equal(true))
		})
	})
})
