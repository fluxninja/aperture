package distcache

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	objectstorage "github.com/fluxninja/aperture/v2/pkg/objectstorage"

	"github.com/buraksezer/olric"
	olricconfig "github.com/buraksezer/olric/config"
	"github.com/buraksezer/olric/events"
	olricstats "github.com/buraksezer/olric/stats"
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
	lock               sync.Mutex
	config             *olricconfig.Config
	olric              *olric.Olric
	objStorage         objectstorage.ObjectStorageIface
	client             olric.Client
	metrics            *DistCacheMetrics
	shutDowner         fx.Shutdowner
	statsFailureCount  uint8
	prometheusRegistry *prometheus.Registry

	// These are periodically update by the scapeMetrics function.
	memberID   string
	memberName string
}

// NewDistCache creates a new instance of DistCache.
func NewDistCache(
	config *olricconfig.Config,
	olric *olric.Olric,
	objStorage objectstorage.ObjectStorageIface,
	metrics *DistCacheMetrics,
	shutDowner fx.Shutdowner,
	prometheusRegistry *prometheus.Registry,
) *DistCache {
	return &DistCache{
		config:             config,
		olric:              olric,
		objStorage:         objStorage,
		client:             olric.NewEmbeddedClient(),
		metrics:            metrics,
		shutDowner:         shutDowner,
		prometheusRegistry: prometheusRegistry,
	}
}

// NewDMap creates a new DMap.
func (dc *DistCache) NewDMap(name string, config olricconfig.DMap, persistent bool) (olric.DMap, error) {
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

	if persistent {
		if dc.objStorage == nil {
			log.Trace().Msg("Object storage not enabled in config, returning non-persistent dmap")
			return d, nil
		}
		return objectstorage.NewPersistentDMap(
			d,
			config.TTLDuration,
			dc.objStorage,
			dc.prometheusRegistry,
			dc.getMetricsLabels,
		)
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

	dc.memberID = strconv.FormatUint(stats.Member.ID, 10)
	dc.memberName = stats.Member.Name
	// We don't care about the second argument as we just set all needed values above.
	metricLabels, _ := dc.getMetricsLabels()

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

	partitionsCountGauge, err := dc.metrics.PartitionsCount.GetMetricWith(metricLabels)
	if err != nil {
		log.Debug().Msgf("Could not extract partitions count gauge metric from olric instance: %v", err)
	} else {
		partitionsCountGauge.Set(float64(countNotEmptyPartitions(stats.Partitions)))
	}

	backupPartitionsCountGauge, err := dc.metrics.BackupPartitionsCount.GetMetricWith(metricLabels)
	if err != nil {
		log.Debug().Msgf("Could not extract backup partitions count gauge metric from olric instance: %v", err)
	} else {
		backupPartitionsCountGauge.Set(float64(countNotEmptyPartitions(stats.Backups)))
	}

	partitionsLength := 0
	for _, v := range stats.Partitions {
		partitionsLength += v.Length
	}
	partitionsLengthGauge, err := dc.metrics.PartitionsLength.GetMetricWith(metricLabels)
	if err != nil {
		log.Debug().Msgf("Could not extract partitions length gauge metric from olric instance: %v", err)
	} else {
		partitionsLengthGauge.Set(float64(partitionsLength))
	}

	backupPartitionsLength := 0
	for _, v := range stats.Backups {
		backupPartitionsLength += v.Length
	}
	backupPartitionsLengthGauge, err := dc.metrics.BackupPartitionsLength.GetMetricWith(metricLabels)
	if err != nil {
		log.Debug().Msgf("Could not extract backup partitions length gauge metric from olric instance: %v", err)
	} else {
		backupPartitionsLengthGauge.Set(float64(backupPartitionsLength))
	}
	return nil, nil
}

// getMetricsLabels returns metric labels based on dc.memberID and dc.memberName.
// If any of those is empty, function will return false as second argument.
func (dc *DistCache) getMetricsLabels() (prometheus.Labels, bool) {
	metricLabels := make(prometheus.Labels)
	metricLabels[metrics.DistCacheMemberIDLabel] = dc.memberID
	metricLabels[metrics.DistCacheMemberNameLabel] = dc.memberName
	return metricLabels, dc.memberID != "" && dc.memberName != ""
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

func (dc *DistCache) startReportingMetricsFromEvents(ctx context.Context) error {
	ps, err := dc.client.NewPubSub()
	if err != nil {
		return fmt.Errorf("failed creating pub sub client: %w", err)
	}
	rps := ps.Subscribe(ctx, events.ClusterEventsChannel)
	msgCh := rps.Channel()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case message, ok := <-msgCh:
				if !ok {
					log.Warn().Msg("Events channel closed")
					return
				}
				event, err := olricEventFromPayload(message.Payload)
				if err != nil {
					log.Debug().Err(err).Msg("Failed getting olric event from payload")
					continue
				}
				metricLabels, ready := dc.getMetricsLabels()
				if !ready {
					log.Warn().Msg("Member ID or Member Name not yet set. Skipping event")
					continue
				}
				switch event.Kind {
				case events.KindFragmentMigrationEvent:
					fragmentMigrationEventsCounter, err := dc.metrics.FragmentMigrationEventsTotal.GetMetricWith(metricLabels)
					if err != nil {
						log.Debug().Msgf("Could not extract fragment migrate counter metric from olric instance: %v", err)
						continue
					}
					fragmentMigrationEventsCounter.Inc()
				case events.KindFragmentReceivedEvent:
					fragmentReceivedEventsCounter, err := dc.metrics.FragmentReceivedEventsTotal.GetMetricWith(metricLabels)
					if err != nil {
						log.Debug().Msgf("Could not extract fragment received counter metric from olric instance: %v", err)
						continue
					}
					fragmentReceivedEventsCounter.Inc()
				}
			}
		}
	}()
	return nil
}

func countNotEmptyPartitions(partitions map[olricstats.PartitionID]olricstats.Partition) int {
	result := 0
	for _, partition := range partitions {
		if partition.Length > 0 {
			result += 1
		}
	}
	return result
}

type olricEvent struct {
	Kind string `json:"Kind"`
}

func olricEventFromPayload(payload string) (*olricEvent, error) {
	var event olricEvent
	err := json.Unmarshal([]byte(payload), &event)
	if err != nil {
		return nil, fmt.Errorf("failed unmarhalling event from payload: %w", err)
	}
	return &event, nil
}
