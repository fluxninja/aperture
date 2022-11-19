package metricsprocessor

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/check/v1"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/metrics"
	"github.com/fluxninja/aperture/pkg/otelcollector"
	"github.com/fluxninja/aperture/pkg/otelcollector/metricsprocessor/internal"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/iface"
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
	err := otelcollector.IterateLogRecords(ld, func(logRecord plog.LogRecord) error {
		retErr := func(sampler zerolog.Sampler, errMsg string) error {
			log.Sample(sampler).Warn().Msg(errMsg)
			return fmt.Errorf(errMsg)
		}
		// Attributes
		attributes := logRecord.Attributes()

		// CheckResponse
		checkResponse := &flowcontrolv1.CheckResponse{}

		// Source specific processing
		source, exists := attributes.Get(otelcollector.ApertureSourceLabel)
		if !exists {
			return retErr(noSourceLabelSampler, "aperture source label not found")
		}
		sourceStr := source.Str()
		if sourceStr == otelcollector.ApertureSourceSDK {
			success := otelcollector.GetStruct(attributes, otelcollector.ApertureCheckResponseLabel, checkResponse, []string{})
			if !success {
				return retErr(
					noSDKCheckResponseSampler,
					"aperture check response label not found in SDK access logs",
				)
			}

			internal.AddSDKSpecificLabels(attributes)
		} else if sourceStr == otelcollector.ApertureSourceEnvoy {
			success := otelcollector.GetStruct(attributes, otelcollector.ApertureCheckResponseLabel, checkResponse, []string{otelcollector.EnvoyMissingAttributeValue})
			if !success {
				return retErr(
					noEnvoyCheckResponseSampler,
					"aperture check response label not found in Envoy access logs",
				)
			}

			internal.AddEnvoySpecificLabels(attributes)
		} else {
			return retErr(
				unrecognizedSourceLabelSampler,
				"aperture source label not recognized",
			)
		}

		statusCode, flowStatus := internal.StatusesFromAttributes(attributes)
		attributes.PutStr(otelcollector.ApertureFlowStatusLabel, internal.FlowStatusForTelemetry(statusCode, flowStatus))
		internal.AddCheckResponseBasedLabels(attributes, checkResponse, sourceStr)

		// Update metrics and enforce include list to eliminate any excess attributes
		if sourceStr == otelcollector.ApertureSourceSDK {
			p.updateMetrics(attributes, checkResponse, []string{})
			internal.EnforceIncludeListSDK(attributes)
		} else if sourceStr == otelcollector.ApertureSourceEnvoy {
			p.updateMetrics(attributes, checkResponse, []string{otelcollector.EnvoyMissingAttributeValue})
			internal.EnforceIncludeListHTTP(attributes)
		}

		// This needs to be called **after** internal.EnforceIncludeList{HTTP,SDK}.
		internal.AddFlowLabels(attributes, checkResponse)
		return nil
	})
	return ld, err
}

var (
	noSourceLabelSampler           = log.NewRatelimitingSampler()
	noSDKCheckResponseSampler      = log.NewRatelimitingSampler()
	noEnvoyCheckResponseSampler    = log.NewRatelimitingSampler()
	unrecognizedSourceLabelSampler = log.NewRatelimitingSampler()
)

func (p *metricsProcessor) updateMetrics(
	attributes pcommon.Map,
	checkResponse *flowcontrolv1.CheckResponse,
	treatAsMissing []string,
) {
	if checkResponse == nil {
		return
	}
	if len(checkResponse.LimiterDecisions) > 0 {
		// Update workload metrics
		latency, latencyFound := otelcollector.GetFloat64(attributes, otelcollector.WorkloadDurationLabel, []string{})
		for _, decision := range checkResponse.LimiterDecisions {
			limiterID := iface.LimiterID{
				PolicyName:     decision.PolicyName,
				PolicyHash:     decision.PolicyHash,
				ComponentIndex: decision.ComponentIndex,
			}

			if cl := decision.GetConcurrencyLimiterInfo(); cl != nil {
				labels := map[string]string{
					metrics.PolicyNameLabel:     decision.PolicyName,
					metrics.PolicyHashLabel:     decision.PolicyHash,
					metrics.ComponentIndexLabel: fmt.Sprintf("%d", decision.ComponentIndex),
					metrics.WorkloadIndexLabel:  cl.GetWorkloadIndex(),
				}

				p.updateMetricsForWorkload(limiterID, labels, checkResponse.DecisionType, latency, latencyFound)
			}

			// Update rate limiter metrics
			if rl := decision.GetRateLimiterInfo(); rl != nil {
				labels := map[string]string{
					metrics.PolicyNameLabel:     decision.PolicyName,
					metrics.PolicyHashLabel:     decision.PolicyHash,
					metrics.ComponentIndexLabel: fmt.Sprintf("%d", decision.ComponentIndex),
				}
				p.updateMetricsForRateLimiter(limiterID, labels, checkResponse.DecisionType)
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

func (p *metricsProcessor) updateMetricsForWorkload(limiterID iface.LimiterID, labels map[string]string, decisionType flowcontrolv1.CheckResponse_DecisionType, latency float64, latencyFound bool) {
	concurrencyLimiter := p.cfg.engine.GetConcurrencyLimiter(limiterID)
	if concurrencyLimiter == nil {
		log.Sample(noConcurrencyLimiterSampler).Warn().
			Str(metrics.PolicyNameLabel, limiterID.PolicyName).
			Str(metrics.PolicyHashLabel, limiterID.PolicyHash).
			Int64(metrics.ComponentIndexLabel, limiterID.ComponentIndex).
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
	requestCounter := concurrencyLimiter.GetRequestCounter(labels)
	if requestCounter != nil {
		requestCounter.Inc()
	}
}

func (p *metricsProcessor) updateMetricsForRateLimiter(limiterID iface.LimiterID, labels map[string]string, decisionType flowcontrolv1.CheckResponse_DecisionType) {
	rateLimiter := p.cfg.engine.GetRateLimiter(limiterID)
	if rateLimiter == nil {
		log.Sample(noRateLimiterSampler).Warn().
			Str(metrics.PolicyNameLabel, limiterID.PolicyName).
			Str(metrics.PolicyHashLabel, limiterID.PolicyHash).
			Int64(metrics.ComponentIndexLabel, limiterID.ComponentIndex).
			Msg("RateLimiter not found")
		return
	}
	// Add decision type label to the request counter metric
	labels[metrics.DecisionTypeLabel] = decisionType.String()
	requestCounter := rateLimiter.GetRequestCounter(labels)
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

var (
	noConcurrencyLimiterSampler = log.NewRatelimitingSampler()
	noRateLimiterSampler        = log.NewRatelimitingSampler()
	noClassifierSampler         = log.NewRatelimitingSampler()
	noFluxMeterSampler          = log.NewRatelimitingSampler()
)
