package distcache

import (
	"fmt"

	"github.com/fluxninja/aperture/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/multierr"
)

// OlricMetrics holds metrics from Olric DMap statistics.
type OlricMetrics struct {
	EntriesTotal *prometheus.GaugeVec
	DeleteHits   *prometheus.GaugeVec
	DeleteMisses *prometheus.GaugeVec
	GetMisses    *prometheus.GaugeVec
	GetHits      *prometheus.GaugeVec
	EvictedTotal *prometheus.GaugeVec
}

func newOlricMetrics() *OlricMetrics {
	olricMetricsLabels := []string{metrics.OlricMemberIDLabel, metrics.OlricMemberNameLabel}
	return &OlricMetrics{
		EntriesTotal: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metrics.OlricEntriesTotalMetricName,
			Help: "Total number of entries in the DMap.",
		}, olricMetricsLabels),
		DeleteHits: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metrics.OlricDeleteHitsMetricName,
			Help: "Number of deletion requests resulting in an item being removed in the DMap.",
		}, olricMetricsLabels),
		DeleteMisses: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metrics.OlricDeleteMissesMetricName,
			Help: "Number of deletion requests for missing keys in the DMap.",
		}, olricMetricsLabels),
		GetMisses: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metrics.OlricGetMissesMetricName,
			Help: "Number of entries that have been requested and not found in the DMap.",
		}, olricMetricsLabels),
		GetHits: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metrics.OlricGetHitsMetricName,
			Help: "Number of entries that have been requested and found present in the DMap.",
		}, olricMetricsLabels),
		EvictedTotal: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metrics.OlricEvictedTotalMetricName,
			Help: "Total number of entries removed from cache to free memory for new entries in the DMap.",
		}, olricMetricsLabels),
	}
}

func (om *OlricMetrics) allMetrics() []prometheus.Collector {
	return []prometheus.Collector{
		om.EntriesTotal,
		om.DeleteHits,
		om.DeleteMisses,
		om.GetMisses,
		om.GetHits,
		om.EvictedTotal,
	}
}

func (om *OlricMetrics) registerMetrics(prometheusRegistry *prometheus.Registry) error {
	var multiErr error
	for _, m := range om.allMetrics() {
		err := prometheusRegistry.Register(m)
		if err != nil {
			if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
				multiErr = multierr.Append(multiErr, err)
			}
		}
	}
	return multiErr
}

func (om *OlricMetrics) unregisterMetrics(prometheusRegistry *prometheus.Registry) error {
	var multiErr error
	if !prometheusRegistry.Unregister(om.EntriesTotal) {
		err := fmt.Errorf("failed to unregister %s metric", metrics.OlricEntriesTotalMetricName)
		multiErr = multierr.Append(multiErr, err)
	}
	if !prometheusRegistry.Unregister(om.DeleteHits) {
		err := fmt.Errorf("failed to unregister %s metric", metrics.OlricDeleteHitsMetricName)
		multiErr = multierr.Append(multiErr, err)
	}
	if !prometheusRegistry.Unregister(om.DeleteMisses) {
		err := fmt.Errorf("failed to unregister %s metric", metrics.OlricDeleteMissesMetricName)
		multiErr = multierr.Append(multiErr, err)
	}
	if !prometheusRegistry.Unregister(om.GetMisses) {
		err := fmt.Errorf("failed to unregister %s metric", metrics.OlricGetMissesMetricName)
		multiErr = multierr.Append(multiErr, err)
	}
	if !prometheusRegistry.Unregister(om.GetHits) {
		err := fmt.Errorf("failed to unregister %s metric", metrics.OlricGetHitsMetricName)
		multiErr = multierr.Append(multiErr, err)
	}
	if !prometheusRegistry.Unregister(om.EvictedTotal) {
		err := fmt.Errorf("failed to unregister %s metric", metrics.OlricEvictedTotalMetricName)
		multiErr = multierr.Append(multiErr, err)
	}
	return multiErr
}
