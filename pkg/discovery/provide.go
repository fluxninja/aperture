package discovery

import (
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/v2/pkg/discovery/entities"
	"github.com/fluxninja/aperture/v2/pkg/discovery/kubernetes"
)

// Module returns an fx.Option that provides the discovery module.
func Module() fx.Option {
	return fx.Options(
		kubernetes.Module(),
		entities.Module(),
	)
}
