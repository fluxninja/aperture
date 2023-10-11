package entities

import (
	"context"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	cmdv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/cmd/v1"
	entitiesv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/discovery/entities/v1"
	"github.com/fluxninja/aperture/v2/pkg/etcd/transport"
	"github.com/fluxninja/aperture/v2/pkg/rpc"
)

// EntitiesService is the implementation of entitiesv1.EntitiesService interface.
type EntitiesService struct {
	entitiesv1.UnimplementedEntitiesServiceServer
	entityCache *Entities
}

// RegisterEntitiesServiceIn bundles and annotates parameters.
type RegisterEntitiesServiceIn struct {
	fx.In
	Server              *grpc.Server `name:"default"`
	Cache               *Entities
	Registry            *rpc.HandlerRegistry
	EtcdTransportClient *transport.EtcdTransportClient
}

// RegisterEntitiesService registers a service for entity cache.
func RegisterEntitiesService(in RegisterEntitiesServiceIn) error {
	svc := &EntitiesService{
		entityCache: in.Cache,
	}
	entitiesv1.RegisterEntitiesServiceServer(in.Server, svc)
	err := transport.RegisterFunction(in.EtcdTransportClient, svc.ListDiscoveryEntities)
	if err != nil {
		return err
	}
	err = transport.RegisterFunction(in.EtcdTransportClient, svc.ListDiscoveryEntity)
	if err != nil {
		return err
	}
	return nil
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
	e, err := c.entityCache.GetByIP(req.GetIpAddress())
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return e.Clone(), nil
}

// GetEntityByName returns an entity by name.
func (c *EntitiesService) GetEntityByName(ctx context.Context, req *entitiesv1.GetEntityByNameRequest) (*entitiesv1.Entity, error) {
	e, err := c.entityCache.GetByName(req.GetName())
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return e.Clone(), nil
}

// ListDiscoveryEntities lists currently discovered entities by IP address.
func (c *EntitiesService) ListDiscoveryEntities(ctx context.Context, _ *cmdv1.ListDiscoveryEntitiesRequest) (*cmdv1.ListDiscoveryEntitiesAgentResponse, error) {
	ec := c.entityCache.GetEntities()
	return &cmdv1.ListDiscoveryEntitiesAgentResponse{
		Entities: ec.EntitiesByIpAddress.Entities,
	}, nil
}

// ListDiscoveryEntity returns an entity by IP address or name.
func (c *EntitiesService) ListDiscoveryEntity(ctx context.Context, req *cmdv1.ListDiscoveryEntityRequest) (*cmdv1.ListDiscoveryEntityAgentResponse, error) {
	// check if request is for an IP address or name
	switch req.By.(type) {
	case *cmdv1.ListDiscoveryEntityRequest_IpAddress:
		entity, err := c.entityCache.GetByIP(req.GetIpAddress())
		if err != nil {
			return nil, err
		}
		return &cmdv1.ListDiscoveryEntityAgentResponse{
			Entity: entity.Clone(),
		}, nil
	case *cmdv1.ListDiscoveryEntityRequest_Name:
		entity, err := c.entityCache.GetByName(req.GetName())
		if err != nil {
			return nil, err
		}
		return &cmdv1.ListDiscoveryEntityAgentResponse{
			Entity: entity.Clone(),
		}, nil
	default:
		return nil, nil
	}
}
