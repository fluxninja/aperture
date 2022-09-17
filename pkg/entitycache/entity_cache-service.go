package entitycache

import (
	"context"
	"errors"

	entitycachev1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/entitycache/v1"
	"google.golang.org/grpc"
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

// GetEntityByIPAddress returns an entity by IP address.
func (c *EntityCacheService) GetEntityByIPAddress(ctx context.Context, req *entitycachev1.GetEntityByIPAddressRequest) (*entitycachev1.Entity, error) {
	entity := c.entityCache.GetByIP(req.GetIpAddress())
	if entity == nil {
		return nil, errors.New("entity not found")
	}
	return entity, nil
}

// GetEntityByName returns an entity by name.
func (c *EntityCacheService) GetEntityByName(ctx context.Context, req *entitycachev1.GetEntityByNameRequest) (*entitycachev1.Entity, error) {
	entity := c.entityCache.GetByName(req.GetName())
	if entity == nil {
		return nil, errors.New("entity not found")
	}
	return entity, nil
}
