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

	"aperture.tech/aperture/pkg/entitycache"
	"aperture.tech/aperture/pkg/log"
	"aperture.tech/aperture/pkg/otelcollector"
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
		ep.enrichAttributes(logRecord.Attributes())
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
		ep.enrichAttributes(span.Attributes())
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

func (ep *enrichmentProcessor) enrichAttributes(attributes pcommon.Map) {
	unpackFlowLabels(attributes)
	var hostIP string
	controlPoint, exists := attributes.Get(otelcollector.ControlPointLabel)
	if !exists {
		log.Warn().Msg("Skipping because 'otelcollector.ControlPointLabel' attribute not found")
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
		if len(hostAddress) == 0 || hostAddress == otelcollector.MissingAttributeSourceValue {
			log.Warn().Msg("Skipping because 'otelcollector.HostAddressLabel' is empty")
			return
		}
		hostIP = ipFromAddress(rawHostAddress.StringVal())
		attributes.Remove(otelcollector.HostAddressLabel)
		attributes.Remove(otelcollector.PeerAddressLabel)
	case otelcollector.ControlPointIngress:
		rawHostIP, exists := attributes.Get(otelcollector.HostIPLabel)
		if !exists {
			log.Warn().Msg("Skipping because 'otelcollector.HostIPLabel' attribute not found")
			return
		}
		hostIP = rawHostIP.StringVal()
		attributes.Remove(otelcollector.HostIPLabel)
		attributes.Remove(otelcollector.PeerIPLabel)
	case otelcollector.ControlPointFeature:
		featureAddress, exists := attributes.Get(otelcollector.FeatureAddressLabel)
		if !exists {
			log.Warn().Msg("Skipping because 'otelcollector.FeatureAddressLabel' attribute not found")
			return
		}
		hostIP = featureAddress.StringVal()
		attributes.Remove(otelcollector.FeatureAddressLabel)
	default:
		log.Warn().Str(otelcollector.ControlPointLabel, controlPoint.AsString()).Msg("Unknown control point")
		return
	}
	hostEntity := ep.cache.GetByIP(hostIP)
	if hostEntity == nil {
		log.Trace().Str("ip", hostIP).Msg("Skipping because entity not found in cache")
		return
	}

	attributes.UpsertString(otelcollector.AgentGroupLabel, hostEntity.AgentGroup)
	attributes.UpsertString(otelcollector.NamespaceLabel, hostEntity.Namespace)
	attributes.UpsertString(otelcollector.ServicesLabel, strings.Join(hostEntity.Services, ","))
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
	attributes.UpsertString(otelcollector.AgentGroupLabel, hostEntity.AgentGroup)
	attributes.UpsertString(otelcollector.NamespaceLabel, hostEntity.Namespace)
	attributes.UpsertString(otelcollector.ServicesLabel, strings.Join(hostEntity.Services, ","))
}

// unpackFlowLabels tries to parse `fn.flow` attribute as json, and adds unmarshalled
// attributes as `fn.flow.<name>`.
func unpackFlowLabels(attributes pcommon.Map) {
	labeled := "false"
	defer func() {
		attributes.UpsertString(otelcollector.LabeledLabel, labeled)
	}()
	rawFlow, exists := attributes.Get(otelcollector.FlowLabel)
	if !exists {
		return
	}
	defer attributes.Remove(otelcollector.FlowLabel)

	var flowAttributes map[string]string
	otelcollector.UnmarshalStringVal(rawFlow, otelcollector.FlowLabel, &flowAttributes)
	for fnKey, fnValue := range flowAttributes {
		labeled = "true"
		// FIXME – this is quadratic (every upsert iterates to search whether label already exists)
		attributes.UpsertString(fmt.Sprintf("%s.%s", otelcollector.FlowLabel, fnKey), fnValue)
	}
}

func ipFromAddress(ip string) string {
	return strings.Split(ip, ":")[0]
}
