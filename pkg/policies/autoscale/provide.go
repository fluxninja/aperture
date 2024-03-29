package autoscale

import (
	"github.com/fluxninja/aperture/v2/pkg/policies/autoscale/kubernetes"
	"go.uber.org/fx"
)

// Module returns the fx module for the autoscale policy.
func Module() fx.Option {
	return fx.Options(
		kubernetes.Module(),
	)
}
