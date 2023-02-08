package kubernetes

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	controlpointcachev1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/controlpointcache/v1"
)

// ControlPointCacheService implements the ControlPointCacheServer interface.
type ControlPointCacheService struct {
	controlpointcachev1.UnimplementedControlPointCacheServer
	ControlPointCache ControlPointCache
}

// RegisterControlPointCacheService returns a new ControlPointCacheService handler.
func RegisterControlPointCacheService(server *grpc.Server, cpc ControlPointCache) {
	svc := &ControlPointCacheService{
		ControlPointCache: cpc,
	}
	controlpointcachev1.RegisterControlPointCacheServer(server, svc)
}

// GetControlPoint returns a ControlPoint from the cache.
func (s ControlPointCacheService) GetControlPoint(ctx context.Context, _ *emptypb.Empty) (*controlpointcachev1.KubernetesControlPoints, error) {
	keys := s.ControlPointCache.Keys()
	controlPoints := make([]*controlpointcachev1.KubernetesControlPoint, len(keys))
	for _, key := range keys {
		cp := key.ToProto()
		controlPoints = append(controlPoints, cp)
	}
	return &controlpointcachev1.KubernetesControlPoints{
		ControlPoints: controlPoints,
	}, nil
}
