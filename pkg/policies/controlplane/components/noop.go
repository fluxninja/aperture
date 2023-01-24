package components

import (
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

// NoOp takes array of signals and emits noOpimum value.
type NoOp struct{}

// Make sure NoOp complies with Component interface.
var _ runtime.Component = (*NoOp)(nil)

// Name implements runtime.Component.
func (*NoOp) Name() string { return "NoOp" }

// Type implements runtime.Component.
func (*NoOp) Type() runtime.ComponentType { return runtime.ComponentTypeSignalProcessor }

// Execute implements runtime.Component.Execute.
func (noOp *NoOp) Execute(inPortReadings runtime.PortToReading, tickInfo runtime.TickInfo) (runtime.PortToReading, error) {
	return inPortReadings, nil
}

// DynamicConfigUpdate is a no-op for NoOp.
func (noOp *NoOp) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {}
