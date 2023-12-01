package components

import (
	"math"

	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
)

// PolynomialRangeFunction .
type PolynomialRangeFunction struct {
	parameters  *policylangv1.PolynomialRangeFunction_Parameters
	componentID runtime.ComponentID
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

// NewPolynomialRangeFunctionAndOptions returns a new PolynomialRangeFunction and its Fx options.
func NewPolynomialRangeFunctionAndOptions(rangeFunctionProto *policylangv1.PolynomialRangeFunction, componentID runtime.ComponentID, _ iface.Policy) (runtime.Component, fx.Option, error) {
	arith := PolynomialRangeFunction{
		parameters:  rangeFunctionProto.Parameters,
		componentID: componentID,
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

	// Helper function to handle outside range behavior.
	handleOutsideRange := func(t float64, output float64) float64 {
		if t < 0 {
			if rangeFunc.parameters.GetClampToDatapoint() {
				return rangeFunc.parameters.Start.Output
			} else if rangeFunc.parameters.GetClampToCustomValues() != nil {
				return rangeFunc.parameters.GetClampToCustomValues().PreStart
			}
		}

		if t > 1 {
			if rangeFunc.parameters.GetClampToDatapoint() {
				return rangeFunc.parameters.End.Output
			} else if rangeFunc.parameters.GetClampToCustomValues() != nil {
				return rangeFunc.parameters.GetClampToCustomValues().PostEnd
			}
		}

		return output
	}

	input := inPortReadings.ReadSingleReadingPort("input")
	if !input.Valid() {
		return returnInvalidReading()
	}

	// Extract datapoint values
	startInput, startOutput := rangeFunc.parameters.Start.Input, rangeFunc.parameters.Start.Output
	endInput, endOutput := rangeFunc.parameters.End.Input, rangeFunc.parameters.End.Output
	inputVal := input.Value()

	// Compute normalized value for interpolation
	var t float64
	if endInput != startInput {
		t = (inputVal - startInput) / (endInput - startInput)
	} else {
		// Handling the edge case where startInput and endInput are the same
		if inputVal <= startInput {
			t = 0.0
		} else {
			t = 1.0
		}
	}

	// Calculate the polynomial output
	polyOutput := float64(startOutput) + (math.Pow(t, rangeFunc.parameters.Degree) * (endOutput - startOutput))

	// Handle outside range behavior
	output := handleOutsideRange(t, polyOutput)

	return returnReading(output)
}

// DynamicConfigUpdate is a no-op for PolynomialRangeFunction.
func (rangeFunc *PolynomialRangeFunction) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
}
