package metricsprocessor

import (
	"context"
	"strconv"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/check/v1"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/metrics"
	"github.com/fluxninja/aperture/pkg/otelcollector"
	otelconsts "github.com/fluxninja/aperture/pkg/otelcollector/consts"
	"github.com/fluxninja/aperture/pkg/otelcollector/metricsprocessor/internal"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/iface"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/selectors"
)

type metricsProcessor struct {
	cfg *Config
}

func newProcessor(cfg *Config) (*metricsProcessor, error) {
	p := &metricsProcessor{
		cfg: cfg,
	}

	return p, nil
}

// Start indicates and logs the start of the metrics processor.
func (p *metricsProcessor) Start(_ context.Context, _ component.Host) error {
	log.Debug().Msg("metrics processor start")
	return nil
}

// Shutdown indicates and logs the shutdown of the metrics processor.
func (p *metricsProcessor) Shutdown(context.Context) error {
	log.Debug().Msg("metrics processor shutdown")
	return nil
}

// Capabilities returns the capabilities of the processor with MutatesData set to true.
func (p *metricsProcessor) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{
		MutatesData: true,
	}
}

// ConsumeLogs receives plog.Logs for consumption then returns updated logs with policy labels and metrics.
func (p *metricsProcessor) ConsumeLogs(ctx context.Context, ld plog.Logs) (plog.Logs, error) {
	otelcollector.IterateLogRecords(ld, func(logRecord plog.LogRecord) otelcollector.IterAction {
		// Attributes
		attributes := logRecord.Attributes()

		// CheckResponse
		checkResponse := &flowcontrolv1.CheckResponse{}

		// Source specific processing
		source, exists := attributes.Get(otelconsts.ApertureSourceLabel)
		if !exists {
			log.Sample(noSourceLabelSampler).Warn().Msg("aperture source label not found")
			return otelcollector.Discard
		}
		sourceStr := source.Str()
		if sourceStr == otelconsts.ApertureSourceSDK {
			success := otelcollector.GetStruct(attributes, otelconsts.ApertureCheckResponseLabel, checkResponse, []string{})
			if !success {
				log.Sample(noSDKCheckResponseSampler).Warn().
					Msg("aperture check response label not found in SDK access logs")
				return otelcollector.Discard
			}

			internal.AddSDKSpecificLabels(attributes)
		} else if sourceStr == otelconsts.ApertureSourceEnvoy {
			success := otelcollector.GetStruct(attributes, otelconsts.ApertureCheckResponseLabel, checkResponse, []string{otelconsts.EnvoyMissingAttributeValue})
			if !success {
				log.Sample(noEnvoyCheckResponseSampler).Warn().
					Msg("aperture check response label not found in Envoy access logs")
				return otelcollector.Discard
			}

			internal.AddEnvoySpecificLabels(attributes)
		} else if sourceStr == otelconsts.ApertureSourceLua {
			success := otelcollector.GetStruct(attributes, otelconsts.ApertureCheckResponseLabel, checkResponse, []string{""})
			if !success {
				log.Sample(noEnvoyCheckResponseSampler).Warn().
					Msg("aperture check response label not found in Lua access logs")
				return otelcollector.Discard
			}

			internal.AddLuaSpecificLabels(attributes)
		} else {
			log.Sample(unrecognizedSourceLabelSampler).Warn().
				Msg("aperture source label not recognized")
			return otelcollector.Discard
		}

		statusCode, flowStatus := internal.StatusesFromAttributes(attributes)
		attributes.PutStr(otelconsts.ApertureFlowStatusLabel, internal.FlowStatusForTelemetry(statusCode, flowStatus))
		internal.AddCheckResponseBasedLabels(attributes, checkResponse, sourceStr)
		controlPointType := ""
		telemetryFlowLabels := checkResponse.GetTelemetryFlowLabels()
		if telemetryFlowLabels == nil {
			log.Sample(noTelemetryFlowLabelsSampler).Debug().Msg("aperture telemetry flow labels not found")
		} else {
			controlPointType, exists = telemetryFlowLabels[otelconsts.ApertureControlPointTypeLabel]
			if !exists {
				log.Sample(noControlPointTypeSampler).Debug().Msg("aperture control point type label not found")
			}
		}
		p.populateControlPointCache(checkResponse, controlPointType)

		// Update metrics and enforce include list to eliminate any excess attributes
		if sourceStr == otelconsts.ApertureSourceSDK {
			p.updateMetrics(attributes, checkResponse, []string{})
			internal.EnforceIncludeListSDK(attributes)
		} else if sourceStr == otelconsts.ApertureSourceEnvoy {
			p.updateMetrics(attributes, checkResponse, []string{otelconsts.EnvoyMissingAttributeValue})
			internal.EnforceIncludeListHTTP(attributes)
		} else if sourceStr == otelconsts.ApertureSourceLua {
			p.updateMetrics(attributes, checkResponse, []string{otelconsts.LuaMissingAttributeValue})
			internal.EnforceIncludeListHTTP(attributes)
		}

		// This needs to be called **after** internal.EnforceIncludeList{HTTP,SDK}.
		internal.AddFlowLabels(attributes, checkResponse)
		return otelcollector.Keep
	})
	return ld, nil
}

