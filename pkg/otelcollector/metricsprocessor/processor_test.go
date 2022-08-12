package metricsprocessor

import (
	"fmt"
	"strings"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"golang.org/x/net/context"
	"k8s.io/apimachinery/pkg/util/json"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	"github.com/fluxninja/aperture/pkg/otelcollector"
	"github.com/fluxninja/aperture/pkg/policies/mocks"
)

var _ = Describe("Metrics Processor", func() {
	var (
		pr        *prometheus.Registry
		cfg       *Config
		processor *metricsProcessor
		engine    *mocks.MockEngineAPI
	)

	BeforeEach(func() {
		pr = prometheus.NewRegistry()
		ctrl := gomock.NewController(GinkgoT())
		engine = mocks.NewMockEngineAPI(ctrl)
		cfg = &Config{
			engine:               engine,
			promRegistry:         pr,
			LatencyBucketStartMS: 0,
			LatencyBucketWidthMS: 10,
			LatencyBucketCount:   3,
		}
		var err error
		processor, err = newProcessor(cfg)
		Expect(err).NotTo(HaveOccurred())
	})

	DescribeTable("Processing logs",
		func(
			controlPoint string,
			checkResponse *flowcontrolv1.CheckResponse,
			authzResponse *flowcontrolv1.AuthzResponse,
			expectedErr error,
			expectedMetrics string,
			expectedLabels map[string]interface{},
		) {
			ctx := context.Background()

			expectEngineCalls(engine, checkResponse)

			logs := someLogs(checkResponse, authzResponse, controlPoint)
			modifiedLogs, err := processor.ConsumeLogs(ctx, logs)
			if expectedErr != nil {
				Expect(err).NotTo(MatchError(expectedErr))
				return
			}
			Expect(err).NotTo(HaveOccurred())
			Expect(modifiedLogs).To(Equal(logs))

			By("sending proper metrics")
			expected := strings.NewReader(expectedMetrics)
			err = testutil.CollectAndCompare(processor.workloadLatencyHistogram, expected, "workload_latency_ms")
			Expect(err).NotTo(HaveOccurred())

			By("adding proper labels")
			logRecords := allLogRecords(modifiedLogs)
			Expect(logRecords).To(HaveLen(1))

			for k, v := range expectedLabels {
				Expect(logRecords[0].Attributes().AsRaw()).To(HaveKeyWithValue(k, v))
			}
		},

		Entry("record with single policy - ingress",
			otelcollector.ControlPointIngress,
			&flowcontrolv1.CheckResponse{
				DecisionType: flowcontrolv1.DecisionType_DECISION_TYPE_ACCEPTED,
				LimiterDecisions: []*flowcontrolv1.LimiterDecision{
					{
						PolicyName:     "foo",
						PolicyHash:     "foo-hash",
						ComponentIndex: 1,
						Dropped:        true,
						Details: &flowcontrolv1.LimiterDecision_ConcurrencyLimiter_{
							ConcurrencyLimiter: &flowcontrolv1.LimiterDecision_ConcurrencyLimiter{
								WorkloadIndex: "0",
							},
						},
					},
				},
				FluxMeters: []*flowcontrolv1.FluxMeter{
					{
						AgentGroup:    "ag",
						PolicyName:    "foo",
						PolicyHash:    "foo-hash",
						FluxMeterName: "bar",
					},
				},
			},
			&flowcontrolv1.AuthzResponse{
				Status: flowcontrolv1.AuthzResponse_STATUS_NO_ERROR,
			},
			nil,
			`# HELP workload_latency_ms Latency histogram of workload
			# TYPE workload_latency_ms histogram
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="0",le="0"} 0
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="0",le="10"} 1
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="0",le="20"} 1
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="0",le="+Inf"} 1
			workload_latency_ms_sum{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="0"} 5
			workload_latency_ms_count{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="0"} 1
			`,
			map[string]interface{}{
				otelcollector.AuthzStatusLabel:                 "STATUS_NO_ERROR",
				otelcollector.DecisionTypeLabel:                "DECISION_TYPE_ACCEPTED",
				otelcollector.DecisionErrorReasonLabel:         "",
				otelcollector.DecisionRejectReasonLabel:        "",
				otelcollector.FluxMetersLabel:                  []interface{}{"policy_name:foo,flux_meter_name:bar,policy_hash:foo-hash"},
				otelcollector.RateLimitersLabel:                []interface{}{},
				otelcollector.DroppingRateLimitersLabel:        []interface{}{},
				otelcollector.ConcurrencyLimitersLabel:         []interface{}{"policy_name:foo,component_index:1,workload_index:0,policy_hash:foo-hash"},
				otelcollector.DroppingConcurrencyLimitersLabel: []interface{}{"policy_name:foo,component_index:1,workload_index:0,policy_hash:foo-hash"},
			},
		),

		Entry("record with single policy - feature",
			otelcollector.ControlPointFeature,
			&flowcontrolv1.CheckResponse{
				DecisionType: flowcontrolv1.DecisionType_DECISION_TYPE_REJECTED,
				DecisionReason: &flowcontrolv1.DecisionReason{
					RejectReason: flowcontrolv1.DecisionReason_REJECT_REASON_RATE_LIMITED,
				},
				LimiterDecisions: []*flowcontrolv1.LimiterDecision{
					{
						PolicyName:     "foo",
						PolicyHash:     "foo-hash",
						ComponentIndex: 1,
						Dropped:        true,
						Details: &flowcontrolv1.LimiterDecision_ConcurrencyLimiter_{
							ConcurrencyLimiter: &flowcontrolv1.LimiterDecision_ConcurrencyLimiter{
								WorkloadIndex: "0",
							},
						},
					},
				},
				FluxMeters: []*flowcontrolv1.FluxMeter{},
			},
			&flowcontrolv1.AuthzResponse{
				Status: flowcontrolv1.AuthzResponse_STATUS_NO_ERROR,
			},
			nil,
			`# HELP workload_latency_ms Latency histogram of workload
			# TYPE workload_latency_ms histogram
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="0",le="0"} 0
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="0",le="10"} 1
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="0",le="20"} 1
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="0",le="+Inf"} 1
			workload_latency_ms_sum{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="0"} 5
			workload_latency_ms_count{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="0"} 1
			`,
			map[string]interface{}{
				otelcollector.AuthzStatusLabel:                 "STATUS_NO_ERROR",
				otelcollector.DecisionTypeLabel:                "DECISION_TYPE_REJECTED",
				otelcollector.DecisionRejectReasonLabel:        "REJECT_REASON_RATE_LIMITED",
				otelcollector.RateLimitersLabel:                []interface{}{},
				otelcollector.DroppingRateLimitersLabel:        []interface{}{},
				otelcollector.ConcurrencyLimitersLabel:         []interface{}{"policy_name:foo,component_index:1,workload_index:0,policy_hash:foo-hash"},
				otelcollector.DroppingConcurrencyLimitersLabel: []interface{}{"policy_name:foo,component_index:1,workload_index:0,policy_hash:foo-hash"},
			},
		),

		Entry("record with two policies",
			otelcollector.ControlPointIngress,
			&flowcontrolv1.CheckResponse{
				DecisionType: flowcontrolv1.DecisionType_DECISION_TYPE_ACCEPTED,
				DecisionReason: &flowcontrolv1.DecisionReason{
					RejectReason: flowcontrolv1.DecisionReason_REJECT_REASON_UNSPECIFIED,
				},
				LimiterDecisions: []*flowcontrolv1.LimiterDecision{
					{
						PolicyName:     "foo",
						PolicyHash:     "foo-hash",
						ComponentIndex: 1,
						Dropped:        true,
						Details: &flowcontrolv1.LimiterDecision_ConcurrencyLimiter_{
							ConcurrencyLimiter: &flowcontrolv1.LimiterDecision_ConcurrencyLimiter{
								WorkloadIndex: "0",
							},
						},
					},
					{
						PolicyName:     "fizz",
						PolicyHash:     "fizz-hash",
						ComponentIndex: 1,
						Dropped:        true,
						Details: &flowcontrolv1.LimiterDecision_ConcurrencyLimiter_{
							ConcurrencyLimiter: &flowcontrolv1.LimiterDecision_ConcurrencyLimiter{
								WorkloadIndex: "1",
							},
						},
					},
					{
						PolicyName:     "fizz",
						PolicyHash:     "fizz-hash",
						ComponentIndex: 2,
						Dropped:        false,
						Details: &flowcontrolv1.LimiterDecision_ConcurrencyLimiter_{
							ConcurrencyLimiter: &flowcontrolv1.LimiterDecision_ConcurrencyLimiter{
								WorkloadIndex: "2",
							},
						},
					},
				},
				FluxMeters: []*flowcontrolv1.FluxMeter{},
			},
			&flowcontrolv1.AuthzResponse{
				Status: flowcontrolv1.AuthzResponse_STATUS_NO_ERROR,
			},
			nil,
			`# HELP workload_latency_ms Latency histogram of workload
			# TYPE workload_latency_ms histogram
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="fizz-hash",policy_name="fizz",workload_index="1",le="0"} 0
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="fizz-hash",policy_name="fizz",workload_index="1",le="10"} 1
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="fizz-hash",policy_name="fizz",workload_index="1",le="20"} 1
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="fizz-hash",policy_name="fizz",workload_index="1",le="+Inf"} 1
			workload_latency_ms_sum{component_index="1",dropped="true",policy_hash="fizz-hash",policy_name="fizz",workload_index="1"} 5
			workload_latency_ms_count{component_index="1",dropped="true",policy_hash="fizz-hash",policy_name="fizz",workload_index="1"} 1
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="0",le="0"} 0
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="0",le="10"} 1
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="0",le="20"} 1
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="0",le="+Inf"} 1
			workload_latency_ms_sum{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="0"} 5
			workload_latency_ms_count{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="0"} 1
			workload_latency_ms_bucket{component_index="2",dropped="false",policy_hash="fizz-hash",policy_name="fizz",workload_index="2",le="0"} 0
			workload_latency_ms_bucket{component_index="2",dropped="false",policy_hash="fizz-hash",policy_name="fizz",workload_index="2",le="10"} 1
			workload_latency_ms_bucket{component_index="2",dropped="false",policy_hash="fizz-hash",policy_name="fizz",workload_index="2",le="20"} 1
			workload_latency_ms_bucket{component_index="2",dropped="false",policy_hash="fizz-hash",policy_name="fizz",workload_index="2",le="+Inf"} 1
			workload_latency_ms_sum{component_index="2",dropped="false",policy_hash="fizz-hash",policy_name="fizz",workload_index="2"} 5
			workload_latency_ms_count{component_index="2",dropped="false",policy_hash="fizz-hash",policy_name="fizz",workload_index="2"} 1
			`,
			map[string]interface{}{
				otelcollector.AuthzStatusLabel:          "STATUS_NO_ERROR",
				otelcollector.DecisionTypeLabel:         "DECISION_TYPE_ACCEPTED",
				otelcollector.DecisionErrorReasonLabel:  "ERROR_REASON_UNSPECIFIED",
				otelcollector.RateLimitersLabel:         []interface{}{},
				otelcollector.DroppingRateLimitersLabel: []interface{}{},
				otelcollector.ConcurrencyLimitersLabel: []interface{}{
					"policy_name:foo,component_index:1,workload_index:0,policy_hash:foo-hash",
					"policy_name:fizz,component_index:1,workload_index:1,policy_hash:fizz-hash",
					"policy_name:fizz,component_index:2,workload_index:2,policy_hash:fizz-hash",
				},
			},
		),
	)

	DescribeTable("Processing traces",
		func(
			controlPoint string,
			checkResponse *flowcontrolv1.CheckResponse,
			expectedErr error,
			expectedMetrics string,
			expectedLabels map[string]interface{},
		) {
			ctx := context.Background()

			expectEngineCalls(engine, checkResponse)

			traces := someTraces(checkResponse, controlPoint)
			modifiedTraces, err := processor.ConsumeTraces(ctx, traces)
			if expectedErr != nil {
				Expect(err).NotTo(MatchError(expectedErr))
				return
			}
			Expect(err).NotTo(HaveOccurred())
			Expect(modifiedTraces).To(Equal(traces))

			By("sending proper metrics")
			expected := strings.NewReader(expectedMetrics)
			err = testutil.CollectAndCompare(processor.workloadLatencyHistogram, expected, "workload_latency_ms")
			Expect(err).NotTo(HaveOccurred())

			By("adding proper labels")
			traceRecords := allTraceRecords(modifiedTraces)
			Expect(traceRecords).To(HaveLen(1))

			for k, v := range expectedLabels {
				Expect(traceRecords[0].Attributes().AsRaw()).To(HaveKeyWithValue(k, v))
			}
		},

		Entry("record with single policy - ingress",
			otelcollector.ControlPointIngress,
			&flowcontrolv1.CheckResponse{
				DecisionType: flowcontrolv1.DecisionType_DECISION_TYPE_ACCEPTED,
				LimiterDecisions: []*flowcontrolv1.LimiterDecision{
					{
						PolicyName:     "foo",
						PolicyHash:     "foo-hash",
						ComponentIndex: 1,
						Dropped:        true,
						Details: &flowcontrolv1.LimiterDecision_ConcurrencyLimiter_{
							ConcurrencyLimiter: &flowcontrolv1.LimiterDecision_ConcurrencyLimiter{
								WorkloadIndex: "0",
							},
						},
					},
				},
				FluxMeters: []*flowcontrolv1.FluxMeter{
					{
						AgentGroup:    "ag",
						PolicyName:    "foo",
						PolicyHash:    "foo-hash",
						FluxMeterName: "bar",
					},
				},
			},
			nil,
			`# HELP workload_latency_ms Latency histogram of workload
			# TYPE workload_latency_ms histogram
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="0",le="0"} 0
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="0",le="10"} 1
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="0",le="20"} 1
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="0",le="+Inf"} 1
			workload_latency_ms_sum{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="0"} 5
			workload_latency_ms_count{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="0"} 1
			`,
			map[string]interface{}{
				otelcollector.DecisionTypeLabel:                "DECISION_TYPE_ACCEPTED",
				otelcollector.DecisionErrorReasonLabel:         "",
				otelcollector.DecisionRejectReasonLabel:        "",
				otelcollector.FluxMetersLabel:                  []interface{}{"policy_name:foo,flux_meter_name:bar,policy_hash:foo-hash"},
				otelcollector.RateLimitersLabel:                []interface{}{},
				otelcollector.DroppingRateLimitersLabel:        []interface{}{},
				otelcollector.ConcurrencyLimitersLabel:         []interface{}{"policy_name:foo,component_index:1,workload_index:0,policy_hash:foo-hash"},
				otelcollector.DroppingConcurrencyLimitersLabel: []interface{}{"policy_name:foo,component_index:1,workload_index:0,policy_hash:foo-hash"},
			},
		),

		Entry("record with single policy - feature",
			otelcollector.ControlPointFeature,
			&flowcontrolv1.CheckResponse{
				DecisionType: flowcontrolv1.DecisionType_DECISION_TYPE_REJECTED,
				DecisionReason: &flowcontrolv1.DecisionReason{
					RejectReason: flowcontrolv1.DecisionReason_REJECT_REASON_RATE_LIMITED,
				},
				LimiterDecisions: []*flowcontrolv1.LimiterDecision{
					{
						PolicyName:     "foo",
						PolicyHash:     "foo-hash",
						ComponentIndex: 1,
						Dropped:        true,
						Details: &flowcontrolv1.LimiterDecision_ConcurrencyLimiter_{
							ConcurrencyLimiter: &flowcontrolv1.LimiterDecision_ConcurrencyLimiter{
								WorkloadIndex: "0",
							},
						},
					},
				},
				FluxMeters: []*flowcontrolv1.FluxMeter{},
			},
			nil,
			`# HELP workload_latency_ms Latency histogram of workload
			# TYPE workload_latency_ms histogram
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="0",le="0"} 0
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="0",le="10"} 1
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="0",le="20"} 1
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="0",le="+Inf"} 1
			workload_latency_ms_sum{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="0"} 5
			workload_latency_ms_count{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="0"} 1
			`,
			map[string]interface{}{
				otelcollector.DecisionTypeLabel:                "DECISION_TYPE_REJECTED",
				otelcollector.DecisionRejectReasonLabel:        "REJECT_REASON_RATE_LIMITED",
				otelcollector.RateLimitersLabel:                []interface{}{},
				otelcollector.DroppingRateLimitersLabel:        []interface{}{},
				otelcollector.ConcurrencyLimitersLabel:         []interface{}{"policy_name:foo,component_index:1,workload_index:0,policy_hash:foo-hash"},
				otelcollector.DroppingConcurrencyLimitersLabel: []interface{}{"policy_name:foo,component_index:1,workload_index:0,policy_hash:foo-hash"},
			},
		),

		Entry("record with two policies",
			otelcollector.ControlPointIngress,
			&flowcontrolv1.CheckResponse{
				DecisionType: flowcontrolv1.DecisionType_DECISION_TYPE_ACCEPTED,
				DecisionReason: &flowcontrolv1.DecisionReason{
					RejectReason: flowcontrolv1.DecisionReason_REJECT_REASON_UNSPECIFIED,
				},
				LimiterDecisions: []*flowcontrolv1.LimiterDecision{
					{
						PolicyName:     "foo",
						PolicyHash:     "foo-hash",
						ComponentIndex: 1,
						Dropped:        true,
						Details: &flowcontrolv1.LimiterDecision_RateLimiter_{
							RateLimiter: &flowcontrolv1.LimiterDecision_RateLimiter{
								Remaining: 10,
								Current:   5,
								Label:     "gold",
							},
						},
					},
					{
						PolicyName:     "fizz",
						PolicyHash:     "fizz-hash",
						ComponentIndex: 1,
						Dropped:        true,
						Details: &flowcontrolv1.LimiterDecision_ConcurrencyLimiter_{
							ConcurrencyLimiter: &flowcontrolv1.LimiterDecision_ConcurrencyLimiter{
								WorkloadIndex: "1",
							},
						},
					},
					{
						PolicyName:     "fizz",
						PolicyHash:     "fizz-hash",
						ComponentIndex: 2,
						Dropped:        true,
						Details: &flowcontrolv1.LimiterDecision_ConcurrencyLimiter_{
							ConcurrencyLimiter: &flowcontrolv1.LimiterDecision_ConcurrencyLimiter{
								WorkloadIndex: "2",
							},
						},
					},
				},
				FluxMeters: []*flowcontrolv1.FluxMeter{},
			},
			nil,
			`# HELP workload_latency_ms Latency histogram of workload
			# TYPE workload_latency_ms histogram
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="fizz-hash",policy_name="fizz",workload_index="1",le="0"} 0
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="fizz-hash",policy_name="fizz",workload_index="1",le="10"} 1
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="fizz-hash",policy_name="fizz",workload_index="1",le="20"} 1
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="fizz-hash",policy_name="fizz",workload_index="1",le="+Inf"} 1
			workload_latency_ms_sum{component_index="1",dropped="true",policy_hash="fizz-hash",policy_name="fizz",workload_index="1"} 5
			workload_latency_ms_count{component_index="1",dropped="true",policy_hash="fizz-hash",policy_name="fizz",workload_index="1"} 1
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="",le="0"} 0
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="",le="10"} 1
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="",le="20"} 1
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="",le="+Inf"} 1
			workload_latency_ms_sum{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index=""} 5
			workload_latency_ms_count{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index=""} 1
			workload_latency_ms_bucket{component_index="2",dropped="true",policy_hash="fizz-hash",policy_name="fizz",workload_index="2",le="0"} 0
			workload_latency_ms_bucket{component_index="2",dropped="true",policy_hash="fizz-hash",policy_name="fizz",workload_index="2",le="10"} 1
			workload_latency_ms_bucket{component_index="2",dropped="true",policy_hash="fizz-hash",policy_name="fizz",workload_index="2",le="20"} 1
			workload_latency_ms_bucket{component_index="2",dropped="true",policy_hash="fizz-hash",policy_name="fizz",workload_index="2",le="+Inf"} 1
			workload_latency_ms_sum{component_index="2",dropped="true",policy_hash="fizz-hash",policy_name="fizz",workload_index="2"} 5
			workload_latency_ms_count{component_index="2",dropped="true",policy_hash="fizz-hash",policy_name="fizz",workload_index="2"} 1
			`,
			map[string]interface{}{
				otelcollector.DecisionTypeLabel:         "DECISION_TYPE_ACCEPTED",
				otelcollector.DecisionRejectReasonLabel: "REJECT_REASON_UNSPECIFIED",
				otelcollector.RateLimitersLabel: []interface{}{
					"policy_name:foo,component_index:1,policy_hash:foo-hash",
				},
				otelcollector.DroppingRateLimitersLabel: []interface{}{
					"policy_name:foo,component_index:1,policy_hash:foo-hash",
				},
				otelcollector.ConcurrencyLimitersLabel: []interface{}{
					"policy_name:fizz,component_index:1,workload_index:1,policy_hash:fizz-hash",
					"policy_name:fizz,component_index:2,workload_index:2,policy_hash:fizz-hash",
				},
				otelcollector.DroppingConcurrencyLimitersLabel: []interface{}{
					"policy_name:fizz,component_index:1,workload_index:1,policy_hash:fizz-hash",
					"policy_name:fizz,component_index:2,workload_index:2,policy_hash:fizz-hash",
				},
			},
		),
	)
})

