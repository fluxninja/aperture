package static

import (
	"context"

	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/discovery/common"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/notifiers"
)

const configKey = common.DiscoveryConfigKey + ".static"

// EntityConfig describes a single entity.
type EntityConfig struct {
	IPAddress string `json:"ip_address"`
	Namespace string `json:"namespace"`
	Prefix    string `json:"prefix"`
	UID       string `json:"uid"`
}

// ServiceConfig describes a service and its entities.
type ServiceConfig struct {
	Name     string          `json:"name"`
	Entities []*EntityConfig `json:"entities"`
}

// StaticDiscoveryConfig for pre-determined list of services.
// swagger:model
type StaticDiscoveryConfig struct {
	Services []*ServiceConfig `json:"services"`
}

// FxIn describes parameters passed to k8s discovery constructor.
type FxIn struct {
	fx.In
	Unmarshaller   config.Unmarshaller
	Lifecycle      fx.Lifecycle
	EntityTrackers notifiers.Trackers `name:"entity_trackers"`
}

// InvokeStaticServiceDiscovery causes statically configured services to be uploaded to the tracker.
func InvokeStaticServiceDiscovery(in FxIn) error {
	var cfg StaticDiscoveryConfig
	if err := in.Unmarshaller.UnmarshalKey(configKey, &cfg); err != nil {
		log.Error().Err(err).Msg("Unable to deserialize static services configuration!")
		return err
	}

	sd, err := newStaticServiceDiscovery(in.EntityTrackers, cfg)
	if err != nil {
		log.Info().Err(err).Msg("Failed to create static discovery service")
		return nil
	}

	in.Lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			return sd.start()
		},
		OnStop: func(_ context.Context) error {
			return nil
		},
	})

	return nil
}
