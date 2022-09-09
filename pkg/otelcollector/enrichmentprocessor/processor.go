package enrichmentprocessor

import (
	"context"
	"fmt"
	"strings"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"

	"github.com/fluxninja/aperture/pkg/entitycache"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/otelcollector"
)

type enrichmentProcessor struct {
	cache *entitycache.EntityCache
}

func newProcessor(cache *entitycache.EntityCache) *enrichmentProcessor {
	return &enrichmentProcessor{
		cache: cache,
	}
}

// Capabilities returns the capabilities of the processor with MutatesData set to true.
func (ep *enrichmentProcessor) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{
		MutatesData: true,
	}
}

// Start indicates and logs the start of the enrichment processor.
func (ep *enrichmentProcessor) Start(_ context.Context, _ component.Host) error {
	log.Debug().Msg("enrichment processor start")
	return nil
}

// Shutdown indicates and logs the shutdown of the enrichment processor.
func (ep *enrichmentProcessor) Shutdown(context.Context) error {
	log.Debug().Msg("enrichment processor shutdown")
	return nil
}

// ConsumeLogs receives plog.Logs for consumption then returns the enriched logs.
func (ep *enrichmentProcessor) ConsumeLogs(ctx context.Context, origLd plog.Logs) (plog.Logs, error) {
	if ep.cache == nil {
		return plog.Logs{}, fmt.Errorf("cache not provided")
	}
	ld := origLd.Clone()
	err := otelcollector.IterateLogRecords(ld, func(logRecord plog.LogRecord) error {
		enforceIncludeListLog(logRecord.Attributes())
		ep.enrichAttributes(logRecord.Attributes(), []string{otelcollector.EnvoyMissingAttributeSourceValue})
		enforceExcludeListLog(logRecord.Attributes())
		return nil
	})
	return ld, err
}

// ConsumeTraces receives ptrace.Traces for consumption then returns the enriched traces.
func (ep *enrichmentProcessor) ConsumeTraces(ctx context.Context, origTd ptrace.Traces) (ptrace.Traces, error) {
	if ep.cache == nil {
		return ptrace.Traces{}, fmt.Errorf("cache not provided")
	}
	td := origTd.Clone()
	err := otelcollector.IterateSpans(td, func(span ptrace.Span) error {
		enforceIncludeListSpan(span.Attributes())
		ep.enrichAttributes(span.Attributes(), []string{})
		enforceExcludeListSpan(span.Attributes())
		return nil
	})
	return td, err
}

// ConsumeMetrics receives pmetric.Metrics for consumption then returns the enriched metrics.
func (ep *enrichmentProcessor) ConsumeMetrics(ctx context.Context, origMd pmetric.Metrics) (pmetric.Metrics, error) {
	if ep.cache == nil {
		return pmetric.Metrics{}, fmt.Errorf("cache not provided")
	}
	md := origMd.Clone()
	err := otelcollector.IterateMetrics(md, func(metric pmetric.Metric) error {
		switch metric.DataType() {
		case pmetric.MetricDataTypeGauge:
			dataPoints := metric.Gauge().DataPoints()
			for dpIt := 0; dpIt < dataPoints.Len(); dpIt++ {
				dp := dataPoints.At(dpIt)
				ep.enrichMetrics(dp.Attributes())
			}
		case pmetric.MetricDataTypeSum:
			dataPoints := metric.Sum().DataPoints()
			for dpIt := 0; dpIt < dataPoints.Len(); dpIt++ {
				dp := dataPoints.At(dpIt)
				ep.enrichMetrics(dp.Attributes())
			}
		case pmetric.MetricDataTypeSummary:
			dataPoints := metric.Summary().DataPoints()
			for dpIt := 0; dpIt < dataPoints.Len(); dpIt++ {
				dp := dataPoints.At(dpIt)
				ep.enrichMetrics(dp.Attributes())
			}
		case pmetric.MetricDataTypeHistogram:
			dataPoints := metric.Histogram().DataPoints()
			for dpIt := 0; dpIt < dataPoints.Len(); dpIt++ {
				dp := dataPoints.At(dpIt)
				ep.enrichMetrics(dp.Attributes())
			}
		case pmetric.MetricDataTypeExponentialHistogram:
			dataPoints := metric.ExponentialHistogram().DataPoints()
			for dpIt := 0; dpIt < dataPoints.Len(); dpIt++ {
				dp := dataPoints.At(dpIt)
				ep.enrichMetrics(dp.Attributes())
			}
		}
		return nil
	})
	return md, err
}

