package static

import (
	"context"

	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/discovery/common"
	"github.com/fluxninja/aperture/pkg/discovery/entities"
	sdconfig "github.com/fluxninja/aperture/pkg/discovery/static/config"
	"github.com/fluxninja/aperture/pkg/log"
)

const (
	configKey                 = common.DiscoveryConfigKey + ".static"
	staticEntityTrackerPrefix = "static_entity"
)

// InvokeStaticServiceDiscovery causes statically configured services to be uploaded to the tracker.
func InvokeStaticServiceDiscovery(
	unmarshaller config.Unmarshaller,
	lifecycle fx.Lifecycle,
	entityTrackers *entities.EntityTrackers,
) error {
	var cfg sdconfig.StaticDiscoveryConfig
	if err := unmarshaller.UnmarshalKey(configKey, &cfg); err != nil {
		log.Error().Err(err).Msg("Unable to deserialize static services configuration!")
		return err
	}

	if len(cfg.Entities) == 0 {
		log.Info().Msg("No services configured, disabling static service discovery")
		return nil
	}

	entityEvents := entityTrackers.RegisterServiceDiscovery(staticEntityTrackerPrefix)
	sd := newStaticServiceDiscovery(entityEvents, &cfg)

	lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			return sd.start()
		},
		OnStop: func(_ context.Context) error {
			return sd.stop()
		},
	})

	return nil
}