var (
	noSourceLabelSampler           = log.NewRatelimitingSampler()
	noSDKCheckResponseSampler      = log.NewRatelimitingSampler()
	noEnvoyCheckResponseSampler    = log.NewRatelimitingSampler()
	unrecognizedSourceLabelSampler = log.NewRatelimitingSampler()
	noTelemetryFlowLabelsSampler   = log.NewRatelimitingSampler()
	noControlPointTypeSampler      = log.NewRatelimitingSampler()
)

func (p *metricsProcessor) updateMetrics(attributes pcommon.Map, checkResponse *flowcontrolv1.CheckResponse, treatAsMissing []string) {
	if checkResponse == nil {
		return
	}
	if len(checkResponse.LimiterDecisions) > 0 {
		// Update workload metrics
		latency, latencyFound := otelcollector.GetFloat64(attributes, otelconsts.WorkloadDurationLabel, []string{})
		for _, decision := range checkResponse.LimiterDecisions {
			limiterID := iface.LimiterID{
				PolicyName:  decision.PolicyName,
				PolicyHash:  decision.PolicyHash,
				ComponentID: decision.ComponentId,
			}

			// Update concurrency limiter metrics.
			if cl := decision.GetConcurrencyLimiterInfo(); cl != nil {
				labels := map[string]string{
					metrics.PolicyNameLabel:    decision.PolicyName,
					metrics.PolicyHashLabel:    decision.PolicyHash,
					metrics.ComponentIDLabel:   decision.ComponentId,
					metrics.WorkloadIndexLabel: cl.GetWorkloadIndex(),
				}

				p.updateMetricsForWorkload(limiterID, labels, decision.Dropped, checkResponse.DecisionType, latency, latencyFound)
			}

			// Update rate limiter metrics.
			if rl := decision.GetRateLimiterInfo(); rl != nil {
				labels := map[string]string{
					metrics.PolicyNameLabel:  decision.PolicyName,
					metrics.PolicyHashLabel:  decision.PolicyHash,
					metrics.ComponentIDLabel: decision.ComponentId,
				}
				p.updateMetricsForRateLimiter(limiterID, labels, decision.Dropped, checkResponse.DecisionType)
			}

			// Update flow regulator metrics.
			if fr := decision.GetFlowRegulatorInfo(); fr != nil {
				labels := map[string]string{
					metrics.PolicyNameLabel:  decision.PolicyName,
					metrics.PolicyHashLabel:  decision.PolicyHash,
					metrics.ComponentIDLabel: decision.ComponentId,
				}
				p.updateMetricsForFlowRegulator(limiterID, labels, decision.Dropped, checkResponse.DecisionType)
			}
		}
	}

	if len(checkResponse.FluxMeterInfos) > 0 {
		// Update flux meter metrics
		statusCode, flowStatus := internal.StatusesFromAttributes(attributes)
		for _, fluxMeter := range checkResponse.FluxMeterInfos {
			p.updateMetricsForFluxMeters(
				fluxMeter,
				checkResponse.DecisionType,
				statusCode, flowStatus,
				attributes,
				treatAsMissing)
		}
	}

	if len(checkResponse.ClassifierInfos) > 0 {
		// Update classifier metrics
		for _, classifierInfo := range checkResponse.ClassifierInfos {
			classifierID := iface.ClassifierID{
				PolicyName:      classifierInfo.PolicyName,
				PolicyHash:      classifierInfo.PolicyHash,
				ClassifierIndex: classifierInfo.ClassifierIndex,
			}
			p.updateMetricsForClassifier(classifierID)
		}
	}
}

func (p *metricsProcessor) updateMetricsForWorkload(limiterID iface.LimiterID, labels map[string]string, dropped bool, decisionType flowcontrolv1.CheckResponse_DecisionType, latency float64, latencyFound bool) {
	concurrencyLimiter := p.cfg.engine.GetConcurrencyLimiter(limiterID)
	if concurrencyLimiter == nil {
		log.Sample(noConcurrencyLimiterSampler).Warn().
			Str(metrics.PolicyNameLabel, limiterID.PolicyName).
			Str(metrics.PolicyHashLabel, limiterID.PolicyHash).
			Str(metrics.ComponentIDLabel, limiterID.ComponentID).
			Msg("ConcurrencyLimiter not found")
		return
	}
	// Observe latency only if the latency is found I.E. the request was allowed and response was received
	if latencyFound {
		latencyObserver := concurrencyLimiter.GetLatencyObserver(labels)
		if latencyObserver != nil {
			latencyObserver.Observe(latency)
		}
	}
	// Add decision type label to the request counter metric
	labels[metrics.DecisionTypeLabel] = decisionType.String()
	labels[metrics.LimiterDroppedLabel] = strconv.FormatBool(dropped)
	requestCounter := concurrencyLimiter.GetRequestCounter(labels)
	if requestCounter != nil {
		requestCounter.Inc()
	}
}

