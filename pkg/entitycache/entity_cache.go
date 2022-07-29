package entitycache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"go.uber.org/fx"

	heartbeatv1 "github.com/FluxNinja/aperture/api/gen/proto/go/aperture/plugins/fluxninja/v1"
	"github.com/FluxNinja/aperture/pkg/agentinfo"
	"github.com/FluxNinja/aperture/pkg/config"
	"github.com/FluxNinja/aperture/pkg/discovery/common"
	"github.com/FluxNinja/aperture/pkg/log"
	"github.com/FluxNinja/aperture/pkg/notifiers"
	"github.com/FluxNinja/aperture/pkg/services"
)

const (
	debugEndpoint = "/debug/entity_cache"
)

// ServiceKey holds key for service.
type ServiceKey struct {
	AgentGroup string `json:"agent_group"`
	Namespace  string `json:"namespace"`
	Name       string `json:"name"`
}

func (sk ServiceKey) lessThan(sk2 ServiceKey) bool {
	if sk.AgentGroup < sk2.AgentGroup {
		return true
	} else if sk.AgentGroup > sk2.AgentGroup {
		return false
	}
	if sk.Namespace < sk2.Namespace {
		return true
	} else if sk.Namespace > sk2.Namespace {
		return false
	}
	if sk.Name < sk2.Name {
		return true
	}
	return false
}

// KeyFromService returns a service key for given service.
func KeyFromService(service *heartbeatv1.Service) *ServiceKey {
	return &ServiceKey{
		AgentGroup: service.AgentGroup,
		Namespace:  service.Namespace,
		Name:       service.Name,
	}
}

// Merge merges `mergedService` into `originalService`. This sums `EnititesCount`.
func Merge(originalService, mergedService *heartbeatv1.Service) {
	originalService.EntitiesCount += mergedService.EntitiesCount
}

// Entity represents a pod, vm, etc.
//
// Entities can be created with NewEntity.
type Entity struct {
	// ID of the entity
	ID EntityID `json:"id,omitempty"`

	// IP Address of this entity
	IPAddress string `json:"ip_address"`
	// Namespace in which this entity belongs
	Namespace string `json:"namespace"`
	// AgentGroup needs to be explicitly set using SetAgentGroup.
	AgentGroup string `json:"agent_group"`
	// Services is a List of names of services this entity is a part of.
	// We store "well-known-labels" (identifying a service) in a separate
	// fields for easier access. Note: we could store `[]agent_core/services/ServiceID`
	// here directly, but right now it'd cause cyclic dependency.
	Services []string `json:"services"`
}

// IP returns IP of this entity.
func (e *Entity) IP() string {
	return e.IPAddress
}

// Name returns Name of this entity. Currently this is based on ID.
func (e *Entity) Name() string {
	return fmt.Sprintf("%v-%v", e.ID.Prefix, e.ID.UID)
}

// SetAgentGroup sets agentGroup.
func (e *Entity) SetAgentGroup(agentGroup string) {
	e.AgentGroup = agentGroup
}

// EntityID is a unique Entity identifier.
type EntityID struct {
	Prefix string `json:"prefix,omitempty"`
	UID    string `json:"uid,omitempty"`
}

// NewEntity creates a new entity from ID and IP address from the tagger.
func NewEntity(id EntityID, namespace, ipAddress string, services []string) *Entity {
	return &Entity{
		ID:        id,
		Namespace: namespace,
		IPAddress: ipAddress,
		Services:  services,
	}
}

// EntityCache maps IP addresses and Entity names to entities.
type EntityCache struct {
	sync.RWMutex
	entitiesByIP   map[string]*Entity
	entitiesByName map[string]*Entity
	agentGroup     string
}

// FxIn are the parameters for ProvideEntityCache.
type FxIn struct {
	fx.In
	Lifecycle      fx.Lifecycle
	Router         *mux.Router
	AgentInfo      *agentinfo.AgentInfo
	EntityTrackers notifiers.Trackers `name:"entity_trackers"`
}

