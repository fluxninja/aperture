package rollupprocessor

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/processor"

	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/metrics"
	"github.com/fluxninja/aperture/pkg/otelcollector"
	otelconsts "github.com/fluxninja/aperture/pkg/otelcollector/consts"
	"github.com/fluxninja/datasketches-go/sketches"
	"github.com/prometheus/client_golang/prometheus"
)

var defaultRollupGroups = []RollupGroup{
	{
		FromField:      otelconsts.WorkloadDurationLabel,
		WithDatasketch: true,
	},
	{
		FromField:      otelconsts.FlowDurationLabel,
		WithDatasketch: true,
	},
	{
		FromField:      otelconsts.ApertureProcessingDurationLabel,
		WithDatasketch: true,
	},
	{
		FromField: otelconsts.HTTPRequestContentLength,
	},
	{
		FromField: otelconsts.HTTPResponseContentLength,
	},
}

type rollupProcessor struct {
	cfg *Config

	logsNextConsumer consumer.Logs
	rollupHistogram  *prometheus.HistogramVec
	rollups          []*Rollup
	rollupFromFields map[string]struct{} // Set of all all FromField names.
}

const (
	defaultAttributeCardinalityLimit = 10
	// RedactedAttributeValue is a value that replaces actual attribute value
	// in case it exceeds cardinality limit.
	RedactedAttributeValue = "REDACTED_VIA_CARDINALITY_LIMIT"
)

var _ consumer.Logs = (*rollupProcessor)(nil)

func newRollupProcessor(set processor.CreateSettings, cfg *Config) (*rollupProcessor, error) {
	rollupHistogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    metrics.RollupMetricName,
		Help:    "Latency of the requests processed by the server",
		Buckets: cfg.RollupBuckets,
	}, []string{})

	rollups := NewRollups(defaultRollupGroups)
	rollupFromFields := map[string]struct{}{}
	for _, rollup := range rollups {
		rollupFromFields[rollup.FromField] = struct{}{}
	}

	err := cfg.promRegistry.Register(rollupHistogram)
	if err != nil {
		// Ignore already registered error, as this is not harmful. Metrics may
		// be registered by other running processor.
		if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
			return nil, fmt.Errorf("could not register prometheus metrics: %w", err)
		}
	}

	return &rollupProcessor{
		cfg:              cfg,
		rollupHistogram:  rollupHistogram,
		rollups:          rollups,
		rollupFromFields: rollupFromFields,
	}, nil
}

// Capabilities returns the capabilities of the processor with MutatesData set to true.
func (rp *rollupProcessor) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: true}
}

// Start is invoked during service startup.
func (rp *rollupProcessor) Start(context.Context, component.Host) error {
	return nil
}

// Shutdown is invoked during service shutdown.
func (rp *rollupProcessor) Shutdown(context.Context) error {
	return nil
}

// ConsumeLogs implements LogsProcessor.
func (rp *rollupProcessor) ConsumeLogs(ctx context.Context, ld plog.Logs) error {
	rp.applyCardinalityLimits(ld, rp.cfg.AttributeCardinalityLimit)

	rollupData := make(map[string]pcommon.Map)
	datasketches := make(map[string]map[string]*sketches.HeapDoublesSketch)

	log.Trace().Int("count", ld.LogRecordCount()).Msg("Before rollup")
	otelcollector.IterateLogRecords(ld, func(logRecord plog.LogRecord) otelcollector.IterAction {
		attributes := logRecord.Attributes()
		key := key(attributes, rp.rollupFromFields)
		_, exists := rollupData[key]
		if !exists {
			rollupData[key] = attributes
			rollupData[key].PutInt(RollupCountKey, 0)
		}
		_, exists = datasketches[key]
		if !exists {
			datasketches[key] = make(map[string]*sketches.HeapDoublesSketch)
		}
		rawCount, _ := rollupData[key].Get(RollupCountKey)
		rollupData[key].PutInt(RollupCountKey, rawCount.Int()+1)
		rp.rollupAttributes(datasketches[key], rollupData[key], attributes)
		return otelcollector.Keep
	})
	for k, v := range datasketches {
		attributes := rollupData[k]
		for toField, sketch := range v {
			serializedBytes, err := sketch.Compact().Serialize()
			if err != nil {
				// Serialize() should never return error, unless there's a bug in sketches library.
				log.Bug().Err(err).Msg("Sketch.Serialize() failed")
			}
			serialized := base64.StdEncoding.EncodeToString(serializedBytes)
			attributes.PutStr(toField, serialized)
		}
	}
	return rp.exportLogs(ctx, rollupData)
}

