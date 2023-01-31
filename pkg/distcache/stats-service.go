package distcache

import (
	"context"

	"github.com/buraksezer/olric"
	"github.com/clarketm/json"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"

	distcachev1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/distcache/v1"
	"github.com/fluxninja/aperture/pkg/log"
)

// DistCacheService implements distcache.v1 service.
type DistCacheService struct {
	distcachev1.UnimplementedDistCacheServiceServer
	Address string
	Olric   *olric.Olric
}

// RegisterDistCacheService returns a new Handler.
func RegisterDistCacheService(distcache *DistCache) *DistCacheService {
	return &DistCacheService{
		Address: distcache.Address,
		Olric:   distcache.Olric,
	}
}

// GetStats returns stats of the current Olric member.
func (svc *DistCacheService) GetStats(ctx context.Context, _ *emptypb.Empty) (*structpb.Struct, error) {
	client := svc.Olric.NewEmbeddedClient()
	stats, err := client.Stats(ctx, svc.Address)
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
	return structpbStats, nil

	// newStats := &distcachev1.Stats{}
	// err = json.Unmarshal(rawStats, newStats)
	// if err != nil {
	// 	log.Error().Err(err).Msgf("Failed to unmarshal Olric statistics")
	// 	return nil, err
	// }

	// // Removing empty partitions to reduce the response size
	// partitions := make(map[uint64]*distcachev1.Partition)
	// for key, partition := range newStats.Partitions {
	// 	if partition.Length != 0 {
	// 		partitions[key] = partition
	// 	}
	// }

	// newStats.Partitions = partitions
	// return newStats, nil
}
