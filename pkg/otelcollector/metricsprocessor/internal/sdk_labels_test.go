package internal_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/collector/pdata/pcommon"

	"github.com/fluxninja/aperture/pkg/otelcollector"
	"github.com/fluxninja/aperture/pkg/otelcollector/metricsprocessor/internal"
)

var _ = DescribeTable("SDK labels", func(before map[string]int64, after map[string]float64) {
	attributes := pcommon.NewMap()
	for k, v := range before {
		attributes.PutInt(k, v)
	}
	internal.AddSDKSpecificLabels(attributes)
	for k, v := range after {
		rawOut, exists := attributes.Get(k)
		Expect(exists).To(BeTrue())
		Expect(rawOut.Double()).To(Equal(v))
	}
},
	Entry("Sets flow duration",
		map[string]int64{
			otelcollector.ApertureFlowStartTimestampLabel: 123e6,
			otelcollector.ApertureFlowEndTimestampLabel:   246e6,
		},
		map[string]float64{otelcollector.FlowDurationLabel: 123},
	),
	Entry("Sets workload duration",
		map[string]int64{
			otelcollector.ApertureWorkloadStartTimestampLabel: 123e6,
			otelcollector.ApertureFlowEndTimestampLabel:       246e6,
		},
		map[string]float64{otelcollector.WorkloadDurationLabel: 123},
	),
)
