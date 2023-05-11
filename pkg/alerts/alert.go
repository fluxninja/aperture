package alerts

import (
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/prometheus/alertmanager/api/v2/models"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"

	"github.com/fluxninja/aperture/v2/pkg/log"
	otelconsts "github.com/fluxninja/aperture/v2/pkg/otelcollector/consts"
)

type alertSeverity string

func (a alertSeverity) String() string { return string(a) }

const (
	// SeverityCrit describes an alert which requires immediate action.
	SeverityCrit alertSeverity = "crit"

	// SeverityWarn describes an alert which requires further observation.
	SeverityWarn alertSeverity = "warn"

	// SeverityInfo describes an alert which has informational purposes.
	SeverityInfo alertSeverity = "info"

	// SeverityUnknown describes an alert which does not have severity set.
	SeverityUnknown alertSeverity = ""
)

// ParseSeverity returns alert severity parsed from string. Returns SeverityUnknown
// if parsing fails.
func ParseSeverity(rawSeverity string) alertSeverity {
	s := alertSeverity(rawSeverity)
	switch s {
	case SeverityCrit, SeverityWarn, SeverityInfo:
		return s
	default:
		return SeverityUnknown
	}
}

// specialLabels are alert labels which are propagated in dedicated fields in OTel logs.
var specialLabels = map[string]struct{}{
	otelconsts.AlertNameLabel:         {},
	otelconsts.AlertSeverityLabel:     {},
	otelconsts.AlertGeneratorURLLabel: {},
}

// AlertOption is a type for constructor options.
type AlertOption func(*Alert)

// NewAlert creates new instance of Alert with StartsAt set to now.
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

// NewAlertFromPostableAlert creates new alert with given PostableAlert.
func NewAlertFromPostableAlert(postabelAlert models.PostableAlert) *Alert {
	return &Alert{
		postableAlert: postabelAlert,
	}
}

// Alert is a wrapper around models.PostableAlert with handy transform methods.
type Alert struct {
	postableAlert models.PostableAlert
}

// Name gets the alert name from labels. Returns empty string if label not found.
func (a *Alert) Name() string {
	return a.postableAlert.Labels[otelconsts.AlertNameLabel]
}

// SetName sets the alert name in labels. Overwrites previous value if exists.
func (a *Alert) SetName(name string) {
	a.postableAlert.Labels[otelconsts.AlertNameLabel] = name
}

// WithName is an option function for constructor.
func WithName(name string) AlertOption {
	return func(a *Alert) {
		a.SetName(name)
	}
}

// Severity gets the alert severity from labels. Returns empty string if label not found.
func (a *Alert) Severity() alertSeverity {
	raw, ok := a.postableAlert.Labels[otelconsts.AlertSeverityLabel]
	if !ok {
		return SeverityUnknown
	}
	return ParseSeverity(raw)
}

// SetSeverity sets the alert severity in labels. Overwrites previous value if exists.
func (a *Alert) SetSeverity(severity alertSeverity) {
	a.postableAlert.Labels[otelconsts.AlertSeverityLabel] = severity.String()
}

// WithSeverity is an option function for constructor.
func WithSeverity(severity alertSeverity) AlertOption {
	return func(a *Alert) {
		a.SetSeverity(severity)
	}
}

// AlertChannels gets the alert channels from labels. Returns empty slice if label not found.
func (a *Alert) AlertChannels() []string {
	channels, ok := a.postableAlert.Labels[otelconsts.AlertChannelsLabel]
	if !ok {
		return []string{}
	}
	return strings.Split(channels, ",")
}

// SetAlertChannels sets the alert channels in labels. Overwrites previous value if exists.
func (a *Alert) SetAlertChannels(alertChannels []string) {
	a.postableAlert.Labels[otelconsts.AlertChannelsLabel] = strings.Join(alertChannels, ",")
}

// WithAlertChannels is an option function for constructor.
func WithAlertChannels(alertChannels []string) AlertOption {
	return func(a *Alert) {
		a.SetAlertChannels(alertChannels)
	}
}

// PostableAlert returns the underlying PostableAlert struct.
func (a *Alert) PostableAlert() models.PostableAlert {
	return a.postableAlert
}