// someLogs will return a plog.Logs instance with single LogRecord
func someLogs(
	checkResponse *flowcontrolv1.CheckResponse,
	authzResponse *flowcontrolv1.AuthzResponse,
	controlPoint string,
) plog.Logs {
	logs := plog.NewLogs()
	logs.ResourceLogs().AppendEmpty()

	resourceLogsSlice := logs.ResourceLogs()
	for i := 0; i < resourceLogsSlice.Len(); i++ {
		resourceLogsSlice.At(i).ScopeLogs().AppendEmpty()

		instrumentationLogsSlice := resourceLogsSlice.At(i).ScopeLogs()
		for j := 0; j < instrumentationLogsSlice.Len(); j++ {
			logRecord := instrumentationLogsSlice.At(j).LogRecords().AppendEmpty()
			marshalledCheckResponse, err := json.Marshal(checkResponse)
			Expect(err).NotTo(HaveOccurred())
			marshalledAuthzResponse, err := json.Marshal(authzResponse)
			Expect(err).NotTo(HaveOccurred())
			logRecord.Attributes().InsertString(otelcollector.MarshalledCheckResponseLabel, string(marshalledCheckResponse))
			logRecord.Attributes().InsertString(otelcollector.MarshalledAuthzResponseLabel, string(marshalledAuthzResponse))
			logRecord.Attributes().InsertString(otelcollector.StatusCodeLabel, "201")
			logRecord.Attributes().InsertString(otelcollector.ControlPointLabel, controlPoint)
			switch controlPoint {
			case otelcollector.ControlPointIngress, otelcollector.ControlPointEgress:
				logRecord.Attributes().InsertString(otelcollector.HTTPDurationLabel, "5")
			case otelcollector.ControlPointFeature:
				logRecord.Attributes().InsertString(otelcollector.FeatureDurationLabel, "5")
			}
		}
	}

	return logs
}

