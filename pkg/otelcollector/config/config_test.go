package config_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	otelconfig "github.com/fluxninja/aperture/pkg/otelcollector/config"
)

var _ = Describe("Config", func() {
	It("Adds extensions properly", func() {
		otelConfig := otelconfig.NewOTelConfig()
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
