package flowcontrol

import (
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/actuators"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/resources/classifier"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/resources/fluxmeter"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/service"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/servicegetter"
)

// Module returns the fx options for dataplane side pieces of policy.
func Module() fx.Option {
	return fx.Options(
		actuators.Module(),
		fluxmeter.Module(),
		classifier.Module(),
		service.Module(),
		servicegetter.Module,
		fx.Provide(
			NewEngine,
		),
	)
}
