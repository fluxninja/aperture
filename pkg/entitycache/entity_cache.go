package entitycache

import (
	"context"
	"sync"

	"go.uber.org/fx"

	entitycachev1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/entitycache/v1"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/discovery/common"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/net/grpcgateway"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/services"
)

// Module sets up EntityCache with Fx.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(provideEntityCache),
		grpcgateway.RegisterHandler{Handler: entitycachev1.RegisterEntityCacheServiceHandlerFromEndpoint}.Annotate(),
		fx.Invoke(RegisterEntityCacheService),
	)
}

// ServiceKey holds key for service.
type ServiceKey struct {
	Name string `json:"name"`
}

func (sk ServiceKey) lessThan(sk2 ServiceKey) bool {
	return sk.Name < sk2.Name
}

// keyFromService returns a service key for given service.
func (c *EntityCache) keyFromService(service *entitycachev1.Service) *ServiceKey {
	return &ServiceKey{
		Name: service.Name,
	}
}

// Merge merges `mergedService` into `originalService`. This sums `EntitiesCount`.
func Merge(originalService, mergedService *entitycachev1.Service) {
	originalService.EntitiesCount += mergedService.EntitiesCount
}

// EntityID is a unique Entity identifier.
type EntityID struct {
	Prefix string `json:"prefix,omitempty"`
	UID    string `json:"uid,omitempty"`
}

// NewEntity creates a new entity from ID and IP address from the tagger.
func NewEntity(id EntityID, ipAddress, name string, services []string) *entitycachev1.Entity {
	return &entitycachev1.Entity{
		EntityId: &entitycachev1.Entity_EntityID{
			Prefix: id.Prefix,
			Uid:    id.UID,
		},
		IpAddress:  ipAddress,
		Services:   services,
		EntityName: name,
	}
}

// EntityCache maps IP addresses and Entity names to entities.
type EntityCache struct {
	sync.RWMutex
	entities *entitycachev1.EntityCache
}

// FxIn are the parameters for ProvideEntityCache.
type FxIn struct {
	fx.In
	Lifecycle      fx.Lifecycle
	EntityTrackers notifiers.Trackers `name:"entity_trackers"`
}

// provideEntityCache creates Entity Cache.
func provideEntityCache(in FxIn) (*EntityCache, error) {
	entityCache := NewEntityCache()

	// create a ConfigPrefixNotifier
	configPrefixNotifier := &notifiers.UnmarshalPrefixNotifier{
		GetUnmarshallerFunc: config.KoanfUnmarshallerConstructor{}.NewKoanfUnmarshaller,
		UnmarshalNotifyFunc: entityCache.processUpdate,
	}

	in.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			err := in.EntityTrackers.AddPrefixNotifier(configPrefixNotifier)
			if err != nil {
				log.Error().Err(err).Msg("failed to add config prefix notifier")
				return err
			}
			return nil
		},
		OnStop: func(context.Context) error {
			err := in.EntityTrackers.RemovePrefixNotifier(configPrefixNotifier)
			if err != nil {
				log.Error().Err(err).Msg("failed to remove prefix notifier")
				return err
			}
			return nil
		},
	})

	return entityCache, nil
}

func (c *EntityCache) processUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
	log.Trace().Str("event", event.String()).Msg("Updating entity")
	var rawEntity common.Entity
	if err := unmarshaller.UnmarshalKey("", &rawEntity); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal entity")
		return
	}
	entity := discoveryEntityToCacheEntity(&rawEntity)
	ip := entity.IpAddress
	name := entity.EntityName

	switch event.Type {
	case notifiers.Write:
		log.Trace().Str("entity", entity.EntityId.Prefix+entity.EntityId.Uid).Str("ip", ip).Str("name", name).Msg("new entity")
		c.Put(entity)
	case notifiers.Remove:
		log.Trace().Str("entity", entity.EntityId.Prefix+entity.EntityId.Uid).Str("ip", ip).Str("name", name).Msg("removing entity")
		c.Remove(entity)
	}
}

func discoveryEntityToCacheEntity(entity *common.Entity) *entitycachev1.Entity {
	return NewEntity(EntityID{
		Prefix: entity.Prefix,
		UID:    entity.UID,
	}, entity.IPAddress, entity.Name, entity.Services)
}

// NewEntityCache creates a new, empty EntityCache.
func NewEntityCache() *EntityCache {
	entities := &entitycachev1.EntityCache{
		EntitiesByIpAddress: &entitycachev1.EntityCache_Entities{
			Entities: make(map[string]*entitycachev1.Entity),
		},
		EntitiesByEntityName: &entitycachev1.EntityCache_Entities{
			Entities: make(map[string]*entitycachev1.Entity),
		},
	}
	return &EntityCache{
		entities: entities,
	}
}

// Put maps given IP address and name to the entity it currently represents.
func (c *EntityCache) Put(entity *entitycachev1.Entity) {
	c.Lock()
	defer c.Unlock()

	entityIP := entity.IpAddress
	if entityIP != "" {
		c.entities.EntitiesByIpAddress.Entities[entityIP] = entity
	}

	entityName := entity.EntityName
	if entityName != "" {
		c.entities.EntitiesByEntityName.Entities[entityName] = entity
	}
}

