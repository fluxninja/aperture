package metricsprocessor

import (
	"strings"

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
)

var _ = Describe("Metrics Processor", func() {
	var (
		pr        *prometheus.Registry
		cfg       *Config
		processor *metricsProcessor
	)

	BeforeEach(func() {
		pr = prometheus.NewRegistry()
		cfg = &Config{
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
		func(controlPoint string, decisions []*flowcontrolv1.LimiterDecision, expectedErr error, expectedMetrics string, expectedLabels map[string]string) {
			ctx := context.Background()

			logs := someLogs(decisions, controlPoint)
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
				Expect(logRecords[0].Attributes().AsRaw()).To(HaveKeyWithValue(k, MatchJSON(v)))
			}
		},

		Entry("record with single policy - ingress",
			otelcollector.ControlPointIngress,
			[]*flowcontrolv1.LimiterDecision{
				{
					PolicyName:     "foo",
					PolicyHash:     "foo-hash",
					ComponentIndex: 1,
					Dropped:        true,
					Details: &flowcontrolv1.LimiterDecision_ConcurrencyLimiter_{
						ConcurrencyLimiter: &flowcontrolv1.LimiterDecision_ConcurrencyLimiter{
							Workload: "workload_key:\"foo\", workload_value:\"bar\"",
						},
					},
				},
			},
			nil,
			`# HELP workload_latency_ms Latency histogram of workload
			# TYPE workload_latency_ms histogram
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="workload_key:\"foo\", workload_value:\"bar\"",le="0"} 0
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="workload_key:\"foo\", workload_value:\"bar\"",le="10"} 1
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="workload_key:\"foo\", workload_value:\"bar\"",le="20"} 1
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="workload_key:\"foo\", workload_value:\"bar\"",le="+Inf"} 1
			workload_latency_ms_sum{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="workload_key:\"foo\", workload_value:\"bar\""} 5
			workload_latency_ms_count{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="workload_key:\"foo\", workload_value:\"bar\""} 1
			`,
			map[string]string{
				"rate_limiters":                 `[]`,
				"dropping_rate_limiters":        `[]`,
				"concurrency_limiters":          `["policy_name:foo,component_index:1,workload_index:workload_key:\"foo\", workload_value:\"bar\",policy_hash:foo-hash"]`,
				"dropping_concurrency_limiters": `["policy_name:foo,component_index:1,workload_index:workload_key:\"foo\", workload_value:\"bar\",policy_hash:foo-hash"]`,
			},
		),

		Entry("record with single policy - feature",
			otelcollector.ControlPointFeature,
			[]*flowcontrolv1.LimiterDecision{
				{
					PolicyName:     "foo",
					PolicyHash:     "foo-hash",
					ComponentIndex: 1,
					Dropped:        true,
					Details: &flowcontrolv1.LimiterDecision_ConcurrencyLimiter_{
						ConcurrencyLimiter: &flowcontrolv1.LimiterDecision_ConcurrencyLimiter{
							Workload: "workload_key:\"foo\", workload_value:\"bar\"",
						},
					},
				},
			},
			nil,
			`# HELP workload_latency_ms Latency histogram of workload
			# TYPE workload_latency_ms histogram
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="workload_key:\"foo\", workload_value:\"bar\"",le="0"} 0
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="workload_key:\"foo\", workload_value:\"bar\"",le="10"} 1
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="workload_key:\"foo\", workload_value:\"bar\"",le="20"} 1
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="workload_key:\"foo\", workload_value:\"bar\"",le="+Inf"} 1
			workload_latency_ms_sum{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="workload_key:\"foo\", workload_value:\"bar\""} 5
			workload_latency_ms_count{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="workload_key:\"foo\", workload_value:\"bar\""} 1
			`,
			map[string]string{
				"rate_limiters":                 `[]`,
				"dropping_rate_limiters":        `[]`,
				"concurrency_limiters":          `["policy_name:foo,component_index:1,workload_index:workload_key:\"foo\", workload_value:\"bar\",policy_hash:foo-hash"]`,
				"dropping_concurrency_limiters": `["policy_name:foo,component_index:1,workload_index:workload_key:\"foo\", workload_value:\"bar\",policy_hash:foo-hash"]`,
			},
		),

		Entry("record with two policies",
			otelcollector.ControlPointIngress,
			[]*flowcontrolv1.LimiterDecision{
				{
					PolicyName:     "foo",
					PolicyHash:     "foo-hash",
					ComponentIndex: 1,
					Dropped:        true,
					Details: &flowcontrolv1.LimiterDecision_ConcurrencyLimiter_{
						ConcurrencyLimiter: &flowcontrolv1.LimiterDecision_ConcurrencyLimiter{
							Workload: "workload_key:\"foo\", workload_value:\"bar\"",
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
							Workload: "workload_key:\"fizz\", workload_value:\"buzz\"",
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
							Workload: "workload_key:\"fizz\", workload_value:\"hoge\"",
						},
					},
				},
			},
			nil,
			`# HELP workload_latency_ms Latency histogram of workload
			# TYPE workload_latency_ms histogram
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="fizz-hash",policy_name="fizz",workload_index="workload_key:\"fizz\", workload_value:\"buzz\"",le="0"} 0
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="fizz-hash",policy_name="fizz",workload_index="workload_key:\"fizz\", workload_value:\"buzz\"",le="10"} 1
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="fizz-hash",policy_name="fizz",workload_index="workload_key:\"fizz\", workload_value:\"buzz\"",le="20"} 1
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="fizz-hash",policy_name="fizz",workload_index="workload_key:\"fizz\", workload_value:\"buzz\"",le="+Inf"} 1
			workload_latency_ms_sum{component_index="1",dropped="true",policy_hash="fizz-hash",policy_name="fizz",workload_index="workload_key:\"fizz\", workload_value:\"buzz\""} 5
			workload_latency_ms_count{component_index="1",dropped="true",policy_hash="fizz-hash",policy_name="fizz",workload_index="workload_key:\"fizz\", workload_value:\"buzz\""} 1
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="workload_key:\"foo\", workload_value:\"bar\"",le="0"} 0
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="workload_key:\"foo\", workload_value:\"bar\"",le="10"} 1
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="workload_key:\"foo\", workload_value:\"bar\"",le="20"} 1
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="workload_key:\"foo\", workload_value:\"bar\"",le="+Inf"} 1
			workload_latency_ms_sum{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="workload_key:\"foo\", workload_value:\"bar\""} 5
			workload_latency_ms_count{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="workload_key:\"foo\", workload_value:\"bar\""} 1
			workload_latency_ms_bucket{component_index="2",dropped="false",policy_hash="fizz-hash",policy_name="fizz",workload_index="workload_key:\"fizz\", workload_value:\"hoge\"",le="0"} 0
			workload_latency_ms_bucket{component_index="2",dropped="false",policy_hash="fizz-hash",policy_name="fizz",workload_index="workload_key:\"fizz\", workload_value:\"hoge\"",le="10"} 1
			workload_latency_ms_bucket{component_index="2",dropped="false",policy_hash="fizz-hash",policy_name="fizz",workload_index="workload_key:\"fizz\", workload_value:\"hoge\"",le="20"} 1
			workload_latency_ms_bucket{component_index="2",dropped="false",policy_hash="fizz-hash",policy_name="fizz",workload_index="workload_key:\"fizz\", workload_value:\"hoge\"",le="+Inf"} 1
			workload_latency_ms_sum{component_index="2",dropped="false",policy_hash="fizz-hash",policy_name="fizz",workload_index="workload_key:\"fizz\", workload_value:\"hoge\""} 5
			workload_latency_ms_count{component_index="2",dropped="false",policy_hash="fizz-hash",policy_name="fizz",workload_index="workload_key:\"fizz\", workload_value:\"hoge\""} 1
			`,
			map[string]string{
				"rate_limiters":          `[]`,
				"dropping_rate_limiters": `[]`,
				"concurrency_limiters": `[
					"policy_name:foo,component_index:1,workload_index:workload_key:\"foo\", workload_value:\"bar\",policy_hash:foo-hash",
					"policy_name:fizz,component_index:1,workload_index:workload_key:\"fizz\", workload_value:\"buzz\",policy_hash:fizz-hash",
					"policy_name:fizz,component_index:2,workload_index:workload_key:\"fizz\", workload_value:\"hoge\",policy_hash:fizz-hash"
				]`,
				"dropping_concurrency_limiters": `[
					"policy_name:foo,component_index:1,workload_index:workload_key:\"foo\", workload_value:\"bar\",policy_hash:foo-hash",
					"policy_name:fizz,component_index:1,workload_index:workload_key:\"fizz\", workload_value:\"buzz\",policy_hash:fizz-hash"
				]`,
			},
		),
	)

	DescribeTable("Processing traces",
		func(controlPoint string, decisions []*flowcontrolv1.LimiterDecision, expectedErr error, expectedMetrics string, expectedLabels map[string]string) {
			ctx := context.Background()

			traces := someTraces(decisions, controlPoint)
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
				Expect(traceRecords[0].Attributes().AsRaw()).To(HaveKeyWithValue(k, MatchJSON(v)))
			}
		},

		Entry("record with single policy - ingress",
			otelcollector.ControlPointIngress,
			[]*flowcontrolv1.LimiterDecision{
				{
					PolicyName:     "foo",
					PolicyHash:     "foo-hash",
					ComponentIndex: 1,
					Dropped:        true,
					Details: &flowcontrolv1.LimiterDecision_ConcurrencyLimiter_{
						ConcurrencyLimiter: &flowcontrolv1.LimiterDecision_ConcurrencyLimiter{
							Workload: "workload_key:\"foo\", workload_value:\"bar\"",
						},
					},
				},
			},
			nil,
			`# HELP workload_latency_ms Latency histogram of workload
			# TYPE workload_latency_ms histogram
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="workload_key:\"foo\", workload_value:\"bar\"",le="0"} 0
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="workload_key:\"foo\", workload_value:\"bar\"",le="10"} 1
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="workload_key:\"foo\", workload_value:\"bar\"",le="20"} 1
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="workload_key:\"foo\", workload_value:\"bar\"",le="+Inf"} 1
			workload_latency_ms_sum{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="workload_key:\"foo\", workload_value:\"bar\""} 5
			workload_latency_ms_count{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="workload_key:\"foo\", workload_value:\"bar\""} 1
			`,
			map[string]string{
				"rate_limiters":                 `[]`,
				"dropping_rate_limiters":        `[]`,
				"concurrency_limiters":          `["policy_name:foo,component_index:1,workload_index:workload_key:\"foo\", workload_value:\"bar\",policy_hash:foo-hash"]`,
				"dropping_concurrency_limiters": `["policy_name:foo,component_index:1,workload_index:workload_key:\"foo\", workload_value:\"bar\",policy_hash:foo-hash"]`,
			},
		),

		Entry("record with single policy - feature",
			otelcollector.ControlPointFeature,
			[]*flowcontrolv1.LimiterDecision{
				{
					PolicyName:     "foo",
					PolicyHash:     "foo-hash",
					ComponentIndex: 1,
					Dropped:        true,
					Details: &flowcontrolv1.LimiterDecision_ConcurrencyLimiter_{
						ConcurrencyLimiter: &flowcontrolv1.LimiterDecision_ConcurrencyLimiter{
							Workload: "workload_key:\"foo\", workload_value:\"bar\"",
						},
					},
				},
			},
			nil,
			`# HELP workload_latency_ms Latency histogram of workload
			# TYPE workload_latency_ms histogram
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="workload_key:\"foo\", workload_value:\"bar\"",le="0"} 0
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="workload_key:\"foo\", workload_value:\"bar\"",le="10"} 1
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="workload_key:\"foo\", workload_value:\"bar\"",le="20"} 1
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="workload_key:\"foo\", workload_value:\"bar\"",le="+Inf"} 1
			workload_latency_ms_sum{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="workload_key:\"foo\", workload_value:\"bar\""} 5
			workload_latency_ms_count{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="workload_key:\"foo\", workload_value:\"bar\""} 1
			`,
			map[string]string{
				"rate_limiters":                 `[]`,
				"dropping_rate_limiters":        `[]`,
				"concurrency_limiters":          `["policy_name:foo,component_index:1,workload_index:workload_key:\"foo\", workload_value:\"bar\",policy_hash:foo-hash"]`,
				"dropping_concurrency_limiters": `["policy_name:foo,component_index:1,workload_index:workload_key:\"foo\", workload_value:\"bar\",policy_hash:foo-hash"]`,
			},
		),

		Entry("record with two policies",
			otelcollector.ControlPointIngress,
			[]*flowcontrolv1.LimiterDecision{
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
							Workload: "workload_key:\"fizz\", workload_value:\"buzz\"",
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
							Workload: "workload_key:\"fizz\", workload_value:\"hoge\"",
						},
					},
				},
			},
			nil,
			`# HELP workload_latency_ms Latency histogram of workload
			# TYPE workload_latency_ms histogram
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="fizz-hash",policy_name="fizz",workload_index="workload_key:\"fizz\", workload_value:\"buzz\"",le="0"} 0
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="fizz-hash",policy_name="fizz",workload_index="workload_key:\"fizz\", workload_value:\"buzz\"",le="10"} 1
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="fizz-hash",policy_name="fizz",workload_index="workload_key:\"fizz\", workload_value:\"buzz\"",le="20"} 1
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="fizz-hash",policy_name="fizz",workload_index="workload_key:\"fizz\", workload_value:\"buzz\"",le="+Inf"} 1
			workload_latency_ms_sum{component_index="1",dropped="true",policy_hash="fizz-hash",policy_name="fizz",workload_index="workload_key:\"fizz\", workload_value:\"buzz\""} 5
			workload_latency_ms_count{component_index="1",dropped="true",policy_hash="fizz-hash",policy_name="fizz",workload_index="workload_key:\"fizz\", workload_value:\"buzz\""} 1
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="",le="0"} 0
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="",le="10"} 1
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="",le="20"} 1
			workload_latency_ms_bucket{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index="",le="+Inf"} 1
			workload_latency_ms_sum{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index=""} 5
			workload_latency_ms_count{component_index="1",dropped="true",policy_hash="foo-hash",policy_name="foo",workload_index=""} 1
			workload_latency_ms_bucket{component_index="2",dropped="true",policy_hash="fizz-hash",policy_name="fizz",workload_index="workload_key:\"fizz\", workload_value:\"hoge\"",le="0"} 0
			workload_latency_ms_bucket{component_index="2",dropped="true",policy_hash="fizz-hash",policy_name="fizz",workload_index="workload_key:\"fizz\", workload_value:\"hoge\"",le="10"} 1
			workload_latency_ms_bucket{component_index="2",dropped="true",policy_hash="fizz-hash",policy_name="fizz",workload_index="workload_key:\"fizz\", workload_value:\"hoge\"",le="20"} 1
			workload_latency_ms_bucket{component_index="2",dropped="true",policy_hash="fizz-hash",policy_name="fizz",workload_index="workload_key:\"fizz\", workload_value:\"hoge\"",le="+Inf"} 1
			workload_latency_ms_sum{component_index="2",dropped="true",policy_hash="fizz-hash",policy_name="fizz",workload_index="workload_key:\"fizz\", workload_value:\"hoge\""} 5
			workload_latency_ms_count{component_index="2",dropped="true",policy_hash="fizz-hash",policy_name="fizz",workload_index="workload_key:\"fizz\", workload_value:\"hoge\""} 1
			`,
			map[string]string{
				"rate_limiters": `[
					"policy_name:foo,component_index:1,policy_hash:foo-hash"
				]`,
				"dropping_rate_limiters": `[
					"policy_name:foo,component_index:1,policy_hash:foo-hash"
				]`,
				"concurrency_limiters": `[
					"policy_name:fizz,component_index:1,workload_index:workload_key:\"fizz\", workload_value:\"buzz\",policy_hash:fizz-hash",
					"policy_name:fizz,component_index:2,workload_index:workload_key:\"fizz\", workload_value:\"hoge\",policy_hash:fizz-hash"
				]`,
				"dropping_concurrency_limiters": `[
					"policy_name:fizz,component_index:1,workload_index:workload_key:\"fizz\", workload_value:\"buzz\",policy_hash:fizz-hash",
					"policy_name:fizz,component_index:2,workload_index:workload_key:\"fizz\", workload_value:\"hoge\",policy_hash:fizz-hash"
				]`,
			},
		),
	)
})

