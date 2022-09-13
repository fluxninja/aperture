package entitycache

import (
	"context"

	entitycachev1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/entitycache/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// EntityCacheService is the implementation of entitycachev1.EntityCacheService interface.
type EntityCacheService struct {
	entitycachev1.UnimplementedEntityCacheServiceServer
	entityCache *EntityCache
}

// RegisterEntityCacheService registers a service for entity cache.
func RegisterEntityCacheService(server *grpc.Server, cache *EntityCache) {
	svc := &EntityCacheService{
		entityCache: cache,
	}
	entitycachev1.RegisterEntityCacheServiceServer(server, svc)
}

// GetServicesList returns a list of services based on entities in cache.
func (c *EntityCacheService) GetServicesList(ctx context.Context, _ *emptypb.Empty) (*entitycachev1.ServicesList, error) {
	return c.entityCache.Services(), nil
}

// GetEntityCache returns *entitycachev1.EntityCache which contains mappings of ip address to entity and entity name to entity.
func (c *EntityCacheService) GetEntityCache(ctx context.Context, _ *emptypb.Empty) (*entitycachev1.EntityCache, error) {
	ec := c.entityCache.Entities()
	return &entitycachev1.EntityCache{
		EntitiesByIpAddress:  ec.EntitiesByIpAddress,
		EntitiesByEntityName: ec.EntitiesByEntityName,
	}, nil
}

// GetEntity returns matching entity in cache based on request field type.
func (c *EntityCacheService) GetEntity(ctx context.Context, req *entitycachev1.GetEntityRequest) (*entitycachev1.Entity, error) {
	switch by := req.By.(type) {
	case *entitycachev1.GetEntityRequest_IpAddress:
		return c.entityCache.GetByIP(req.GetIpAddress()), nil
	case *entitycachev1.GetEntityRequest_EntityName:
		return c.entityCache.GetByName(req.GetEntityName()), nil
	default:
		return nil, status.Errorf(codes.InvalidArgument, "unsupported by: %v", by)
	}
}
