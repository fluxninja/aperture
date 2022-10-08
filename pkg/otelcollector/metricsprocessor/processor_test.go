package metricsprocessor

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"go.opentelemetry.io/collector/pdata/plog"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/timestamppb"
	"k8s.io/apimachinery/pkg/util/json"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	m "github.com/fluxninja/aperture/pkg/metrics"
	oc "github.com/fluxninja/aperture/pkg/otelcollector"
	"github.com/fluxninja/aperture/pkg/policies/mocks"
)

var _ = Describe("Metrics Processor", func() {
	var (
		pr                *prometheus.Registry
		cfg               *Config
		processor         *metricsProcessor
		engine            *mocks.MockEngine
		clasEngine        *mocks.MockClassificationEngine
		conLimiter        *mocks.MockConcurrencyLimiter
		rateLimiter       *mocks.MockRateLimiter
		classifier        *mocks.MockClassifier
		summaryVec        *prometheus.SummaryVec
		rateCounter       prometheus.Counter
		classifierCounter prometheus.Counter
		baseCheckResp     *flowcontrolv1.CheckResponse
		labelsFoo1        map[string]string
		labelsFizz1       map[string]string
		labelsFizz2       map[string]string
		expectedLabels    map[string]interface{}
		expectedMetrics   string
		source            string
	)

	BeforeEach(func() {
		pr = prometheus.NewRegistry()
		ctrl := gomock.NewController(GinkgoT())
		engine = mocks.NewMockEngine(ctrl)
		clasEngine = mocks.NewMockClassificationEngine(ctrl)
		conLimiter = mocks.NewMockConcurrencyLimiter(ctrl)
		rateLimiter = mocks.NewMockRateLimiter(ctrl)
		classifier = mocks.NewMockClassifier(ctrl)
		expectedLabels = make(map[string]interface{})
		cfg = &Config{
			engine:               engine,
			classificationEngine: clasEngine,
			promRegistry:         pr,
		}

		summaryVec = prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Name: m.WorkloadLatencyMetricName,
			Help: "dummy",
		}, []string{m.PolicyNameLabel, m.PolicyHashLabel, m.ComponentIndexLabel,
			m.DecisionTypeLabel, m.WorkloadIndexLabel,
		})
		rateCounter = prometheus.NewCounter(prometheus.CounterOpts{
			Name: m.RateLimiterCounterMetricName,
			Help: "dummy",
			ConstLabels: prometheus.Labels{
				m.PolicyNameLabel:     "foo",
				m.PolicyHashLabel:     "foo-hash",
				m.ComponentIndexLabel: "2",
			},
		})
		classifierCounter = prometheus.NewCounter(prometheus.CounterOpts{
			Name: m.ClassifierCounterMetricName,
			Help: "dummy",
			ConstLabels: prometheus.Labels{
				m.PolicyNameLabel:      "foo",
				m.PolicyHashLabel:      "foo-hash",
				m.ClassifierIndexLabel: "1",
			},
		})

		var err error
		processor, err = newProcessor(cfg)
		Expect(err).NotTo(HaveOccurred())

		start := time.Date(1969, time.Month(7), 20, 17, 0, 0, 0, time.UTC)
		end := time.Date(1969, time.Month(7), 20, 17, 0, 1, 0, time.UTC)
		baseCheckResp = &flowcontrolv1.CheckResponse{
			Start: timestamppb.New(start),
			End:   timestamppb.New(end),
			ControlPointInfo: &flowcontrolv1.ControlPointInfo{
				Type: flowcontrolv1.ControlPointInfo_TYPE_INGRESS,
			},
			DecisionType: flowcontrolv1.CheckResponse_DECISION_TYPE_REJECTED,
			LimiterDecisions: []*flowcontrolv1.LimiterDecision{
				{
					PolicyName:     "foo",
					PolicyHash:     "foo-hash",
					ComponentIndex: 1,
					Dropped:        true,
					Details: &flowcontrolv1.LimiterDecision_ConcurrencyLimiterInfo_{
						ConcurrencyLimiterInfo: &flowcontrolv1.LimiterDecision_ConcurrencyLimiterInfo{
							WorkloadIndex: "0",
						},
					},
				},
			},
		}
		labelsFoo1 = map[string]string{
			m.PolicyNameLabel:     "foo",
			m.PolicyHashLabel:     "foo-hash",
			m.ComponentIndexLabel: "1",
			m.DecisionTypeLabel:   flowcontrolv1.CheckResponse_DECISION_TYPE_REJECTED.String(),
			m.WorkloadIndexLabel:  "0",
		}
		labelsFizz1 = map[string]string{
			m.PolicyNameLabel:     "fizz",
			m.PolicyHashLabel:     "fizz-hash",
			m.ComponentIndexLabel: "1",
			m.DecisionTypeLabel:   flowcontrolv1.CheckResponse_DECISION_TYPE_REJECTED.String(),
			m.WorkloadIndexLabel:  "1",
		}
		labelsFizz2 = map[string]string{
			m.PolicyNameLabel:     "fizz",
			m.PolicyHashLabel:     "fizz-hash",
			m.ComponentIndexLabel: "2",
			m.DecisionTypeLabel:   flowcontrolv1.CheckResponse_DECISION_TYPE_REJECTED.String(),
			m.WorkloadIndexLabel:  "2",
		}
		engine.EXPECT().GetConcurrencyLimiter(gomock.Any()).Return(conLimiter).AnyTimes()
		engine.EXPECT().GetRateLimiter(gomock.Any()).Return(rateLimiter).AnyTimes()
		clasEngine.EXPECT().GetClassifier(gomock.Any()).Return(classifier, nil).AnyTimes()
	})

	AfterEach(func() {
		logs := someLogs(engine, baseCheckResp, source)
		modifiedLogs, err := processor.ConsumeLogs(context.Background(), logs)
		Expect(err).NotTo(HaveOccurred())
		Expect(modifiedLogs).To(Equal(logs))

		By("sending proper metrics")
		splitMetrics := strings.Split(expectedMetrics, "<split>")
		for _, expectedMetrics := range splitMetrics {
			if strings.Contains(expectedMetrics, m.WorkloadLatencyMetricName) {
				expected := strings.NewReader(expectedMetrics)
				err = testutil.CollectAndCompare(summaryVec, expected, m.WorkloadLatencyMetricName)
				Expect(err).NotTo(HaveOccurred())
			}

			if strings.Contains(expectedMetrics, m.RateLimiterCounterMetricName) {
				expected2 := strings.NewReader(expectedMetrics)
				err = testutil.CollectAndCompare(rateCounter, expected2, m.RateLimiterCounterMetricName)
				Expect(err).NotTo(HaveOccurred())
			}
		}

		By("adding proper labels")
		logRecords := allLogRecords(modifiedLogs)
		Expect(logRecords).To(HaveLen(1))

		for k, v := range expectedLabels {
			Expect(logRecords[0].Attributes().AsRaw()).To(HaveKeyWithValue(k, v))
		}
	})

	It("Processes logs for single policy - ingress", func() {
		rateLimiterDecision := &flowcontrolv1.LimiterDecision{
			PolicyName:     "foo",
			PolicyHash:     "foo-hash",
			ComponentIndex: 2,
			Dropped:        true,
			Details: &flowcontrolv1.LimiterDecision_RateLimiterInfo_{
				RateLimiterInfo: &flowcontrolv1.LimiterDecision_RateLimiterInfo{
					Remaining: 1,
					Current:   1,
					Label:     "test",
				}}}
		baseCheckResp.LimiterDecisions = append(baseCheckResp.LimiterDecisions, rateLimiterDecision)
		baseCheckResp.ClassifierInfos = []*flowcontrolv1.ClassifierInfo{{
			PolicyName:      "foo",
			PolicyHash:      "foo-hash",
			ClassifierIndex: 1,
			Error:           flowcontrolv1.ClassifierInfo_ERROR_EMPTY_RESULTSET,
		}}
		baseCheckResp.FluxMeterInfos = []*flowcontrolv1.FluxMeterInfo{{FluxMeterName: "bar"}}
		baseCheckResp.FlowLabelKeys = []string{"someLabel"}
		baseCheckResp.Services = []string{"svc1", "svc2"}

		// <split> is a workaround until PR https://github.com/prometheus/client_golang/pull/1143 is released
		expectedMetrics = `# HELP classifier_counter dummy
# TYPE classifier_counter counter
classifier_counter{component_index="1",policy_hash="foo-hash",policy_name="foo"} 1
<split># HELP rate_limiter_counter dummy
# TYPE rate_limiter_counter counter
rate_limiter_counter{component_index="2",policy_hash="foo-hash",policy_name="foo"} 1
<split># HELP workload_latency_ms dummy
# TYPE workload_latency_ms summary
workload_latency_ms_sum{component_index="1",decision_type="DECISION_TYPE_REJECTED",policy_hash="foo-hash",policy_name="foo",workload_index="0"} 5
workload_latency_ms_count{component_index="1",decision_type="DECISION_TYPE_REJECTED",policy_hash="foo-hash",policy_name="foo",workload_index="0"} 1
`

		expectedLabels = map[string]interface{}{
			oc.ApertureDecisionTypeLabel: flowcontrolv1.CheckResponse_DECISION_TYPE_REJECTED.String(),
			oc.ApertureErrorLabel:        flowcontrolv1.CheckResponse_ERROR_NONE.String(),
			oc.ApertureRejectReasonLabel: flowcontrolv1.CheckResponse_REJECT_REASON_NONE.String(),
			oc.ApertureClassifiersLabel:  []interface{}{"policy_name:foo,classifier_index:1"},

			oc.ApertureClassifierErrorsLabel: []interface{}{fmt.Sprintf("%s,policy_name:foo,classifier_index:1,policy_hash:foo-hash",
				flowcontrolv1.ClassifierInfo_ERROR_EMPTY_RESULTSET.String())},

			oc.ApertureFluxMetersLabel:                  []interface{}{"bar"},
			oc.ApertureFlowLabelKeysLabel:               []interface{}{"someLabel"},
			oc.ApertureRateLimitersLabel:                []interface{}{"policy_name:foo,component_index:2,policy_hash:foo-hash"},
			oc.ApertureDroppingRateLimitersLabel:        []interface{}{"policy_name:foo,component_index:2,policy_hash:foo-hash"},
			oc.ApertureConcurrencyLimitersLabel:         []interface{}{"policy_name:foo,component_index:1,policy_hash:foo-hash"},
			oc.ApertureDroppingConcurrencyLimitersLabel: []interface{}{"policy_name:foo,component_index:1,policy_hash:foo-hash"},
			oc.ApertureWorkloadsLabel:                   []interface{}{"policy_name:foo,component_index:1,workload_index:0,policy_hash:foo-hash"},
			oc.ApertureDroppingWorkloadsLabel:           []interface{}{"policy_name:foo,component_index:1,workload_index:0,policy_hash:foo-hash"},

			oc.ApertureProcessingDurationLabel: float64(1000),
			oc.ApertureServicesLabel:           []interface{}{"svc1", "svc2"},
			oc.ApertureControlPointLabel:       "type:TYPE_INGRESS",
		}
		source = oc.ApertureSourceEnvoy

		summary, err := summaryVec.GetMetricWith(labelsFoo1)
		Expect(err).NotTo(HaveOccurred())
		conLimiter.EXPECT().GetObserver(labelsFoo1).Return(summary).Times(1)

		rateLimiter.EXPECT().GetCounter().Return(rateCounter).Times(1)
		classifier.EXPECT().GetCounter().Return(classifierCounter).Times(1)
	})

	It("Processes logs for single policy - feature", func() {
		baseCheckResp.ControlPointInfo = &flowcontrolv1.ControlPointInfo{
			Type:    flowcontrolv1.ControlPointInfo_TYPE_FEATURE,
			Feature: "featureX",
		}
		baseCheckResp.RejectReason = flowcontrolv1.CheckResponse_REJECT_REASON_RATE_LIMITED

		expectedMetrics = `# HELP workload_latency_ms dummy
# TYPE workload_latency_ms summary
workload_latency_ms_sum{component_index="1",decision_type="DECISION_TYPE_REJECTED",policy_hash="foo-hash",policy_name="foo",workload_index="0"} 5
workload_latency_ms_count{component_index="1",decision_type="DECISION_TYPE_REJECTED",policy_hash="foo-hash",policy_name="foo",workload_index="0"} 1
`

		expectedLabels = map[string]interface{}{
			oc.ApertureDecisionTypeLabel:                "DECISION_TYPE_REJECTED",
			oc.ApertureRejectReasonLabel:                "REJECT_REASON_RATE_LIMITED",
			oc.ApertureRateLimitersLabel:                []interface{}{},
			oc.ApertureDroppingRateLimitersLabel:        []interface{}{},
			oc.ApertureConcurrencyLimitersLabel:         []interface{}{"policy_name:foo,component_index:1,policy_hash:foo-hash"},
			oc.ApertureDroppingConcurrencyLimitersLabel: []interface{}{"policy_name:foo,component_index:1,policy_hash:foo-hash"},
			oc.ApertureWorkloadsLabel:                   []interface{}{"policy_name:foo,component_index:1,workload_index:0,policy_hash:foo-hash"},
			oc.ApertureDroppingWorkloadsLabel:           []interface{}{"policy_name:foo,component_index:1,workload_index:0,policy_hash:foo-hash"},
		}
		source = oc.ApertureSourceSDK

		summary, err := summaryVec.GetMetricWith(labelsFoo1)
		Expect(err).NotTo(HaveOccurred())
		conLimiter.EXPECT().GetObserver(labelsFoo1).Return(summary).Times(1)
	})

	It("Processes logs for two policies - ingress", func() {
		baseCheckResp.LimiterDecisions = append(baseCheckResp.LimiterDecisions,
			&flowcontrolv1.LimiterDecision{
				PolicyName:     "fizz",
				PolicyHash:     "fizz-hash",
				ComponentIndex: 1,
				Dropped:        true,
				Details: &flowcontrolv1.LimiterDecision_ConcurrencyLimiterInfo_{
					ConcurrencyLimiterInfo: &flowcontrolv1.LimiterDecision_ConcurrencyLimiterInfo{
						WorkloadIndex: "1",
					},
				},
			},
			&flowcontrolv1.LimiterDecision{
				PolicyName:     "fizz",
				PolicyHash:     "fizz-hash",
				ComponentIndex: 2,
				Dropped:        false,
				Details: &flowcontrolv1.LimiterDecision_ConcurrencyLimiterInfo_{
					ConcurrencyLimiterInfo: &flowcontrolv1.LimiterDecision_ConcurrencyLimiterInfo{
						WorkloadIndex: "2",
					},
				},
			})

		expectedMetrics = `# HELP workload_latency_ms dummy
# TYPE workload_latency_ms summary
workload_latency_ms_sum{component_index="1",decision_type="DECISION_TYPE_REJECTED",policy_hash="fizz-hash",policy_name="fizz",workload_index="1"} 5
workload_latency_ms_count{component_index="1",decision_type="DECISION_TYPE_REJECTED",policy_hash="fizz-hash",policy_name="fizz",workload_index="1"} 1
workload_latency_ms_sum{component_index="1",decision_type="DECISION_TYPE_REJECTED",policy_hash="foo-hash",policy_name="foo",workload_index="0"} 5
workload_latency_ms_count{component_index="1",decision_type="DECISION_TYPE_REJECTED",policy_hash="foo-hash",policy_name="foo",workload_index="0"} 1
workload_latency_ms_sum{component_index="2",decision_type="DECISION_TYPE_REJECTED",policy_hash="fizz-hash",policy_name="fizz",workload_index="2"} 5
workload_latency_ms_count{component_index="2",decision_type="DECISION_TYPE_REJECTED",policy_hash="fizz-hash",policy_name="fizz",workload_index="2"} 1
`

		expectedLabels = map[string]interface{}{
			oc.ApertureDecisionTypeLabel:         flowcontrolv1.CheckResponse_DECISION_TYPE_REJECTED.String(),
			oc.ApertureErrorLabel:                flowcontrolv1.CheckResponse_ERROR_NONE.String(),
			oc.ApertureRateLimitersLabel:         []interface{}{},
			oc.ApertureDroppingRateLimitersLabel: []interface{}{},
			oc.ApertureConcurrencyLimitersLabel: []interface{}{
				"policy_name:foo,component_index:1,policy_hash:foo-hash",
				"policy_name:fizz,component_index:1,policy_hash:fizz-hash",
				"policy_name:fizz,component_index:2,policy_hash:fizz-hash",
			},
			oc.ApertureDroppingConcurrencyLimitersLabel: []interface{}{
				"policy_name:foo,component_index:1,policy_hash:foo-hash",
				"policy_name:fizz,component_index:1,policy_hash:fizz-hash",
			},
			oc.ApertureWorkloadsLabel: []interface{}{
				"policy_name:foo,component_index:1,workload_index:0,policy_hash:foo-hash",
				"policy_name:fizz,component_index:1,workload_index:1,policy_hash:fizz-hash",
				"policy_name:fizz,component_index:2,workload_index:2,policy_hash:fizz-hash",
			},
			oc.ApertureDroppingWorkloadsLabel: []interface{}{
				"policy_name:foo,component_index:1,workload_index:0,policy_hash:foo-hash",
				"policy_name:fizz,component_index:1,workload_index:1,policy_hash:fizz-hash",
			},
		}
		source = oc.ApertureSourceEnvoy

		summaryFoo, err := summaryVec.GetMetricWith(labelsFoo1)
		Expect(err).NotTo(HaveOccurred())
		conLimiter.EXPECT().GetObserver(labelsFoo1).Return(summaryFoo).Times(1)

		summaryFizz1, err := summaryVec.GetMetricWith(labelsFizz1)
		Expect(err).NotTo(HaveOccurred())
		conLimiter.EXPECT().GetObserver(labelsFizz1).Return(summaryFizz1).Times(1)

		summaryFizz2, err := summaryVec.GetMetricWith(labelsFizz2)
		Expect(err).NotTo(HaveOccurred())
		conLimiter.EXPECT().GetObserver(labelsFizz2).Return(summaryFizz2).Times(1)

	})
})