func (ep *enrichmentProcessor) enrichAttributes(attributes pcommon.Map, treatAsMissing []string) {
	// TODO tgill: split this into multiple functions, one for each source
	unpackFlowLabels(attributes, treatAsMissing)
	var hostIP string
	controlPoint, exists := attributes.Get(otelcollector.ControlPointLabel)
	if !exists {
		otelcollector.LogSampled.Warn().Msg("Skipping because 'otelcollector.ControlPointLabel' attribute not found")
		return
	}
	switch controlPoint.AsString() {
	case otelcollector.ControlPointEgress:
		rawHostAddress, exists := attributes.Get(otelcollector.HostAddressLabel)
		if !exists {
			log.Warn().Msg("Skipping because 'otelcollector.HostAddressLabel' attribute not found")
			return
		}
		hostAddress := rawHostAddress.StringVal()
		if len(hostAddress) == 0 || hostAddress == otelcollector.EnvoyMissingAttributeSourceValue {
			otelcollector.LogSampled.Warn().Msg("Skipping because 'otelcollector.HostAddressLabel' is empty")
			return
		}
		hostIP = ipFromAddress(rawHostAddress.StringVal())
	case otelcollector.ControlPointIngress:
		rawHostIP, exists := attributes.Get(otelcollector.HostIPLabel)
		if !exists {
			otelcollector.LogSampled.Warn().Msg("Skipping because 'otelcollector.HostIPLabel' attribute not found")
			return
		}
		hostIP = rawHostIP.StringVal()
	case otelcollector.ControlPointFeature:
		featureAddress, exists := attributes.Get(otelcollector.FeatureAddressLabel)
		if !exists {
			otelcollector.LogSampled.Warn().Msg("Skipping because 'otelcollector.FeatureAddressLabel' attribute not found")
			return
		}
		hostIP = featureAddress.StringVal()
	default:
		otelcollector.LogSampled.Warn().Str(otelcollector.ControlPointLabel, controlPoint.AsString()).Msg("Unknown control point")
		return
	}

	hostEntity := ep.cache.GetByIP(hostIP)
	if hostEntity == nil {
		otelcollector.LogSampled.Trace().Str("ip", hostIP).Msg("Skipping because entity not found in cache")
		return
	}
	servicesValue := pcommon.NewValueSlice()
	for _, service := range hostEntity.Services {
		servicesValue.SliceVal().AppendEmpty().SetStringVal(service)
	}
	attributes.Upsert(otelcollector.ServicesLabel, servicesValue)
}

func (ep *enrichmentProcessor) enrichMetrics(attributes pcommon.Map) {
	hostNamex, ok := attributes.Get(otelcollector.EntityNameLabel)
	if !ok {
		return
	}
	hostName := hostNamex.StringVal()
	hostEntity := ep.cache.GetByName(hostName)
	attributes.Remove(otelcollector.EntityNameLabel)
	if hostEntity == nil {
		log.Trace().Str("name", hostName).Msg("Skipping because entity not found in cache")
		return
	}
	servicesValue := pcommon.NewValueSlice()
	for _, service := range hostEntity.Services {
		servicesValue.SliceVal().AppendEmpty().SetStringVal(service)
	}
	attributes.Upsert(otelcollector.ServicesLabel, servicesValue)
}

