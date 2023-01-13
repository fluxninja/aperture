package infra

import (
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/policies/infra/actuators"
)

// Module returns the fx options for infra integrations of policy.
func Module() fx.Option {
	return fx.Options(
		actuators.Module(),
	)
}
