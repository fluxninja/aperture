package dataplane

import (
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/policies/dataplane/actuators"
	"github.com/fluxninja/aperture/pkg/policies/dataplane/resources/fluxmeter"
)

// PolicyModule returns the fx options for dataplane side pieces of policy.
func PolicyModule() fx.Option {
	return fx.Options(
		actuators.Module(),
		fluxmeter.Module(),
		fx.Provide(
			ProvideEngineAPI,
		),
	)
}
