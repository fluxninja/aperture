package component

import (
	"math"

	"go.uber.org/fx"

	policylangv1 "aperture.tech/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"aperture.tech/aperture/pkg/policies/apis/policyapi"
	"aperture.tech/aperture/pkg/policies/controlplane/reading"
	"aperture.tech/aperture/pkg/policies/controlplane/runtime"
)

// Sqrt takes an input signal and emits Square Root of it multiplied by scale as output.
type Sqrt struct {
	scale float64
}

// Make sure Sqrt complies with Component interface.
var _ runtime.Component = (*Sqrt)(nil)

// NewSqrtAndOptions creates a new Sqrt Component.
func NewSqrtAndOptions(sqrtProto *policylangv1.Sqrt, componentIndex int, policyReadAPI policyapi.PolicyReadAPI) (runtime.Component, fx.Option, error) {
	sqrt := Sqrt{
		scale: sqrtProto.Scale,
	}
	return &sqrt, fx.Options(), nil
}

// Execute implements runtime.Component.Execute.
func (sqrt *Sqrt) Execute(inPortReadings runtime.PortToValue, tickInfo runtime.TickInfo) (runtime.PortToValue, error) {
	input := inPortReadings.ReadSingleValuePort("input")
	output := reading.NewInvalid()

	if input.Valid {
		// Square root the input and scale it.
		sqrtIn := math.Sqrt(input.Value)
		if !math.IsNaN(sqrtIn) {
			output = reading.New(sqrt.scale * sqrtIn)
		}
	}

	return runtime.PortToValue{
		"output": []reading.Reading{output},
	}, nil
}
