package entities

import (
	"context"
	"errors"
	"sync"

	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"

	entitiesv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/discovery/entities/v1"
	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/net/grpcgateway"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
)

// Module sets up Entities with Fx.
func Module() fx.Option {
	return fx.Options(
		notifiers.TrackersConstructor{Name: "entity_trackers_private"}.Annotate(),
		fx.Provide(provideEntities),
		grpcgateway.RegisterHandler{Handler: entitiesv1.RegisterEntitiesServiceHandlerFromEndpoint}.Annotate(),
		fx.Invoke(RegisterEntitiesService),
	)
}

// Entities maps IP addresses and Entity names to entities.
type Entities struct {
	sync.RWMutex
	entities *entitiesv1.Entities
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
// Keys passed to EventWriter should not be prefixed.
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

// FxIn are the parameters for ProvideEntities.
type FxIn struct {
	fx.In
	Lifecycle      fx.Lifecycle
	EntityTrackers notifiers.Trackers `name:"entity_trackers_private"`
}

// provideEntities creates Entity Cache.
func provideEntities(in FxIn) (*Entities, *EntityTrackers, error) {
	entityCache := NewEntities()

	// create a ConfigPrefixNotifier
	configPrefixNotifier, err := notifiers.NewUnmarshalPrefixNotifier("",
		entityCache.processUpdate,
		config.KoanfUnmarshallerConstructor{}.NewKoanfUnmarshaller,
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to create config prefix notifier")
		return nil, nil, err
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

	return entityCache, &EntityTrackers{trackers: in.EntityTrackers}, nil
}

func (c *Entities) processUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
	log.Trace().Str("event", event.String()).Msg("Updating entity")
	entity := &entitiesv1.Entity{}
	if err := unmarshaller.UnmarshalKey("", entity); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal entity")
		return
	}
	ip := entity.IpAddress
	name := entity.Name

	switch event.Type {
	case notifiers.Write:
		log.Trace().Str("entity", entity.Uid).Str("ip", ip).Str("name", name).Msg("new entity")
		c.Put(entity)
	case notifiers.Remove:
		log.Trace().Str("entity", entity.Uid).Str("ip", ip).Str("name", name).Msg("removing entity")
		c.Remove(entity)
	}
}

// NewEntities creates a new, empty Entities.
func NewEntities() *Entities {
	entities := &entitiesv1.Entities{
		EntitiesByIpAddress: &entitiesv1.Entities_Entities{
			Entities: make(map[string]*entitiesv1.Entity),
		},
		EntitiesByName: &entitiesv1.Entities_Entities{
			Entities: make(map[string]*entitiesv1.Entity),
		},
	}
	return &Entities{
		entities: entities,
	}
}

// Put maps given IP address and name to the entity it currently represents.
func (c *Entities) Put(entity *entitiesv1.Entity) {
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
func (c *Entities) GetByIP(entityIP string) (*entitiesv1.Entity, error) {
	c.RLock()
	defer c.RUnlock()

	if len(c.entities.EntitiesByIpAddress.Entities) == 0 {
		return nil, errNoEntities
	}

	v, ok := c.entities.EntitiesByIpAddress.Entities[entityIP]
	if !ok {
		return nil, errNotFound
	}

	return proto.Clone(v).(*entitiesv1.Entity), nil
}

// GetByName retrieves entity with a given name.
func (c *Entities) GetByName(entityName string) (*entitiesv1.Entity, error) {
	c.RLock()
	defer c.RUnlock()

	if len(c.entities.EntitiesByName.Entities) == 0 {
		return nil, errNoEntities
	}

	v, ok := c.entities.EntitiesByName.Entities[entityName]
	if !ok {
		return nil, errNotFound
	}

	return proto.Clone(v).(*entitiesv1.Entity), nil
}

var (
	errNotFound   = errors.New("entity not found")
	errNoEntities = errors.New("entity not found (empty cache)")
)

// Clear removes all entities from the cache.
func (c *Entities) Clear() {
	c.RLock()
	defer c.RUnlock()
	c.entities.EntitiesByIpAddress = &entitiesv1.Entities_Entities{
		Entities: make(map[string]*entitiesv1.Entity),
	}
	c.entities.EntitiesByName = &entitiesv1.Entities_Entities{
		Entities: make(map[string]*entitiesv1.Entity),
	}
}

// Remove removes entity from the cache and returns `true` if any of IP address
// or name mapping exists.
// If no such entity was found, returns `false`.
func (c *Entities) Remove(entity *entitiesv1.Entity) bool {
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

// GetEntities returns *entitiesv1.EntitiyCache entities.
func (c *Entities) GetEntities() *entitiesv1.Entities {
	c.RLock()
	defer c.RUnlock()

	return proto.Clone(c.entities).(*entitiesv1.Entities)
}
