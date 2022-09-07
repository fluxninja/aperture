package rollupprocessor

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/fluxninja/datasketches-go/sketches"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/ptrace"

	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/otelcollector"
)

const (
	rollupCountKey = "rollup_count"
)

type rollupProcessor struct {
	cfg *Config

	logsNextConsumer   consumer.Logs
	tracesNextConsumer consumer.Traces
}

var (
	_ consumer.Traces = (*rollupProcessor)(nil)
	_ consumer.Logs   = (*rollupProcessor)(nil)
)

func newRollupProcessor(set component.ProcessorCreateSettings, cfg *Config) (*rollupProcessor, error) {
	return &rollupProcessor{
		cfg: cfg,
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

// ConsumeTraces implements TracesProcessor.
func (rp *rollupProcessor) ConsumeTraces(ctx context.Context, td ptrace.Traces) error {
	rollupData := make(map[string]pcommon.Map)
	datasketches := make(map[string]map[string]*sketches.HeapDoublesSketch)
	err := otelcollector.IterateSpans(td, func(span ptrace.Span) error {
		key := rp.key(span.Attributes())
		_, exists := rollupData[key]
		if !exists {
			rollupData[key] = span.Attributes()
			rollupData[key].InsertInt(rollupCountKey, 0)
		}
		_, exists = datasketches[key]
		if !exists {
			datasketches[key] = make(map[string]*sketches.HeapDoublesSketch)
		}
		rawCount, _ := rollupData[key].Get(rollupCountKey)
		rollupData[key].UpdateInt(rollupCountKey, rawCount.IntVal()+1)
		rp.rollupAttributes(datasketches[key], rollupData[key], span.Attributes(), rp.cfg.RollupsSpan)
		return nil
	})
	if err != nil {
		return err
	}
	for k, v := range datasketches {
		attributes := rollupData[k]
		for toField, sketch := range v {
			serializedBytes, err := sketch.Compact().Serialize()
			if err != nil {
				return err
			}
			serialized := base64.StdEncoding.EncodeToString(serializedBytes)
			attributes.UpsertString(toField, serialized)
		}
	}
	return rp.exportTraces(ctx, rollupData)
}

// ConsumeLogs implements LogsProcessor.
func (rp *rollupProcessor) ConsumeLogs(ctx context.Context, ld plog.Logs) error {
	rollupData := make(map[string]pcommon.Map)
	datasketches := make(map[string]map[string]*sketches.HeapDoublesSketch)
	err := otelcollector.IterateLogRecords(ld, func(logRecord plog.LogRecord) error {
		key := rp.key(logRecord.Attributes())
		_, exists := rollupData[key]
		if !exists {
			rollupData[key] = logRecord.Attributes()
			rollupData[key].InsertInt(rollupCountKey, 0)
		}
		_, exists = datasketches[key]
		if !exists {
			datasketches[key] = make(map[string]*sketches.HeapDoublesSketch)
		}
		rawCount, _ := rollupData[key].Get(rollupCountKey)
		rollupData[key].UpdateInt(rollupCountKey, rawCount.IntVal()+1)
		rp.rollupAttributes(datasketches[key], rollupData[key], logRecord.Attributes(), rp.cfg.RollupsLog)
		return nil
	})
	if err != nil {
		return err
	}
	for k, v := range datasketches {
		attributes := rollupData[k]
		for toField, sketch := range v {
			serializedBytes, err := sketch.Compact().Serialize()
			if err != nil {
				return err
			}
			serialized := base64.StdEncoding.EncodeToString(serializedBytes)
			attributes.UpsertString(toField, serialized)
		}
	}
	return rp.exportLogs(ctx, rollupData)
}

func (rp *rollupProcessor) rollupAttributes(datasketches map[string]*sketches.HeapDoublesSketch, baseAttributes, attributes pcommon.Map, rollups []*Rollup) {
	for _, rollup := range rollups {
		switch rollup.Type {
		case RollupSum:
			newValue, found := rollup.GetFromFieldValue(attributes)
			if !found {
				continue
			}
			rollupSum, exists := rollup.GetToFieldValue(baseAttributes)
			if !exists {
				rollupSum = 0
				baseAttributes.UpsertDouble(rollup.ToField, rollupSum)
			}
			baseAttributes.UpdateDouble(rollup.ToField, rollupSum+newValue)
		case RollupSumOfSquares:
			newValue, found := rollup.GetFromFieldValue(attributes)
			if !found {
				continue
			}
			rollupSos, exists := rollup.GetToFieldValue(baseAttributes)
			if !exists {
				rollupSos = 0
				baseAttributes.UpsertDouble(rollup.ToField, rollupSos)
			}
			baseAttributes.UpdateDouble(rollup.ToField, rollupSos+newValue*newValue)
		case RollupMin:
			newValue, found := rollup.GetFromFieldValue(attributes)
			if !found {
				continue
			}
			rollupMin, exists := rollup.GetToFieldValue(baseAttributes)
			if !exists {
				rollupMin = newValue
				baseAttributes.UpsertDouble(rollup.ToField, rollupMin)
			}
			newMin := otelcollector.Min(rollupMin, newValue)
			baseAttributes.UpdateDouble(rollup.ToField, newMin)
		case RollupMax:
			newValue, found := rollup.GetFromFieldValue(attributes)
			if !found {
				continue
			}
			rollupMax, exists := rollup.GetToFieldValue(baseAttributes)
			if !exists {
				rollupMax = newValue
				baseAttributes.UpsertDouble(rollup.ToField, rollupMax)
			}
			newMax := otelcollector.Max(rollupMax, newValue)
			baseAttributes.UpdateDouble(rollup.ToField, newMax)
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
	for _, rollup := range rp.cfg.RollupsLog {
		baseAttributes.Remove(rollup.FromField)
	}
}

func (rp *rollupProcessor) setLogsNextConsumer(c consumer.Logs) {
	rp.logsNextConsumer = c
}

func (rp *rollupProcessor) setTracesNextConsumer(c consumer.Traces) {
	rp.tracesNextConsumer = c
}

func (rp *rollupProcessor) exportLogs(ctx context.Context, rollupData map[string]pcommon.Map) error {
	ld := plog.NewLogs()
	logs := ld.ResourceLogs().AppendEmpty().ScopeLogs().AppendEmpty().LogRecords()
	for _, v := range rollupData {
		logRecord := logs.AppendEmpty()
		logRecord.SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))
		v.CopyTo(logRecord.Attributes())
	}
	return rp.logsNextConsumer.ConsumeLogs(ctx, ld)
}

func (rp *rollupProcessor) exportTraces(ctx context.Context, rollupData map[string]pcommon.Map) error {
	ld := ptrace.NewTraces()
	spans := ld.ResourceSpans().AppendEmpty().ScopeSpans().AppendEmpty().Spans()
	for _, v := range rollupData {
		spanRecord := spans.AppendEmpty()
		spanRecord.SetStartTimestamp(pcommon.NewTimestampFromTime(time.Now()))
		v.CopyTo(spanRecord.Attributes())
	}
	return rp.tracesNextConsumer.ConsumeTraces(ctx, ld)
}

// key returns string key used in the hashmap. Current implementations marshals
// the map to JSON. This might be suboptimal.
func (rp *rollupProcessor) key(am pcommon.Map) string {
	raw := am.AsRaw()
	for _, rollup := range rp.cfg.RollupsLog {
		// Removing all fields from which we will get rolled up values, as those
		// are dimensions not to be considered as "key".
		delete(raw, rollup.FromField)
	}
	key, err := json.Marshal(raw)
	if err != nil {
		log.Error().Err(err).Msg("Failed marshaling map to JSON")
		return ""
	}
	return string(key)
}

// newRollupLogsProcessor creates a new rollup processor that rollupes logs.
func newRollupLogsProcessor(set component.ProcessorCreateSettings, next consumer.Logs, cfg *Config) (*rollupProcessor, error) {
	rp, err := newRollupProcessor(set, cfg)
	if err != nil {
		return nil, err
	}
	rp.setLogsNextConsumer(next)
	return rp, nil
}

// newRollupTracesProcessor creates a new rollup processor that rollupes traces.
func newRollupTracesProcessor(set component.ProcessorCreateSettings, next consumer.Traces, cfg *Config) (*rollupProcessor, error) {
	rp, err := newRollupProcessor(set, cfg)
	if err != nil {
		return nil, err
	}
	rp.setTracesNextConsumer(next)
	return rp, nil
}