// someLogs will return a plog.Logs instance with single LogRecord
func someLogs(
	engine *mocks.MockEngine,
	checkResponse *flowcontrolv1.CheckResponse,
	source string,
) plog.Logs {
	logs := plog.NewLogs()
	logs.ResourceLogs().AppendEmpty()

	expectedCalls := make([]*gomock.Call, len(checkResponse.FluxMeterInfos))
	resourceLogsSlice := logs.ResourceLogs()
	for i := 0; i < resourceLogsSlice.Len(); i++ {
		resourceLogsSlice.At(i).ScopeLogs().AppendEmpty()

		instrumentationLogsSlice := resourceLogsSlice.At(i).ScopeLogs()
		for j := 0; j < instrumentationLogsSlice.Len(); j++ {
			logRecord := instrumentationLogsSlice.At(j).LogRecords().AppendEmpty()
			marshalledCheckResponse, err := json.Marshal(checkResponse)
			Expect(err).NotTo(HaveOccurred())
			logRecord.Attributes().InsertString(oc.ApertureSourceLabel, source)
			logRecord.Attributes().InsertString(oc.ApertureCheckResponseLabel, string(marshalledCheckResponse))
			logRecord.Attributes().InsertString(oc.HTTPStatusCodeLabel, "201")
			logRecord.Attributes().InsertDouble(oc.WorkloadDurationLabel, 5)
			logRecord.Attributes().InsertDouble(oc.EnvoyAuthzDurationLabel, 1)
			for i, fm := range checkResponse.FluxMeterInfos {
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
