package components

import (
	"math"

	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
)

// Max takes array of signals and emits maximum value.
type Max struct{}

// Name implements runtime.Component.
func (*Max) Name() string { return "Max" }

// Type implements runtime.Component.
func (*Max) Type() runtime.ComponentType { return runtime.ComponentTypeSignalProcessor }

// ShortDescription implements runtime.Component.
func (*Max) ShortDescription() string { return "" }

// IsActuator implements runtime.Component.
func (*Max) IsActuator() bool { return false }

// Make sure Max complies with Component interface.
var _ runtime.Component = (*Max)(nil)

// NewMaxAndOptions creates a new Max Component.
func NewMaxAndOptions(_ *policylangv1.Max, _ runtime.ComponentID, _ iface.Policy) (runtime.Component, fx.Option, error) {
	max := Max{}
	return &max, fx.Options(), nil
}

// Execute implements runtime.Component.Execute.
func (max *Max) Execute(inPortReadings runtime.PortToReading, tickInfo runtime.TickInfo) (runtime.PortToReading, error) {
	maxValue := -math.MaxFloat64
	inputs := inPortReadings.ReadRepeatedReadingPort("inputs")
	output := runtime.InvalidReading()

	if len(inputs) > 0 {
		for _, singleInput := range inputs {
			if !singleInput.Valid() {
				return runtime.PortToReading{
					"output": []runtime.Reading{output},
				}, nil
			}
			if singleInput.Value() > maxValue {
				maxValue = singleInput.Value()
			}
		}
		output = runtime.NewReading(maxValue)
	} else {
		output = runtime.InvalidReading()
	}

	return runtime.PortToReading{
		"output": []runtime.Reading{output},
	}, nil
}

// DynamicConfigUpdate is a no-op for Max.
func (max *Max) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {}
