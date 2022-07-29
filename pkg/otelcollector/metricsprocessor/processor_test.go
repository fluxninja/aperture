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

	"github.com/FluxNinja/aperture/pkg/otelcollector"
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
		func(controlPoint string, policies []policy, expectedErr error, expectedMetrics string) {
			ctx := context.Background()

			logs := someLogs(policies, controlPoint)
			modifiedLogs, err := processor.ConsumeLogs(ctx, logs)
			if expectedErr != nil {
				Expect(err).NotTo(MatchError(expectedErr))
				return
			}
			Expect(err).NotTo(HaveOccurred())
			Expect(modifiedLogs).To(Equal(logs))

			By("sending proper metrics")
			expected := strings.NewReader(expectedMetrics)
			err = testutil.CollectAndCompare(
				processor.requestLatencyHistogram,
				expected,
				"request_latency_ms")
			Expect(err).NotTo(HaveOccurred())

			By("adding proper labels")
			logRecords := allLogRecords(modifiedLogs)
			Expect(logRecords).To(HaveLen(1))

			expectedMatched, expectedDropped := getIDs(policies)
			Expect(logRecords[0].Attributes().AsRaw()).To(
				HaveKeyWithValue(otelcollector.PoliciesMatchedLabel, expectedMatched))
			Expect(logRecords[0].Attributes().AsRaw()).To(
				HaveKeyWithValue(otelcollector.PoliciesDroppedLabel, expectedDropped))
		},
		Entry("record with single policy - ingress",
			otelcollector.ControlPointIngress,
			[]policy{{
				ID:       "foo",
				Dropped:  true,
				Workload: "workload_key:\"foo\", workload_value:\"bar\"",
			}},
			nil,
			`# HELP request_latency_ms Latency of requests histogram
			# TYPE request_latency_ms histogram
			request_latency_ms_bucket{dropped="true",metric_id="foo",workload_key_name="foo",workload_key_value="bar",le="0"} 0
			request_latency_ms_bucket{dropped="true",metric_id="foo",workload_key_name="foo",workload_key_value="bar",le="10"} 1
			request_latency_ms_bucket{dropped="true",metric_id="foo",workload_key_name="foo",workload_key_value="bar",le="20"} 1
			request_latency_ms_bucket{dropped="true",metric_id="foo",workload_key_name="foo",workload_key_value="bar",le="+Inf"} 1
			request_latency_ms_sum{dropped="true",metric_id="foo",workload_key_name="foo",workload_key_value="bar",} 5
			request_latency_ms_count{dropped="true",metric_id="foo",workload_key_name="foo",workload_key_value="bar"} 1
			`,
		),
		Entry("record with single policy - feature",
			otelcollector.ControlPointFeature,
			[]policy{{
				ID:       "foo",
				Dropped:  true,
				Workload: "workload_key:\"foo\", workload_value:\"bar\"",
			}},
			nil,
			`# HELP request_latency_ms Latency of requests histogram
			# TYPE request_latency_ms histogram
			request_latency_ms_bucket{dropped="true",metric_id="foo",workload_key_name="foo",workload_key_value="bar",le="0"} 0
			request_latency_ms_bucket{dropped="true",metric_id="foo",workload_key_name="foo",workload_key_value="bar",le="10"} 1
			request_latency_ms_bucket{dropped="true",metric_id="foo",workload_key_name="foo",workload_key_value="bar",le="20"} 1
			request_latency_ms_bucket{dropped="true",metric_id="foo",workload_key_name="foo",workload_key_value="bar",le="+Inf"} 1
			request_latency_ms_sum{dropped="true",metric_id="foo",workload_key_name="foo",workload_key_value="bar",} 5
			request_latency_ms_count{dropped="true",metric_id="foo",workload_key_name="foo",workload_key_value="bar"} 1
			`,
		),
		Entry("record with two policies",
			otelcollector.ControlPointIngress,
			[]policy{
				{
					ID:       "foo",
					Dropped:  true,
					Workload: "workload_key:\"foo\", workload_value:\"bar\"",
				},
				{
					ID:       "fizz",
					Dropped:  false,
					Workload: "workload_key:\"fizz\", workload_value:\"buzz\"",
				},
				{
					ID:       "fizz",
					Dropped:  false,
					Workload: "workload_key:\"fizz\", workload_value:\"hoge\"",
				},
			},
			nil,
			`# HELP request_latency_ms Latency of requests histogram
			# TYPE request_latency_ms histogram
			request_latency_ms_bucket{dropped="false",metric_id="fizz",workload_key_name="fizz",workload_key_value="buzz",le="0"} 0
			request_latency_ms_bucket{dropped="false",metric_id="fizz",workload_key_name="fizz",workload_key_value="buzz",le="10"} 1
			request_latency_ms_bucket{dropped="false",metric_id="fizz",workload_key_name="fizz",workload_key_value="buzz",le="20"} 1
			request_latency_ms_bucket{dropped="false",metric_id="fizz",workload_key_name="fizz",workload_key_value="buzz",le="+Inf"} 1
			request_latency_ms_sum{dropped="false",metric_id="fizz",workload_key_name="fizz",workload_key_value="buzz"} 5
			request_latency_ms_count{dropped="false",metric_id="fizz",workload_key_name="fizz",workload_key_value="buzz"} 1
			request_latency_ms_bucket{dropped="false",metric_id="fizz",workload_key_name="fizz",workload_key_value="hoge",le="0"} 0
			request_latency_ms_bucket{dropped="false",metric_id="fizz",workload_key_name="fizz",workload_key_value="hoge",le="10"} 1
			request_latency_ms_bucket{dropped="false",metric_id="fizz",workload_key_name="fizz",workload_key_value="hoge",le="20"} 1
			request_latency_ms_bucket{dropped="false",metric_id="fizz",workload_key_name="fizz",workload_key_value="hoge",le="+Inf"} 1
			request_latency_ms_sum{dropped="false",metric_id="fizz",workload_key_name="fizz",workload_key_value="hoge"} 5
			request_latency_ms_count{dropped="false",metric_id="fizz",workload_key_name="fizz",workload_key_value="hoge"} 1
			request_latency_ms_bucket{dropped="true",metric_id="foo",workload_key_name="foo",workload_key_value="bar",le="0"} 0
			request_latency_ms_bucket{dropped="true",metric_id="foo",workload_key_name="foo",workload_key_value="bar",le="10"} 1
			request_latency_ms_bucket{dropped="true",metric_id="foo",workload_key_name="foo",workload_key_value="bar",le="20"} 1
			request_latency_ms_bucket{dropped="true",metric_id="foo",workload_key_name="foo",workload_key_value="bar",le="+Inf"} 1
			request_latency_ms_sum{dropped="true",metric_id="foo",workload_key_name="foo",workload_key_value="bar"} 5
			request_latency_ms_count{dropped="true",metric_id="foo",workload_key_name="foo",workload_key_value="bar"} 1
			`,
		),
		Entry("policy without priority keys",
			otelcollector.ControlPointIngress,
			[]policy{
				{
					ID:      "foo",
					Dropped: true,
				},
			},
			nil,
			`# HELP request_latency_ms Latency of requests histogram
			# TYPE request_latency_ms histogram
			request_latency_ms_bucket{dropped="true",metric_id="foo",workload_key_name="default_workload_key",workload_key_value="default_workload_value",le="0"} 0
			request_latency_ms_bucket{dropped="true",metric_id="foo",workload_key_name="default_workload_key",workload_key_value="default_workload_value",le="10"} 1
			request_latency_ms_bucket{dropped="true",metric_id="foo",workload_key_name="default_workload_key",workload_key_value="default_workload_value",le="20"} 1
			request_latency_ms_bucket{dropped="true",metric_id="foo",workload_key_name="default_workload_key",workload_key_value="default_workload_value",le="+Inf"} 1
			request_latency_ms_sum{dropped="true",metric_id="foo",workload_key_name="default_workload_key",workload_key_value="default_workload_value"} 5
			request_latency_ms_count{dropped="true",metric_id="foo",workload_key_name="default_workload_key",workload_key_value="default_workload_value"} 1
			`,
		),
	)

	DescribeTable("Processing traces",
		func(controlPoint string, policies []policy, expectedErr error, expectedMetrics string) {
			ctx := context.Background()

			traces := someTraces(policies, controlPoint)
			modifiedTraces, err := processor.ConsumeTraces(ctx, traces)
			if expectedErr != nil {
				Expect(err).NotTo(MatchError(expectedErr))
				return
			}
			Expect(err).NotTo(HaveOccurred())
			Expect(modifiedTraces).To(Equal(traces))

			By("sending proper metrics")
			expected := strings.NewReader(expectedMetrics)
			err = testutil.CollectAndCompare(
				processor.requestLatencyHistogram,
				expected,
				"request_latency_ms")
			Expect(err).NotTo(HaveOccurred())

			By("adding proper labels")
			traceRecords := allTraceRecords(modifiedTraces)
			Expect(traceRecords).To(HaveLen(1))

			expectedMatched, expectedDropped := getIDs(policies)
			Expect(traceRecords[0].Attributes().AsRaw()).To(
				HaveKeyWithValue(otelcollector.PoliciesMatchedLabel, expectedMatched))
			Expect(traceRecords[0].Attributes().AsRaw()).To(
				HaveKeyWithValue(otelcollector.PoliciesDroppedLabel, expectedDropped))
		},
		Entry("record with single policy - ingress",
			otelcollector.ControlPointIngress,
			[]policy{{
				ID:       "foo",
				Dropped:  true,
				Workload: "workload_key:\"foo\", workload_value:\"bar\"",
			}},
			nil,
			`# HELP request_latency_ms Latency of requests histogram
			# TYPE request_latency_ms histogram
			request_latency_ms_bucket{dropped="true",metric_id="foo",workload_key_name="foo",workload_key_value="bar",le="0"} 0
			request_latency_ms_bucket{dropped="true",metric_id="foo",workload_key_name="foo",workload_key_value="bar",le="10"} 1
			request_latency_ms_bucket{dropped="true",metric_id="foo",workload_key_name="foo",workload_key_value="bar",le="20"} 1
			request_latency_ms_bucket{dropped="true",metric_id="foo",workload_key_name="foo",workload_key_value="bar",le="+Inf"} 1
			request_latency_ms_sum{dropped="true",metric_id="foo",workload_key_name="foo",workload_key_value="bar",} 5
			request_latency_ms_count{dropped="true",metric_id="foo",workload_key_name="foo",workload_key_value="bar"} 1
			`,
		),
		Entry("record with single policy - feature",
			otelcollector.ControlPointFeature,
			[]policy{{
				ID:       "foo",
				Dropped:  true,
				Workload: "workload_key:\"foo\", workload_value:\"bar\"",
			}},
			nil,
			`# HELP request_latency_ms Latency of requests histogram
			# TYPE request_latency_ms histogram
			request_latency_ms_bucket{dropped="true",metric_id="foo",workload_key_name="foo",workload_key_value="bar",le="0"} 0
			request_latency_ms_bucket{dropped="true",metric_id="foo",workload_key_name="foo",workload_key_value="bar",le="10"} 1
			request_latency_ms_bucket{dropped="true",metric_id="foo",workload_key_name="foo",workload_key_value="bar",le="20"} 1
			request_latency_ms_bucket{dropped="true",metric_id="foo",workload_key_name="foo",workload_key_value="bar",le="+Inf"} 1
			request_latency_ms_sum{dropped="true",metric_id="foo",workload_key_name="foo",workload_key_value="bar",} 5
			request_latency_ms_count{dropped="true",metric_id="foo",workload_key_name="foo",workload_key_value="bar"} 1
			`,
		),
		Entry("record with two policies",
			otelcollector.ControlPointIngress,
			[]policy{
				{
					ID:       "foo",
					Dropped:  true,
					Workload: "workload_key:\"foo\", workload_value:\"bar\"",
				},
				{
					ID:       "fizz",
					Dropped:  false,
					Workload: "workload_key:\"fizz\", workload_value:\"buzz\"",
				},
				{
					ID:       "fizz",
					Dropped:  false,
					Workload: "workload_key:\"fizz\", workload_value:\"hoge\"",
				},
			},
			nil,
			`# HELP request_latency_ms Latency of requests histogram
			# TYPE request_latency_ms histogram
			request_latency_ms_bucket{dropped="false",metric_id="fizz",workload_key_name="fizz",workload_key_value="buzz",le="0"} 0
			request_latency_ms_bucket{dropped="false",metric_id="fizz",workload_key_name="fizz",workload_key_value="buzz",le="10"} 1
			request_latency_ms_bucket{dropped="false",metric_id="fizz",workload_key_name="fizz",workload_key_value="buzz",le="20"} 1
			request_latency_ms_bucket{dropped="false",metric_id="fizz",workload_key_name="fizz",workload_key_value="buzz",le="+Inf"} 1
			request_latency_ms_sum{dropped="false",metric_id="fizz",workload_key_name="fizz",workload_key_value="buzz"} 5
			request_latency_ms_count{dropped="false",metric_id="fizz",workload_key_name="fizz",workload_key_value="buzz"} 1
			request_latency_ms_bucket{dropped="false",metric_id="fizz",workload_key_name="fizz",workload_key_value="hoge",le="0"} 0
			request_latency_ms_bucket{dropped="false",metric_id="fizz",workload_key_name="fizz",workload_key_value="hoge",le="10"} 1
			request_latency_ms_bucket{dropped="false",metric_id="fizz",workload_key_name="fizz",workload_key_value="hoge",le="20"} 1
			request_latency_ms_bucket{dropped="false",metric_id="fizz",workload_key_name="fizz",workload_key_value="hoge",le="+Inf"} 1
			request_latency_ms_sum{dropped="false",metric_id="fizz",workload_key_name="fizz",workload_key_value="hoge"} 5
			request_latency_ms_count{dropped="false",metric_id="fizz",workload_key_name="fizz",workload_key_value="hoge"} 1
			request_latency_ms_bucket{dropped="true",metric_id="foo",workload_key_name="foo",workload_key_value="bar",le="0"} 0
			request_latency_ms_bucket{dropped="true",metric_id="foo",workload_key_name="foo",workload_key_value="bar",le="10"} 1
			request_latency_ms_bucket{dropped="true",metric_id="foo",workload_key_name="foo",workload_key_value="bar",le="20"} 1
			request_latency_ms_bucket{dropped="true",metric_id="foo",workload_key_name="foo",workload_key_value="bar",le="+Inf"} 1
			request_latency_ms_sum{dropped="true",metric_id="foo",workload_key_name="foo",workload_key_value="bar"} 5
			request_latency_ms_count{dropped="true",metric_id="foo",workload_key_name="foo",workload_key_value="bar"} 1
			`,
		),
		Entry("policy without priority keys",
			otelcollector.ControlPointIngress,
			[]policy{
				{
					ID:      "foo",
					Dropped: true,
				},
			},
			nil,
			`# HELP request_latency_ms Latency of requests histogram
			# TYPE request_latency_ms histogram
			request_latency_ms_bucket{dropped="true",metric_id="foo",workload_key_name="default_workload_key",workload_key_value="default_workload_value",le="0"} 0
			request_latency_ms_bucket{dropped="true",metric_id="foo",workload_key_name="default_workload_key",workload_key_value="default_workload_value",le="10"} 1
			request_latency_ms_bucket{dropped="true",metric_id="foo",workload_key_name="default_workload_key",workload_key_value="default_workload_value",le="20"} 1
			request_latency_ms_bucket{dropped="true",metric_id="foo",workload_key_name="default_workload_key",workload_key_value="default_workload_value",le="+Inf"} 1
			request_latency_ms_sum{dropped="true",metric_id="foo",workload_key_name="default_workload_key",workload_key_value="default_workload_value"} 5
			request_latency_ms_count{dropped="true",metric_id="foo",workload_key_name="default_workload_key",workload_key_value="default_workload_value"} 1
			`,
		),
	)
})

