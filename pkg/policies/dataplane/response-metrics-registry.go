package dataplane

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/metrics"
	"github.com/fluxninja/aperture/pkg/policies/dataplane/iface"
)

type responseMetricsRegistry struct {
	pr                        *prometheus.Registry
	workloadLatencySummaryVec *prometheus.SummaryVec
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

	lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
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
					return err
				}
			}

			return nil
		},
		OnStop: func(_ context.Context) error {
			unregistered := registry.Unregister(rm.workloadLatencySummaryVec)
			if !unregistered {
				log.Error().Msgf("Failed to unregister workload metric from Prometheus registry")
			}
			return nil
		},
	})

	return rm
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

// DeleteTokenLatencyHistogram - deletes histogram from token latency vector with specified labels.
func (rm *responseMetricsRegistry) DeleteTokenLatencyHistogram(labels map[string]string) bool {
	deleted := rm.workloadLatencySummaryVec.DeletePartialMatch(labels)
	log.Info().Msgf("Deleted %d metrics for token latency", deleted)

	return deleted != 0
}
