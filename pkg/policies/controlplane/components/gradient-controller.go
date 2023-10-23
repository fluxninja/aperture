package components

import (
	"fmt"
	"math"

	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
)

// GradientController describes gradient values.
type GradientController struct {
	minGradient       float64
	maxGradient       float64
	slope             float64
	componentID       string
	dynamicConfigKey  string
	defaultManualMode bool
	manualMode        bool
}

// Make sure Gradient complies with Component interface.
var _ runtime.Component = (*GradientController)(nil)

// NewGradientControllerAndOptions creates a Gradient Controller Component and its fx options.
func NewGradientControllerAndOptions(
	gradientControllerProto *policylangv1.GradientController,
	componentID runtime.ComponentID,
	policyReadAPI iface.Policy,
) (runtime.Component, fx.Option, error) {
	gradientParameters := gradientControllerProto.GetParameters()
	// Make sure max is greater than min
	if gradientParameters.MaxGradient < gradientParameters.MinGradient {
		return nil, nil, fmt.Errorf("max_gradient must be greater than min_gradient")
	}

	gradient := &GradientController{
		slope:             gradientParameters.Slope,
		minGradient:       gradientParameters.MinGradient,
		maxGradient:       gradientParameters.MaxGradient,
		componentID:       componentID.String(),
		dynamicConfigKey:  gradientControllerProto.ManualModeConfigKey,
		defaultManualMode: gradientControllerProto.ManualMode,
		manualMode:        gradientControllerProto.ManualMode,
	}

	return gradient, fx.Options(), nil
}

// Name implements runtime.Component.
func (g *GradientController) Name() string {
	return "Gradient"
}

// Type implements runtime.Component.
func (g *GradientController) Type() runtime.ComponentType {
	return runtime.ComponentTypeSignalProcessor
}

// ShortDescription implements runtime.Component.
func (g *GradientController) ShortDescription() string {
	return "Gradient Controller"
}

// IsActuator implements runtime.Component.
func (g *GradientController) IsActuator() bool {
	return false
}

// Execute implements runtime.Component.Execute.
func (g *GradientController) Execute(inPortReadings runtime.PortToReading, circuitAPI runtime.CircuitAPI) (outPortReadings runtime.PortToReading, err error) {
	signal := inPortReadings.ReadSingleReadingPort("signal")
	setpoint := inPortReadings.ReadSingleReadingPort("setpoint")
	max := inPortReadings.ReadSingleReadingPort("max")
	min := inPortReadings.ReadSingleReadingPort("min")
	controlVariable := inPortReadings.ReadSingleReadingPort("control_variable")
	var output runtime.Reading

	// ComputeOutput
	if setpoint.Valid() && setpoint.Value() != 0 && signal.Valid() && controlVariable.Valid() {
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

	if output.Valid() {
		outputReading := output
		// Constrain output
		if max.Valid() {
			outputReading = runtime.NewReading(math.Min(output.Value(), max.Value()))
		}
		if min.Valid() {
			outputReading = runtime.NewReading(math.Max(outputReading.Value(), min.Value()))
		}
		if outputReading.Value() != output.Value() {
			output = outputReading
		}
	}

	// Set output to control variable in-case of Manual mode
	if g.manualMode {
		output = controlVariable
	}

	return runtime.PortToReading{
		"output": []runtime.Reading{output},
	}, nil
}

// DynamicConfigUpdate handles setting of controller.ControllerMode.
func (g *GradientController) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
	g.manualMode = config.GetBoolValue(unmarshaller, g.dynamicConfigKey, g.defaultManualMode)
}
