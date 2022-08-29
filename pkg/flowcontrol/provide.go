package flowcontrol

import (
	"github.com/fluxninja/aperture/pkg/flowcontrol/common"
	"github.com/fluxninja/aperture/pkg/flowcontrol/envoy"
	"go.uber.org/fx"
)

// Module is a set of default providers for flowcontrol components
//
// Note that the handler needs to be Registered for flowcontrol to be available
// externally.
var Module = fx.Options(
	common.Module,
	envoy.Module,
)
