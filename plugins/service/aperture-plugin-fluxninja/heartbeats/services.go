package heartbeats

import (
	entitycachev1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/entitycache/v1"
	heartbeatv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/plugins/fluxninja/v1"
	"github.com/fluxninja/aperture/pkg/entitycache"
	"github.com/fluxninja/aperture/pkg/log"
)

// ServiceKey holds key for service.
type ServiceKey struct {
	Name string `json:"name"`
}

func (sk ServiceKey) lessThan(sk2 ServiceKey) bool {
	return sk.Name < sk2.Name
}

// keyFromService returns a service key for given service.
func keyFromService(service *heartbeatv1.Service) *ServiceKey {
	return &ServiceKey{
		Name: service.Name,
	}
}

// Merge merges `mergedService` into `originalService`. This sums `EntitiesCount`.
func Merge(originalService, mergedService *heartbeatv1.Service) {
	originalService.EntitiesCount += mergedService.EntitiesCount
}

// populateServicesList returns a list of populateServicesList based on entities in cache.
//
// Each service is identified by 2 values:
// - agent group
// - service name
//
// This shouldn't happen in real world, but entities which have multiple values
// for an agent group is ignored.
// Entities which have multiple values for service name will create one service
// for each of them.
func populateServicesList(c *entitycache.EntityCache) *heartbeatv1.ServicesList {
	services := map[ServiceKey]*heartbeatv1.Service{}
	overlapping := make(map[pair]int)
	entities := c.GetEntities()

	for _, entity := range entities.EntitiesByIpAddress.Entities {
		entityServices, err := servicesFromEntity(entity)
		if err != nil {
			log.Trace().Err(err).Str("entity", entity.Uid).Msg("Failed getting services from entity. Skipping")
			continue
		}
		var serviceKeys []ServiceKey
		for _, es := range entityServices {
			key := *keyFromService(es)
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

	entityCache := &heartbeatv1.ServicesList{
		Services:            make([]*heartbeatv1.Service, 0, len(services)),
		OverlappingServices: make([]*heartbeatv1.OverlappingService, 0, len(overlapping)),
	}

	for _, svc := range services {
		entityCache.Services = append(entityCache.Services, svc)
	}
	for k, v := range overlapping {
		entityCache.OverlappingServices = append(entityCache.OverlappingServices, &heartbeatv1.OverlappingService{
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

func servicesFromEntity(entity *entitycachev1.Entity) ([]*heartbeatv1.Service, error) {
	svcIDs := entity.Services
	svcs := make([]*heartbeatv1.Service, 0, len(svcIDs))
	for _, svc := range svcIDs {
		svcs = append(svcs, &heartbeatv1.Service{
			Name:          svc,
			EntitiesCount: 1,
		})
	}
	return svcs, nil
}
