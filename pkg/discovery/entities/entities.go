package entities

import (
	"context"
	"errors"
	"sync"

	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"

	entitiesv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/discovery/entities/v1"
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

// Entity is an immutable wrapper over *entitiesv1.Entity.
type Entity struct {
	immutableEntity *entitiesv1.Entity
}

// NewEntity creates a new immutable entity from the copy of given entity.
func NewEntity(entity *entitiesv1.Entity) Entity {
	return Entity{immutableEntity: proto.Clone(entity).(*entitiesv1.Entity)}
}

// NewEntityFromImmutable creates a new immutable entity, assuming given entity is immutable.
//
// This allows avoiding a copy compared to NewEntity.
func NewEntityFromImmutable(entity *entitiesv1.Entity) Entity {
	return Entity{immutableEntity: entity}
}

// Clone returns a mutable copy of the entity.
func (e Entity) Clone() *entitiesv1.Entity {
	return proto.Clone(e.immutableEntity).(*entitiesv1.Entity)
}

// UID returns the entity's UID.
func (e Entity) UID() string { return e.immutableEntity.Uid }

// IPAddress returns the entity's IP address.
func (e Entity) IPAddress() string { return e.immutableEntity.IpAddress }

// Name returns the entity's name.
func (e Entity) Name() string { return e.immutableEntity.Name }

// Namespace returns the entity's namespace.
func (e Entity) Namespace() string { return e.immutableEntity.Namespace }

// NodeName returns the entity's node name.
func (e Entity) NodeName() string { return e.immutableEntity.NodeName }

// Services returns list of services the entity belongs to.
//
// The returned slice must not be modified.
func (e Entity) Services() []string { return e.immutableEntity.Services }

// Entities maps IP addresses and Entity names to entities.
type Entities struct {
	sync.RWMutex
	byIP   map[string]Entity
	byName map[string]Entity
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
		config.NewProtobufUnmarshaller,
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

func (e *Entities) processUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
	log.Trace().Str("event", event.String()).Msg("Updating entity")
	entityProto := &entitiesv1.Entity{}
	if err := unmarshaller.Unmarshal(entityProto); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal entity")
		return
	}
	entity := NewEntityFromImmutable(entityProto)
	ip := entity.IPAddress()
	name := entity.Name()

	switch event.Type {
	case notifiers.Write:
		log.Trace().Str("entity", entity.UID()).Str("ip", ip).Str("name", name).Msg("new entity")
		e.Put(entity)
	case notifiers.Remove:
		log.Trace().Str("entity", entity.UID()).Str("ip", ip).Str("name", name).Msg("removing entity")
		e.Remove(entity)
	}
}

// NewEntities creates a new, empty Entities.
func NewEntities() *Entities {
	return &Entities{
		byIP:   make(map[string]Entity),
		byName: make(map[string]Entity),
	}
}

// PutForTest maps given IP address and name to the entity it currently represents.
func (e *Entities) PutForTest(entity *entitiesv1.Entity) {
	e.Put(NewEntity(entity))
}

// Put maps given IP address and name to the entity it currently represents.
func (e *Entities) Put(entity Entity) {
	e.Lock()
	defer e.Unlock()

	entityIP := entity.IPAddress()
	if entityIP != "" {
		// FIXME: would be nice to store Entity directly in the map, but that
		// would require removing the reusal of proto-generated structs as
		// containers.
		e.byIP[entityIP] = entity
	}

	entityName := entity.Name()
	if entityName != "" {
		e.byName[entityName] = entity
	}
}

// GetByIP retrieves entity with a given IP address.
func (e *Entities) GetByIP(entityIP string) (Entity, error) {
	return e.getFromMap(e.byIP, entityIP)
}

// GetByName retrieves entity with a given name.
func (e *Entities) GetByName(entityName string) (Entity, error) {
	return e.getFromMap(e.byName, entityName)
}

func (e *Entities) getFromMap(m map[string]Entity, k string) (Entity, error) {
	e.RLock()
	defer e.RUnlock()

	if len(m) == 0 {
		return Entity{}, errNoEntities
	}

	v, ok := m[k]
	if !ok {
		return Entity{}, errNotFound
	}

	return v, nil
}

var (
	errNotFound   = errors.New("entity not found")
	errNoEntities = errors.New("entity not found (empty cache)")
)

// Clear removes all entities from the cache.
func (e *Entities) Clear() {
	e.Lock()
	defer e.Unlock()
	e.byIP = make(map[string]Entity)
	e.byName = make(map[string]Entity)
}

// Remove removes entity from the cache and returns `true` if any of IP address
// or name mapping exists.
// If no such entity was found, returns `false`.
func (e *Entities) Remove(entity Entity) bool {
	e.Lock()
	defer e.Unlock()

	entityIP := entity.IPAddress()
	_, okByIP := e.byIP[entityIP]
	if okByIP {
		delete(e.byIP, entityIP)
	}

	entityName := entity.Name()
	_, okByName := e.byName[entityName]
	if okByName {
		delete(e.byName, entityName)
	}
	return okByIP || okByName
}

// GetEntities returns *entitiesv1.EntitiyCache entities.
func (e *Entities) GetEntities() *entitiesv1.Entities {
	e.RLock()
	defer e.RUnlock()

	// Not sure what caller will do with the result, let's clone
	return &entitiesv1.Entities{
		EntitiesByIpAddress: cloneEntitiesMap(e.byIP),
		EntitiesByName:      cloneEntitiesMap(e.byName),
	}
}

func cloneEntitiesMap(m map[string]Entity) *entitiesv1.Entities_Entities {
	clones := make(map[string]*entitiesv1.Entity, len(m))
	for k, entity := range m {
		clones[k] = entity.Clone()
	}
	return &entitiesv1.Entities_Entities{
		Entities: clones,
	}
}
