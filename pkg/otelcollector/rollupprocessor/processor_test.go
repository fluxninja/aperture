package rollupprocessor

import (
	"context"
	"encoding/base64"
	"strconv"

	"github.com/fluxninja/datasketches-go/sketches"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.opentelemetry.io/collector/processor"

	otelconsts "github.com/fluxninja/aperture/pkg/otelcollector/consts"
)

var _ = Describe("Rollup processor", func() {
	var (
		config       *Config
		testConsumer *fakeConsumer
	)

	BeforeEach(func() {
		config = &Config{
			AttributeCardinalityLimit: 10,
			RollupBuckets:             []float64{10, 20, 30},
			promRegistry:              prometheus.NewRegistry(),
		}
		testConsumer = &fakeConsumer{
			receivedLogs:    []plog.Logs{},
			receivedMetrics: []pmetric.Metrics{},
			receivedTraces:  []ptrace.Traces{},
		}
	})

	Describe("Logs", func() {
		var logsProcessor processor.Logs

		JustBeforeEach(func() {
			var err error
			logsProcessor, err = CreateLogsProcessor(
				context.Background(), processor.CreateSettings{}, config, testConsumer)
			Expect(err).NotTo(HaveOccurred())
		})

		It("works for single log record", func() {
			attributeValues := []int{5}
			expectedSerializedDatasketch, err := serializedDatasketchFromAttributeValues(attributeValues)
			Expect(err).NotTo(HaveOccurred())

			input := plog.NewLogs()
			logs := input.ResourceLogs().AppendEmpty().
				ScopeLogs().AppendEmpty().
				LogRecords()
			logRecord := logs.AppendEmpty()
			logRecord.Attributes().PutStr("fizz", "buzz")
			logRecord.Attributes().PutStr(otelconsts.ApertureSourceLabel, otelconsts.ApertureSourceEnvoy)
			logRecord.Attributes().PutStr(otelconsts.WorkloadDurationLabel, strconv.Itoa(attributeValues[0]))

			err = logsProcessor.ConsumeLogs(context.Background(), input)
			Expect(err).NotTo(HaveOccurred())

			Expect(testConsumer.receivedLogs).To(HaveLen(1))
			attributes := testConsumer.receivedLogs[0].ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0).Attributes().AsRaw()
			Expect(attributes).To(HaveLen(8))
			Expect(attributes).To(HaveKeyWithValue(RollupCountKey, int64(1)))
			Expect(attributes).To(HaveKeyWithValue(AggregateField(otelconsts.WorkloadDurationLabel, RollupDatasketch), expectedSerializedDatasketch))
			Expect(attributes).To(HaveKeyWithValue(AggregateField(otelconsts.WorkloadDurationLabel, RollupSum), float64(5)))
			Expect(attributes).To(HaveKeyWithValue(AggregateField(otelconsts.WorkloadDurationLabel, RollupMin), float64(5)))
			Expect(attributes).To(HaveKeyWithValue(AggregateField(otelconsts.WorkloadDurationLabel, RollupMax), float64(5)))
			Expect(attributes).To(HaveKeyWithValue(AggregateField(otelconsts.WorkloadDurationLabel, RollupSumOfSquares), float64(25)))
			Expect(attributes).To(HaveKeyWithValue("fizz", "buzz"))
		})

		It("works for multiple log records", func() {
			attributeValues := []int{5, 6, 7}
			expectedSerializedDatasketch, err := serializedDatasketchFromAttributeValues(attributeValues)
			Expect(err).NotTo(HaveOccurred())

			input := plog.NewLogs()
			logs := input.ResourceLogs().AppendEmpty().ScopeLogs().AppendEmpty().LogRecords()
			logRecord := logs.AppendEmpty()
			logRecord.Attributes().PutStr(otelconsts.WorkloadDurationLabel, strconv.Itoa(attributeValues[0]))
			logRecord = logs.AppendEmpty()
			logRecord.Attributes().PutStr(otelconsts.WorkloadDurationLabel, strconv.Itoa(attributeValues[1]))
			logRecord = logs.AppendEmpty()
			logRecord.Attributes().PutStr(otelconsts.WorkloadDurationLabel, strconv.Itoa(attributeValues[2]))
			logRecord = logs.AppendEmpty()
			logRecord.Attributes().PutStr(otelconsts.HTTPRequestContentLength, strconv.Itoa(1234))

			err = logsProcessor.ConsumeLogs(context.Background(), input)
			Expect(err).NotTo(HaveOccurred())

			Expect(testConsumer.receivedLogs).To(HaveLen(1))
			attributes := testConsumer.receivedLogs[0].ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0).Attributes().AsRaw()
			Expect(attributes).To(HaveLen(10))
			Expect(attributes).To(HaveKeyWithValue(RollupCountKey, int64(4)))
			Expect(attributes).To(HaveKeyWithValue(AggregateField(otelconsts.WorkloadDurationLabel, RollupDatasketch), expectedSerializedDatasketch))
			Expect(attributes).To(HaveKeyWithValue(AggregateField(otelconsts.WorkloadDurationLabel, RollupSum), float64(18)))
			Expect(attributes).To(HaveKeyWithValue(AggregateField(otelconsts.WorkloadDurationLabel, RollupMin), float64(5)))
			Expect(attributes).To(HaveKeyWithValue(AggregateField(otelconsts.WorkloadDurationLabel, RollupMax), float64(7)))
			Expect(attributes).To(HaveKeyWithValue(AggregateField(otelconsts.WorkloadDurationLabel, RollupSumOfSquares), float64(110)))
			Expect(attributes).NotTo(HaveKey(AggregateField(otelconsts.HTTPRequestContentLength, RollupDatasketch)))
		})

		It("applies cardinality limits", func() {
			input := plog.NewLogs()
			logs := input.ResourceLogs().AppendEmpty().ScopeLogs().AppendEmpty().LogRecords()

			for i := 0; i < 30; i++ {
				logRecord := logs.AppendEmpty()
				logRecord.Attributes().PutStr(otelconsts.WorkloadDurationLabel, strconv.Itoa(i))
				logRecord.Attributes().PutStr(otelconsts.ApertureSourceLabel, otelconsts.ApertureSourceEnvoy)
				logRecord.Attributes().PutStr("low-cardinality", strconv.Itoa(i%2))
				logRecord.Attributes().PutStr("almost-high-cardinality", strconv.Itoa(i%10))
				logRecord.Attributes().PutStr("high-cardinality", strconv.Itoa(i))
			}

			err := logsProcessor.ConsumeLogs(context.Background(), input)
			Expect(err).NotTo(HaveOccurred())

			Expect(testConsumer.receivedLogs).To(HaveLen(1))
			receivedLogRecords := testConsumer.receivedLogs[0].ResourceLogs().At(0).ScopeLogs().At(0).LogRecords()
			Expect(receivedLogRecords.Len()).To(Equal(20))

			uniqueValues := map[string]map[string]struct{}{}

			for i := 0; i < receivedLogRecords.Len(); i++ {
				attrs := receivedLogRecords.At(i).Attributes()
				attrs.Range(func(k string, v pcommon.Value) bool {
					value := v.AsString()
					values, exist := uniqueValues[k]
					if !exist {
						values = map[string]struct{}{}
						uniqueValues[k] = values
					}
					values[value] = struct{}{}
					return true
				})
			}
			Expect(uniqueValues[otelconsts.WorkloadDurationLabel+"_max"]).To(HaveLen(20))
			Expect(uniqueValues[otelconsts.WorkloadDurationLabel+"_max"]).To(HaveKey("29"))
			Expect(uniqueValues["low-cardinality"]).To(HaveLen(2))
			Expect(uniqueValues["almost-high-cardinality"]).To(HaveLen(10))
			Expect(uniqueValues["high-cardinality"]).To(HaveLen(11))
			Expect(uniqueValues["high-cardinality"]).To(HaveKey(RedactedAttributeValue))
		})
	})
})

type fakeConsumer struct {
	receivedLogs    []plog.Logs
	receivedMetrics []pmetric.Metrics
	receivedTraces  []ptrace.Traces
}

func (fc *fakeConsumer) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{}
}

func (fc *fakeConsumer) ConsumeLogs(_ context.Context, ld plog.Logs) error {
	fc.receivedLogs = append(fc.receivedLogs, ld)
	return nil
}

func (fc *fakeConsumer) ConsumeMetrics(_ context.Context, ld pmetric.Metrics) error {
	fc.receivedMetrics = append(fc.receivedMetrics, ld)
	return nil
}

func (fc *fakeConsumer) ConsumeTraces(_ context.Context, ld ptrace.Traces) error {
	fc.receivedTraces = append(fc.receivedTraces, ld)
	return nil
}

func serializedDatasketchFromAttributeValues(values []int) (string, error) {
	sketch, err := sketches.NewDoublesSketch(128)
	if err != nil {
		return "", err
	}
	for _, v := range values {
		err = sketch.Update(float64(v))
		if err != nil {
			return "", err
		}
	}
	sketchBytes, _ := sketch.Compact().Serialize()
	return base64.StdEncoding.EncodeToString(sketchBytes), nil
}
