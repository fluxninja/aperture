package entitycache

import (
	"context"
	"errors"
	"sync"

	"go.uber.org/fx"

	entitycachev1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/entitycache/v1"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/net/grpcgateway"
	"github.com/fluxninja/aperture/pkg/notifiers"
)

// Module sets up EntityCache with Fx.
func Module() fx.Option {
	return fx.Options(
		notifiers.TrackersConstructor{Name: "entity_trackers_private"}.Annotate(),
		fx.Provide(provideEntityCache),
		grpcgateway.RegisterHandler{Handler: entitycachev1.RegisterEntityCacheServiceHandlerFromEndpoint}.Annotate(),
		fx.Invoke(RegisterEntityCacheService),
	)
}

// EntityCache maps IP addresses and Entity names to entities.
type EntityCache struct {
	sync.RWMutex
	entities *entitycachev1.EntityCache
}

// EntityTrackers allows to register a service discovery for entity cache
//
// Intended to be used during FX initialization.
type EntityTrackers struct {
	trackers     notifiers.Trackers
	hasDiscovery bool
}

// RegisterServiceDiscovery registers service discovery for entity cache and
// returns an EventWriter to push discovery events into.
//
// Keys passed to EventWriter shouldn't be prefixed.
//
// Should be called at FX provide/invoke stage.
func (et *EntityTrackers) RegisterServiceDiscovery(name string) notifiers.EventWriter {
	et.hasDiscovery = true
	return notifiers.NewPrefixedEventWriter(name+".", et.trackers)
}

// HasDiscovery returns whether RegisterServiceDiscovery was called before.
func (et *EntityTrackers) HasDiscovery() bool { return et.hasDiscovery }

// Watcher returns watcher that watches all events from registered service discoveries.
func (et *EntityTrackers) Watcher() notifiers.Watcher { return et.trackers }

// FxIn are the parameters for ProvideEntityCache.
type FxIn struct {
	fx.In
	Lifecycle      fx.Lifecycle
	EntityTrackers notifiers.Trackers `name:"entity_trackers_private"`
}

// provideEntityCache creates Entity Cache.
func provideEntityCache(in FxIn) (*EntityCache, *EntityTrackers) {
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

	return entityCache, &EntityTrackers{trackers: in.EntityTrackers}
}

func (c *EntityCache) processUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
	log.Trace().Str("event", event.String()).Msg("Updating entity")
	entity := &entitycachev1.Entity{}
	if err := unmarshaller.UnmarshalKey("", entity); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal entity")
		return
	}
	ip := entity.IpAddress
	name := entity.Name

	switch event.Type {
	case notifiers.Write:
		log.Trace().Str("entity", entity.Prefix+entity.Uid).Str("ip", ip).Str("name", name).Msg("new entity")
		c.Put(entity)
	case notifiers.Remove:
		log.Trace().Str("entity", entity.Prefix+entity.Uid).Str("ip", ip).Str("name", name).Msg("removing entity")
		c.Remove(entity)
	}
}

// NewEntityCache creates a new, empty EntityCache.
func NewEntityCache() *EntityCache {
	entities := &entitycachev1.EntityCache{
		EntitiesByIpAddress: &entitycachev1.EntityCache_Entities{
			Entities: make(map[string]*entitycachev1.Entity),
		},
		EntitiesByName: &entitycachev1.EntityCache_Entities{
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

	entityName := entity.Name
	if entityName != "" {
		c.entities.EntitiesByName.Entities[entityName] = entity
	}
}

// GetByIP retrieves entity with a given IP address.
func (c *EntityCache) GetByIP(entityIP string) (*entitycachev1.Entity, error) {
	c.RLock()
	defer c.RUnlock()

	if len(c.entities.EntitiesByIpAddress.Entities) == 0 {
		return nil, errNoEntities
	}

	v, ok := c.entities.EntitiesByIpAddress.Entities[entityIP]
	if !ok {
		return nil, errNotFound
	}

	return v.DeepCopy(), nil
}

// GetByName retrieves entity with a given name.
func (c *EntityCache) GetByName(entityName string) (*entitycachev1.Entity, error) {
	c.RLock()
	defer c.RUnlock()

	if len(c.entities.EntitiesByName.Entities) == 0 {
		return nil, errNoEntities
	}

	v, ok := c.entities.EntitiesByName.Entities[entityName]
	if !ok {
		return nil, errNotFound
	}
	return v.DeepCopy(), nil
}

var (
	errNotFound   = errors.New("entity not found")
	errNoEntities = errors.New("entity not found (empty cache)")
)

// Clear removes all entities from the cache.
func (c *EntityCache) Clear() {
	c.RLock()
	defer c.RUnlock()
	c.entities.EntitiesByIpAddress = &entitycachev1.EntityCache_Entities{
		Entities: make(map[string]*entitycachev1.Entity),
	}
	c.entities.EntitiesByName = &entitycachev1.EntityCache_Entities{
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

	entityName := entity.Name
	_, okByName := c.entities.EntitiesByName.Entities[entityName]
	if okByName {
		delete(c.entities.EntitiesByName.Entities, entityName)
	}
	return okByIP || okByName
}

// GetEntities returns *entitycachev1.EntitiyCache entities.
func (c *EntityCache) GetEntities() *entitycachev1.EntityCache {
	c.RLock()
	defer c.RUnlock()
	return c.entities.DeepCopy()
}
