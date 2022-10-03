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

// Min takes array of signals and emits minimum value.
type Min struct{}

// Make sure Min complies with Component interface.
var _ runtime.Component = (*Min)(nil)

// NewMinAndOptions creates a new Min Component.
func NewMinAndOptions(minProto *policylangv1.Min, componentIndex int, policyReadAPI iface.Policy) (runtime.Component, fx.Option, error) {
	min := Min{}
	return &min, fx.Options(), nil
}

// Execute implements runtime.Component.Execute.
func (min *Min) Execute(inPortReadings runtime.PortToValue, tickInfo runtime.TickInfo) (runtime.PortToValue, error) {
	minValue := math.MaxFloat64
	inputs := inPortReadings.ReadRepeatedValuePort("inputs")
	output := runtime.InvalidReading()

	if len(inputs) > 0 {
		for _, singleInput := range inputs {
			if !singleInput.Valid() {
				return runtime.PortToValue{
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

	return runtime.PortToValue{
		"output": []runtime.Reading{output},
	}, nil
}

// DynamicConfigUpdate is a no-op for Min.
func (min *Min) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {}