func (p *metricsProcessor) updateMetricsForRateLimiter(limiterID iface.LimiterID, labels map[string]string, dropped bool, decisionType flowcontrolv1.CheckResponse_DecisionType) {
	rateLimiter := p.cfg.engine.GetRateLimiter(limiterID)
	if rateLimiter == nil {
		log.Sample(noRateLimiterSampler).Warn().
			Str(metrics.PolicyNameLabel, limiterID.PolicyName).
			Str(metrics.PolicyHashLabel, limiterID.PolicyHash).
			Str(metrics.ComponentIDLabel, limiterID.ComponentID).
			Msg("RateLimiter not found")
		return
	}
	// Add decision type label to the request counter metric
	labels[metrics.DecisionTypeLabel] = decisionType.String()
	labels[metrics.LimiterDroppedLabel] = strconv.FormatBool(dropped)
	requestCounter := rateLimiter.GetRequestCounter(labels)
	if requestCounter != nil {
		requestCounter.Inc()
	}
}

func (p *metricsProcessor) updateMetricsForFlowRegulator(limiterID iface.LimiterID, labels map[string]string, dropped bool, decisionType flowcontrolv1.CheckResponse_DecisionType) {
	flowRegulator := p.cfg.engine.GetFlowRegulator(limiterID)
	if flowRegulator == nil {
		log.Sample(noFlowRegulatorSampler).Warn().
			Str(metrics.PolicyNameLabel, limiterID.PolicyName).
			Str(metrics.PolicyHashLabel, limiterID.PolicyHash).
			Str(metrics.ComponentIDLabel, limiterID.ComponentID).
			Msg("FlowRegulator not found")
		return
	}
	// Add decision type label to the request counter metric
	labels[metrics.DecisionTypeLabel] = decisionType.String()
	labels[metrics.RegulatorDroppedLabel] = strconv.FormatBool(dropped)
	requestCounter := flowRegulator.GetRequestCounter(labels)
	if requestCounter != nil {
		requestCounter.Inc()
	}
}

func (p *metricsProcessor) updateMetricsForClassifier(classifierID iface.ClassifierID) {
	classifier := p.cfg.classificationEngine.GetClassifier(classifierID)
	if classifier == nil {
		log.Sample(noClassifierSampler).Warn().
			Str(metrics.PolicyNameLabel, classifierID.PolicyName).
			Str(metrics.PolicyHashLabel, classifierID.PolicyHash).
			Int64(metrics.ClassifierIndexLabel, classifierID.ClassifierIndex).
			Msg("Classifier not found")
		return
	}

	requestCounter := classifier.GetRequestCounter()
	if requestCounter != nil {
		requestCounter.Inc()
	}
}

func (p *metricsProcessor) updateMetricsForFluxMeters(
	fluxMeterMessage *flowcontrolv1.FluxMeterInfo,
	decisionType flowcontrolv1.CheckResponse_DecisionType,
	statusCode string,
	flowStatus string,
	attributes pcommon.Map,
	treatAsMissing []string,
) {
	fluxMeter := p.cfg.engine.GetFluxMeter(fluxMeterMessage.FluxMeterName)
	if fluxMeter == nil {
		log.Sample(noFluxMeterSampler).Warn().Str(metrics.FluxMeterNameLabel, fluxMeterMessage.GetFluxMeterName()).
			Str(metrics.DecisionTypeLabel, decisionType.String()).
			Str(metrics.StatusCodeLabel, statusCode).
			Str(metrics.FlowStatusLabel, flowStatus).
			Msg("FluxMeter not found")
		return
	}

	labels := internal.StatusLabelsForMetrics(decisionType, statusCode, flowStatus)

	// metricValue is the value at fluxMeter's AttributeKey
	metricValue, found := otelcollector.GetFloat64(attributes, fluxMeter.GetAttributeKey(), treatAsMissing)

	// Add attribute found label to the flux meter metric
	if found {
		labels[metrics.ValidLabel] = metrics.ValidTrue
	} else {
		labels[metrics.ValidLabel] = metrics.ValidFalse
	}
	fluxMeterHistogram := fluxMeter.GetHistogram(labels)
	if fluxMeterHistogram != nil {
		fluxMeterHistogram.Observe(metricValue)
	}
}

func (p *metricsProcessor) populateControlPointCache(checkResponse *flowcontrolv1.CheckResponse, controlPointType string) {
	for _, service := range checkResponse.GetServices() {
		p.cfg.controlPointCache.Put(selectors.NewTypedControlPointID(service, checkResponse.GetControlPoint(), controlPointType))
	}
}

var (
	noConcurrencyLimiterSampler = log.NewRatelimitingSampler()
	noRateLimiterSampler        = log.NewRatelimitingSampler()
	noFlowRegulatorSampler      = log.NewRatelimitingSampler()
	noClassifierSampler         = log.NewRatelimitingSampler()
	noFluxMeterSampler          = log.NewRatelimitingSampler()
)
