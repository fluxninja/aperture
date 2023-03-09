package entities

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	entitiesv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/discovery/entities/v1"
)

// EntitiesService is the implementation of entitiesv1.EntitiesService interface.
type EntitiesService struct {
	entitiesv1.UnimplementedEntitiesServiceServer
	entityCache *Entities
}

// RegisterEntitiesService registers a service for entity cache.
func RegisterEntitiesService(server *grpc.Server, cache *Entities) {
	svc := &EntitiesService{
		entityCache: cache,
	}
	entitiesv1.RegisterEntitiesServiceServer(server, svc)
}

// GetEntities returns *entitiesv1.Entities which contains mappings of ip address to entity and entity name to entity.
func (c *EntitiesService) GetEntities(ctx context.Context, _ *emptypb.Empty) (*entitiesv1.Entities, error) {
	ec := c.entityCache.GetEntities()
	return &entitiesv1.Entities{
		EntitiesByIpAddress: ec.EntitiesByIpAddress,
		EntitiesByName:      ec.EntitiesByName,
	}, nil
}

// GetEntityByIPAddress returns an entity by IP address.
func (c *EntitiesService) GetEntityByIPAddress(ctx context.Context, req *entitiesv1.GetEntityByIPAddressRequest) (*entitiesv1.Entity, error) {
	return c.entityCache.GetByIP(req.GetIpAddress())
}

// GetEntityByName returns an entity by name.
func (c *EntitiesService) GetEntityByName(ctx context.Context, req *entitiesv1.GetEntityByNameRequest) (*entitiesv1.Entity, error) {
	return c.entityCache.GetByName(req.GetName())
}