// GetByIP retrieves entity with a given IP address.
func (c *EntityCache) GetByIP(entityIP string) *entitycachev1.Entity {
	c.RLock()
	defer c.RUnlock()

	v, ok := c.entities.EntitiesByIpAddress.Entities[entityIP]
	if !ok {
		return nil
	}
	return v.DeepCopy()
}

// GetByName retrieves entity with a given name.
func (c *EntityCache) GetByName(entityName string) *entitycachev1.Entity {
	c.RLock()
	defer c.RUnlock()

	v, ok := c.entities.EntitiesByEntityName.Entities[entityName]
	if !ok {
		return nil
	}
	return v.DeepCopy()
}

// Clear removes all entities from the cache.
func (c *EntityCache) Clear() {
	c.RLock()
	defer c.RUnlock()
	c.entities.EntitiesByIpAddress = &entitycachev1.EntityCache_Entities{
		Entities: make(map[string]*entitycachev1.Entity),
	}
	c.entities.EntitiesByEntityName = &entitycachev1.EntityCache_Entities{
		Entities: make(map[string]*entitycachev1.Entity),
	}
}

// Remove removes entity from the cache and returns `true` if any of IP address
// or name mapping exists.
// If no such entity was found, returns `false`.
func (c *EntityCache) Remove(entity *entitycachev1.Entity) bool {
	c.Lock()
	defer c.Unlock()

	entityIP := entity.IpAddress
	_, okByIP := c.entities.EntitiesByIpAddress.Entities[entityIP]
	if okByIP {
		delete(c.entities.EntitiesByIpAddress.Entities, entityIP)
	}
	entityName := entity.EntityName
	_, okByName := c.entities.EntitiesByEntityName.Entities[entityName]
	if okByName {
		delete(c.entities.EntitiesByEntityName.Entities, entityName)
	}
	return okByIP || okByName
}

// Entities returns *entitycachev1.EntitiyCache entities.
func (c *EntityCache) Entities() *entitycachev1.EntityCache {
	c.RLock()
	defer c.RUnlock()
	return c.entities.DeepCopy()
}

// Services returns a list of services based on entities in cache.
//
// Each service is identified by 2 values:
// - agent group
// - service name
//
// This shouldn't happen in real world, but entities which have multiple values
// for an agent group is ignored.
// Entities which have multiple values for service name will create one service
// for each of them.
func (c *EntityCache) Services() *entitycachev1.ServicesList {
	c.RLock()
	defer c.RUnlock()

	services := map[ServiceKey]*entitycachev1.Service{}
	overlapping := make(map[pair]int)

	for _, entity := range c.entities.EntitiesByIpAddress.Entities {
		entityServices, err := servicesFromEntity(entity)
		if err != nil {
			log.Trace().Err(err).Str("entity", entity.EntityId.Uid).Msg("Failed getting services from entity. Skipping")
			continue
		}
		var serviceKeys []ServiceKey
		for _, es := range entityServices {
			key := *c.keyFromService(es)
			serviceKeys = append(serviceKeys, key)
			if _, ok := services[key]; !ok {
				services[key] = es
				continue
			}
			Merge(services[key], es)
		}
		// for each pair in entityServices count number of overlapping entities
		for _, pair := range eachPair(serviceKeys) {
			overlapping[pair]++
		}

	}

	entityCache := &entitycachev1.ServicesList{
		Services:            make([]*entitycachev1.Service, 0, len(services)),
		OverlappingServices: make([]*entitycachev1.OverlappingService, 0, len(overlapping)),
	}

	for _, svc := range services {
		entityCache.Services = append(entityCache.Services, svc)
	}
	for k, v := range overlapping {
		entityCache.OverlappingServices = append(entityCache.OverlappingServices, &entitycachev1.OverlappingService{
			Service1:      k.x.Name,
			Service2:      k.y.Name,
			EntitiesCount: int32(v),
		})
	}
	return entityCache
}

type pair struct {
	x, y ServiceKey
}

// eachPair returns each pair of elements in a slice. Elements in the pair are sorted so that
// x < y.
func eachPair(services []ServiceKey) []pair {
	n := len(services)
	pairs := make([]pair, 0, n*n)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if services[i].lessThan(services[j]) {
				pairs = append(pairs, pair{
					x: services[i],
					y: services[j],
				})
			} else {
				pairs = append(pairs, pair{
					x: services[j],
					y: services[i],
				})
			}
		}
	}
	return pairs
}

// ServiceIDsFromEntity returns a list of services the entity is a part of.
func ServiceIDsFromEntity(entity *entitycachev1.Entity) []services.ServiceID {
	var svcs []services.ServiceID
	if entity != nil {
		svcs = make([]services.ServiceID, 0, len(entity.Services))
		for _, service := range entity.Services {
			svcs = append(svcs, services.ServiceID{
				Service: service,
			})
		}
	}
	return svcs
}

func servicesFromEntity(entity *entitycachev1.Entity) ([]*entitycachev1.Service, error) {
	svcIDs := ServiceIDsFromEntity(entity)
	svcs := make([]*entitycachev1.Service, 0, len(svcIDs))
	for _, svc := range svcIDs {
		svcs = append(svcs, &entitycachev1.Service{
			Name:          svc.Service,
			EntitiesCount: 1,
		})
	}
	return svcs, nil
}
