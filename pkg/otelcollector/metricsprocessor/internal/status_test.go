package internal_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/collector/pdata/pcommon"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/check/v1"
	"github.com/fluxninja/aperture/pkg/metrics"
	"github.com/fluxninja/aperture/pkg/otelcollector"
	"github.com/fluxninja/aperture/pkg/otelcollector/metricsprocessor/internal"
)

var _ = Describe("Status", func() {
	Context("StatusesFromAttributes", func() {
		It("Read both statuses if exist", func() {
			attributes := pcommon.NewMap()
			attributes.PutStr(otelcollector.HTTPStatusCodeLabel, "201")
			attributes.PutStr(otelcollector.ApertureFeatureStatusLabel, otelcollector.ApertureResponseStatusOK)
			statusCode, featureStatus := internal.StatusesFromAttributes(attributes)
			Expect(statusCode).To(Equal("201"))
			Expect(featureStatus).To(Equal(otelcollector.ApertureResponseStatusOK))
		})

		It("Defaults to empty if not exist", func() {
			attributes := pcommon.NewMap()
			statusCode, featureStatus := internal.StatusesFromAttributes(attributes)
			Expect(statusCode).To(Equal(""))
			Expect(featureStatus).To(Equal(""))
		})
	})

	DescribeTable("StatusLabelsForMetrics", func(
		decisionType flowcontrolv1.CheckResponse_DecisionType,
		statusCode string,
		featureStatus string,
		expectedResponseStatus string,
	) {
		result := internal.StatusLabelsForMetrics(decisionType, statusCode, featureStatus)
		Expect(result).To(HaveLen(4))
		Expect(result).To(HaveKeyWithValue(metrics.ResponseStatusLabel, expectedResponseStatus))
		Expect(result).To(HaveKeyWithValue(metrics.DecisionTypeLabel, decisionType.String()))
		Expect(result).To(HaveKeyWithValue(metrics.StatusCodeLabel, statusCode))
		Expect(result).To(HaveKeyWithValue(metrics.FeatureStatusLabel, featureStatus))
	},
		Entry("Works for HTTP status OK",
			flowcontrolv1.CheckResponse_DECISION_TYPE_ACCEPTED,
			"201",
			"",
			metrics.ResponseStatusOK,
		),
		Entry("Works for HTTP status Error",
			flowcontrolv1.CheckResponse_DECISION_TYPE_ACCEPTED,
			"404",
			"",
			metrics.ResponseStatusError,
		),
		Entry("Works for Feature status OK",
			flowcontrolv1.CheckResponse_DECISION_TYPE_ACCEPTED,
			"",
			metrics.FeatureStatusOK,
			metrics.ResponseStatusOK,
		),
		Entry("Works for Feature status Error",
			flowcontrolv1.CheckResponse_DECISION_TYPE_ACCEPTED,
			"",
			metrics.FeatureStatusError,
			metrics.ResponseStatusError,
		),
	)
})
