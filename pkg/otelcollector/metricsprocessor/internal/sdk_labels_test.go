package internal_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/collector/pdata/pcommon"

	otelconsts "github.com/fluxninja/aperture/v2/pkg/otelcollector/consts"
	"github.com/fluxninja/aperture/v2/pkg/otelcollector/metricsprocessor/internal"
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
			otelconsts.ApertureFlowStartTimestampLabel: 123e6,
			otelconsts.ApertureFlowEndTimestampLabel:   246e6,
		},
		map[string]float64{otelconsts.FlowDurationLabel: 123},
	),
	Entry("Sets workload duration",
		map[string]int64{
			otelconsts.ApertureWorkloadStartTimestampLabel: 123e6,
			otelconsts.ApertureFlowEndTimestampLabel:       246e6,
		},
		map[string]float64{otelconsts.WorkloadDurationLabel: 123},
	),
)
