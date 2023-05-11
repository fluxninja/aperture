// +kubebuilder:validation:Optional
package kubernetes

import (
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/v2/pkg/policies/autoscale/kubernetes/actuators"
	"github.com/fluxninja/aperture/v2/pkg/policies/autoscale/kubernetes/config"
	"github.com/fluxninja/aperture/v2/pkg/policies/autoscale/kubernetes/discovery"
	"github.com/fluxninja/aperture/v2/pkg/policies/autoscale/kubernetes/service"
)

// Module returns the fx options for infra integrations of policy.
func Module() fx.Option {
	return fx.Options(
		actuators.Module(),
		service.Module(),
		discovery.Module(),
		config.Module(),
	)
}
