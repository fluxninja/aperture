package static

import (
	"encoding/json"

	"github.com/fluxninja/aperture/v2/pkg/discovery/static/config"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
)

// StaticDiscovery reads entities from config and writes them to tracker.
type StaticDiscovery struct {
	entityEvents notifiers.EventWriter
	config       *config.StaticDiscoveryConfig
}

func newStaticServiceDiscovery(
	entityEvents notifiers.EventWriter,
	config *config.StaticDiscoveryConfig,
) *StaticDiscovery {
	return &StaticDiscovery{
		entityEvents: entityEvents,
		config:       config,
	}
}

// start loads all configured entities into the tracker.
func (sd *StaticDiscovery) start() error {
	log.Debug().Msgf("Uploading %v pre-configured entities to tracker", len(sd.config.Entities))
	for i := 0; i < len(sd.config.Entities); i++ {
		entity := &sd.config.Entities[i]
		key := notifiers.Key(entity.GetUid())
		value, err := json.Marshal(entity)
		if err != nil {
			log.Error().Msgf("Error marshaling entity: %v", err)
			return err
		}
		sd.entityEvents.WriteEvent(key, value)
	}
	log.Info().Msgf("Uploaded %v pre-configured entities to tracker", len(sd.config.Entities))
	return nil
}

func (sd *StaticDiscovery) stop() error {
	sd.entityEvents.Purge("")
	return nil
}