// ProvideEntityCache creates Entity Cache.
func ProvideEntityCache(in FxIn) (*EntityCache, error) {
	entityCache := NewEntityCache()
	agentGroup := in.AgentInfo.GetAgentGroup()
	entityCache.agentGroup = agentGroup

	in.Router.HandleFunc(debugEndpoint, entityCache.DumpHandler)

	// create a ConfigPrefixNotifier
	configPrefixNotifier := &notifiers.UnmarshalPrefixNotifier{
		GetUnmarshallerFunc: config.KoanfUnmarshallerConstructor{}.NewKoanfUnmarshaller,
		UnmarshalNotifyFunc: entityCache.processUpdate,
	}

	in.Lifecycle.Append(
		fx.Hook{
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
	entity.SetAgentGroup(c.agentGroup)
	ip := entity.IP()
	name := entity.Name()

	switch event.Type {
	case notifiers.Write:
		log.Trace().Str("entity", entity.ID.Prefix+entity.ID.UID).Str("ip", ip).Str("name", name).Msg("new entity")
		c.Put(entity)
	case notifiers.Remove:
		log.Trace().Str("entity", entity.ID.Prefix+entity.ID.UID).Str("ip", ip).Str("name", name).Msg("removing entity")
		c.Remove(entity)
	}
}

func discoveryEntityToCacheEntity(entity *common.Entity) *Entity {
	return NewEntity(EntityID{
		Prefix: entity.Prefix,
		UID:    entity.UID,
	}, entity.Namespace, entity.IPAddress, entity.Services)
}

// NewEntityCache creates a new, empty EntityCache.
func NewEntityCache() *EntityCache {
	entitiesByIP := make(map[string]*Entity)
	entitiesByName := make(map[string]*Entity)
	return &EntityCache{
		entitiesByIP:   entitiesByIP,
		entitiesByName: entitiesByName,
	}
}

// Put maps given IP address and name to the entity it currently represents.
func (c *EntityCache) Put(entity *Entity) {
	c.Lock()
	defer c.Unlock()

	entityIP := entity.IP()
	if entityIP != "" {
		c.entitiesByIP[entityIP] = entity
	}

	entityName := entity.Name()
	if entityName != "" {
		c.entitiesByName[entityName] = entity
	}
}

// GetByIP retrieves entity with a given IP address.
func (c *EntityCache) GetByIP(entityIP string) *Entity {
	c.RLock()
	defer c.RUnlock()
	v, ok := c.entitiesByIP[entityIP]
	if !ok {
		return nil
	}
	return v
}

// GetByName retrieves entity with a given name.
func (c *EntityCache) GetByName(entityName string) *Entity {
	c.RLock()
	defer c.RUnlock()
	v, ok := c.entitiesByName[entityName]
	if !ok {
		return nil
	}
	return v
}

// Clear removes all entities from the cache.
func (c *EntityCache) Clear() {
	c.RLock()
	defer c.RUnlock()
	c.entitiesByIP = make(map[string]*Entity)
	c.entitiesByName = make(map[string]*Entity)
}

// Remove removes entity from the cache and returns `true` if any of IP address
// or name mapping exists.
// If no such entity was found, returns `false`.
func (c *EntityCache) Remove(entity *Entity) bool {
	c.Lock()
	defer c.Unlock()

	entityIP := entity.IP()
	_, okByIP := c.entitiesByIP[entityIP]
	if okByIP {
		delete(c.entitiesByIP, entityIP)
	}
	entityName := entity.Name()
	_, okByName := c.entitiesByName[entityName]
	if okByName {
		delete(c.entitiesByName, entityName)
	}
	return okByIP || okByName
}

// Services returns a list of services based on entities in cache.
//
// Each service is identified by 3 values:
// - agent group
// - namespace
// - service name
//
// This shouldn't happen in real world, but entities which have multiple values
// either for agent group or namespace are ignored.
// Entities which have multiple values for service name will create one service
// for each of them.
func (c *EntityCache) Services() ([]*heartbeatv1.Service, []*heartbeatv1.OverlappingService) {
	c.RLock()
	defer c.RUnlock()
	services := map[ServiceKey]*heartbeatv1.Service{}
	overlapping := make(map[pair]int)

	for _, entity := range c.entitiesByIP {
		entityServices, err := servicesFromEntity(entity)
		if err != nil {
			log.Trace().Err(err).Str("entity", entity.ID.UID).Msg("Failed getting services from entity. Skipping")
			continue
		}
		var serviceKeys []ServiceKey
		for _, es := range entityServices {
			key := *KeyFromService(es)
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
	ret := make([]*heartbeatv1.Service, 0, len(services))
	for _, svc := range services {
		ret = append(ret, svc)
	}
	retOverlapping := make([]*heartbeatv1.OverlappingService, 0, len(overlapping))
	for k, v := range overlapping {
		retOverlapping = append(retOverlapping, &heartbeatv1.OverlappingService{
			Service1: &heartbeatv1.ServiceKey{
				AgentGroup: k.x.AgentGroup,
				Namespace:  k.x.Namespace,
				Name:       k.x.Name,
			},
			Service2: &heartbeatv1.ServiceKey{
				AgentGroup: k.y.AgentGroup,
				Namespace:  k.y.Namespace,
				Name:       k.y.Name,
			},
			EntitiesCount: int32(v),
		})
	}
	return ret, retOverlapping
}

// DumpHandler is used to return entity cache data in JSON format.
func (c *EntityCache) DumpHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(c.entitiesByIP)
	if err != nil {
		log.Error().Err(err).Msg("Error writing entity cache response body")
		http.Error(w, "", http.StatusInternalServerError)
	}
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
func ServiceIDsFromEntity(entity *Entity) ([]services.ServiceID, error) {
	if entity.Namespace == "" {
		return nil, errors.New("missing namespace")
	}

	if len(entity.Services) == 0 {
		return nil, errors.New("missing services")
	}

	svcs := make([]services.ServiceID, 0, len(entity.Services))
	for _, service := range entity.Services {
		svcs = append(svcs, services.ServiceID{
			AgentGroup: entity.AgentGroup,
			Namespace:  entity.Namespace,
			Service:    service,
		})
	}
	return svcs, nil
}

func servicesFromEntity(entity *Entity) ([]*heartbeatv1.Service, error) {
	if entity.AgentGroup == "" {
		return nil, errors.New("missing agent group")
	}

	svcIDs, err := ServiceIDsFromEntity(entity)
	if err != nil {
		return nil, err
	}

	svcs := make([]*heartbeatv1.Service, 0, len(svcIDs))
	for _, svc := range svcIDs {
		svcs = append(svcs, &heartbeatv1.Service{
			AgentGroup:    entity.AgentGroup,
			Namespace:     svc.Namespace,
			Name:          svc.Service,
			EntitiesCount: 1,
		})
	}
	return svcs, nil
}
