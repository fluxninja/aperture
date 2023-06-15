package internal_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/collector/pdata/pcommon"

	otelconsts "github.com/fluxninja/aperture/v2/pkg/otelcollector/consts"
	"github.com/fluxninja/aperture/v2/pkg/otelcollector/metricsprocessor/internal"
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
		map[string]float64{otelconsts.BytesSentLabel: 123},
		map[string]float64{otelconsts.HTTPRequestContentLength: 123},
	),
	Entry("Sets response content length",
		map[string]float64{otelconsts.BytesReceivedLabel: 123},
		map[string]float64{otelconsts.HTTPResponseContentLength: 123},
	),
	Entry("Sets flow duration",
		map[string]float64{otelconsts.ResponseDurationLabel: 123},
		map[string]float64{otelconsts.FlowDurationLabel: 123},
	),
	Entry("Sets workload duration",
		map[string]float64{
			otelconsts.ResponseDurationLabel:   123,
			otelconsts.EnvoyAuthzDurationLabel: 23,
		},
		map[string]float64{otelconsts.WorkloadDurationLabel: 100},
	),
)
