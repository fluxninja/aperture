package otelcollector_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/fluxninja/aperture/pkg/otelcollector"
)

var _ = Describe("Config", func() {
	It("Adds extensions properly", func() {
		otelConfig := otelcollector.NewOTELConfig()
		otelConfig.AddExtension("foo", map[string]interface{}{"bar": "baz"})
		otelConfig.AddExtension("empty", map[string]interface{}{})
		marshalledConfig, err := json.Marshal(otelConfig)
		Expect(err).NotTo(HaveOccurred())
		Expect(marshalledConfig).To(MatchJSON(`{
			"extensions": {
				"foo": {
					"bar": "baz"
				},
				"empty": {}
			},
			"service": {
				"Telemetry": {
					"logs": {
						"level": "INFO"
					}
				},
				"Pipelines": {},
				"Extensions": ["empty", "foo"]
			}
		}`))
	})
})
