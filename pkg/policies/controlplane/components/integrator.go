package components

import (
	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

// Integrator is a component that accumulates sum of signal every tick.
type Integrator struct {
	sum float64
}

// Name implements runtime.Component.
func (*Integrator) Name() string { return "Integrator" }

// Type implements runtime.Component.
func (*Integrator) Type() runtime.ComponentType { return runtime.ComponentTypeSignalProcessor }

// NewIntegrator creates an integrator component.
func NewIntegrator() runtime.Component {
	integrator := &Integrator{
		sum: 0,
	}
	return integrator
}

// NewIntegratorAndOptions creates an integrator component and its fx options.
func NewIntegratorAndOptions(_ *policylangv1.Integrator, _ int, _ iface.Policy) (runtime.Component, fx.Option, error) {
	return NewIntegrator(), fx.Options(), nil
}

// Execute implements runtime.Component.Execute.
func (in *Integrator) Execute(inPortReadings runtime.PortToValue, tickInfo runtime.TickInfo) (runtime.PortToValue, error) {
	inputVal := inPortReadings.ReadSingleValuePort("input")
	resetVal := inPortReadings.ReadSingleValuePort("reset")
	if resetVal.Valid() && resetVal.Value() > 0 {
		in.sum = 0
	} else if inputVal.Valid() {
		minVal := inPortReadings.ReadSingleValuePort("min")
		maxVal := inPortReadings.ReadSingleValuePort("max")

		value := inputVal.Value()
		if value < minVal.Value() {
			value = minVal.Value()
		} else if value > maxVal.Value() {
			value = maxVal.Value()
		}

		in.sum += value
	}

	return runtime.PortToValue{
		"output": []runtime.Reading{runtime.NewReading(in.sum)},
	}, nil
}

// DynamicConfigUpdate is a no-op for Integrator.
func (in *Integrator) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {}