// unpackFlowLabels tries to parse `LabelsLabel` attribute as json, and adds
// unmarshalled attributes to given map.
func unpackFlowLabels(attributes pcommon.Map, treatAsMissing []string) {
	labeled := "false"
	defer func() {
		attributes.UpsertString(otelcollector.LabeledLabel, labeled)
	}()

	var flowAttributes map[string]string
	otelcollector.GetStruct(attributes, otelcollector.MarshalledLabelsLabel, &flowAttributes, treatAsMissing)
	for k, v := range flowAttributes {
		labeled = "true"
		// FIXME – this is quadratic (every upsert iterates to search whether label already exists)
		attributes.UpsertString(k, v)
	}
}

func ipFromAddress(ip string) string {
	return strings.Split(ip, ":")[0]
}

/*
 * IncludeList: This IncludeList is applied to logs and spans at the beginning of enrichment process.
 */
var (
	_includeAttributesCommon = []string{
		otelcollector.MarshalledLabelsLabel,
		otelcollector.ControlPointLabel,
		otelcollector.MarshalledCheckResponseLabel,
	}

	_includeAttributesLog = []string{
		otelcollector.DurationLabel,
		otelcollector.MarshalledAuthzResponseLabel,
		otelcollector.HTTPStatusCodeLabel,
		otelcollector.HTTPRequestContentLength,
		otelcollector.HTTPResponseContentLength,
		otelcollector.HTTPMethodLabel,
		otelcollector.HTTPTargetLabel,
		otelcollector.HTTPFlavorLabel,
		otelcollector.HTTPUserAgentLabel,
		otelcollector.HTTPHostLabel,
		otelcollector.HostAddressLabel,
		otelcollector.PeerAddressLabel,
		otelcollector.HostIPLabel,
		otelcollector.PeerIPLabel,
		otelcollector.EnvoyDurationLabel,
		otelcollector.EnvoyRequestDurationLabel,
		otelcollector.EnvoyRequestTxDurationLabel,
		otelcollector.EnvoyResponseDurationLabel,
		otelcollector.EnvoyResponseTxDurationLabel,
		otelcollector.EnvoyCallerLabel,
	}

	_includeAttributesSpan = []string{
		otelcollector.FeatureAddressLabel,
		otelcollector.FeatureIDLabel,
		otelcollector.FeatureStatusLabel,
	}

	includeListLog  = otelcollector.FormIncludeList(append(_includeAttributesCommon, _includeAttributesLog...))
	includeListSpan = otelcollector.FormIncludeList(append(_includeAttributesCommon, _includeAttributesSpan...))
)

func enforceIncludeListLog(attributes pcommon.Map) {
	otelcollector.EnforceIncludeList(attributes, includeListLog)
}

func enforceIncludeListSpan(attributes pcommon.Map) {
	otelcollector.EnforceIncludeList(attributes, includeListSpan)
}

var (
	_excludeAttributesCommon = []string{
		otelcollector.MarshalledLabelsLabel,
	}

	_excludeAttributesLog = []string{
		otelcollector.HostAddressLabel,
		otelcollector.PeerAddressLabel,
		otelcollector.HostIPLabel,
		otelcollector.PeerIPLabel,
		otelcollector.EnvoyCallerLabel,
	}

	_excludeAttributesSpan = []string{
		otelcollector.FeatureAddressLabel,
	}

	excludeListLog  = otelcollector.FormExcludeList(append(_excludeAttributesCommon, _excludeAttributesLog...))
	excludeListSpan = otelcollector.FormExcludeList(append(_excludeAttributesCommon, _excludeAttributesSpan...))
)

func enforceExcludeListLog(attributes pcommon.Map) {
	otelcollector.EnforceExcludeList(attributes, excludeListLog)
}

func enforceExcludeListSpan(attributes pcommon.Map) {
	otelcollector.EnforceExcludeList(attributes, excludeListSpan)
}