// someTraces will return a ptrace.Traces instance with single SpanRecord
func someTraces(
	checkResponse *flowcontrolv1.CheckResponse,
	controlPoint string,
) ptrace.Traces {
	traces := ptrace.NewTraces()
	traces.ResourceSpans().AppendEmpty()

	resourceSpansSlice := traces.ResourceSpans()
	for i := 0; i < resourceSpansSlice.Len(); i++ {
		resourceSpansSlice.At(i).ScopeSpans().AppendEmpty()

		instrumentationSpansSlice := resourceSpansSlice.At(i).ScopeSpans()
		for j := 0; j < instrumentationSpansSlice.Len(); j++ {
			span := instrumentationSpansSlice.At(j).Spans().AppendEmpty()
			marshalledCheckResponse, err := json.Marshal(checkResponse)
			Expect(err).NotTo(HaveOccurred())
			span.Attributes().InsertString(otelcollector.MarshalledCheckResponseLabel, string(marshalledCheckResponse))
			span.Attributes().InsertString(otelcollector.StatusCodeLabel, "201")
			span.Attributes().InsertString(otelcollector.ControlPointLabel, controlPoint)
			switch controlPoint {
			case otelcollector.ControlPointIngress, otelcollector.ControlPointEgress:
				span.Attributes().InsertString(otelcollector.HTTPDurationLabel, "5")
			case otelcollector.ControlPointFeature:
				span.Attributes().InsertString(otelcollector.FeatureDurationLabel, "5")
			}
		}
	}

	return traces
}

