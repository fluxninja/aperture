package distcache

import (
	"context"
	"strconv"
	"sync"

	"github.com/buraksezer/olric"
	olricconfig "github.com/buraksezer/olric/config"
	"github.com/clarketm/json"
	distcachev1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/distcache/v1"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

// DistCache is a peer to peer distributed cache.
type DistCache struct {
	distcachev1.UnimplementedDistCacheServiceServer
	lock    sync.Mutex
	config  *olricconfig.Config
	olric   *olric.Olric
	metrics *DistCacheMetrics
}

// NewDistCache creates a new instance of DistCache.
func NewDistCache(config *olricconfig.Config,
	olric *olric.Olric,
	metrics *DistCacheMetrics,
) *DistCache {
	return &DistCache{
		config:  config,
		olric:   olric,
		metrics: metrics,
	}
}

// NewDMap creates a new Distributed Map.
func (dc *DistCache) NewDMap(name string, config olricconfig.DMap) (*olric.DMap, error) {
	dc.lock.Lock()
	defer dc.lock.Unlock()
	dc.config.DMaps.Custom[name] = config
	defer delete(dc.config.DMaps.Custom, name)
	return dc.olric.NewDMap(name)
}

func (dc *DistCache) scrapeMetrics(context.Context) (proto.Message, error) {
	stats, err := dc.olric.Stats()
	if err != nil {
		log.Error().Err(err).Msgf("Failed to scrape Olric statistics")
		return nil, err
	}

	memberID := stats.Member.ID
	memberName := stats.Member.Name
	metricLabels := make(prometheus.Labels)
	metricLabels[metrics.DistCacheMemberIDLabel] = strconv.FormatUint(memberID, 10)
	metricLabels[metrics.DistCacheMemberNameLabel] = memberName

	entriesTotalGauge, err := dc.metrics.EntriesTotal.GetMetricWith(metricLabels)
	if err != nil {
		log.Debug().Msgf("Could not extract entries total gauge metric from olric instance: %v", err)
	} else {
		entriesTotalGauge.Set(float64(stats.DMaps.EntriesTotal))
	}

	deleteHitsGauge, err := dc.metrics.DeleteHits.GetMetricWith(metricLabels)
	if err != nil {
		log.Debug().Msgf("Could not extract delete hits gauge metric from olric instance: %v", err)
	} else {
		deleteHitsGauge.Set(float64(stats.DMaps.DeleteHits))
	}

	deleteMissesGauge, err := dc.metrics.DeleteMisses.GetMetricWith(metricLabels)
	if err != nil {
		log.Debug().Msgf("Could not extract delete misses gauge metric from olric instance: %v", err)
	} else {
		deleteMissesGauge.Set(float64(stats.DMaps.DeleteMisses))
	}

	getMissesGauge, err := dc.metrics.GetMisses.GetMetricWith(metricLabels)
	if err != nil {
		log.Debug().Msgf("Could not extract get misses gauge metric from olric instance: %v", err)
	} else {
		getMissesGauge.Set(float64(stats.DMaps.GetMisses))
	}

	getHitsGauge, err := dc.metrics.GetHits.GetMetricWith(metricLabels)
	if err != nil {
		log.Debug().Msgf("Could not extract get hits gauge metric from olric instance: %v", err)
	} else {
		getHitsGauge.Set(float64(stats.DMaps.GetHits))
	}

	evictedTotalGauge, err := dc.metrics.EvictedTotal.GetMetricWith(metricLabels)
	if err != nil {
		log.Debug().Msgf("Could not extract evicted total gauge metric from olric instance: %v", err)
	} else {
		evictedTotalGauge.Set(float64(stats.DMaps.EvictedTotal))
	}
	return nil, nil
}

// GetStats returns stats of the current Olric member.
func (dc *DistCache) GetStats(ctx context.Context, _ *emptypb.Empty) (*distcachev1.Stats, error) {
	stats, err := dc.olric.Stats()
	if err != nil {
		log.Error().Err(err).Msgf("Failed to scrape Olric statistics")
		return nil, err
	}

	rawStats, err := json.Marshal(stats)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to marshal Olric statistics")
		return nil, err
	}

	newStats := &distcachev1.Stats{}
	err = json.Unmarshal(rawStats, newStats)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to unmarshal Olric statistics")
		return nil, err
	}

	// Removing empty partitions to reduce the response size
	partitions := make(map[uint64]*distcachev1.Partition)
	for key, partition := range newStats.Partitions {
		if partition.Length != 0 {
			partitions[key] = partition
		}
	}

	newStats.Partitions = partitions
	return newStats, nil
}
