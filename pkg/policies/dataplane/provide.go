package dataplane

import (
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/policies/dataplane/actuator"
	"github.com/fluxninja/aperture/pkg/policies/dataplane/fluxmeter"
)

// PolicyModule returns the fx options for dataplane side pieces of policy.
func PolicyModule() fx.Option {
	return fx.Options(
		actuator.Module(),
		fluxmeter.Module(),
		fx.Provide(
			ProvideEngineAPI,
		),
	)
}
