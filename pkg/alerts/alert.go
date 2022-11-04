package alerts

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"

	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/otelcollector"
)

// Alert TODO as we will switch to prom client struct anyway.
type Alert struct {
	StartsAt     string            `json:"startsAt"`
	Annotations  map[string]string `json:"annotations"`
	Labels       map[string]string `json:"labels"`
	GeneratorURL string            `json:"generatorURL"`
}

// AlertsFromLogs gets slice of alerts from OTEL Logs.
func AlertsFromLogs(ld plog.Logs) []*Alert {
	// We can't preallocate size, as we don't know how many of those log records
	// has incorrect data and will be dropped.
	alerts := []*Alert{}
	err := otelcollector.IterateLogRecords(ld, func(lr plog.LogRecord) error {
		a := &Alert{}
		attributes := lr.Attributes()
		startsAt, exists := attributes.Get(otelcollector.AlertStartsAtLabel)
		if !exists {
			log.Sample(zerolog.Sometimes).Trace().
				Str("key", otelcollector.AlertStartsAtLabel).Msg("Key not found")
			return nil
		}
		a.StartsAt = startsAt.AsString()
		generatorURL, exists := attributes.Get(otelcollector.AlertGeneratorURLLabel)
		if !exists {
			log.Sample(zerolog.Sometimes).Trace().
				Str("key", otelcollector.AlertGeneratorURLLabel).Msg("Key not found")
			return nil
		}
		a.GeneratorURL = generatorURL.AsString()
		annotations := map[string]string{}
		labels := map[string]string{}
		attributes.Range(func(k string, v pcommon.Value) bool {
			if strings.HasPrefix(k, otelcollector.AlertAnnotationsLabelPrefix) {
				trimmed := strings.TrimPrefix(k, otelcollector.AlertAnnotationsLabelPrefix)
				annotations[trimmed] = v.AsString()
			}
			if strings.HasPrefix(k, otelcollector.AlertLabelsLabelPrefix) {
				trimmed := strings.TrimPrefix(k, otelcollector.AlertLabelsLabelPrefix)
				labels[trimmed] = v.AsString()
			}
			return true
		})
		a.Annotations = annotations
		a.Labels = labels
		alerts = append(alerts, a)
		return nil
	})
	if err != nil {
		// This should not happen, as we don't return error anywhere in the IterateLogs func.
		log.Sample(zerolog.Sometimes).Trace().Err(err).Msg("Getting alerts from logs")
	}
	return alerts
}

// AsLogs returns alert as OTEL Logs.
func (a *Alert) AsLogs() plog.Logs {
	ld := plog.NewLogs()
	logRecord := ld.
		ResourceLogs().AppendEmpty().
		ScopeLogs().AppendEmpty().
		LogRecords().AppendEmpty()
	attributes := logRecord.Attributes()
	attributes.PutStr(otelcollector.AlertStartsAtLabel, a.StartsAt)
	attributes.PutStr(otelcollector.AlertGeneratorURLLabel, a.GeneratorURL)
	populateMapWithPrefix(attributes, a.Labels, otelcollector.AlertLabelsLabelPrefix)
	populateMapWithPrefix(attributes, a.Annotations, otelcollector.AlertAnnotationsLabelPrefix)
	return ld
}

func populateMapWithPrefix(attributes pcommon.Map, values map[string]string, prefix string) {
	for k, v := range values {
		attributes.PutStr(fmt.Sprintf("%v%v", prefix, k), v)
	}
}
