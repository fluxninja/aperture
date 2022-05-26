package component

import (
	"math"

	"go.uber.org/fx"

	policylangv1 "aperture.tech/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"aperture.tech/aperture/pkg/policies/apis/policyapi"
	"aperture.tech/aperture/pkg/policies/controlplane/reading"
	"aperture.tech/aperture/pkg/policies/controlplane/runtime"
)

// Min takes array of signals and emits minimum value.
type Min struct{}

// Make sure Min complies with Component interface.
var _ runtime.Component = (*Min)(nil)

// NewMinAndOptions creates a new Min Component.
func NewMinAndOptions(minProto *policylangv1.Min, componentIndex int, policyReadAPI policyapi.PolicyReadAPI) (runtime.Component, fx.Option, error) {
	min := Min{}
	return &min, fx.Options(), nil
}

// Execute implements runtime.Component.Execute.
func (min *Min) Execute(inPortReadings runtime.PortToValue, tickInfo runtime.TickInfo) (runtime.PortToValue, error) {
	minValue := math.MaxFloat64
	inputs := inPortReadings.ReadRepeatedValuePort("inputs")
	output := reading.NewInvalid()

	if len(inputs) > 0 {
		for _, singleInput := range inputs {
			if !singleInput.Valid {
				return runtime.PortToValue{
					"output": []reading.Reading{output},
				}, nil
			}
			if singleInput.Value < minValue {
				minValue = singleInput.Value
			}
		}
		output = reading.New(minValue)
	} else {
		output = reading.NewInvalid()
	}

	return runtime.PortToValue{
		"output": []reading.Reading{output},
	}, nil
}
