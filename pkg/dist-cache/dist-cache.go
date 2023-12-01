package distcache

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/buraksezer/olric"
	olricconfig "github.com/buraksezer/olric/config"
	"github.com/clarketm/json"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"

	distcachev1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/distcache/v1"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/metrics"
)

// DistCache is a peer to peer distributed cache.
type DistCache struct {
	distcachev1.UnimplementedDistCacheServiceServer
	lock              sync.Mutex
	config            *olricconfig.Config
	olric             *olric.Olric
	client            olric.Client
	metrics           *DistCacheMetrics
	shutDowner        fx.Shutdowner
	statsFailureCount uint8
}

// NewDistCache creates a new instance of DistCache.
func NewDistCache(config *olricconfig.Config, olric *olric.Olric, metrics *DistCacheMetrics, shutDowner fx.Shutdowner) *DistCache {
	return &DistCache{
		config:     config,
		olric:      olric,
		client:     olric.NewEmbeddedClient(),
		metrics:    metrics,
		shutDowner: shutDowner,
	}
}

// NewDMap creates a new DMap.
func (dc *DistCache) NewDMap(name string, config olricconfig.DMap) (olric.DMap, error) {
	dc.lock.Lock()
	defer dc.lock.Unlock()
	dc.config.DMaps.Custom[name] = config
	d, err := dc.client.NewDMap(name)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to create new DMap: %s, shutting down", name)
		// shutdown
		_ = dc.shutDowner.Shutdown()
		return nil, err
	}
	return d, nil
}

// DeleteDMap deletes a DMap.
func (dc *DistCache) DeleteDMap(name string) error {
	dc.lock.Lock()
	defer dc.lock.Unlock()
	defer delete(dc.config.DMaps.Custom, name)
	return dc.client.DeleteDMap(name)
}

func (dc *DistCache) scrapeMetrics(ctx context.Context) (proto.Message, error) {
	stats, err := dc.client.Stats(ctx, "")
	if err != nil {
		dc.statsFailureCount++
		if dc.statsFailureCount > 10 {
			log.Error().Err(err).Msgf("Failed to scrape Olric statistics 10 times in a row, shutting down")
			_ = dc.shutDowner.Shutdown()
		}
		log.Error().Err(err).Msgf("Failed to scrape Olric statistics")
		return nil, err
	}

	dc.statsFailureCount = 0

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
func (dc *DistCache) GetStats(ctx context.Context, _ *emptypb.Empty) (*structpb.Struct, error) {
	// create a new context with a timeout to avoid hanging
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	stats, err := dc.client.Stats(ctx, "")
	if err != nil {
		log.Error().Err(err).Msgf("Failed to scrape Olric statistics")
		return nil, err
	}

	rawStats, err := json.Marshal(stats)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to marshal Olric statistics")
		return nil, err
	}

	structpbStats := &structpb.Struct{}
	err = json.Unmarshal(rawStats, structpbStats)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to unmarshal Olric statistics")
		return nil, err
	}

	// remove empty partitions from the stats
	for k, v := range structpbStats.GetFields()["partitions"].GetStructValue().GetFields() {
		if v.GetStructValue().GetFields()["length"].GetNumberValue() == 0 {
			delete(structpbStats.GetFields()["partitions"].GetStructValue().GetFields(), k)
		}
	}

	return structpbStats, nil
}
