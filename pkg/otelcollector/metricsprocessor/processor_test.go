package metricsprocessor

import (
	"strings"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"go.opentelemetry.io/collector/pdata/plog"
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
		engine    *mocks.MockEngine
	)

	BeforeEach(func() {
		pr = prometheus.NewRegistry()
		ctrl := gomock.NewController(GinkgoT())
		engine = mocks.NewMockEngine(ctrl)
		cfg = &Config{
			engine:       engine,
			promRegistry: pr,
		}
		var err error
		processor, err = newProcessor(cfg)
		Expect(err).NotTo(HaveOccurred())
	})

	DescribeTable("Processing logs",
		func(
			checkResponse *flowcontrolv1.CheckResponse,
			expectedErr error,
			expectedMetrics string,
			expectedLabels map[string]interface{},
		) {
			ctx := context.Background()

			logs := someLogs(engine, checkResponse)
			modifiedLogs, err := processor.ConsumeLogs(ctx, logs)
			if expectedErr != nil {
				Expect(err).NotTo(MatchError(expectedErr))
				return
			}
			Expect(err).NotTo(HaveOccurred())
			Expect(modifiedLogs).To(Equal(logs))

			By("sending proper metrics")
			expected := strings.NewReader(expectedMetrics)
			err = testutil.CollectAndCompare(processor.workloadLatencySummary, expected, "workload_latency_ms")
			Expect(err).NotTo(HaveOccurred())

			By("adding proper labels")
			logRecords := allLogRecords(modifiedLogs)
			Expect(logRecords).To(HaveLen(1))

			for k, v := range expectedLabels {
				Expect(logRecords[0].Attributes().AsRaw()).To(HaveKeyWithValue(k, v))
			}
		},

		Entry("record with single policy - ingress",
			&flowcontrolv1.CheckResponse{
				ControlPoint: &flowcontrolv1.ControlPoint{
					Type: flowcontrolv1.ControlPoint_TYPE_INGRESS,
				},
				DecisionType: flowcontrolv1.CheckResponse_DECISION_TYPE_REJECTED,
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
				Classifiers: []*flowcontrolv1.Classifier{
					{
						PolicyName:      "foo",
						PolicyHash:      "foo-hash",
						ClassifierIndex: 1,
					},
				},
				FluxMeters: []*flowcontrolv1.FluxMeter{
					{
						FluxMeterName: "bar",
					},
				},
				FlowLabelKeys: []string{
					"someLabel",
				},
			},
			nil,
			`# HELP workload_latency_ms Latency summary of workload
			# TYPE workload_latency_ms summary
			workload_latency_ms_sum{component_index="1",decision_type="DECISION_TYPE_REJECTED",policy_hash="foo-hash",policy_name="foo",workload_index="0"} 5
			workload_latency_ms_count{component_index="1",decision_type="DECISION_TYPE_REJECTED",policy_hash="foo-hash",policy_name="foo",workload_index="0"} 1
			`,
			map[string]interface{}{
				otelcollector.ApertureDecisionTypeLabel:                "DECISION_TYPE_REJECTED",
				otelcollector.ApertureErrorLabel:                       "",
				otelcollector.ApertureRejectReasonLabel:                "",
				otelcollector.ApertureClassifiersLabel:                 []interface{}{"policy_name:foo,classifier_index:1"},
				otelcollector.ApertureFluxMetersLabel:                  []interface{}{"bar"},
				otelcollector.ApertureFlowLabelKeysLabel:               []interface{}{"someLabel"},
				otelcollector.ApertureRateLimitersLabel:                []interface{}{},
				otelcollector.ApertureDroppingRateLimitersLabel:        []interface{}{},
				otelcollector.ApertureConcurrencyLimitersLabel:         []interface{}{"policy_name:foo,component_index:1,policy_hash:foo-hash"},
				otelcollector.ApertureDroppingConcurrencyLimitersLabel: []interface{}{"policy_name:foo,component_index:1,policy_hash:foo-hash"},
				otelcollector.ApertureWorkloadsLabel:                   []interface{}{"policy_name:foo,component_index:1,workload_index:0,policy_hash:foo-hash"},
				otelcollector.ApertureDroppingWorkloadsLabel:           []interface{}{"policy_name:foo,component_index:1,workload_index:0,policy_hash:foo-hash"},
			},
		),

		Entry("record with single policy - feature",
			&flowcontrolv1.CheckResponse{
				ControlPoint: &flowcontrolv1.ControlPoint{
					Type:    flowcontrolv1.ControlPoint_TYPE_FEATURE,
					Feature: "featureX",
				},
				DecisionType: flowcontrolv1.CheckResponse_DECISION_TYPE_REJECTED,
				RejectReason: flowcontrolv1.CheckResponse_REJECT_REASON_RATE_LIMITED,
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
				FluxMeters:    []*flowcontrolv1.FluxMeter{},
				FlowLabelKeys: []string{},
			},
			nil,
			`# HELP workload_latency_ms Latency summary of workload
			# TYPE workload_latency_ms summary
			workload_latency_ms_sum{component_index="1",decision_type="DECISION_TYPE_REJECTED",policy_hash="foo-hash",policy_name="foo",workload_index="0"} 5
			workload_latency_ms_count{component_index="1",decision_type="DECISION_TYPE_REJECTED",policy_hash="foo-hash",policy_name="foo",workload_index="0"} 1
			`,
			map[string]interface{}{
				otelcollector.ApertureDecisionTypeLabel:                "DECISION_TYPE_REJECTED",
				otelcollector.ApertureRejectReasonLabel:                "REJECT_REASON_RATE_LIMITED",
				otelcollector.ApertureRateLimitersLabel:                []interface{}{},
				otelcollector.ApertureDroppingRateLimitersLabel:        []interface{}{},
				otelcollector.ApertureConcurrencyLimitersLabel:         []interface{}{"policy_name:foo,component_index:1,policy_hash:foo-hash"},
				otelcollector.ApertureDroppingConcurrencyLimitersLabel: []interface{}{"policy_name:foo,component_index:1,policy_hash:foo-hash"},
				otelcollector.ApertureWorkloadsLabel:                   []interface{}{"policy_name:foo,component_index:1,workload_index:0,policy_hash:foo-hash"},
				otelcollector.ApertureDroppingWorkloadsLabel:           []interface{}{"policy_name:foo,component_index:1,workload_index:0,policy_hash:foo-hash"},
			},
		),

		Entry("record with two policies",
			&flowcontrolv1.CheckResponse{
				ControlPoint: &flowcontrolv1.ControlPoint{
					Type: flowcontrolv1.ControlPoint_TYPE_INGRESS,
				},
				DecisionType: flowcontrolv1.CheckResponse_DECISION_TYPE_REJECTED,
				RejectReason: flowcontrolv1.CheckResponse_REJECT_REASON_NONE,
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
				FluxMeters:    []*flowcontrolv1.FluxMeter{},
				FlowLabelKeys: []string{},
			},
			nil,
			`# HELP workload_latency_ms Latency summary of workload
			# TYPE workload_latency_ms summary
			workload_latency_ms_sum{component_index="1",decision_type="DECISION_TYPE_REJECTED",policy_hash="fizz-hash",policy_name="fizz",workload_index="1"} 5
			workload_latency_ms_count{component_index="1",decision_type="DECISION_TYPE_REJECTED",policy_hash="fizz-hash",policy_name="fizz",workload_index="1"} 1
			workload_latency_ms_sum{component_index="1",decision_type="DECISION_TYPE_REJECTED",policy_hash="foo-hash",policy_name="foo",workload_index="0"} 5
			workload_latency_ms_count{component_index="1",decision_type="DECISION_TYPE_REJECTED",policy_hash="foo-hash",policy_name="foo",workload_index="0"} 1
			workload_latency_ms_sum{component_index="2",decision_type="DECISION_TYPE_REJECTED",policy_hash="fizz-hash",policy_name="fizz",workload_index="2"} 5
			workload_latency_ms_count{component_index="2",decision_type="DECISION_TYPE_REJECTED",policy_hash="fizz-hash",policy_name="fizz",workload_index="2"} 1
			`,
			map[string]interface{}{
				otelcollector.ApertureDecisionTypeLabel:         "DECISION_TYPE_REJECTED",
				otelcollector.ApertureErrorLabel:                "ERROR_REASON_UNSPECIFIED",
				otelcollector.ApertureRateLimitersLabel:         []interface{}{},
				otelcollector.ApertureDroppingRateLimitersLabel: []interface{}{},
				otelcollector.ApertureConcurrencyLimitersLabel: []interface{}{
					"policy_name:foo,component_index:1,policy_hash:foo-hash",
					"policy_name:fizz,component_index:1,policy_hash:fizz-hash",
					"policy_name:fizz,component_index:2,policy_hash:fizz-hash",
				},
				otelcollector.ApertureDroppingConcurrencyLimitersLabel: []interface{}{
					"policy_name:foo,component_index:1,policy_hash:foo-hash",
					"policy_name:fizz,component_index:1,policy_hash:fizz-hash",
				},
				otelcollector.ApertureWorkloadsLabel: []interface{}{
					"policy_name:foo,component_index:1,workload_index:0,policy_hash:foo-hash",
					"policy_name:fizz,component_index:1,workload_index:1,policy_hash:fizz-hash",
					"policy_name:fizz,component_index:2,workload_index:2,policy_hash:fizz-hash",
				},
				otelcollector.ApertureDroppingWorkloadsLabel: []interface{}{
					"policy_name:foo,component_index:1,workload_index:0,policy_hash:foo-hash",
					"policy_name:fizz,component_index:1,workload_index:1,policy_hash:fizz-hash",
				},
			},
		),
	)
})

