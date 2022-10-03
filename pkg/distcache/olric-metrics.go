package distcache

import (
	"github.com/fluxninja/aperture/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// OlricMetrics holds metrics from Olric DMap statistics.
type OlricMetrics struct {
	EntriesTotal prometheus.Gauge
	DeleteHits   prometheus.Gauge
	DeleteMisses prometheus.Gauge
	GetMisses    prometheus.Gauge
	GetHits      prometheus.Gauge
	EvictedTotal prometheus.Gauge
}

func newOlricMetrics() *OlricMetrics {
	return &OlricMetrics{
		EntriesTotal: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: metrics.OlricEntriesTotalMetricName,
			Help: "Total number of entries in the DMap.",
		}),
		DeleteHits: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: metrics.OlricDeleteHitsMetricName,
			Help: "Number of entries deleted from the DMap.",
		}),
		DeleteMisses: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: metrics.OlricDeleteMissesMetricName,
			Help: "Number of entries not deleted from the DMap.",
		}),
		GetMisses: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: metrics.OlricGetMissesMetricName,
			Help: "Number of entries not found in the DMap.",
		}),
		GetHits: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: metrics.OlricGetHitsMetricName,
			Help: "Number of entries found in the DMap.",
		}),
		EvictedTotal: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: metrics.OlricEvictedTotalMetricName,
			Help: "Total number of evicted entries in the DMap.",
		}),
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
