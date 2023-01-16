package discovery

import (
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/discovery/kubernetes"
	"github.com/fluxninja/aperture/pkg/discovery/static"
)

// Module returns an fx.Option that provides the discovery module.
func Module() fx.Option {
	return fx.Options(
		kubernetes.Module(),
		fx.Invoke(
			static.InvokeStaticServiceDiscovery,
		),
	)
}
