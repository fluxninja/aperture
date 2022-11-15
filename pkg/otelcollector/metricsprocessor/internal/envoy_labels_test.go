package internal_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/collector/pdata/pcommon"

	"github.com/fluxninja/aperture/pkg/otelcollector"
	"github.com/fluxninja/aperture/pkg/otelcollector/metricsprocessor/internal"
)

var _ = DescribeTable("Envoy labels", func(before, after map[string]float64) {
	attributes := pcommon.NewMap()
	for k, v := range before {
		attributes.PutDouble(k, v)
	}
	internal.AddEnvoySpecificLabels(attributes)
	for k, v := range after {
		rawOut, exists := attributes.Get(k)
		Expect(exists).To(BeTrue())
		Expect(rawOut.Double()).To(Equal(v))
	}
},
	Entry("Sets request content length",
		map[string]float64{otelcollector.EnvoyBytesSentLabel: 123},
		map[string]float64{otelcollector.HTTPRequestContentLength: 123},
	),
	Entry("Sets response content length",
		map[string]float64{otelcollector.EnvoyBytesReceivedLabel: 123},
		map[string]float64{otelcollector.HTTPResponseContentLength: 123},
	),
	Entry("Sets flow duration",
		map[string]float64{otelcollector.EnvoyResponseDurationLabel: 123},
		map[string]float64{otelcollector.FlowDurationLabel: 123},
	),
	Entry("Sets workload duration",
		map[string]float64{
			otelcollector.EnvoyResponseDurationLabel: 123,
			otelcollector.EnvoyAuthzDurationLabel:    23,
		},
		map[string]float64{otelcollector.WorkloadDurationLabel: 100},
	),
)