// someLogs will return a plog.Logs instance with single LogRecord
func someLogs(policies []policy, controlPoint string) plog.Logs {
	logs := plog.NewLogs()
	logs.ResourceLogs().AppendEmpty()

	resourceLogsSlice := logs.ResourceLogs()
	for i := 0; i < resourceLogsSlice.Len(); i++ {
		resourceLogsSlice.At(i).ScopeLogs().AppendEmpty()

		instrumentationLogsSlice := resourceLogsSlice.At(i).ScopeLogs()
		for j := 0; j < instrumentationLogsSlice.Len(); j++ {
			logRecord := instrumentationLogsSlice.At(j).LogRecords().AppendEmpty()
			marshalled, err := json.Marshal(policies)
			Expect(err).NotTo(HaveOccurred())
			logRecord.Attributes().InsertString(otelcollector.PoliciesLabel, string(marshalled))
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
func someTraces(policies []policy, controlPoint string) ptrace.Traces {
	traces := ptrace.NewTraces()
	traces.ResourceSpans().AppendEmpty()

	resourceSpansSlice := traces.ResourceSpans()
	for i := 0; i < resourceSpansSlice.Len(); i++ {
		resourceSpansSlice.At(i).ScopeSpans().AppendEmpty()

		instrumentationSpansSlice := resourceSpansSlice.At(i).ScopeSpans()
		for j := 0; j < instrumentationSpansSlice.Len(); j++ {
			span := instrumentationSpansSlice.At(j).Spans().AppendEmpty()
			marshalled, err := json.Marshal(policies)
			Expect(err).NotTo(HaveOccurred())
			span.Attributes().InsertString(otelcollector.PoliciesLabel, string(marshalled))
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
