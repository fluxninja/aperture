package controller

import (
	"fmt"
	"math"

	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

// GradientController describes gradient values.
type GradientController struct {
	minGradient float64
	maxGradient float64
	slope       float64
}

// Make sure Gradient complies with Controller interface.
var _ Controller = (*GradientController)(nil)

// NewGradientControllerAndOptions creates a Gradient Controller Component and its fx options.
func NewGradientControllerAndOptions(gradientControllerProto *policylangv1.GradientController, componentIndex int, policyReadAPI iface.Policy) (runtime.Component, fx.Option, error) {
	// Make sure max is greater than min
	if gradientControllerProto.MaxGradient < gradientControllerProto.MinGradient {
		return nil, nil, fmt.Errorf("max_gradient must be greater than min_gradient")
	}

	gradient := &GradientController{
		slope:       gradientControllerProto.Slope,
		minGradient: gradientControllerProto.MinGradient,
		maxGradient: gradientControllerProto.MaxGradient,
	}

	controller := NewControllerComponent(gradient, componentIndex, policyReadAPI)

	return controller, fx.Options(), nil
}

// ComputeOutput based on previous and current signal reading.
func (g *GradientController) ComputeOutput(signal, setpoint, controlVariable runtime.Reading, controllerStateReadAPI ControllerStateReadAPI, tickInfo runtime.TickInfo) (runtime.Reading, error) {
	var output runtime.Reading
	if setpoint.Valid() && signal.Valid() && controlVariable.Valid() {
		var gradient float64
		gradient = math.Pow(signal.Value()/setpoint.Value(), g.slope)
		// clamp to min/max
		if gradient < g.minGradient {
			gradient = g.minGradient
		}
		if gradient > g.maxGradient {
			gradient = g.maxGradient
		}
		if math.IsNaN(gradient) {
			output = runtime.InvalidReading()
		} else {
			output = runtime.NewReading(controlVariable.Value() * gradient)
		}
	} else {
		output = runtime.InvalidReading()
	}

	return output, nil
}

// MaintainOutput - Gradient Controller is stateless, so a bump is inevitable.
func (g *GradientController) MaintainOutput(prevSetpoint, currentSetpoint runtime.Reading, _ ControllerStateReadAPI, tickInfo runtime.TickInfo) error {
	return nil
}

// WindOutput - Gradient Controller relies on ControllerComponent to store the last output, returning targetOutput should wind the output.
func (g *GradientController) WindOutput(currentOutput, targetOutput runtime.Reading, _ ControllerStateReadAPI, tickInfo runtime.TickInfo) (runtime.Reading, error) {
	return targetOutput, nil
}
