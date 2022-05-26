package component

import (
	"math"

	"go.uber.org/fx"

	policylangv1 "aperture.tech/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"aperture.tech/aperture/pkg/policies/apis/policyapi"
	"aperture.tech/aperture/pkg/policies/controlplane/reading"
	"aperture.tech/aperture/pkg/policies/controlplane/runtime"
)

// Max takes array of signals and emits maximum value.
type Max struct{}

// Make sure Max complies with Component interface.
var _ runtime.Component = (*Max)(nil)

// NewMaxAndOptions creates a new Max Component.
func NewMaxAndOptions(maxProto *policylangv1.Max, componentIndex int, policyReadAPI policyapi.PolicyReadAPI) (runtime.Component, fx.Option, error) {
	max := Max{}
	return &max, fx.Options(), nil
}

// Execute implements runtime.Component.Execute.
func (max *Max) Execute(inPortReadings runtime.PortToValue, tickInfo runtime.TickInfo) (runtime.PortToValue, error) {
	maxValue := -math.MaxFloat64
	inputs := inPortReadings.ReadRepeatedValuePort("inputs")
	output := reading.NewInvalid()

	if len(inputs) > 0 {
		for _, singleInput := range inputs {
			if !singleInput.Valid {
				return runtime.PortToValue{
					"output": []reading.Reading{output},
				}, nil
			}
			if singleInput.Value > maxValue {
				maxValue = singleInput.Value
			}
		}
		output = reading.New(maxValue)
	} else {
		output = reading.NewInvalid()
	}

	return runtime.PortToValue{
		"output": []reading.Reading{output},
	}, nil
}
