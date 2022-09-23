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
	"context"
	_ "embed"
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

//go:embed hook_test.tpl
var controllerSampleYAML string

//go:embed hook_invalid_test.tpl
var controllerInvalidSampleYAML string

var _ = Describe("Controller Hook Tests", Ordered, func() {
	Context("testing Handle", func() {
		It("should add defaults in spec when valid Controller instance is provided", func() {
			controllerHook := ControllerHooks{}

			res := controllerHook.Handle(context.Background(), admission.Request{
				AdmissionRequest: v1.AdmissionRequest{
					Object: runtime.RawExtension{Raw: []byte(controllerSampleYAML)},
				},
			})

			Expect(res.Allowed).To(Equal(true))
			Expect(len(res.Patches) > 0).To(Equal(true))
		})

		It("should not add defaults in spec when invalid Controller instance is provided", func() {
			controllerHook := ControllerHooks{}

			res := controllerHook.Handle(context.Background(), admission.Request{
				AdmissionRequest: v1.AdmissionRequest{
					Object: runtime.RawExtension{Raw: []byte(controllerInvalidSampleYAML)},
				},
			})

			Expect(res.Allowed).To(Equal(false))
			Expect(int(res.Result.Code)).To(Equal(http.StatusBadRequest))
			Expect(strings.Contains(res.Result.Message, "PullPolicy")).To(Equal(true))
		})
	})
})
