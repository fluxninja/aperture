package circuitfactory

import (
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
)

// Module for circuit and component factory run via the main app.
func Module() fx.Option {
	return fx.Options(
		runtime.CircuitModule(),
		FactoryModule(),
	)
}