// someLogs will return a plog.Logs instance with single LogRecord
func someLogs(
	engine *mocks.MockEngine,
	checkResponse *flowcontrolv1.CheckResponse,
) plog.Logs {
	logs := plog.NewLogs()
	logs.ResourceLogs().AppendEmpty()

	expectedCalls := make([]*gomock.Call, len(checkResponse.FluxMeters))
	resourceLogsSlice := logs.ResourceLogs()
	for i := 0; i < resourceLogsSlice.Len(); i++ {
		resourceLogsSlice.At(i).ScopeLogs().AppendEmpty()

		instrumentationLogsSlice := resourceLogsSlice.At(i).ScopeLogs()
		for j := 0; j < instrumentationLogsSlice.Len(); j++ {
			logRecord := instrumentationLogsSlice.At(j).LogRecords().AppendEmpty()
			marshalledCheckResponse, err := json.Marshal(checkResponse)
			Expect(err).NotTo(HaveOccurred())
			logRecord.Attributes().InsertString(otelcollector.ApertureCheckResponseLabel, string(marshalledCheckResponse))
			logRecord.Attributes().InsertString(otelcollector.HTTPStatusCodeLabel, "201")
			logRecord.Attributes().InsertString(otelcollector.WorkloadDurationLabel, "5")
			for i, fm := range checkResponse.FluxMeters {
				// TODO actually return some Histogram
				expectedCalls[i] = engine.EXPECT().GetFluxMeter(fm.GetFluxMeterName()).Return(nil)
			}
		}
	}
	gomock.InOrder(expectedCalls...)

	return logs
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
