package flowcontrol

import (
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/actuators"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/api"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/resources/classifier"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/resources/fluxmeter"
)

// Module returns the fx options for dataplane side pieces of policy.
func Module() fx.Option {
	return fx.Options(
		actuators.Module(),
		fluxmeter.Module(),
		classifier.Module(),
		api.Module(),
		fx.Provide(
			NewEngine,
		),
	)
}
