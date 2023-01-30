package internal_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/collector/pdata/pcommon"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/check/v1"
	"github.com/fluxninja/aperture/pkg/metrics"
	otelconsts "github.com/fluxninja/aperture/pkg/otelcollector/consts"
	"github.com/fluxninja/aperture/pkg/otelcollector/metricsprocessor/internal"
)

var _ = Describe("Status", func() {
	Context("StatusesFromAttributes", func() {
		It("Read both statuses if exist", func() {
			attributes := pcommon.NewMap()
			attributes.PutStr(otelconsts.HTTPStatusCodeLabel, "201")
			attributes.PutStr(otelconsts.ApertureFlowStatusLabel, otelconsts.ApertureFlowStatusOK)
			statusCode, flowStatus := internal.StatusesFromAttributes(attributes)
			Expect(statusCode).To(Equal("201"))
			Expect(flowStatus).To(Equal(otelconsts.ApertureFlowStatusOK))
		})

		It("Defaults to empty if not exist", func() {
			attributes := pcommon.NewMap()
			statusCode, flowStatus := internal.StatusesFromAttributes(attributes)
			Expect(statusCode).To(Equal(""))
			Expect(flowStatus).To(Equal(""))
		})
	})

	DescribeTable("StatusLabelsForMetrics", func(
		decisionType flowcontrolv1.CheckResponse_DecisionType,
		statusCode string,
		flowStatus string,
		expectedResponseStatus string,
	) {
		result := internal.StatusLabelsForMetrics(decisionType, statusCode, flowStatus)
		Expect(result).To(HaveLen(3))
		Expect(result).To(HaveKeyWithValue(metrics.FlowStatusLabel, expectedResponseStatus))
		Expect(result).To(HaveKeyWithValue(metrics.DecisionTypeLabel, decisionType.String()))
		Expect(result).To(HaveKeyWithValue(metrics.StatusCodeLabel, statusCode))
	},
		Entry("Works for HTTP status OK",
			flowcontrolv1.CheckResponse_DECISION_TYPE_ACCEPTED,
			"201",
			"",
			metrics.FlowStatusOK,
		),
		Entry("Works for HTTP status Error",
			flowcontrolv1.CheckResponse_DECISION_TYPE_ACCEPTED,
			"404",
			"",
			metrics.FlowStatusError,
		),
		Entry("Works for Flow status OK",
			flowcontrolv1.CheckResponse_DECISION_TYPE_ACCEPTED,
			"",
			metrics.FlowStatusOK,
			metrics.FlowStatusOK,
		),
		Entry("Works for Flow status Error",
			flowcontrolv1.CheckResponse_DECISION_TYPE_ACCEPTED,
			"",
			metrics.FlowStatusError,
			metrics.FlowStatusError,
		),
	)
})