// SetAnnotation sets a single annotation. It overwrites the previous value if exists.
func (a *Alert) SetAnnotation(key, value string) {
	a.postableAlert.Annotations[key] = value
}

// WithAnnotation is an option function for constructor.
func WithAnnotation(key, value string) AlertOption {
	return func(a *Alert) {
		a.SetAnnotation(key, value)
	}
}

// SetLabel sets a single label. It overwrites the previous value if exists.
func (a *Alert) SetLabel(key, value string) {
	a.postableAlert.Labels[key] = value
}

// WithLabel is an option function for constructor.
func WithLabel(key, value string) AlertOption {
	return func(a *Alert) {
		a.SetLabel(key, value)
	}
}

// SetGeneratorURL sets a generator URL. It overwrites the previous value if exists.
func (a *Alert) SetGeneratorURL(value string) {
	a.postableAlert.GeneratorURL = strfmt.URI(value)
}

// WithGeneratorURL is an option function for constructor.
func WithGeneratorURL(value string) AlertOption {
	return func(a *Alert) {
		a.SetGeneratorURL(value)
	}
}

// SetResolveTimeout sets a resolve timeout which says when given alert becomes resolved.
func (a *Alert) SetResolveTimeout(t time.Duration) {
	a.postableAlert.EndsAt = strfmt.DateTime(time.Time(a.postableAlert.StartsAt).Add(t))
}

// WithResolveTimeout is an option function for constructor.
func WithResolveTimeout(t time.Duration) AlertOption {
	return func(a *Alert) {
		a.SetResolveTimeout(t)
	}
}

// AlertsFromLogs gets slice of alerts from OTel Logs.
func AlertsFromLogs(ld plog.Logs) []*Alert {
	// We cannot preallocate size, as we do not know how many of those log records
	// has incorrect data and will be dropped.
	alerts := []*Alert{}
	resourceLogsSlice := ld.ResourceLogs()
	for resourceLogsIt := 0; resourceLogsIt < resourceLogsSlice.Len(); resourceLogsIt++ {
		resourceLogs := resourceLogsSlice.At(resourceLogsIt)
		resourceAttributes := resourceLogs.Resource().Attributes()
		generatorURL, exists := resourceAttributes.Get(otelconsts.AlertGeneratorURLLabel)
		if !exists {
			log.Trace().
				Str("key", otelconsts.AlertGeneratorURLLabel).Msg("Key not found")
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
				a.postableAlert.EndsAt = strfmt.DateTime(logRecord.ObservedTimestamp().AsTime())
				a.postableAlert.GeneratorURL = strfmt.URI(generatorURL.AsString())
				a.postableAlert.Labels = models.LabelSet(mapFromAttributes(resourceAttributes, specialLabels))
				a.SetSeverity(ParseSeverity(logRecord.SeverityText()))
				a.SetName(logRecord.Body().AsString())
				attributes := logRecord.Attributes()
				a.postableAlert.Annotations = models.LabelSet(mapFromAttributes(attributes, map[string]struct{}{}))
				alerts = append(alerts, a)
			}
		}
	}
	return alerts
}

// AsLogs returns alert as OTel Logs.
func (a *Alert) AsLogs() plog.Logs {
	ld := plog.NewLogs()
	resource := ld.ResourceLogs().AppendEmpty()
	resourceAttributes := resource.Resource().Attributes()
	// Labels in AM are used to identify identical instances of an alert. This corresponds
	// with the resource notion in OTLP protocol, which describes the source of a log.
	populateAttributesFromMap(resourceAttributes, a.postableAlert.Labels, specialLabels)
	resourceAttributes.PutStr(otelconsts.AlertGeneratorURLLabel, string(a.postableAlert.GeneratorURL))
	resourceAttributes.PutBool(otelconsts.IsAlertLabel, true)

	logRecord := resource.ScopeLogs().AppendEmpty().LogRecords().AppendEmpty()
	logRecord.SetTimestamp(pcommon.NewTimestampFromTime(time.Time(a.postableAlert.StartsAt)))
	logRecord.SetObservedTimestamp(pcommon.NewTimestampFromTime(time.Time(a.postableAlert.EndsAt)))
	logRecord.SetSeverityText(a.Severity().String())
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