// firstLogRecord extracts the only log record from one-record logs created by someLogs()
func allLogRecords(logs plog.Logs) []plog.LogRecord {
	var logRecords []plog.LogRecord

	resourceLogsSlice := logs.ResourceLogs()
	for i := 0; i < resourceLogsSlice.Len(); i++ {
		instrumentationLogsSlice := resourceLogsSlice.At(i).ScopeLogs()
		for j := 0; j < instrumentationLogsSlice.Len(); j++ {
			records := instrumentationLogsSlice.At(j).LogRecords()
			for k := 0; k < records.Len(); k++ {
				record := records.At(k)
				logRecords = append(logRecords, record)
			}
		}
	}

	return logRecords
}

// firstTraceRecord extracts the only span record from one-record traces created by someTraces()
func allTraceRecords(traces ptrace.Traces) []ptrace.Span {
	var spanRecords []ptrace.Span

	resourceSpansSlice := traces.ResourceSpans()
	for i := 0; i < resourceSpansSlice.Len(); i++ {
		instrumentationSpansSlice := resourceSpansSlice.At(i).ScopeSpans()
		for j := 0; j < instrumentationSpansSlice.Len(); j++ {
			records := instrumentationSpansSlice.At(j).Spans()
			for k := 0; k < records.Len(); k++ {
				record := records.At(k)
				spanRecords = append(spanRecords, record)
			}
		}
	}

	return spanRecords
}

func expectEngineCalls(engine *mocks.MockEngineAPI, checkResponse *flowcontrolv1.CheckResponse) {
	expectedCalls := make([]*gomock.Call, len(checkResponse.FluxMeters))
	for i, fm := range checkResponse.FluxMeters {
		// TODO actually return some Histogram
		fmID := fmt.Sprintf(
			"agent_group-%v-policy-%v-flux_meter-%v-policy_hash-%v",
			fm.GetAgentGroup(),
			fm.GetPolicyName(),
			fm.GetFluxMeterName(),
			fm.GetPolicyHash(),
		)
		expectedCalls[i] = engine.EXPECT().GetFluxMeterHist(fmID).Return(nil)
	}
	gomock.InOrder(expectedCalls...)
}
