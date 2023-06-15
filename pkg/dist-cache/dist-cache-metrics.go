package distcache

import (
	"fmt"

	"github.com/fluxninja/aperture/v2/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// DistCacheMetrics holds metrics from DistCache, Olric, DMap statistics.
type DistCacheMetrics struct {
	EntriesTotal *prometheus.GaugeVec
	DeleteHits   *prometheus.GaugeVec
	DeleteMisses *prometheus.GaugeVec
	GetMisses    *prometheus.GaugeVec
	GetHits      *prometheus.GaugeVec
	EvictedTotal *prometheus.GaugeVec
}

func newDistCacheMetrics() *DistCacheMetrics {
	distCacheMetricsLabels := []string{metrics.DistCacheMemberIDLabel, metrics.DistCacheMemberNameLabel}
	return &DistCacheMetrics{
		EntriesTotal: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metrics.DistCacheEntriesTotalMetricName,
			Help: "Total number of entries in the DMap.",
		}, distCacheMetricsLabels),
		DeleteHits: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metrics.DistCacheDeleteHitsMetricName,
			Help: "Number of deletion requests resulting in an item being removed in the DMap.",
		}, distCacheMetricsLabels),
		DeleteMisses: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metrics.DistCacheDeleteMissesMetricName,
			Help: "Number of deletion requests for missing keys in the DMap.",
		}, distCacheMetricsLabels),
		GetMisses: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metrics.DistCacheGetMissesMetricName,
			Help: "Number of entries that have been requested and not found in the DMap.",
		}, distCacheMetricsLabels),
		GetHits: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metrics.DistCacheGetHitsMetricName,
			Help: "Number of entries that have been requested and found present in the DMap.",
		}, distCacheMetricsLabels),
		EvictedTotal: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metrics.DistCacheEvictedTotalMetricName,
			Help: "Total number of entries removed from cache to free memory for new entries in the DMap.",
		}, distCacheMetricsLabels),
	}
}

func (dm *DistCacheMetrics) allMetrics() []prometheus.Collector {
	return []prometheus.Collector{
		dm.EntriesTotal,
		dm.DeleteHits,
		dm.DeleteMisses,
		dm.GetMisses,
		dm.GetHits,
		dm.EvictedTotal,
	}
}

func (dm *DistCacheMetrics) registerMetrics(prometheusRegistry *prometheus.Registry) error {
	for _, m := range dm.allMetrics() {
		err := prometheusRegistry.Register(m)
		if err != nil {
			if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
				return fmt.Errorf("unable to register distcache metrics: %v", err)
			}
		}
	}
	return nil
}

func (dm *DistCacheMetrics) unregisterMetrics(prometheusRegistry *prometheus.Registry) error {
	for _, m := range dm.allMetrics() {
		if !prometheusRegistry.Unregister(m) {
			return fmt.Errorf("unable to unregister distcache metrics")
		}
	}
	return nil
}
