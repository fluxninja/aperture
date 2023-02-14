package components

import (
	"math"

	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/components/tristate"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

// Min takes array of signals and emits minimum value.
type Min struct{}

// Name implements runtime.Component.
func (*Min) Name() string { return "Min" }

// Type implements runtime.Component.
func (*Min) Type() runtime.ComponentType { return runtime.ComponentTypeSignalProcessor }

// ShortDescription implements runtime.Component.
func (*Min) ShortDescription() string { return "" }

// Make sure Min complies with Component interface.
var _ runtime.Component = (*Min)(nil)

// NewMinAndOptions creates a new Min Component.
func NewMinAndOptions(_ *policylangv1.Min, _ string, _ iface.Policy) (runtime.Component, fx.Option, error) {
	min := Min{}
	return &min, fx.Options(), nil
}

// Execute implements runtime.Component.Execute.
func (min *Min) Execute(inPortReadings runtime.PortToReading, tickInfo runtime.TickInfo) (runtime.PortToReading, error) {
	minValue := math.MaxFloat64
	inputs := inPortReadings.ReadRepeatedReadingPort("inputs")
	output := runtime.InvalidReading()

	if len(inputs) > 0 {
		for _, singleInput := range inputs {
			if tristate.FromReading(singleInput) == tristate.False {
				return runtime.PortToReading{
					"output": []runtime.Reading{output},
				}, nil
			}
			if singleInput.Value() < minValue {
				minValue = singleInput.Value()
			}
		}
		output = runtime.NewReading(minValue)
	} else {
		output = runtime.InvalidReading()
	}

	return runtime.PortToReading{
		"output": []runtime.Reading{output},
	}, nil
}

// DynamicConfigUpdate is a no-op for Min.
func (min *Min) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {}
