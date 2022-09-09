// +kubebuilder:validation:Optional
package static

import (
	"context"

	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/discovery/common"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/notifiers"
)

const (
	configKey                 = common.DiscoveryConfigKey + ".static"
	staticEntityTrackerPrefix = "static_entity"
)

// EntityConfig describes a single entity.
// swagger:model
// +kubebuilder:object:generate=true
type EntityConfig struct {
	// IP address of the entity.
	IPAddress string `json:"ip_address" validate:"required,ip"`
	// UID of the entity.
	UID string `json:"uid"`
	// Name of the entity.
	Name string `json:"name"`
}

// ServiceConfig describes a service and its entities.
// swagger:model
// +kubebuilder:object:generate=true
type ServiceConfig struct {
	// Name of the service.
	Name string `json:"name" validate:"required"`
	// Entities of the service.
	Entities []*EntityConfig `json:"entities"`
}

// StaticDiscoveryConfig for pre-determined list of services.
// swagger:model
// +kubebuilder:object:generate=true
type StaticDiscoveryConfig struct {
	// Services list.
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
			return sd.stop()
		},
	})

	return nil
}
