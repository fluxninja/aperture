package discovery

import (
	"go.uber.org/fx"

	"github.com/FluxNinja/aperture/pkg/discovery/kubernetes"
	"github.com/FluxNinja/aperture/pkg/discovery/static"
)

// Module returns an fx.Option that provides the discovery module.
func Module() fx.Option {
	return fx.Options(
		fx.Invoke(
			kubernetes.InvokeKubernetesServiceDiscovery,
			static.InvokeStaticServiceDiscovery,
		),
	)
}