// someLogs will return a plog.Logs instance with single LogRecord
func someLogs(decisions []*flowcontrolv1.LimiterDecision, controlPoint string) plog.Logs {
	logs := plog.NewLogs()
	logs.ResourceLogs().AppendEmpty()

	resourceLogsSlice := logs.ResourceLogs()
	for i := 0; i < resourceLogsSlice.Len(); i++ {
		resourceLogsSlice.At(i).ScopeLogs().AppendEmpty()

		instrumentationLogsSlice := resourceLogsSlice.At(i).ScopeLogs()
		for j := 0; j < instrumentationLogsSlice.Len(); j++ {
			logRecord := instrumentationLogsSlice.At(j).LogRecords().AppendEmpty()
			marshalled, err := json.Marshal(decisions)
			Expect(err).NotTo(HaveOccurred())
			logRecord.Attributes().InsertString(otelcollector.LimiterDecisionsLabel, string(marshalled))
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
func someTraces(decisions []*flowcontrolv1.LimiterDecision, controlPoint string) ptrace.Traces {
	traces := ptrace.NewTraces()
	traces.ResourceSpans().AppendEmpty()

	resourceSpansSlice := traces.ResourceSpans()
	for i := 0; i < resourceSpansSlice.Len(); i++ {
		resourceSpansSlice.At(i).ScopeSpans().AppendEmpty()

		instrumentationSpansSlice := resourceSpansSlice.At(i).ScopeSpans()
		for j := 0; j < instrumentationSpansSlice.Len(); j++ {
			span := instrumentationSpansSlice.At(j).Spans().AppendEmpty()
			marshalled, err := json.Marshal(decisions)
			Expect(err).NotTo(HaveOccurred())
			span.Attributes().InsertString(otelcollector.LimiterDecisionsLabel, string(marshalled))
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
