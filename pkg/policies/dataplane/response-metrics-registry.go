package dataplane

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/multierr"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/metrics"
	"github.com/fluxninja/aperture/pkg/policies/dataplane/iface"
)

// TODO: Bring metrics to this file
// 1. per component metrics from actuator, fluxmeter.

type responseMetricsRegistry struct {
	pr                        *prometheus.Registry
	workloadLatencySummaryVec *prometheus.SummaryVec
	fluxmeterHistogramVec     *prometheus.HistogramVec
}

// ProvideResponseMetricsAPI returns API for getting metrics observers.
func ProvideResponseMetricsAPI(registry *prometheus.Registry, lifecycle fx.Lifecycle) iface.ResponseMetricsAPI {
	rm := &responseMetricsRegistry{
		pr: registry,
	}

	// WORKLOAD LATENCY METRICS
	rm.workloadLatencySummaryVec = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: metrics.WorkloadLatencyMetricName,
		Help: "Latency summary of workload",
	}, []string{
		metrics.PolicyNameLabel,
		metrics.PolicyHashLabel,
		metrics.ComponentIndexLabel,
		metrics.DecisionTypeLabel,
		metrics.WorkloadIndexLabel,
	})

	// FLUXMETER METRICS
	rm.fluxmeterHistogramVec = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: metrics.FluxMeterMetricName,
		Help: "Latency histogram of fluxmeters",
	}, []string{
		metrics.FluxMeterNameLabel,
		metrics.DecisionTypeLabel,
		metrics.StatusCodeLabel,
		metrics.FeatureStatusLabel,
	})

	lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			var errs error
			err := registry.Register(rm.workloadLatencySummaryVec)
			if err != nil {
				if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
					// We're registering this histogram vec from multiple processors
					// (logs processor and traces processor), so if both processors are
					// enabled, it's expected that whichever processor is created
					// second, it will see that the histogram vec was already
					// registered. Use the existing histogram vec from now on.
					rm.workloadLatencySummaryVec = are.ExistingCollector.(*prometheus.SummaryVec)
				} else {
					errs = multierr.Append(errs, err)
				}
			}

			err = registry.Register(rm.fluxmeterHistogramVec)
			if err != nil {
				if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
					rm.fluxmeterHistogramVec = are.ExistingCollector.(*prometheus.HistogramVec)
				} else {
					errs = multierr.Append(errs, err)
				}
			}

			return errs
		},
		OnStop: func(_ context.Context) error {
			unregistered := registry.Unregister(rm.workloadLatencySummaryVec)
			if !unregistered {
				log.Error().Msgf("Failed to unregister workload metric from Prometheus registry")
			}
			unregistered = registry.Unregister(rm.fluxmeterHistogramVec)
			if !unregistered {
				log.Error().Msgf("Failed to unregister fluxmeters metric from Prometheus registry")
			}
			return nil
		},
	})

	return rm
}

// GetFluxMeterHistogram - get prometheus Observer for fluxmeter.
func (rm *responseMetricsRegistry) GetFluxMeterHistogram(
	fluxmeterID, statusCode, featureStatus string,
	decisionType flowcontrolv1.CheckResponse_DecisionType,
) (prometheus.Observer, error) {
	labels := make(map[string]string)
	labels[metrics.DecisionTypeLabel] = decisionType.String()
	labels[metrics.StatusCodeLabel] = statusCode
	labels[metrics.FluxMeterNameLabel] = fluxmeterID
	labels[metrics.FeatureStatusLabel] = featureStatus

	fluxmeterHistogram, err := rm.fluxmeterHistogramVec.GetMetricWith(labels)
	if err != nil {
		log.Warn().Err(err).Msg("Getting fluxmeter histogram")
		return nil, err
	}

	return fluxmeterHistogram, nil
}

// GetTokenLatencyHistogram - get prometheus Observer for tokens.
func (rm *responseMetricsRegistry) GetTokenLatencyHistogram(labels map[string]string) (prometheus.Observer, error) {
	tokenHistogram, err := rm.workloadLatencySummaryVec.GetMetricWith(labels)
	if err != nil {
		log.Warn().Err(err).Msg("Getting latency histogram")
		return nil, err
	}

	return tokenHistogram, nil
}

// DeleteFluxmeterHistogram - deletes histogram from fluxmeter vector with specified labels.
func (rm *responseMetricsRegistry) DeleteFluxmeterHistogram(fluxmeterID string) bool {
	labels := make(map[string]string)
	labels[metrics.FluxMeterNameLabel] = fluxmeterID

	deleted := rm.fluxmeterHistogramVec.DeletePartialMatch(labels)
	log.Info().Msgf("Deleted %d metrics for fluxmeter: %+v", deleted, fluxmeterID)

	return deleted != 0
}

// DeleteTokenLatencyHistogram - deletes histogram from token latency vector with specified labels.
func (rm *responseMetricsRegistry) DeleteTokenLatencyHistogram(labels map[string]string) bool {
	deleted := rm.workloadLatencySummaryVec.DeletePartialMatch(labels)
	log.Info().Msgf("Deleted %d metrics for token latency", deleted)

	return deleted != 0
}
