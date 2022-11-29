package alerts

import (
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/prometheus/alertmanager/api/v2/models"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"

	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/otelcollector"
)

// specialLabels are alert labels which are propagated in dedicated fields in OTEL logs.
var specialLabels = map[string]struct{}{
	otelcollector.AlertNameLabel:         {},
	otelcollector.AlertSeverityLabel:     {},
	otelcollector.AlertGeneratorURLLabel: {},
}

// AlertOption is a type for constructor options.
type AlertOption func(*Alert)

// NewAlert created new instance of Alert with StartsAt set to now.
func NewAlert(opts ...AlertOption) *Alert {
	newAlert := &Alert{
		postableAlert: models.PostableAlert{
			Alert: models.Alert{
				Labels: models.LabelSet(map[string]string{}),
			},
			Annotations: models.LabelSet(map[string]string{}),
			StartsAt:    strfmt.DateTime(time.Now().UTC()),
		},
	}

	for _, opt := range opts {
		opt(newAlert)
	}
	return newAlert
}

// Alert is a wrapper around models.PostableAlert with handy transform methods.
type Alert struct {
	postableAlert models.PostableAlert
}

// Name gets the alert name from labels. Returns empty string if label not found.
func (a *Alert) Name() string {
	return a.postableAlert.Labels[otelcollector.AlertNameLabel]
}

// SetName sets the alert name in labels. Overwrites previous value if exists.
func (a *Alert) SetName(name string) {
	a.postableAlert.Labels[otelcollector.AlertNameLabel] = name
}

// WithName is an option function for constructor.
func WithName(name string) AlertOption {
	return func(a *Alert) {
		a.postableAlert.Labels[otelcollector.AlertNameLabel] = name
	}
}

// Severity gets the alert severity from labels. Returns empty string if label not found.
func (a *Alert) Severity() string {
	return a.postableAlert.Labels[otelcollector.AlertSeverityLabel]
}

// SetSeverity sets the alert severity in labels. Overwrites previous value if exists.
func (a *Alert) SetSeverity(severity string) {
	a.postableAlert.Labels[otelcollector.AlertSeverityLabel] = severity
}

// WithSeverity is an option function for constructor.
func WithSeverity(severity string) AlertOption {
	return func(a *Alert) {
		a.postableAlert.Labels[otelcollector.AlertSeverityLabel] = severity
	}
}

// SetAnnotation sets a single annotation. It overwrites the previous value if exists.
func (a *Alert) SetAnnotation(key, value string) {
	a.postableAlert.Annotations[key] = value
}

// WithAnnotation is an option function for constructor.
func WithAnnotation(key, value string) AlertOption {
	return func(a *Alert) {
		a.postableAlert.Annotations[key] = value
	}
}

// SetLabel sets a single label. It overwrites the previous value if exists.
func (a *Alert) SetLabel(key, value string) {
	a.postableAlert.Labels[key] = value
}

// WithLabel is an option function for constructor.
func WithLabel(key, value string) AlertOption {
	return func(a *Alert) {
		a.postableAlert.Labels[key] = value
	}
}

// AlertsFromLogs gets slice of alerts from OTEL Logs.
func AlertsFromLogs(ld plog.Logs) []*Alert {
	// We can't preallocate size, as we don't know how many of those log records
	// has incorrect data and will be dropped.
	alerts := []*Alert{}
	resourceLogsSlice := ld.ResourceLogs()
	for resourceLogsIt := 0; resourceLogsIt < resourceLogsSlice.Len(); resourceLogsIt++ {
		resourceLogs := resourceLogsSlice.At(resourceLogsIt)
		resourceAttributes := resourceLogs.Resource().Attributes()
		generatorURL, exists := resourceAttributes.Get(otelcollector.AlertGeneratorURLLabel)
		if !exists {
			log.Trace().
				Str("key", otelcollector.AlertGeneratorURLLabel).Msg("Key not found")
			return nil
		}
		scopeLogsSlice := resourceLogs.ScopeLogs()
		for scopeLogsIt := 0; scopeLogsIt < scopeLogsSlice.Len(); scopeLogsIt++ {
			scopeLogs := scopeLogsSlice.At(scopeLogsIt)
			logsSlice := scopeLogs.LogRecords()
			for logsIt := 0; logsIt < logsSlice.Len(); logsIt++ {
				logRecord := logsSlice.At(logsIt)
				a := &Alert{}
				a.postableAlert.StartsAt = strfmt.DateTime(logRecord.Timestamp().AsTime())
				a.postableAlert.GeneratorURL = strfmt.URI(generatorURL.AsString())
				a.postableAlert.Labels = models.LabelSet(mapFromAttributes(resourceAttributes, specialLabels))
				a.SetSeverity(logRecord.SeverityText())
				a.SetName(logRecord.Body().AsString())
				attributes := logRecord.Attributes()
				a.postableAlert.Annotations = models.LabelSet(mapFromAttributes(attributes, map[string]struct{}{}))
				alerts = append(alerts, a)
			}
		}
	}
	return alerts
}

// AsLogs returns alert as OTEL Logs.
func (a *Alert) AsLogs() plog.Logs {
	ld := plog.NewLogs()
	resource := ld.ResourceLogs().AppendEmpty()
	resourceAttributes := resource.Resource().Attributes()
	// Labels in AM are used to identify identical instances of an alert. This corresponds
	// with the resource notion in OTLP protocol, which describes the source of a log.
	populateAttributesFromMap(resourceAttributes, a.postableAlert.Labels, specialLabels)
	resourceAttributes.PutStr(otelcollector.AlertGeneratorURLLabel, string(a.postableAlert.GeneratorURL))

	logRecord := resource.ScopeLogs().AppendEmpty().LogRecords().AppendEmpty()
	logRecord.SetTimestamp(pcommon.NewTimestampFromTime(time.Time(a.postableAlert.StartsAt)))
	logRecord.SetSeverityText(a.Severity())
	pcommon.NewValueStr(a.Name()).CopyTo(logRecord.Body())

	attributes := logRecord.Attributes()
	populateAttributesFromMap(attributes, a.postableAlert.Annotations, map[string]struct{}{})
	return ld
}

func populateAttributesFromMap(attributes pcommon.Map, values map[string]string, ignore map[string]struct{}) {
	for k, v := range values {
		if _, ok := ignore[k]; ok {
			continue
		}
		attributes.PutStr(k, v)
	}
}

func mapFromAttributes(attributes pcommon.Map, ignore map[string]struct{}) map[string]string {
	result := map[string]string{}
	attributes.Range(func(k string, v pcommon.Value) bool {
		if _, exists := ignore[k]; exists {
			return true
		}
		result[k] = v.AsString()
		return true
	})
	return result
}
