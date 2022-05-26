package controller

import (
	"fmt"
	"math"

	"go.uber.org/fx"

	policylangv1 "aperture.tech/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"aperture.tech/aperture/pkg/policies/apis/policyapi"
	"aperture.tech/aperture/pkg/policies/controlplane/reading"
	"aperture.tech/aperture/pkg/policies/controlplane/runtime"
)

// GradientController describes gradient values.
type GradientController struct {
	minGradient float64
	maxGradient float64
	tolerance   float64
}

// Make sure Gradient complies with Controller interface.
var _ Controller = (*GradientController)(nil)

// NewGradientControllerAndOptions creates a Gradient Controller Component and its fx options.
func NewGradientControllerAndOptions(gradientControllerProto *policylangv1.GradientController, componentIndex int, policyReadAPI policyapi.PolicyReadAPI) (runtime.Component, fx.Option, error) {
	// Make sure max is greater than min
	if gradientControllerProto.MaxGradient < gradientControllerProto.MinGradient {
		return nil, nil, fmt.Errorf("max_gradient must be greater than min_gradient")
	}

	gradient := &GradientController{
		tolerance:   gradientControllerProto.Tolerance,
		minGradient: gradientControllerProto.MinGradient,
		maxGradient: gradientControllerProto.MaxGradient,
	}

	controller := NewControllerComponent(gradient, componentIndex, policyReadAPI)

	return controller, fx.Options(), nil
}

// ComputeOutput based on previous and current signal reading.
func (g *GradientController) ComputeOutput(signal, setpoint, controlVariable reading.Reading, controllerStateReadAPI ControllerStateReadAPI, tickInfo runtime.TickInfo) (reading.Reading, error) {
	var output reading.Reading
	if setpoint.Valid && signal.Valid && controlVariable.Valid {
		var gradient float64
		gradient = setpoint.Value / signal.Value * g.tolerance
		// clamp to min/max
		if gradient < g.minGradient {
			gradient = g.minGradient
		}
		if gradient > g.maxGradient {
			gradient = g.maxGradient
		}
		if math.IsNaN(gradient) {
			output = reading.NewInvalid()
		} else {
			output = reading.New(controlVariable.Value * gradient)
		}
	} else {
		output = reading.NewInvalid()
	}

	return output, nil
}

// MaintainOutput - Gradient Controller is stateless, so a bump is inevitable.
func (g *GradientController) MaintainOutput(prevSetpoint, currentSetpoint reading.Reading, _ ControllerStateReadAPI, tickInfo runtime.TickInfo) error {
	return nil
}

// WindOutput - Gradient Controller relies on ControllerComponent to store the last output, returning targetOutput should wind the output.
func (g *GradientController) WindOutput(currentOutput, targetOutput reading.Reading, _ ControllerStateReadAPI, tickInfo runtime.TickInfo) (reading.Reading, error) {
	return targetOutput, nil
}
