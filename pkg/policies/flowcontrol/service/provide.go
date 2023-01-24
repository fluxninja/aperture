package service

import (
	"go.uber.org/fx"

	awsgateway "github.com/fluxninja/aperture/pkg/policies/flowcontrol/service/aws-gateway"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/service/check"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/service/envoy"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/service/preview"
)

// Module is a set of default providers for flowcontrol components.
func Module() fx.Option {
	return fx.Options(
		check.Module(),
		envoy.Module(),
		awsgateway.Module(),
		preview.Module(),
	)
}
