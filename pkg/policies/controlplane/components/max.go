package components

import (
	"math"

	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

// Max takes array of signals and emits maximum value.
type Max struct{}

// Make sure Max complies with Component interface.
var _ runtime.Component = (*Max)(nil)

// NewMaxAndOptions creates a new Max Component.
func NewMaxAndOptions(maxProto *policylangv1.Max, componentIndex int, policyReadAPI iface.Policy) (runtime.Component, fx.Option, error) {
	max := Max{}
	return &max, fx.Options(), nil
}

// Execute implements runtime.Component.Execute.
func (max *Max) Execute(inPortReadings runtime.PortToValue, tickInfo runtime.TickInfo) (runtime.PortToValue, error) {
	maxValue := -math.MaxFloat64
	inputs := inPortReadings.ReadRepeatedValuePort("inputs")
	output := runtime.InvalidReading()

	if len(inputs) > 0 {
		for _, singleInput := range inputs {
			if !singleInput.Valid() {
				return runtime.PortToValue{
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

	return runtime.PortToValue{
		"output": []runtime.Reading{output},
	}, nil
}

// DynamicConfigUpdate is a no-op for Max.
func (max *Max) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {}
