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

// PolynomialRangeFunction .
type PolynomialRangeFunction struct {
	parameters *policylangv1.PolynomialRangeFunction_Parameters
}

// Name implements runtime.Component.
func (*PolynomialRangeFunction) Name() string { return "PolynomialRangeFunction" }

// Type implements runtime.Component.
func (*PolynomialRangeFunction) Type() runtime.ComponentType {
	return runtime.ComponentTypeSignalProcessor
}

// ShortDescription implements runtime.Component.
func (rangeFunc *PolynomialRangeFunction) ShortDescription() string { return "" }

// IsActuator implements runtime.Component.
func (*PolynomialRangeFunction) IsActuator() bool { return false }

// Make sure PolynomialRangeFunction complies with Component interface.
var _ runtime.Component = (*PolynomialRangeFunction)(nil)

// NewRangeFunctionAndOptions returns a new PolynomialRangeFunction and its Fx options.
func NewRangeFunctionAndOptions(rangeFunctionProto *policylangv1.PolynomialRangeFunction, _ runtime.ComponentID, _ iface.Policy) (runtime.Component, fx.Option, error) {
	arith := PolynomialRangeFunction{
		parameters: rangeFunctionProto.Parameters,
	}
	return &arith, fx.Options(), nil
}

// Execute implements runtime.Component.Execute.
func (rangeFunc *PolynomialRangeFunction) Execute(inPortReadings runtime.PortToReading, circuitAPI runtime.CircuitAPI) (runtime.PortToReading, error) {
	// Helper function to return a runtime.PortToReading with given value.
	returnReading := func(value float64) (runtime.PortToReading, error) {
		return runtime.PortToReading{
			"output": []runtime.Reading{runtime.NewReading(value)},
		}, nil
	}

	// Helper function to return an invalid reading.
	returnInvalidReading := func() (runtime.PortToReading, error) {
		return runtime.PortToReading{
			"output": []runtime.Reading{runtime.InvalidReading()},
		}, nil
	}

	input := inPortReadings.ReadSingleReadingPort("input")
	if !input.Valid() {
		return returnInvalidReading()
	}

	// Extract datapoint values
	startInput, startOutput := rangeFunc.parameters.Start.Input, rangeFunc.parameters.Start.Output
	endInput, endOutput := rangeFunc.parameters.End.Input, rangeFunc.parameters.End.Output
	inputVal := input.Value()

	// Edge case: both input values are equal
	if endInput == startInput {
		if inputVal <= startInput {
			return returnReading(startOutput)
		}
		return returnReading(endOutput)
	}

	// Compute normalized value for interpolation
	t := (inputVal - startInput) / (endInput - startInput)

	// Check bounds using t and order of start and end inputs
	if t <= 0 {
		return returnReading(startOutput)
	}

	if t >= 1 {
		return returnReading(endOutput)
	}

	// Compute the range function value based on the degree of the polynomial curve
	return returnReading(startOutput + math.Pow(t, rangeFunc.parameters.Degree)*(endOutput-startOutput))
}

// DynamicConfigUpdate is a no-op for PolynomialRangeFunction.
func (rangeFunc *PolynomialRangeFunction) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
}