func (rp *rollupProcessor) applyCardinalityLimits(ld plog.Logs, limit int) {
	attributeValues := map[string]map[string]struct{}{}
	otelcollector.IterateLogRecords(ld, func(logRecord plog.LogRecord) otelcollector.IterAction {
		logRecord.Attributes().Range(func(k string, v pcommon.Value) bool {
			// Applying the cardinality limits only to string attributes, as
			// all user-created attributes where we risk high-cardinality are
			// string attributes.
			if v.Type() != pcommon.ValueTypeStr {
				return true
			}
			if _, toRollup := rp.rollupFromFields[k]; toRollup {
				return true
			}
			value := v.Str()
			values, exist := attributeValues[k]
			if !exist {
				attributeValues[k] = map[string]struct{}{value: {}}
				return true
			}
			if _, exists := values[value]; !exists {
				if len(values) < limit {
					values[value] = struct{}{}
				} else {
					v.SetStr(RedactedAttributeValue)
				}
			}
			return true
		})
		return otelcollector.Keep
	})
}

func (rp *rollupProcessor) rollupAttributes(
	datasketches map[string]*sketches.HeapDoublesSketch,
	baseAttributes,
	attributes pcommon.Map,
) {
	// TODO tgill: need to track latest timestamp from attributes as the timestamp in baseAttributes
	for _, rollup := range rp.rollups {
		switch rollup.Type {
		case RollupSum:
			newValue, found := rollup.GetFromFieldValue(attributes)
			if !found {
				continue
			}
			rollupSum, exists := rollup.GetToFieldValue(baseAttributes)
			if !exists {
				rollupSum = 0
				baseAttributes.PutDouble(rollup.ToField, rollupSum)
			}
			baseAttributes.PutDouble(rollup.ToField, rollupSum+newValue)
		case RollupSumOfSquares:
			newValue, found := rollup.GetFromFieldValue(attributes)
			if !found {
				continue
			}
			rollupSos, exists := rollup.GetToFieldValue(baseAttributes)
			if !exists {
				rollupSos = 0
				baseAttributes.PutDouble(rollup.ToField, rollupSos)
			}
			baseAttributes.PutDouble(rollup.ToField, rollupSos+newValue*newValue)
		case RollupMin:
			newValue, found := rollup.GetFromFieldValue(attributes)
			if !found {
				continue
			}
			rollupMin, exists := rollup.GetToFieldValue(baseAttributes)
			if !exists {
				rollupMin = newValue
				baseAttributes.PutDouble(rollup.ToField, rollupMin)
			}
			newMin := otelcollector.Min(rollupMin, newValue)
			baseAttributes.PutDouble(rollup.ToField, newMin)
		case RollupMax:
			newValue, found := rollup.GetFromFieldValue(attributes)
			if !found {
				continue
			}
			rollupMax, exists := rollup.GetToFieldValue(baseAttributes)
			if !exists {
				rollupMax = newValue
				baseAttributes.PutDouble(rollup.ToField, rollupMax)
			}
			newMax := otelcollector.Max(rollupMax, newValue)
			baseAttributes.PutDouble(rollup.ToField, newMax)
		case RollupDatasketch:
			newValue, found := rollup.GetFromFieldValue(attributes)
			if !found {
				continue
			}
			_, exists := datasketches[rollup.ToField]
			if !exists {
				ds, err := sketches.NewDoublesSketch(128)
				if err != nil {
					log.Warn().Msg("Failed creating an empty HeapDoublesSketch")
					continue
				}
				datasketches[rollup.ToField] = ds
			}
			datasketch := datasketches[rollup.ToField]
			err := datasketch.Update(newValue)
			if err != nil {
				log.Warn().Float64("value", newValue).Msg("Failed updating datasketch with value")
			}
		}
	}
	// Exclude list
	for fromField := range rp.rollupFromFields {
		baseAttributes.Remove(fromField)
	}
}

func (rp *rollupProcessor) setLogsNextConsumer(c consumer.Logs) {
	rp.logsNextConsumer = c
}

func (rp *rollupProcessor) exportLogs(ctx context.Context, rollupData map[string]pcommon.Map) error {
	ld := plog.NewLogs()
	logs := ld.ResourceLogs().AppendEmpty().ScopeLogs().AppendEmpty().LogRecords()
	for _, v := range rollupData {
		logRecord := logs.AppendEmpty()
		// TODO tgill: need to get timestamp from v
		logRecord.SetTimestamp(pcommon.NewTimestampFromTime(time.Now().UTC()))
		v.CopyTo(logRecord.Attributes())
		rawCount, _ := v.Get(RollupCountKey)
		hist, err := rp.rollupHistogram.GetMetricWith(nil)
		if err != nil {
			log.Debug().Msgf("Could not extract rollup histogram metric from registry: %v", err)
		} else {
			hist.Observe(float64(rawCount.Int()))
		}
	}
	log.Trace().Int("count", ld.LogRecordCount()).Msg("After rollup")
	return rp.logsNextConsumer.ConsumeLogs(ctx, ld)
}

// newRollupLogsProcessor creates a new rollup processor that rollupes logs.
func newRollupLogsProcessor(set processor.CreateSettings, next consumer.Logs, cfg *Config) (*rollupProcessor, error) {
	rp, err := newRollupProcessor(set, cfg)
	if err != nil {
		return nil, err
	}
	rp.setLogsNextConsumer(next)
	return rp, nil
}
