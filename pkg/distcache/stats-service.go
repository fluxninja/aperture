package distcache

import (
	"context"

	"github.com/buraksezer/olric"
	"github.com/clarketm/json"
	distcachev1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/distcache/v1"
	"github.com/fluxninja/aperture/pkg/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

// DistCacheService implements distcache.v1 service.
type DistCacheService struct {
	distcachev1.UnimplementedDistCacheServiceServer
	Olric *olric.Olric
}

// RegisterDistCacheService returns a new Handler.
func RegisterDistCacheService(distcahce *DistCache) *DistCacheService {
	return &DistCacheService{
		Olric: distcahce.Olric,
	}
}

// GetStats returns stats of the current Olric member.
func (svc *DistCacheService) GetStats(ctx context.Context, _ *emptypb.Empty) (*distcachev1.Stats, error) {
	stats, err := svc.Olric.Stats()
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
